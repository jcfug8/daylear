package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// CircleFromCoreModel converts a core model.Circle to a GORM model.Circle
func CircleFromCoreModel(m cmodel.Circle) (gmodel.Circle, error) {
	return gmodel.Circle{
		CircleId: m.Id.CircleId,
		Title:    m.Title,
		IsPublic: m.IsPublic,
	}, nil
}

// CircleToCoreModel converts a GORM model.Circle to a core model.Circle
func CircleToCoreModel(g gmodel.Circle) (cmodel.Circle, error) {
	return cmodel.Circle{
		Id:       cmodel.CircleId{CircleId: g.CircleId},
		Title:    g.Title,
		IsPublic: g.IsPublic,
	}, nil
}
