package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	accessMaxPageSize     int32 = 1000
	accessDefaultPageSize int32 = 100
)

var calendarAccessFieldMap = map[string][]string{
	"name":      {model.CalendarAccessField_Parent, model.CalendarAccessField_Id},
	"level":     {model.CalendarAccessField_PermissionLevel},
	"state":     {model.CalendarAccessField_State},
	"requester": {model.CalendarAccessField_Requester},
	"recipient": {model.CalendarAccessField_Recipient},
}

// CreateAccess creates a new calendar access
func (s *CalendarService) CreateAccess(ctx context.Context, request *pb.CreateAccessRequest) (*pb.Access, error) {
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
	modelAccess, err := s.ProtoToCalendarAccess(pbAccess)
	if err != nil {
		log.Warn().Err(err).Msg("invalid access proto")
		return nil, status.Errorf(codes.InvalidArgument, "invalid access: %v", err)
	}

	// parse parent
	_, err = s.calendarAccessNamer.ParseParent(request.Parent, &modelAccess)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// create access
	createdAccess, err := s.domain.CreateCalendarAccess(ctx, authAccount, modelAccess)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateCalendarAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err = s.CalendarAccessToProto(createdAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC CreateAccess success")
	return pbAccess, nil
}

// DeleteAccess deletes a calendar access
func (s *CalendarService) DeleteAccess(ctx context.Context, request *pb.DeleteAccessRequest) (*emptypb.Empty, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse name
	calendarAccess := &model.CalendarAccess{}
	_, err = s.calendarAccessNamer.Parse(request.Name, calendarAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", request.Name).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.Name)
	}

	// delete access
	err = s.domain.DeleteCalendarAccess(ctx, authAccount, calendarAccess.CalendarAccessParent, calendarAccess.CalendarAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteCalendarAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC DeleteAccess returning successfully")
	return &emptypb.Empty{}, nil
}

// GetAccess retrieves a calendar access
func (s *CalendarService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse name
	var calendarAccess model.CalendarAccess
	_, err = s.calendarAccessNamer.Parse(request.Name, &calendarAccess)
	if err != nil {
		log.Warn().Err(err).Str("name", request.Name).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.Name)
	}

	// get access
	access, err := s.domain.GetCalendarAccess(ctx, authAccount, calendarAccess.CalendarAccessParent, calendarAccess.CalendarAccessId, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetCalendarAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err := s.CalendarAccessToProto(access)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC GetAccess returning successfully")
	return pbAccess, nil
}

// ListAccesses lists calendar accesses
func (s *CalendarService) ListAccesses(ctx context.Context, request *pb.ListAccessesRequest) (*pb.ListAccessesResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListAccesses called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse parent
	var calendarAccessParent model.CalendarAccessParent
	_, err = s.calendarAccessNamer.ParseParent(request.Parent, &calendarAccessParent)
	if err != nil {
		log.Warn().Err(err).Str("parent", request.Parent).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.Parent)
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: accessDefaultPageSize,
		MaxPageSize:     accessMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}
	request.PageSize = pageSize

	// list accesses
	accesses, err := s.domain.ListCalendarAccesses(ctx, authAccount, calendarAccessParent, request.GetPageSize(), pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListCalendarAccesses failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	pbAccesses := make([]*pb.Access, len(accesses))
	for i, access := range accesses {
		pbAccess, err := s.CalendarAccessToProto(access)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		pbAccesses[i] = pbAccess
	}

	// check field behavior
	for _, pbAccess := range pbAccesses {
		grpc.ProcessResponseFieldBehavior(pbAccess)
	}

	response := &pb.ListAccessesResponse{
		Accesses: pbAccesses,
	}

	if len(pbAccesses) > 0 {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListAccesses returning successfully")
	return response, nil
}

// UpdateAccess updates a calendar access
func (s *CalendarService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessUpdateRequestFieldBehavior(request)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, err
	}

	// convert proto to model
	modelAccess, err := s.ProtoToCalendarAccess(request.Access)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	updateMask := s.calendarAccessFieldMasker.Convert(request.GetUpdateMask().GetPaths())

	// update access
	updatedAccess, err := s.domain.UpdateCalendarAccess(ctx, authAccount, modelAccess, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateCalendarAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err := s.CalendarAccessToProto(updatedAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC UpdateAccess returning successfully")
	return pbAccess, nil
}

// AcceptAccess accepts a calendar access
func (s *CalendarService) AcceptAccess(ctx context.Context, request *pb.AcceptAccessRequest) (*pb.Access, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC AcceptAccess called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse name to get the parent and access id
	var access model.CalendarAccess
	_, err = s.calendarAccessNamer.Parse(request.GetName(), &access)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	acceptedAccess, err := s.domain.AcceptCalendarAccess(ctx, authAccount, access.CalendarAccessParent, access.CalendarAccessId)
	if err != nil {
		log.Error().Err(err).Msg("domain.AcceptCalendarAccess failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err := s.CalendarAccessToProto(acceptedAccess)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)
	log.Info().Msg("gRPC AcceptAccess returning successfully")
	return pbAccess, nil
}

// UTILS

// ProtoToCalendarAccess converts a proto Access to a model CalendarAccess
func (s *CalendarService) ProtoToCalendarAccess(pbAccess *pb.Access) (model.CalendarAccess, error) {
	modelAccess := model.CalendarAccess{
		PermissionLevel: pbAccess.GetLevel(),
		State:           pbAccess.GetState(),
	}

	if pbAccess.GetName() != "" {
		_, err := s.calendarAccessNamer.Parse(pbAccess.GetName(), &modelAccess)
		if err != nil {
			return model.CalendarAccess{}, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		}
	}

	// Handle requester (only user is supported for calendars)
	switch pbrequester := pbAccess.GetRequester().GetName().(type) {
	case *pb.Access_RequesterOrRecipient_User:
		modelAccess.Requester = model.CalendarRecipientOrRequester{}
		if pbrequester.User != nil {
			modelAccess.Requester = model.CalendarRecipientOrRequester{}
			_, err := s.userNamer.Parse(pbrequester.User.Name, &modelAccess.Requester)
			if err != nil {
				return model.CalendarAccess{}, status.Errorf(codes.InvalidArgument, "invalid requester: %v", err)
			}
		}
	case *pb.Access_RequesterOrRecipient_Circle:
		modelAccess.Requester = model.CalendarRecipientOrRequester{}
		if pbrequester.Circle != nil {
			_, err := s.circleNamer.Parse(pbrequester.Circle.Name, &modelAccess.Requester)
			if err != nil {
				return model.CalendarAccess{}, status.Errorf(codes.InvalidArgument, "invalid requester: %v", err)
			}
		}
	}

	// Handle recipient (only user is supported for calendars)
	switch pbRecipient := pbAccess.GetRecipient().GetName().(type) {
	case *pb.Access_RequesterOrRecipient_User:
		modelAccess.Recipient = model.CalendarRecipientOrRequester{}
		if pbRecipient.User != nil {
			modelAccess.Recipient = model.CalendarRecipientOrRequester{}
			_, err := s.userNamer.Parse(pbRecipient.User.Name, &modelAccess.Recipient)
			if err != nil {
				return model.CalendarAccess{}, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", err)
			}
		}
	case *pb.Access_RequesterOrRecipient_Circle:
		modelAccess.Recipient = model.CalendarRecipientOrRequester{}
		if pbRecipient.Circle != nil {
			_, err := s.circleNamer.Parse(pbRecipient.Circle.Name, &modelAccess.Recipient)
			if err != nil {
				return model.CalendarAccess{}, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", err)
			}
		}
	}

	return modelAccess, nil
}

// CalendarAccessToProto converts a model CalendarAccess to a proto Access
func (s *CalendarService) CalendarAccessToProto(modelAccess model.CalendarAccess) (*pb.Access, error) {
	pbAccess := &pb.Access{
		Level: modelAccess.PermissionLevel,
		State: modelAccess.State,
	}

	if modelAccess.CalendarId != 0 && modelAccess.CalendarAccessId.CalendarAccessId != 0 {
		name, err := s.calendarAccessNamer.Format(modelAccess)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format access: %v", err)
		}
		pbAccess.Name = name
	}

	// Handle requester (only user is supported for calendars)
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

	// Handle recipient (only user is supported for calendars)
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
