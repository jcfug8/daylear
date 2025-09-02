package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/lists/list/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	accessMaxPageSize     int32 = 1000
	accessDefaultPageSize int32 = 100
)

var listAccessFieldMap = map[string][]string{
	"name":          {model.ListAccessField_Parent, model.ListAccessField_Id},
	"level":         {model.ListAccessField_PermissionLevel},
	"state":         {model.ListAccessField_State},
	"accept_target": {model.ListAccessField_AcceptTarget},
	"requester":     {model.ListAccessField_Requester},
	"recipient":     {model.ListAccessField_Recipient},
}

func (s *ListService) CreateAccess(ctx context.Context, request *pb.CreateAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Warn().Err(err).Msg("invalid access proto")
		return nil, err
	}

	// convert proto to model
	pbAccess := request.GetAccess()
	pbAccess.Name = ""
	modelAccess, err := s.ProtoToListAccess(pbAccess)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	_, err = s.accessNamer.ParseParent(request.GetParent(), &modelAccess.ListAccessParent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// create access
	modelAccess, err = s.domain.CreateListAccess(ctx, authAccount, modelAccess)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateListAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err = s.ListAccessToProto(modelAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC CreateAccess success")
	return pbAccess, nil
}

func (s *ListService) DeleteAccess(ctx context.Context, request *pb.DeleteAccessRequest) (*emptypb.Empty, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	modelAccess := model.ListAccess{}
	_, err = s.accessNamer.Parse(request.GetName(), &modelAccess)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	err = s.domain.DeleteListAccess(ctx, authAccount, modelAccess.ListAccessParent, modelAccess.ListAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteListAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC DeleteAccess success")
	return &emptypb.Empty{}, nil
}

func (s *ListService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	modelAccess := model.ListAccess{}
	_, err = s.accessNamer.Parse(request.GetName(), &modelAccess)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	modelAccess, err = s.domain.GetListAccess(ctx, authAccount, modelAccess.ListAccessParent, modelAccess.ListAccessId, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetListAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbAccess, err := s.ListAccessToProto(modelAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC GetAccess success")
	return pbAccess, nil
}

func (s *ListService) ListAccesses(ctx context.Context, request *pb.ListAccessesRequest) (*pb.ListAccessesResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListAccesses called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	modelAccessParent := model.ListAccessParent{}
	_, err = s.accessNamer.ParseParent(request.GetParent(), &modelAccessParent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: accessDefaultPageSize,
		MaxPageSize:     accessMaxPageSize,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to setup pagination")
		return nil, err
	}
	request.PageSize = pageSize

	res, err := s.domain.ListListAccesses(ctx, authAccount, modelAccessParent, request.GetPageSize(), pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListListAccesses failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	accesses := make([]*pb.Access, len(res))
	for i, access := range res {
		accessProto, err := s.ListAccessToProto(access)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		accesses[i] = accessProto
	}

	// check field behavior
	for _, accessProto := range accesses {
		grpc.ProcessResponseFieldBehavior(accessProto)
	}

	response := &pb.ListAccessesResponse{
		Accesses: accesses,
	}

	if len(accesses) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListAccesses success")
	return response, nil
}

func (s *ListService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	fieldMask := request.GetUpdateMask()
	updateMask := s.accessFieldMasker.Convert(fieldMask.GetPaths())

	accessProto := request.GetAccess()
	modelAccess, err := s.ProtoToListAccess(accessProto)
	if err != nil {
		log.Error().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.Internal, err.Error())
	}

	modelAccess, err = s.domain.UpdateListAccess(ctx, authAccount, modelAccess, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateListAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessProto, err = s.ListAccessToProto(modelAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC UpdateAccess success")
	return accessProto, nil
}

func (s *ListService) AcceptListAccess(ctx context.Context, request *pb.AcceptListAccessRequest) (*pb.AcceptListAccessResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC AcceptAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	modelAccess := model.ListAccess{}
	_, err = s.accessNamer.Parse(request.GetName(), &modelAccess)
	if err != nil {
		log.Warn().Err(err).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	modelAccess, err = s.domain.AcceptListAccess(ctx, authAccount, modelAccess.ListAccessParent, modelAccess.ListAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.AcceptListAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC AcceptListAccess success")
	return &pb.AcceptListAccessResponse{}, nil
}

func (s *ListService) ProtoToListAccess(pbAccess *pb.Access) (model.ListAccess, error) {
	modelAccess := model.ListAccess{
		PermissionLevel: pbAccess.GetLevel(),
		State:           pbAccess.GetState(),
		AcceptTarget:    pbAccess.GetAcceptTarget(),
	}

	if pbAccess.GetName() != "" {
		_, err := s.accessNamer.Parse(pbAccess.GetName(), &modelAccess)
		if err != nil {
			return model.ListAccess{}, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		}
	}

	switch pbrequester := pbAccess.GetRequester().GetName().(type) {
	case *pb.Access_RequesterOrRecipient_User:
		modelAccess.Requester = model.ListRecipientOrRequester{}
		if pbrequester.User != nil {
			_, err := s.userNamer.Parse(pbrequester.User.Name, &modelAccess.Requester)
			if err != nil {
				return model.ListAccess{}, status.Errorf(codes.InvalidArgument, "invalid requester: %v", err)
			}
		}
	case *pb.Access_RequesterOrRecipient_Circle:
		modelAccess.Requester = model.ListRecipientOrRequester{}
		if pbrequester.Circle != nil {
			_, err := s.circleNamer.Parse(pbrequester.Circle.Name, &modelAccess.Requester)
			if err != nil {
				return model.ListAccess{}, status.Errorf(codes.InvalidArgument, "invalid requester: %v", err)
			}
		}
	}

	switch pbRecipient := pbAccess.GetRecipient().GetName().(type) {
	case *pb.Access_RequesterOrRecipient_User:
		modelAccess.Recipient = model.ListRecipientOrRequester{}
		if pbRecipient.User != nil {
			_, err := s.userNamer.Parse(pbRecipient.User.Name, &modelAccess.Recipient)
			if err != nil {
				return model.ListAccess{}, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", err)
			}
		}
	case *pb.Access_RequesterOrRecipient_Circle:
		modelAccess.Recipient = model.ListRecipientOrRequester{}
		if pbRecipient.Circle != nil {
			_, err := s.circleNamer.Parse(pbRecipient.Circle.Name, &modelAccess.Recipient)
			if err != nil {
				return model.ListAccess{}, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", err)
			}
		}
	}

	return modelAccess, nil
}

func (s *ListService) ListAccessToProto(modelAccess model.ListAccess) (*pb.Access, error) {
	pbAccess := &pb.Access{
		Level:        modelAccess.PermissionLevel,
		State:        modelAccess.State,
		AcceptTarget: modelAccess.AcceptTarget,
	}

	if modelAccess.ListId.ListId != 0 && modelAccess.ListAccessId.ListAccessId != 0 {
		name, err := s.accessNamer.Format(modelAccess)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format access: %v", err)
		}
		pbAccess.Name = name
	}

	if modelAccess.Requester.UserId != 0 {
		userName, err := s.userNamer.Format(modelAccess.Requester)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format requester: %v", err)
		}
		pbAccess.Requester = &pb.Access_RequesterOrRecipient{
			Name: &pb.Access_RequesterOrRecipient_User{
				User: &pb.Access_User{
					Name: userName,
				},
			},
		}
	} else if modelAccess.Requester.CircleId != 0 {
		circleName, err := s.circleNamer.Format(modelAccess.Requester)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format requester: %v", err)
		}
		pbAccess.Requester = &pb.Access_RequesterOrRecipient{
			Name: &pb.Access_RequesterOrRecipient_Circle{
				Circle: &pb.Access_Circle{
					Name: circleName,
				},
			},
		}
	}

	if modelAccess.Recipient.UserId != 0 {
		userName, err := s.userNamer.Format(modelAccess.Recipient)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format recipient: %v", err)
		}
		pbAccess.Recipient = &pb.Access_RequesterOrRecipient{
			Name: &pb.Access_RequesterOrRecipient_User{
				User: &pb.Access_User{
					Name:       userName,
					Username:   modelAccess.RecipientUsername,
					GivenName:  modelAccess.RecipientGivenName,
					FamilyName: modelAccess.RecipientFamilyName,
				},
			},
		}
	} else if modelAccess.Recipient.CircleId != 0 {
		circleName, err := s.circleNamer.Format(modelAccess.Recipient)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format recipient: %v", err)
		}
		pbAccess.Recipient = &pb.Access_RequesterOrRecipient{
			Name: &pb.Access_RequesterOrRecipient_Circle{
				Circle: &pb.Access_Circle{
					Name:   circleName,
					Title:  modelAccess.RecipientCircleTitle,
					Handle: modelAccess.RecipientCircleHandle,
				},
			},
		}
	}

	return pbAccess, nil
}
