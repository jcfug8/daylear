package namer

import (
	"fmt"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
	"go.uber.org/fx"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

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

type NewRecipeNamerResult struct {
	fx.Out

	UserRecipeNamer       RecipeNamer `group:"recipeNamer"`
	UserCircleRecipeNamer RecipeNamer `group:"recipeNamer"`
}

// NewRecipeNamer creates a new RecipeNamer.
func NewRecipeNamer() (NewRecipeNamerResult, error) {
	t := new(pb.Recipe)
	resourceOption := proto.GetExtension(
		t.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)
	patterns := resourceOption.Pattern
	if len(patterns) == 0 {
		return NewRecipeNamerResult{}, fmt.Errorf("no resource pattern found in %T", t)
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

	return NewRecipeNamerResult{
		UserRecipeNamer: &userRecipeNamer{
			namer: namers[0],
		},
		UserCircleRecipeNamer: &userCircleRecipeNamer{
			namer: namers[1],
		},
	}, nil
}
