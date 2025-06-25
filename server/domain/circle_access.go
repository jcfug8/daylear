package domain

import (
	"context"
	"errors"

	"github.com/jcfug8/daylear/server/core/model"
	permPb "github.com/jcfug8/daylear/server/genapi/api/types"
	domain "github.com/jcfug8/daylear/server/ports/domain"
)

func (d *Domain) verifyCircleAccess(ctx context.Context, authAccount model.AuthAccount) (permPb.PermissionLevel, error) {
	// verify auth account is set
	if authAccount.UserId == 0 {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "auth user is required"}
	}

	// verify circle id is set
	if authAccount.CircleId == 0 {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrInvalidArgument{Msg: "circle id is required"}
	}

	// verify circle exists
	circle, err := d.repo.GetCircle(ctx, model.Circle{Id: model.CircleId{CircleId: authAccount.CircleId}}, []string{model.CircleFields.IsPublic})
	if err != nil {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
	}

	// check if user is a member of the circle
	permissionLevel, err := d.repo.GetCircleUserPermission(ctx, authAccount.UserId, authAccount.CircleId)
	if errors.Is(err, domain.ErrNotFound{}) && !circle.IsPublic {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, domain.ErrPermissionDenied{Msg: "user not a member of the circle"}
	} else if err != nil && !errors.Is(err, domain.ErrNotFound{}) {
		return permPb.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED, err
	}

	return permissionLevel, nil
}
