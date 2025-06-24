package domain

import (
	"context"
	"fmt"

	model "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

func (d *Domain) CreateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	// verify recipe is set
	if access.RecipeAccessParent.RecipeId.RecipeId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify issuer.user is set
	if access.RecipeAccessParent.Issuer.UserId == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "issuer required"}
	}

	// if issuer.circle is set, verify issuer has access to circle
	if access.RecipeAccessParent.Issuer.CircleId != 0 {
		permissionLevel, err := d.repo.GetCircleUserPermission(ctx, access.RecipeAccessParent.Issuer.UserId, access.RecipeAccessParent.Issuer.CircleId)
		if err != nil {
			return model.RecipeAccess{}, err
		}
		if permissionLevel < permPb.PermissionLevel_RESOURCE_PERMISSION_WRITE {
			return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "issuer does not have access to circle"}
		}
	}

	var filter string
	// based on recipient, set state and verify that the issuer has access to recipient
	if access.RecipeAccessParent.Recipient.UserId != 0 {
		filter = fmt.Sprintf("user_id = %d", access.RecipeAccessParent.Recipient.UserId)
		access.State = pb.Access_STATE_PENDING
		// TODO: verify issuer has correct permission to recipient user
	} else if access.RecipeAccessParent.Recipient.CircleId != 0 {
		filter = fmt.Sprintf("circle_id = %d", access.RecipeAccessParent.Recipient.CircleId)
		// TODO: verify issuer has correct permission to recipient circle
		access.State = pb.Access_STATE_ACCEPTED
	} else {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "recipient is required"}
	}

	// verify issuer has access to recipe
	recipeAccess, err := d.repo.ListRecipeAccesses(ctx, access.RecipeAccessParent, 1, 0, filter)
	if err != nil {
		return model.RecipeAccess{}, err
	}
	if len(recipeAccess) == 0 {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "issuer does not have access to recipe"}
	}

	if recipeAccess[0].Level < pb.Access_LEVEL_WRITE {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "issuer does not have correct access to recipe"}
	}

	// verify permission is set
	if access.Level == pb.Access_LEVEL_UNSPECIFIED {
		return model.RecipeAccess{}, domain.ErrInvalidArgument{Msg: "permission is required"}
	}

	// create access
	access, err = d.repo.CreateRecipeAccess(ctx, access)
	if err != nil {
		return model.RecipeAccess{}, err
	}

	return access, nil
}

func (d *Domain) DeleteRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) error {
	// verify recipe is set
	if parent.RecipeId.RecipeId == 0 {
		return domain.ErrInvalidArgument{Msg: "recipe id is required"}
	}

	// verify access id is set
	if id.RecipeAccessId == 0 {
		return domain.ErrInvalidArgument{Msg: "access id is required"}
	}

	// verify access is for the given recipe
	access, err := d.repo.GetRecipeAccess(ctx, parent, id)
	if err != nil {
		return err
	}

	if access.RecipeAccessParent.RecipeId.RecipeId != parent.RecipeId.RecipeId {
		return domain.ErrInvalidArgument{Msg: "access is not for the given recipe"}
	}

	// verify issuer is set
	if access.RecipeAccessParent.Issuer.UserId == 0 {
		return domain.ErrInvalidArgument{Msg: "issuer is required"}
	}

	// verify issuer has access to recipe
	recipeAccess, err := d.repo.ListRecipeAccesses(ctx, access.RecipeAccessParent, 1, 0, "")
	if err != nil {
		return err
	}
	if len(recipeAccess) == 0 {
		return domain.ErrInvalidArgument{Msg: "issuer does not have access to recipe"}
	}

	if recipeAccess[0].Level < pb.Access_LEVEL_WRITE {
		return domain.ErrInvalidArgument{Msg: "issuer does not have correct access to recipe"}
	}

	return d.repo.DeleteRecipeAccess(ctx, parent, id)
}

func (d *Domain) GetRecipeAccess(ctx context.Context, parent model.RecipeAccessParent, id model.RecipeAccessId) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}

func (d *Domain) ListRecipeAccesses(ctx context.Context, parent model.RecipeAccessParent, pageSize int32, pageOffset int32, filter string) ([]model.RecipeAccess, error) {
	if parent.Issuer.UserId == 0 && parent.Issuer.CircleId == 0 {
		return nil, domain.ErrInvalidArgument{Msg: "issuer is required"}
	}

	// TODO: verify issuer has access to recipe

	return d.repo.ListRecipeAccesses(ctx, parent, int64(pageSize), int64(pageOffset), filter)
}

func (d *Domain) UpdateRecipeAccess(ctx context.Context, access model.RecipeAccess) (model.RecipeAccess, error) {
	return model.RecipeAccess{}, nil
}
