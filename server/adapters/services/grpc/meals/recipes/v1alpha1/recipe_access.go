package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *RecipeService) CreateAccess(ctx context.Context, request *pb.CreateAccessRequest) (*pb.Access, error) {
	// check field behavior
	fieldbehavior.ClearFields(request, annotations.FieldBehavior_OUTPUT_ONLY)
	err := fieldbehavior.ValidateRequiredFields(request)
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
	_, err = s.accessNamer.Parse(request.Parent, &modelAccess)
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

	// return response
	return pbAccess, nil
}

func (s *RecipeService) DeleteAccess(ctx context.Context, request *pb.DeleteAccessRequest) (*pb.Access, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAccess not implemented")
}

func (s *RecipeService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAccess not implemented")
}

func (s *RecipeService) ListAccesses(ctx context.Context, request *pb.ListAccessesRequest) (*pb.ListAccessesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListAccesses not implemented")
}

func (s *RecipeService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAccess not implemented")
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
		_, err := s.recipeNamer_User.Parse(pbIssuer.User, &modelAccess.Issuer)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid issuer: %v", err)
		}
	case *pb.Access_IssuerOrRecipient_Circle:
		modelAccess.Issuer = model.RecipeParent{}
		_, err := s.recipeNamer_Circle.Parse(pbIssuer.Circle, &modelAccess.Issuer)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid issuer: %v", err)
		}
	}

	switch pbRecipient := pbAccess.GetRecipient().GetName().(type) {
	case *pb.Access_IssuerOrRecipient_User:
		modelAccess.Recipient = model.RecipeParent{}
		_, err := s.recipeNamer_User.Parse(pbRecipient.User, &modelAccess.Recipient)
		if err != nil {
			return model.RecipeAccess{}, status.Errorf(codes.InvalidArgument, "invalid recipient: %v", err)
		}
	case *pb.Access_IssuerOrRecipient_Circle:
		modelAccess.Recipient = model.RecipeParent{}
		_, err := s.recipeNamer_Circle.Parse(pbRecipient.Circle, &modelAccess.Recipient)
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
		userName, err := s.recipeNamer_User.Format(modelAccess.Issuer)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format issuer: %v", err)
		}
		pbAccess.Issuer = &pb.Access_IssuerOrRecipient{
			Name: &pb.Access_IssuerOrRecipient_User{
				User: userName,
			},
		}
	} else if modelAccess.Issuer.CircleId != 0 {
		circleName, err := s.recipeNamer_Circle.Format(modelAccess.Issuer)
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
		userName, err := s.recipeNamer_User.Format(modelAccess.Recipient)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to format recipient: %v", err)
		}
		pbAccess.Recipient = &pb.Access_IssuerOrRecipient{
			Name: &pb.Access_IssuerOrRecipient_User{
				User: userName,
			},
		}
	} else if modelAccess.Recipient.CircleId != 0 {
		circleName, err := s.recipeNamer_Circle.Format(modelAccess.Recipient)
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
