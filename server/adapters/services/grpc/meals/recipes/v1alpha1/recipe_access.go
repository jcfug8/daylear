package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	accessMaxPageSize     int32 = 1000
	accessDefaultPageSize int32 = 100
)

func (s *RecipeService) CreateAccess(ctx context.Context, request *pb.CreateAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		return nil, err
	}

	// convert proto to model
	pbAccess := request.GetAccess()
	pbAccess.Name = ""
	modelAccess, err := s.ProtoToRecipeAccess(pbAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid access: %v", err)
	}

	// parse parent
	_, err = s.accessNamer.ParseParent(request.Parent, &modelAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// create access
	createdAccess, err := s.domain.CreateRecipeAccess(ctx, authAccount, modelAccess)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err = s.RecipeAccessToProto(createdAccess)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)

	return pbAccess, nil
}

func (s *RecipeService) DeleteAccess(ctx context.Context, request *pb.DeleteAccessRequest) (*emptypb.Empty, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	// parse name
	recipeAccess := &model.RecipeAccess{}
	_, err = s.accessNamer.Parse(request.Name, recipeAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.Name)
	}

	// delete access
	err = s.domain.DeleteRecipeAccess(ctx, authAccount, recipeAccess.RecipeAccessParent, recipeAccess.RecipeAccessId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *RecipeService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	// parse name
	var recipeAccess model.RecipeAccess
	_, err = s.accessNamer.Parse(request.Name, &recipeAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.Name)
	}

	// get access
	access, err := s.domain.GetRecipeAccess(ctx, authAccount, recipeAccess.RecipeAccessParent, recipeAccess.RecipeAccessId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err := s.RecipeAccessToProto(access)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)

	return pbAccess, nil
}

func (s *RecipeService) ListAccesses(ctx context.Context, request *pb.ListAccessesRequest) (*pb.ListAccessesResponse, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	// parse parent
	var recipeAccessParent model.RecipeAccessParent
	_, err = s.accessNamer.ParseParent(request.Parent, &recipeAccessParent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.Parent)
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: accessDefaultPageSize,
		MaxPageSize:     accessMaxPageSize,
	})
	if err != nil {
		return nil, err
	}
	request.PageSize = pageSize

	// list accesses
	accesses, err := s.domain.ListRecipeAccesses(ctx, authAccount, recipeAccessParent, request.GetPageSize(), pageToken.Offset, request.GetFilter())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	pbAccesses := make([]*pb.Access, len(accesses))
	for i, access := range accesses {
		pbAccess, err := s.RecipeAccessToProto(access)
		if err != nil {
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

	return response, nil
}

func (s *RecipeService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		return nil, err
	}

	// convert proto to model
	modelAccess, err := s.ProtoToRecipeAccess(request.Access)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	// TODO: update mask

	// update access
	updatedAccess, err := s.domain.UpdateRecipeAccess(ctx, authAccount, modelAccess)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err := s.RecipeAccessToProto(updatedAccess)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)

	return pbAccess, nil
}

func (s *RecipeService) AcceptRecipeAccess(ctx context.Context, request *pb.AcceptRecipeAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		return nil, err
	}

	// parse name to get the parent and access id
	var access model.RecipeAccess
	_, err = s.accessNamer.Parse(request.GetName(), &access)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	// accept the access
	updatedAccess, err := s.domain.AcceptRecipeAccess(ctx, authAccount, access.RecipeAccessParent, access.RecipeAccessId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	pbAccess, err := s.RecipeAccessToProto(updatedAccess)
	if err != nil {
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(pbAccess)

	return pbAccess, nil
}

// UTILS

func (s *RecipeService) ProtoToRecipeAccess(pbAccess *pb.Access) (model.RecipeAccess, error) {
	modelAccess := model.RecipeAccess{
		Level: pbAccess.GetLevel(),
		State: pbAccess.GetState(),
	}

	if pbAccess.GetName() != "" {
		_, err := s.accessNamer.Parse(pbAccess.GetName(), &modelAccess)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		}
	}

	switch pbrequester := pbAccess.GetRequester().GetName().(type) {
	case *pb.Access_RequesterOrRecipient_User:
		modelAccess.Requester = model.AuthAccount{}
		_, err := s.userNamer.Parse(pbrequester.User, &modelAccess.Requester)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid requester: %v", err)
		}
	case *pb.Access_RequesterOrRecipient_Circle:
		modelAccess.Requester = model.AuthAccount{}
		_, err := s.circleNamer.Parse(pbrequester.Circle, &modelAccess.Requester)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid requester: %v", err)
		}
	}

	switch pbRecipient := pbAccess.GetRecipient().GetName().(type) {
	case *pb.Access_RequesterOrRecipient_User:
		modelAccess.Recipient = model.AuthAccount{}
		_, err := s.userNamer.Parse(pbRecipient.User, &modelAccess.Recipient)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", err)
		}
	case *pb.Access_RequesterOrRecipient_Circle:
		modelAccess.Recipient = model.AuthAccount{}
		_, err := s.circleNamer.Parse(pbRecipient.Circle, &modelAccess.Recipient)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", err)
		}
	}

	return modelAccess, nil
}

func (s *RecipeService) RecipeAccessToProto(modelAccess model.RecipeAccess) (*pb.Access, error) {
	pbAccess := &pb.Access{
		Level: modelAccess.Level,
		State: modelAccess.State,
	}

	if modelAccess.RecipeId.RecipeId != 0 && modelAccess.RecipeAccessId.RecipeAccessId != 0 {
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
				User: userName,
			},
		}
	} else if modelAccess.Requester.CircleId != 0 {
		circleName, err := s.circleNamer.Format(modelAccess.Requester)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format requester: %v", err)
		}
		pbAccess.Requester = &pb.Access_RequesterOrRecipient{
			Name: &pb.Access_RequesterOrRecipient_Circle{
				Circle: circleName,
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
				User: userName,
			},
		}
	} else if modelAccess.Recipient.CircleId != 0 {
		circleName, err := s.circleNamer.Format(modelAccess.Recipient)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format recipient: %v", err)
		}
		pbAccess.Recipient = &pb.Access_RequesterOrRecipient{
			Name: &pb.Access_RequesterOrRecipient_Circle{
				Circle: circleName,
			},
		}
	}

	return pbAccess, nil
}
