package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// CircleFromCoreModel converts a core model.Circle to a GORM model.Circle
func CircleFromCoreModel(m cmodel.Circle) (gmodel.Circle, error) {
	return gmodel.Circle{
		CircleId:        m.Id.CircleId,
		Title:           m.Title,
		Description:     m.Description,
		Handle:          m.Handle,
		ImageURI:        m.ImageURI,
		VisibilityLevel: m.VisibilityLevel,
	}, nil
}

// CircleToCoreModel converts a GORM model.Circle to a core model.Circle
func CircleToCoreModel(g gmodel.Circle) (cmodel.Circle, error) {
	if g.VisibilityLevel == types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC && g.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		g.PermissionLevel = types.PermissionLevel_PERMISSION_LEVEL_PUBLIC
	}

	return cmodel.Circle{
		Id:              cmodel.CircleId{CircleId: g.CircleId},
		Title:           g.Title,
		Description:     g.Description,
		Handle:          g.Handle,
		ImageURI:        g.ImageURI,
		VisibilityLevel: g.VisibilityLevel,
		CircleAccess: cmodel.CircleAccess{
			CircleAccessParent: cmodel.CircleAccessParent{CircleId: cmodel.CircleId{CircleId: g.CircleId}},
			CircleAccessId:     cmodel.CircleAccessId{CircleAccessId: g.CircleAccessId},
			PermissionLevel:    g.PermissionLevel,
			State:              g.State,
		},
	}, nil
}
