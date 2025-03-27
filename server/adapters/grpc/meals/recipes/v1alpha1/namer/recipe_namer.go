package namer

import (
	"fmt"
	"strconv"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"

	"github.com/jcfug8/daylear/server/core/errz"
	"go.einride.tech/aip/resourcename"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	recipeParentSegmentCount = 2
)

var _ RecipeNamer = &defaultRecipeNamer{}

// RecipeNamer is an interface for creating and validating recipe names.
type RecipeNamer interface {
	Format(model.RecipeParent, model.RecipeId) (string, error)
	IsMatch(name string) bool
	IsParent(parent string) bool
	Parse(string) (model.RecipeParent, model.RecipeId, error)
	ParseParent(string) (model.RecipeParent, error)
}

type recipeNamer struct {
	pattern  string
	varCount int
}

type defaultRecipeNamer struct {
	// IRIOMO:CUSTOM_CODE_SLOT_START RecipeNamerFields
	namer *recipeNamer
	// IRIOMO:CUSTOM_CODE_SLOT_END
}

// NewRecipeNamer creates a new RecipeNamer.
func NewRecipeNamer() (RecipeNamer, error) {
	t := new(pb.Recipe)
	resourceOption := proto.GetExtension(
		t.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)
	patterns := resourceOption.Pattern
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", t)
	}

	namers := make([]*recipeNamer, 0, len(patterns))
	for _, pattern := range patterns {
		var patternScanner resourcename.Scanner
		var varCount int
		patternScanner.Init(patterns[0])
		for patternScanner.Scan() {
			segment := patternScanner.Segment()
			if segment.IsVariable() {
				varCount++
			}
		}

		namers = append(namers, &recipeNamer{
			pattern:  pattern,
			varCount: varCount,
		})
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START RecipeNamers
	return &defaultRecipeNamer{
		namer: namers[0],
	}, nil
	// IRIOMO:CUSTOM_CODE_SLOT_END
}

// Format formats a recipe name.
func (n *defaultRecipeNamer) Format(parent model.RecipeParent, id model.RecipeId) (string, error) {
	return resourcename.Sprint(n.namer.pattern, fmt.Sprintf("%v", parent.UserId), fmt.Sprintf("%v", id.RecipeId)), nil
}

// IsMatch checks if a name matches the recipe pattern.
func (n *defaultRecipeNamer) IsMatch(name string) bool {
	return resourcename.Match(n.namer.pattern, name)
}

// IsParent checks if a name matches the recipe parent pattern.
func (n *defaultRecipeNamer) IsParent(parent string) bool {
	isParent := false
	foundSegments := 1
	resourcename.RangeParents(n.namer.pattern, func(p string) bool {
		if resourcename.Match(p, parent) && recipeParentSegmentCount == foundSegments {
			isParent = true
			return false
		}
		foundSegments++
		return true
	})
	return isParent
}

// Parse parses a recipe name.
func (n *defaultRecipeNamer) Parse(name string) (parent model.RecipeParent, id model.RecipeId, err error) {

	var userIdStr string

	var recipeIdStr string

	err = resourcename.Sscan(name, n.namer.pattern, &userIdStr, &recipeIdStr)
	if err != nil {
		return parent, id, err
	}

	parent.UserId, err = strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return parent, id, errz.NewInvalidArgument("invalid parent format")
	}
	id.RecipeId, err = strconv.ParseInt(recipeIdStr, 10, 64)
	if err != nil {
		return parent, id, errz.NewInvalidArgument("invalid format")
	}

	return parent, id, nil
}

// ParseParent parses a recipe parent name.
func (n *defaultRecipeNamer) ParseParent(name string) (parent model.RecipeParent, err error) {
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
		return parent, errz.NewInvalidArgument("invalid parent format")
	}

	return parent, err
}
