package convert

import (
	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
)

// ListItemFromCoreModel converts a core model to a gorm model.
func ListItemFromCoreModel(m cmodel.ListItem) (gmodel.ListItem, error) {
	listItem := gmodel.ListItem{
		ListItemId:     m.Id.ListItemId,
		ListId:         m.Parent.ListId.ListId,
		Title:          m.Title,
		Points:         m.Points,
		RecurrenceRule: m.RecurrenceRule,
		ListSectionId:  m.ListSectionId,
		CreateTime:     m.CreateTime,
		UpdateTime:     m.UpdateTime,
	}

	return listItem, nil
}

// ListItemToCoreModel converts a gorm model to a core model.
func ListItemToCoreModel(m gmodel.ListItem) (cmodel.ListItem, error) {
	listItem := cmodel.ListItem{
		Id: cmodel.ListItemId{
			ListItemId: m.ListItemId,
		},
		Parent: cmodel.ListItemParent{
			ListId: cmodel.ListId{
				ListId: m.ListId,
			},
		},
		Title:          m.Title,
		Points:         m.Points,
		RecurrenceRule: m.RecurrenceRule,
		ListSectionId:  m.ListSectionId,
		CreateTime:     m.CreateTime,
		UpdateTime:     m.UpdateTime,
	}

	return listItem, nil
}
