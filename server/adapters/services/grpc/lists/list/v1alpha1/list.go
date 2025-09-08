package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/lists/list/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	listMaxPageSize     int32 = 1000
	listDefaultPageSize int32 = 100
)

var listFieldMap = map[string][]string{
	"name":           {model.ListField_Parent, model.ListField_Id},
	"title":          {model.ListField_Title},
	"description":    {model.ListField_Description},
	"show_completed": {model.ListField_ShowCompleted},
	"visibility":     {model.ListField_VisibilityLevel},
	"sections":       {model.ListField_Sections},
	"create_time":    {model.ListField_CreateTime},
	"update_time":    {model.ListField_UpdateTime},
	"favorited":      {model.ListField_Favorited},

	"list_access": {model.ListField_ListAccess},
}

// CreateList -
func (s *ListService) CreateList(ctx context.Context, request *pb.CreateListRequest) (response *pb.List, err error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateList called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to process request field behavior")
		return nil, err
	}

	// convert proto to model
	pbList := request.GetList()
	pbList.Name = ""
	_, mList, err := s.ProtoToList(pbList)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	_, err = s.listNamer.ParseParent(request.GetParent(), &mList.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// create list
	mList, err = s.domain.CreateList(ctx, authAccount, mList)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateList failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbList, err = s.ListToProto(mList)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbList)
	log.Info().Msg("gRPC CreateList success")
	return pbList, nil
}

