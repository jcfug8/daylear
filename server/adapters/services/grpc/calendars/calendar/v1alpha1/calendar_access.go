package v1alpha1

import (
	"github.com/jcfug8/daylear/server/core/model"
)

var calendarAccessFieldMap = map[string][]string{
	"name":             {model.CalendarAccessField_Parent, model.CalendarAccessField_Id},
	"permission_level": {model.CalendarAccessField_PermissionLevel},
	"state":            {model.CalendarAccessField_State},
	"requester":        {model.CalendarAccessField_Requester},
	"recipient":        {model.CalendarAccessField_Recipient},
}
