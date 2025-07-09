package convert

import (
	"encoding/json"
	"time"

	gmodel "github.com/jcfug8/daylear/server/adapters/clients/gorm/model"
	cmodel "github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/genapi/api/types"
)

// RecipeFromCoreModel converts a core model to a gorm model.
func RecipeFromCoreModel(m cmodel.Recipe) (gmodel.Recipe, error) {
	var err error
	recipe := gmodel.Recipe{
		RecipeId:             m.Id.RecipeId,
		Title:                m.Title,
		Description:          m.Description,
		ImageURI:             m.ImageURI,
		VisibilityLevel:      m.Visibility,
		Citation:             m.Citation,
		PrepDurationSeconds:  int64(m.PrepDuration.Seconds()),
		CookDurationSeconds:  int64(m.CookDuration.Seconds()),
		TotalDurationSeconds: int64(m.TotalDuration.Seconds()),
		CookingMethod:        m.CookingMethod,
		YieldAmount:          m.YieldAmount,
		CreateTime:           m.CreateTime,
		UpdateTime:           m.UpdateTime,
	}

	recipe.Directions, err = json.Marshal(m.Directions)
	if err != nil {
		return gmodel.Recipe{}, err
	}

	recipe.IngredientGroups, err = json.Marshal(m.IngredientGroups)
	if err != nil {
		return gmodel.Recipe{}, err
	}

	// Marshal Categories and Cuisines ([]string) to []byte (jsonb)
	if m.Categories != nil {
		recipe.Categories, err = json.Marshal(m.Categories)
		if err != nil {
			return gmodel.Recipe{}, err
		}
	}
	if m.Cuisines != nil {
		recipe.Cuisines, err = json.Marshal(m.Cuisines)
		if err != nil {
			return gmodel.Recipe{}, err
		}
	}

	return recipe, nil
}

// RecipeToCoreModel converts a gorm model to a core model.
func RecipeToCoreModel(m gmodel.Recipe) (cmodel.Recipe, error) {
	permissionLevel := m.PermissionLevel
	if m.VisibilityLevel == types.VisibilityLevel_VISIBILITY_LEVEL_PUBLIC && m.PermissionLevel == types.PermissionLevel_PERMISSION_LEVEL_UNSPECIFIED {
		permissionLevel = types.PermissionLevel_PERMISSION_LEVEL_PUBLIC
	}

	var err error
	recipe := cmodel.Recipe{
		Id: cmodel.RecipeId{
			RecipeId: m.RecipeId,
		},
		Title:         m.Title,
		Description:   m.Description,
		ImageURI:      m.ImageURI,
		Visibility:    m.VisibilityLevel,
		Citation:      m.Citation,
		PrepDuration:  time.Duration(m.PrepDurationSeconds) * time.Second,
		CookDuration:  time.Duration(m.CookDurationSeconds) * time.Second,
		TotalDuration: time.Duration(m.TotalDurationSeconds) * time.Second,
		CookingMethod: m.CookingMethod,
		YieldAmount:   m.YieldAmount,
		CreateTime:    m.CreateTime,
		UpdateTime:    m.UpdateTime,
	}

	// Populate RecipeAccess if permission or state is set (i.e., join succeeded)
	if m.PermissionLevel != 0 || m.State != 0 {
		recipe.RecipeAccess = cmodel.RecipeAccess{
			RecipeAccessParent: cmodel.RecipeAccessParent{
				RecipeId: cmodel.RecipeId{RecipeId: m.RecipeId},
			},
			RecipeAccessId:  cmodel.RecipeAccessId{RecipeAccessId: m.RecipeAccessId},
			PermissionLevel: permissionLevel,
			State:           m.State,
		}
	}

	if m.Directions != nil {
		err = json.Unmarshal(m.Directions, &recipe.Directions)
		if err != nil {
			return cmodel.Recipe{}, err
		}
	}

	if m.IngredientGroups != nil {
		err = json.Unmarshal(m.IngredientGroups, &recipe.IngredientGroups)
		if err != nil {
			return cmodel.Recipe{}, err
		}
	}

	// Unmarshal Categories and Cuisines ([]byte) to []string
	if m.Categories != nil {
		err = json.Unmarshal(m.Categories, &recipe.Categories)
		if err != nil {
			return cmodel.Recipe{}, err
		}
	}
	if m.Cuisines != nil {
		err = json.Unmarshal(m.Cuisines, &recipe.Cuisines)
		if err != nil {
			return cmodel.Recipe{}, err
		}
	}

	return recipe, nil
}

// RecipeListFromCoreModel converts a list of core models to a list of gorm models.
func RecipeListFromCoreModel(m []cmodel.Recipe) (res []gmodel.Recipe, err error) {
	res = make([]gmodel.Recipe, len(m))
	for i, v := range m {
		res[i], err = RecipeFromCoreModel(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

// RecipeListToCoreModel converts a list of gorm models to a list of core models.
func RecipeListToCoreModel(m []gmodel.Recipe) (res []cmodel.Recipe, err error) {
	res = make([]cmodel.Recipe, len(m))
	for i, v := range m {
		res[i], err = RecipeToCoreModel(v)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}
