package convert

import (
	"encoding/json"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// ListFromCoreModel converts a core model to a gorm model.
func ListFromCoreModel(m cmodel.List) (gmodel.List, error) {
	var err error
	list := gmodel.List{
		ListId:          m.Id.ListId,
		Title:           m.Title,
		Description:     m.Description,
		ShowCompleted:   m.ShowCompleted,
		VisibilityLevel: m.VisibilityLevel,
		CreateTime:      m.CreateTime,
		UpdateTime:      m.UpdateTime,
	}

	// Marshal Sections ([]ListSection) to []byte (jsonb)
	if m.Sections != nil {
		list.Sections, err = json.Marshal(m.Sections)
		if err != nil {
			return gmodel.List{}, err
		}
	}

	return list, nil
}

// ListToCoreModel converts a gorm model to a core model.
func ListToCoreModel(m gmodel.List) (cmodel.List, error) {
	permissionLevel := m.PermissionLevel
	if m.VisibilityLevel == types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC && m.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		permissionLevel = types.PermissionLevel_PERMISSION_LEVEL_PUBLIC
	}

	var err error
	list := cmodel.List{
		Id: cmodel.ListId{
			ListId: m.ListId,
		},
		Title:           m.Title,
		Description:     m.Description,
		ShowCompleted:   m.ShowCompleted,
		VisibilityLevel: m.VisibilityLevel,
		CreateTime:      m.CreateTime,
		UpdateTime:      m.UpdateTime,
	}

	if m.ListFavoriteId != 0 {
		list.Favorited = true
	}

	// Populate ListAccess if permission or state is set (i.e., join succeeded)
	if m.PermissionLevel != 0 || m.State != 0 {
		list.ListAccess = cmodel.ListAccess{
			ListAccessParent: cmodel.ListAccessParent{
				ListId: cmodel.ListId{ListId: m.ListId},
			},
			ListAccessId:    cmodel.ListAccessId{ListAccessId: m.ListAccessId},
			PermissionLevel: permissionLevel,
			State:           m.State,
			AcceptTarget:    m.AcceptTarget,
		}
	}

	// Unmarshal Sections ([]byte) to []ListSection
	if m.Sections != nil {
		err = json.Unmarshal(m.Sections, &list.Sections)
		if err != nil {
			return cmodel.List{}, err
		}
	}

	return list, nil
}
