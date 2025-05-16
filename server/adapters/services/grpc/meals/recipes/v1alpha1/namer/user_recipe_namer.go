package namer

import (
	"fmt"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/ports/fileretriever"

	"go.einride.tech/aip/resourcename"
)

const (
	userRecipeParentSegmentCount = 2
)

var _ RecipeNamer = &userRecipeNamer{}

type userRecipeNamer struct {
	namer *recipeNamer
}

// Format formats a recipe name.
func (n *userRecipeNamer) Format(parent model.RecipeParent, id model.RecipeId) (string, error) {
	return resourcename.Sprint(n.namer.pattern, fmt.Sprintf("%v", parent.UserId), fmt.Sprintf("%v", id.RecipeId)), nil
}

// IsMatch checks if a name matches the recipe pattern.
func (n *userRecipeNamer) IsMatch(name string) bool {
	return resourcename.Match(n.namer.pattern, name)
}

// IsParent checks if a name matches the recipe parent pattern.
func (n *userRecipeNamer) IsParent(parent string) bool {
	isParent := false
	foundSegments := 1
	resourcename.RangeParents(n.namer.pattern, func(p string) bool {
		if resourcename.Match(p, parent) && userRecipeParentSegmentCount == foundSegments {
			isParent = true
			return false
		}
		foundSegments++
		return true
	})
	return isParent
}

// Parse parses a recipe name.
func (n *userRecipeNamer) Parse(name string) (parent model.RecipeParent, id model.RecipeId, err error) {

	var userIdStr string

	var recipeIdStr string

	err = resourcename.Sscan(name, n.namer.pattern, &userIdStr, &recipeIdStr)
	if err != nil {
		return parent, id, err
	}

	parent.UserId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return parent, id, fileretriever.ErrInvalidArgument{Msg: "invalid parent format"}
	}
	id.RecipeId, err = strconv.ParseInt(recipeIdStr, 10, 64)
	if err != nil {
		return parent, id, fileretriever.ErrInvalidArgument{Msg: "invalid format"}
	}

	return parent, id, nil
}

// ParseParent parses a recipe parent name.
func (n *userRecipeNamer) ParseParent(name string) (parent model.RecipeParent, err error) {
	if !n.IsParent(name) {
		return parent, fmt.Errorf("invalid parent %s", name)
	}

	var userIdStr string

	resourcename.RangeParents(n.namer.pattern, func(p string) bool {
		if !resourcename.Match(p, name) {
			return true
		}

		err = resourcename.Sscan(name, p, &userIdStr)
		if err != nil {
			return false
		}

		return false
	})

	parent.UserId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return parent, fileretriever.ErrInvalidArgument{Msg: "invalid parent format"}
	}

	return parent, err
}
