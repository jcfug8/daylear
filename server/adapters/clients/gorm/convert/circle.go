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
		VisibilityLevel: m.VisibilityLevel,
		PermissionLevel: m.PermissionLevel,
		AccessState:     m.AccessState,
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
		VisibilityLevel: g.VisibilityLevel,
		PermissionLevel: g.PermissionLevel,
		AccessState:     g.AccessState,
	}, nil
}
