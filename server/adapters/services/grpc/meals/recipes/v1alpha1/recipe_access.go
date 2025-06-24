package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc/pagination"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"go.einride.tech/aip/fieldbehavior"
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
	user, circleID, err := headers.ParseAuthData(ctx, s.circleNamer)
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

	// set issuer
	modelAccess.RecipeAccessParent.Issuer = model.RecipeParent{
		UserId:   user.Id.UserId,
		CircleId: circleID.CircleId,
	}

	// parse parent
	_, err = s.accessNamer.ParseParent(request.Parent, &modelAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// create access
	createdAccess, err := s.domain.CreateRecipeAccess(ctx, modelAccess)
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
	user, circleID, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// parse name
	recipeAccessParent := model.RecipeAccessParent{
		Issuer: model.RecipeParent{
			UserId:   user.Id.UserId,
			CircleId: circleID.CircleId,
		},
	}
	accessId := model.RecipeAccessId{}
	_, err = s.accessNamer.Parse(request.Name, &model.RecipeAccess{
		RecipeAccessParent: recipeAccessParent,
		RecipeAccessId:     accessId,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// delete access
	err = s.domain.DeleteRecipeAccess(ctx, recipeAccessParent, accessId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete access: %v", err)
	}

	return nil, nil
}

func (s *RecipeService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	user, circleID, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// set issuer
	recipeAccessParent := model.RecipeAccessParent{
		Issuer: model.RecipeParent{
			UserId:   user.Id.UserId,
			CircleId: circleID.CircleId,
		},
	}

	// parse name
	recipeAccessId := model.RecipeAccessId{}
	_, err = s.accessNamer.Parse(request.Name, &model.RecipeAccess{
		RecipeAccessParent: recipeAccessParent,
		RecipeAccessId:     recipeAccessId,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// get access
	access, err := s.domain.GetRecipeAccess(ctx, recipeAccessParent, recipeAccessId)
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
	user, circleID, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// set issuer
	recipeAccessParent := model.RecipeAccessParent{
		Issuer: model.RecipeParent{
			UserId:   user.Id.UserId,
			CircleId: circleID.CircleId,
		},
	}

	// parse parent
	_, err = s.accessNamer.ParseParent(request.Parent, &recipeAccessParent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// pagination
	pageToken, err := pagination.ParsePageToken[model.RecipeAccess](request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	if pageToken.PageSize == 0 {
		pageToken.PageSize = accessDefaultPageSize
	}
	pageToken.PageSize = min(pageToken.PageSize, accessMaxPageSize)

	// list accesses
	accesses, err := s.domain.ListRecipeAccesses(ctx, recipeAccessParent, pageToken.PageSize, pageToken.Skip, request.Filter)
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
		NextPageToken: pagination.EncodePageToken(pageToken.Next(accesses)),
		Accesses:      pbAccesses,
	}, nil
}

func (s *RecipeService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	user, circleID, err := headers.ParseAuthData(ctx, s.circleNamer)
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

	// set issuer
	modelAccess.RecipeAccessParent.Issuer = model.RecipeParent{
		UserId:   user.Id.UserId,
		CircleId: circleID.CircleId,
	}

	// parse name
	_, err = s.accessNamer.Parse(request.Access.GetName(), modelAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// update access
	updatedAccess, err := s.domain.UpdateRecipeAccess(ctx, modelAccess)
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

	switch pbIssuer := pbAccess.GetIssuer().GetName().(type) {
	case *pb.Access_IssuerOrRecipient_User:
		modelAccess.Issuer = model.RecipeParent{}
		_, err := s.userNamer.Parse(pbIssuer.User, &modelAccess.Issuer)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid issuer: %v", err)
		}
	case *pb.Access_IssuerOrRecipient_Circle:
		modelAccess.Issuer = model.RecipeParent{}
		_, err := s.circleNamer.Parse(pbIssuer.Circle, &modelAccess.Issuer)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid issuer: %v", err)
		}
	}

	switch pbRecipient := pbAccess.GetRecipient().GetName().(type) {
	case *pb.Access_IssuerOrRecipient_User:
		modelAccess.Recipient = model.RecipeParent{}
		_, err := s.userNamer.Parse(pbRecipient.User, &modelAccess.Recipient)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", err)
		}
	case *pb.Access_IssuerOrRecipient_Circle:
		modelAccess.Recipient = model.RecipeParent{}
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

	if modelAccess.Issuer.UserId != 0 {
		userName, err := s.userNamer.Format(modelAccess.Issuer)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format issuer: %v", err)
		}
		pbAccess.Issuer = &pb.Access_IssuerOrRecipient{
			Name: &pb.Access_IssuerOrRecipient_User{
				User: userName,
			},
		}
	} else if modelAccess.Issuer.CircleId != 0 {
		circleName, err := s.circleNamer.Format(modelAccess.Issuer)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format issuer: %v", err)
		}
		pbAccess.Issuer = &pb.Access_IssuerOrRecipient{
			Name: &pb.Access_IssuerOrRecipient_Circle{
				Circle: circleName,
			},
		}
	}

	if modelAccess.Recipient.UserId != 0 {
		userName, err := s.userNamer.Format(modelAccess.Recipient)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format recipient: %v", err)
		}
		pbAccess.Recipient = &pb.Access_IssuerOrRecipient{
			Name: &pb.Access_IssuerOrRecipient_User{
				User: userName,
			},
		}
	} else if modelAccess.Recipient.CircleId != 0 {
		circleName, err := s.circleNamer.Format(modelAccess.Recipient)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format recipient: %v", err)
		}
		pbAccess.Recipient = &pb.Access_IssuerOrRecipient{
			Name: &pb.Access_IssuerOrRecipient_Circle{
				Circle: circleName,
			},
		}
	}

	return pbAccess, nil
}
