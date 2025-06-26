package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"go.einride.tech/aip/fieldbehavior"
	"go.einride.tech/aip/pagination"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	accessMaxPageSize     int32 = 1000
	accessDefaultPageSize int32 = 100
)

func (s *RecipeService) CreateAccess(ctx context.Context, request *pb.CreateAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// check field behavior
	fieldbehavior.ClearFields(request, annotations.FieldBehavior_OUTPUT_ONLY)
	err = fieldbehavior.ValidateRequiredFields(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	// convert proto to model
	pbAccess := request.GetAccess()
	pbAccess.Name = ""
	modelAccess, err := s.convertProtoToModel(pbAccess)
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
		return nil, status.Errorf(codes.Internal, "failed to create access: %v", err)
	}

	// convert model to proto
	pbAccess, err = s.convertModelToProto(createdAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
	}

	// check field behavior
	fieldbehavior.ClearFields(pbAccess, annotations.FieldBehavior_INPUT_ONLY)

	return pbAccess, nil
}

func (s *RecipeService) DeleteAccess(ctx context.Context, request *pb.DeleteAccessRequest) (*emptypb.Empty, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// parse name
	recipeAccess := &model.RecipeAccess{}

	_, err = s.accessNamer.Parse(request.Name, recipeAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// delete access
	err = s.domain.DeleteRecipeAccess(ctx, authAccount, recipeAccess.RecipeAccessParent, recipeAccess.RecipeAccessId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete access: %v", err)
	}

	return nil, nil
}

func (s *RecipeService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// set requester
	var recipeAccess model.RecipeAccess

	// parse name
	_, err = s.accessNamer.Parse(request.Name, &recipeAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// get access
	access, err := s.domain.GetRecipeAccess(ctx, authAccount, recipeAccess.RecipeAccessParent, recipeAccess.RecipeAccessId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get access: %v", err)
	}

	// convert model to proto
	pbAccess, err := s.convertModelToProto(access)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
	}

	// check field behavior
	fieldbehavior.ClearFields(pbAccess, annotations.FieldBehavior_INPUT_ONLY)

	return pbAccess, nil
}

func (s *RecipeService) ListAccesses(ctx context.Context, request *pb.ListAccessesRequest) (*pb.ListAccessesResponse, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// set requester
	var recipeAccessParent model.RecipeAccessParent

	// parse parent
	_, err = s.accessNamer.ParseParent(request.Parent, &recipeAccessParent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// pagination
	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	if request.GetPageSize() == 0 {
		request.PageSize = accessDefaultPageSize
	}
	request.PageSize = min(request.PageSize, accessMaxPageSize)

	// list accesses
	accesses, err := s.domain.ListRecipeAccesses(ctx, authAccount, recipeAccessParent, request.GetPageSize(), pageToken.Offset, request.GetFilter())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list accesses: %v", err)
	}

	// convert models to protos
	pbAccesses := make([]*pb.Access, len(accesses))
	for i, access := range accesses {
		pbAccess, err := s.convertModelToProto(access)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
		}
		pbAccesses[i] = pbAccess
	}

	// return response
	return &pb.ListAccessesResponse{
		NextPageToken: pageToken.Next(request).String(),
		Accesses:      pbAccesses,
	}, nil
}

func (s *RecipeService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// check field behavior
	fieldbehavior.ClearFields(request, annotations.FieldBehavior_OUTPUT_ONLY)
	err = fieldbehavior.ValidateRequiredFields(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	// convert proto to model
	modelAccess, err := s.convertProtoToModel(request.Access)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid access: %v", err)
	}

	// update access
	updatedAccess, err := s.domain.UpdateRecipeAccess(ctx, authAccount, modelAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update access: %v", err)
	}

	// convert model to proto
	pbAccess, err := s.convertModelToProto(updatedAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
	}

	// check field behavior
	fieldbehavior.ClearFields(pbAccess, annotations.FieldBehavior_INPUT_ONLY)

	return pbAccess, nil
}

func (s *RecipeService) AcceptRecipeAccess(ctx context.Context, request *pb.AcceptRecipeAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// check field behavior
	err = fieldbehavior.ValidateRequiredFields(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	// create access parent with requester (user accepting the access)
	var access model.RecipeAccess

	// parse name to get the parent and access id
	_, err = s.accessNamer.Parse(request.GetName(), &access)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// accept the access
	updatedAccess, err := s.domain.AcceptRecipeAccess(ctx, authAccount, access.RecipeAccessParent, access.RecipeAccessId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to accept access: %v", err)
	}

	// convert model to proto
	pbAccess, err := s.convertModelToProto(updatedAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
	}

	// check field behavior
	fieldbehavior.ClearFields(pbAccess, annotations.FieldBehavior_INPUT_ONLY)

	return pbAccess, nil
}

// UTILS

func (s *RecipeService) convertProtoToModel(pbAccess *pb.Access) (model.RecipeAccess, error) {
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

	switch pbrequester := pbAccess.Getrequester().GetName().(type) {
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

func (s *RecipeService) convertModelToProto(modelAccess model.RecipeAccess) (*pb.Access, error) {
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
