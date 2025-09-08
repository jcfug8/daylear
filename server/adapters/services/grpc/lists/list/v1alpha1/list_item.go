package v1alpha1

import (
	"context"
	"time"

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
	listItemMaxPageSize     int32 = 1000
	listItemDefaultPageSize int32 = 100
)

var listItemFieldMap = map[string][]string{
	"name":            {model.ListItemField_Parent, model.ListItemField_Id},
	"title":           {model.ListItemField_Title},
	"points":          {model.ListItemField_Points},
	"recurrence_rule": {model.ListItemField_RecurrenceRule},
	"list_section":    {model.ListItemField_ListSectionId},
	"create_time":     {model.ListItemField_CreateTime},
	"update_time":     {model.ListItemField_UpdateTime},
}

// CreateListItem creates a new list item
func (s *ListService) CreateListItem(ctx context.Context, request *pb.CreateListItemRequest) (*pb.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateListItem called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// Check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to process request field behavior")
		return nil, err
	}

	// Parse parent
	var mListItem model.ListItem
	_, err = s.listItemNamer.ParseParent(request.GetParent(), &mListItem.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// Convert proto to model
	mListItem, err = s.ListItemFromProto(request.GetListItem(), mListItem.Parent)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid list item data")
	}

	// Create list item
	dbListItem, err := s.domain.CreateListItem(ctx, authAccount, mListItem)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateListItem failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert model to proto
	pbListItem, err := s.ListItemToProto(dbListItem)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// Check field behavior
	grpc.ProcessResponseFieldBehavior(pbListItem)

	log.Info().Msg("gRPC CreateListItem success")
	return pbListItem, nil
}

// GetListItem retrieves a list item
func (s *ListService) GetListItem(ctx context.Context, request *pb.GetListItemRequest) (*pb.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetListItem called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// Parse name to get parent and ID
	var mListItem model.ListItem
	_, err = s.listItemNamer.Parse(request.GetName(), &mListItem)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	// Get list item
	dbListItem, err := s.domain.GetListItem(ctx, authAccount, mListItem.Parent, mListItem.Id, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetListItem failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert model to proto
	pbListItem, err := s.ListItemToProto(dbListItem)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// Check field behavior
	grpc.ProcessResponseFieldBehavior(pbListItem)

	log.Info().Msg("gRPC GetListItem success")
	return pbListItem, nil
}

// ListListItems lists list items with pagination and filtering
func (s *ListService) ListListItems(ctx context.Context, request *pb.ListListItemsRequest) (*pb.ListListItemsResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListListItems called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// Parse parent
	var mListItem model.ListItem
	_, err = s.listItemNamer.ParseParent(request.GetParent(), &mListItem.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// Setup pagination
	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: listItemDefaultPageSize,
		MaxPageSize:     listItemMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}

	// List list items
	dbListItems, err := s.domain.ListListItems(ctx, authAccount, mListItem.Parent, pageSize, int32(pageToken.Offset), request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListListItems failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert models to protos
	listItemProtos := make([]*pb.ListItem, len(dbListItems))
	for i, dbListItem := range dbListItems {
		listItemProto, err := s.ListItemToProto(dbListItem)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		listItemProtos[i] = listItemProto
	}

	// Check field behavior
	for _, listItemProto := range listItemProtos {
		grpc.ProcessResponseFieldBehavior(listItemProto)
	}

	// Create response
	response := &pb.ListListItemsResponse{
		ListItems: listItemProtos,
	}

	// Add next page token if there are more results
	if len(dbListItems) == int(pageSize) {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListListItems success")
	return response, nil
}

// UpdateListItem updates an existing list item
func (s *ListService) UpdateListItem(ctx context.Context, request *pb.UpdateListItemRequest) (*pb.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateListItem called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// Check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Error().Err(err).Msg("failed to process request field behavior")
		return nil, err
	}

	// Parse name to get parent and ID
	var mListItem model.ListItem
	_, err = s.listItemNamer.Parse(request.GetListItem().GetName(), &mListItem)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetListItem().GetName())
	}

	// Convert proto to model
	mListItem, err = s.ListItemFromProto(request.GetListItem(), mListItem.Parent)
	if err != nil {
		log.Error().Err(err).Msg("failed to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid list item data")
	}

	// Process field mask
	fieldMask := request.GetUpdateMask()
	fields := s.listItemFieldMasker.Convert(fieldMask.GetPaths())

	// Update list item
	dbListItem, err := s.domain.UpdateListItem(ctx, authAccount, mListItem, fields)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateListItem failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert model to proto
	pbListItem, err := s.ListItemToProto(dbListItem)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// Check field behavior
	grpc.ProcessResponseFieldBehavior(pbListItem)

	log.Info().Msg("gRPC UpdateListItem success")
	return pbListItem, nil
}

// DeleteListItem deletes a list item
func (s *ListService) DeleteListItem(ctx context.Context, request *pb.DeleteListItemRequest) (*pb.ListItem, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteListItem called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// Parse name to get parent and ID
	var mListItem model.ListItem
	_, err = s.listItemNamer.Parse(request.GetName(), &mListItem)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	// Delete list item
	dbListItem, err := s.domain.DeleteListItem(ctx, authAccount, mListItem.Parent, mListItem.Id)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteListItem failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert model to proto
	pbListItem, err := s.ListItemToProto(dbListItem)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// Check field behavior
	grpc.ProcessResponseFieldBehavior(pbListItem)

	log.Info().Msg("gRPC DeleteListItem success")
	return pbListItem, nil
}

// ListItemToProto converts a domain model to a proto message
func (s *ListService) ListItemToProto(mListItem model.ListItem) (*pb.ListItem, error) {
	pbListItem := &pb.ListItem{
		Title:          mListItem.Title,
		Points:         mListItem.Points,
		RecurrenceRule: mListItem.RecurrenceRule,
		CreateTime:     timestamppb.New(mListItem.CreateTime),
		UpdateTime:     timestamppb.New(mListItem.UpdateTime),
	}

	// Generate the name using the namer if ID is set
	if mListItem.Id.ListItemId != 0 {
		name, err := s.listItemNamer.Format(mListItem)
		if err != nil {
			return nil, err
		}
		pbListItem.Name = name
	}

	if mListItem.ListSectionId != 0 {
		name, err := s.listSectionNamer.Format(mListItem, namer.AsPatternIndex(-1))
		if err != nil {
			return nil, err
		}
		pbListItem.ListSection = name
	}

	return pbListItem, nil
}

// ListItemFromProto converts a proto message to a domain model
func (s *ListService) ListItemFromProto(pbListItem *pb.ListItem, parent model.ListItemParent) (model.ListItem, error) {
	mListItem := model.ListItem{
		Parent:         parent,
		Title:          pbListItem.GetTitle(),
		Points:         pbListItem.GetPoints(),
		RecurrenceRule: pbListItem.GetRecurrenceRule(),
		CreateTime:     time.Now(), // Will be set by database
		UpdateTime:     time.Now(), // Will be set by database
	}

	// If this is an update operation, parse the ID from the name
	if pbListItem.GetName() != "" {
		_, err := s.listItemNamer.Parse(pbListItem.GetName(), &mListItem)
		if err != nil {
			return model.ListItem{}, err
		}
	}

	if pbListItem.GetListSection() != "" {
		_, err := s.listSectionNamer.Parse(pbListItem.GetListSection(), &mListItem)
		if err != nil {
			return model.ListItem{}, err
		}
	}

	mListItem.Parent = parent

	return mListItem, nil
}
