package domain

import (
	"context"

	"github.com/jcfug8/daylear/server/core/model"
)

// CreateUserAccess -
func (d *Domain) CreateUserAccess(ctx context.Context, authAccount model.AuthAccount, access model.UserAccess) (model.UserAccess, error) {
	return model.UserAccess{}, nil
}

// DeleteUserAccess -
func (d *Domain) DeleteUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) error {
	return nil
}

// GetUserAccess -
func (d *Domain) GetUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error) {
	return model.UserAccess{}, nil
}

// ListUserAccesses -
func (d *Domain) ListUserAccesses(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, pageSize int32, pageOffset int64, filter string) ([]model.UserAccess, error) {
	return nil, nil
}

// UpdateUserAccess -
func (d *Domain) UpdateUserAccess(ctx context.Context, authAccount model.AuthAccount, access model.UserAccess) (model.UserAccess, error) {
	return model.UserAccess{}, nil
}

// AcceptUserAccess -
func (d *Domain) AcceptUserAccess(ctx context.Context, authAccount model.AuthAccount, parent model.UserAccessParent, id model.UserAccessId) (model.UserAccess, error) {
	return model.UserAccess{}, nil
}
