package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	// IRIOMO:CUSTOM_CODE_SLOT_START resourceConvert
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// UserFromCoreModel converts a core model to a gorm model.
func UserFromCoreModel(m cmodel.User) (gmodel.User, error) {
	// IRIOMO:CUSTOM_CODE_SLOT_START userFromCoreModel
	u := gmodel.User{
		UserId:   m.Id.UserId,
		Email:    m.Email,
		Username: m.Username,
	}

	if m.AmazonId != "" {
		u.AmazonId = &m.AmazonId
	}

	if m.FacebookId != "" {
		u.FacebookId = &m.FacebookId
	}

	if m.GoogleId != "" {
		u.GoogleId = &m.GoogleId
	}

	return u, nil
	// IRIOMO:CUSTOM_CODE_SLOT_END
}

// UserToCoreModel converts a gorm model to a core model.
func UserToCoreModel(m gmodel.User) (cmodel.User, error) {
	// IRIOMO:CUSTOM_CODE_SLOT_START userToCoreModel
	u := cmodel.User{
		Id: cmodel.UserId{
			UserId: m.UserId,
		},
		Email: m.Email,
		Username: m.Username,
	}

	if m.AmazonId != nil {
		u.AmazonId = *m.AmazonId
	}

	if m.FacebookId != nil {
		u.FacebookId = *m.FacebookId
	}

	if m.GoogleId != nil {
		u.GoogleId = *m.GoogleId
	}

	return u, nil
	// IRIOMO:CUSTOM_CODE_SLOT_END
}

// UserListFromCoreModel converts a list of core models to a list of gorm models.
func UserListFromCoreModel(m []cmodel.User) (res []gmodel.User, err error) {
	res = make([]gmodel.User, len(m))
	for i, v := range m {
		res[i], err = UserFromCoreModel(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// UserListToCoreModel converts a list of gorm models to a list of core models.
func UserListToCoreModel(m []gmodel.User) (res []cmodel.User, err error) {
	res = make([]cmodel.User, len(m))
	for i, v := range m {
		res[i], err = UserToCoreModel(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