// DeleteList -
func (s *ListService) DeleteList(ctx context.Context, request *pb.DeleteListRequest) (*pb.List, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteList called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mList := model.List{}
	_, err = s.listNamer.Parse(request.GetName(), &mList)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	err = s.domain.DeleteList(ctx, authAccount, mList.Parent, mList.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteList failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Get the deleted list to return it
	mList, err = s.domain.GetList(ctx, authAccount, mList.Parent, mList.Id, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetList failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbList, err := s.ListToProto(mList)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC DeleteList success")
	return pbList, nil
}

// GetList -
func (s *ListService) GetList(ctx context.Context, request *pb.GetListRequest) (*pb.List, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetList called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mList := model.List{}
	_, err = s.listNamer.Parse(request.GetName(), &mList)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mList, err = s.domain.GetList(ctx, authAccount, mList.Parent, mList.Id, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetList failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbList, err := s.ListToProto(mList)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC GetList success")
	return pbList, nil
}

// UpdateList -
func (s *ListService) UpdateList(ctx context.Context, request *pb.UpdateListRequest) (*pb.List, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateList called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	fieldMask := request.GetUpdateMask()
	updateMask := s.listFieldMasker.Convert(fieldMask.GetPaths())

	listProto := request.GetList()
	_, mList, err := s.ProtoToList(listProto)
	if err != nil {
		log.Error().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.Internal, err.Error())
	}

	mList, err = s.domain.UpdateList(ctx, authAccount, mList, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateList failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	listProto, err = s.ListToProto(mList)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC UpdateList success")
	return listProto, nil
}

// ListLists -
func (s *ListService) ListLists(ctx context.Context, request *pb.ListListsRequest) (*pb.ListListsResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListLists called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mListParent := model.ListParent{}
	_, err = s.listNamer.ParseParent(request.GetParent(), &mListParent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: listDefaultPageSize,
		MaxPageSize:     listMaxPageSize,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to setup pagination")
		return nil, err
	}
	request.PageSize = pageSize

	res, err := s.domain.ListLists(ctx, authAccount, mListParent, request.GetPageSize(), int32(pageToken.Offset), request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListLists failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	lists := make([]*pb.List, len(res))
	for i, list := range res {
		listProto, err := s.ListToProto(list)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		lists[i] = listProto
	}

	// check field behavior
	for _, listProto := range lists {
		grpc.ProcessResponseFieldBehavior(listProto)
	}

	response := &pb.ListListsResponse{
		Lists: lists,
	}

	if len(lists) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListLists success")
	return response, nil
}

// FavoriteList favorites a list for the authenticated user.
func (s *ListService) FavoriteList(ctx context.Context, request *pb.FavoriteListRequest) (*pb.FavoriteListResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC FavoriteList called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// Parse the list name to get the ID
	list := model.List{}
	_, err = s.listNamer.Parse(request.GetName(), &list)
	if err != nil {
		log.Warn().Err(err).Msg("invalid list name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid list name: %v", request.GetName())
	}

	// Call domain method to favorite the list
	err = s.domain.FavoriteList(ctx, authAccount, list.Parent, list.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.FavoriteList failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC FavoriteList success")
	return &pb.FavoriteListResponse{}, nil
}

// UnfavoriteList removes a list from the authenticated user's favorites.
func (s *ListService) UnfavoriteList(ctx context.Context, request *pb.UnfavoriteListRequest) (*pb.UnfavoriteListResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UnfavoriteList called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// Parse the list name to get the ID
	list := model.List{}
	_, err = s.listNamer.Parse(request.GetName(), &list)
	if err != nil {
		log.Warn().Err(err).Msg("invalid list name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid list name: %v", request.GetName())
	}

	// Call domain method to unfavorite the list
	err = s.domain.UnfavoriteList(ctx, authAccount, list.Parent, list.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.UnfavoriteList failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC UnfavoriteList success")
	return &pb.UnfavoriteListResponse{}, nil
}

// ProtoToList converts a protobuf List to a model List
func (s *ListService) ProtoToList(proto *pb.List) (int, model.List, error) {
	list := model.List{}
	var nameIndex int
	var err error
	if proto.Name != "" {
		nameIndex, err = s.listNamer.Parse(proto.Name, &list)
		if err != nil {
			return nameIndex, list, err
		}
	}

	list.Title = proto.Title
	list.Description = proto.Description
	list.ShowCompleted = proto.ShowCompleted
	list.VisibilityLevel = proto.Visibility
	list.Sections, err = s.ProtosToSections(proto.Sections)
	if err != nil {
		return nameIndex, list, err
	}
	if proto.CreateTime != nil {
		list.CreateTime = proto.CreateTime.AsTime()
	}
	if proto.UpdateTime != nil {
		list.UpdateTime = proto.UpdateTime.AsTime()
	}

	return nameIndex, list, nil
}

// ListToProto converts a model List to a protobuf List
func (s *ListService) ListToProto(list model.List) (*pb.List, error) {
	proto := &pb.List{}

	if list.Id.ListId != 0 {
		name, err := s.listNamer.Format(list)
		if err != nil {
			return proto, err
		}
		proto.Name = name
	}

	proto.Title = list.Title
	proto.Description = list.Description
	proto.ShowCompleted = list.ShowCompleted
	proto.Visibility = list.VisibilityLevel
	sections, err := s.SectionsToProtos(list.Id, list.Sections)
	if err != nil {
		return proto, err
	}
	proto.Sections = sections
	proto.CreateTime = timestamppb.New(list.CreateTime)
	proto.UpdateTime = timestamppb.New(list.UpdateTime)
	proto.Favorited = list.Favorited

	// Handle list_access field if present
	if (list.ListAccess != model.ListAccess{}) {
		name, err := s.accessNamer.Format(list.ListAccess)
		if err == nil {
			proto.ListAccess = &pb.List_ListAccess{
				Name:            name,
				PermissionLevel: list.ListAccess.PermissionLevel,
				State:           list.ListAccess.State,
				AcceptTarget:    list.ListAccess.AcceptTarget,
			}
		}
	}

	return proto, nil
}

// ProtoToSection converts a proto ListSection to a domain ListSection.
func (s *ListService) ProtoToSection(proto *pb.List_ListSection) (model.ListSection, error) {
	section := model.ListSection{}

	section.Title = proto.Title
	if proto.Name != "" {
		_, err := s.listSectionNamer.Parse(proto.Name, &section)
		if err != nil {
			return section, err
		}
	}
	return section, nil
}

// SectionToProto converts a domain ListSection to a proto ListSection.
func (s *ListService) SectionToProto(section model.ListSection) (*pb.List_ListSection, error) {
	proto := &pb.List_ListSection{
		Title: section.Title,
	}
	if section.Id != 0 {
		name, err := s.listSectionNamer.Format(section, namer.AsPatternIndex(-1))
		if err != nil {
			return proto, err
		}
		proto.Name = name
	}
	return proto, nil
}

// ProtosToSections converts a slice of proto ListSections to a slice of domain ListSections.
func (s *ListService) ProtosToSections(protos []*pb.List_ListSection) ([]model.ListSection, error) {
	sections := make([]model.ListSection, len(protos))
	for i, proto := range protos {
		section, err := s.ProtoToSection(proto)
		if err != nil {
			return nil, err
		}
		sections[i] = section
	}
	return sections, nil
}

// SectionsToProtos converts a slice of domain ListSections to a slice of proto ListSections.
func (s *ListService) SectionsToProtos(listId model.ListId, sections []model.ListSection) ([]*pb.List_ListSection, error) {
	protos := make([]*pb.List_ListSection, len(sections))
	for i, section := range sections {
		section.ListId = listId.ListId
		proto, err := s.SectionToProto(section)
		if err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}
