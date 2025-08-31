package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/calendars/calendar/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	eventRecipeMaxPageSize     int32 = 1000
	eventRecipeDefaultPageSize int32 = 100
)

var eventRecipeFieldMap = map[string][]string{
	"name":        {model.EventRecipeField_Parent, model.EventRecipeField_EventRecipeId},
	"recipe":      {model.EventRecipeField_RecipeId},
	"create_time": {model.EventRecipeField_CreateTime},
}

// CreateEventRecipe creates a new eventRecipe
func (s *CalendarService) CreateEventRecipe(ctx context.Context, request *pb.CreateEventRecipeRequest) (response *pb.EventRecipe, err error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateEventRecipe called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, err
	}

	// convert proto to model
	eventRecipeProto := request.GetEventRecipe()
	eventRecipeProto.Name = ""
	_, mEventRecipe, err := s.ProtoToEventRecipe(eventRecipeProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	_, err = s.eventRecipeNamer.ParseParent(request.GetParent(), &mEventRecipe.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// create eventRecipe
	mEventRecipe, err = s.domain.CreateEventRecipe(ctx, authAccount, mEventRecipe)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateEventRecipe failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	eventRecipeProto, err = s.EventRecipeToProto(mEventRecipe)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(eventRecipeProto)
	log.Info().Msg("gRPC CreateEventRecipe returning successfully")
	return eventRecipeProto, nil
}

// DeleteEventRecipe deletes a eventRecipe
func (s *CalendarService) DeleteEventRecipe(ctx context.Context, request *pb.DeleteEventRecipeRequest) (*pb.EventRecipe, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteEventRecipe called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mEventRecipe model.EventRecipe
	_, err = s.eventRecipeNamer.Parse(request.GetName(), &mEventRecipe)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mEventRecipe, err = s.domain.DeleteEventRecipe(ctx, authAccount, mEventRecipe.Parent, mEventRecipe.EventRecipeId)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteEventRecipe failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	eventRecipeProto, err := s.EventRecipeToProto(mEventRecipe)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(eventRecipeProto)
	log.Info().Msg("gRPC DeleteEventRecipe returning successfully")
	return eventRecipeProto, nil
}

// GetEventRecipe retrieves a eventRecipe
func (s *CalendarService) GetEventRecipe(ctx context.Context, request *pb.GetEventRecipeRequest) (*pb.EventRecipe, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetEventRecipe called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mEventRecipe model.EventRecipe
	_, err = s.eventRecipeNamer.Parse(request.GetName(), &mEventRecipe)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mEventRecipe, err = s.domain.GetEventRecipe(ctx, authAccount, mEventRecipe.Parent, mEventRecipe.EventRecipeId, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetEventRecipe failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	eventRecipeProto, err := s.EventRecipeToProto(mEventRecipe)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(eventRecipeProto)
	log.Info().Msg("gRPC GetEventRecipe returning successfully")
	return eventRecipeProto, nil
}

// ListEventRecipes lists eventRecipes
func (s *CalendarService) ListEventRecipes(ctx context.Context, request *pb.ListEventRecipesRequest) (*pb.ListEventRecipesResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListEventRecipes called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse parent
	var mEventRecipe model.EventRecipe
	_, err = s.eventRecipeNamer.ParseParent(request.GetParent(), &mEventRecipe)
	if err != nil {
		log.Warn().Err(err).Str("parent", request.GetParent()).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: eventRecipeDefaultPageSize,
		MaxPageSize:     eventRecipeMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}

	// list eventRecipes
	mEventRecipes, err := s.domain.ListEventRecipes(ctx, authAccount, mEventRecipe.Parent, pageSize, pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListEventRecipes failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	eventRecipeProtos := make([]*pb.EventRecipe, len(mEventRecipes))
	for i, mEventRecipe := range mEventRecipes {
		eventRecipeProto, err := s.EventRecipeToProto(mEventRecipe)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		eventRecipeProtos[i] = eventRecipeProto
	}

	// check field behavior
	for _, eventRecipeProto := range eventRecipeProtos {
		grpc.ProcessResponseFieldBehavior(eventRecipeProto)
	}

	// create response
	response := &pb.ListEventRecipesResponse{
		EventRecipes: eventRecipeProtos,
	}

	// add next page token if there are more results
	if len(mEventRecipes) == int(pageSize) {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListEventRecipes returning successfully")
	return response, nil
}

// ProtoToEventRecipe converts a proto EventRecipe to a model EventRecipe
func (s *CalendarService) ProtoToEventRecipe(proto *pb.EventRecipe) (nameIndex int, eventRecipe model.EventRecipe, err error) {
	eventRecipe = model.EventRecipe{
		CreateTime: proto.GetCreateTime().AsTime(),
	}

	if proto.GetRecipe() != "" {
		nameIndex, err = s.recipeNamer.Parse(proto.GetRecipe(), &eventRecipe)
		if err != nil {
			return 0, model.EventRecipe{}, err
		}
	}

	// Parse parent from name if provided
	if proto.GetName() != "" {
		nameIndex, err = s.eventRecipeNamer.Parse(proto.GetName(), &eventRecipe)
		if err != nil {
			return 0, model.EventRecipe{}, err
		}
	}

	return nameIndex, eventRecipe, nil
}

// EventRecipeToProto converts a model EventRecipe to a proto EventRecipe
func (s *CalendarService) EventRecipeToProto(eventRecipe model.EventRecipe, options ...namer.FormatReflectNamerOption) (*pb.EventRecipe, error) {
	proto := &pb.EventRecipe{
		CreateTime: timestamppb.New(eventRecipe.CreateTime),
	}

	if eventRecipe.RecipeId.RecipeId != 0 {
		name, err := s.recipeNamer.Format(eventRecipe, options...)
		if err != nil {
			return nil, err
		}
		proto.Recipe = name
	}

	// Generate name
	if eventRecipe.EventRecipeId.EventRecipeId != 0 {
		name, err := s.eventRecipeNamer.Format(eventRecipe, options...)
		if err != nil {
			return nil, err
		}
		proto.Name = name
	}

	return proto, nil
}
