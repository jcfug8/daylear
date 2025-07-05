package fieldmasker

import (
	"context"
	"fmt"
	"strings"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

var _ RecipeFieldMasker = &defaultRecipeFieldMasker{}

// RecipeFieldMasker is an interface for handling field masks for recipe.
type RecipeFieldMasker interface {
	GetFieldMaskFromCtx(ctx context.Context) *fieldmaskpb.FieldMask
	GetReadMask(*fieldmaskpb.FieldMask) ([]string, error)
	GetWriteMask(*fieldmaskpb.FieldMask) ([]string, error)
}

type defaultRecipeFieldMasker struct {
	fieldMaskFields map[string]fieldMaskField
}

// NewRecipeFieldMasker creates a new RecipeFieldMasker.
func NewRecipeFieldMasker() RecipeFieldMasker {
	t := new(pb.Recipe)
	fm := &defaultRecipeFieldMasker{
		fieldMaskFields: make(map[string]fieldMaskField),
	}

	// Recursively collect field masks
	collectFieldMasks(t.ProtoReflect(), "", fm.fieldMaskFields)

	fm.mapFieldMaskPathToDomainMasks("name", model.RecipeFields.Id)
	fm.mapFieldMaskPathToDomainMasks("title", model.RecipeFields.Title)
	fm.mapFieldMaskPathToDomainMasks("description", model.RecipeFields.Description)
	fm.mapFieldMaskPathToDomainMasks("directions", model.RecipeFields.Directions)
	fm.mapFieldMaskPathToDomainMasks("ingredient_groups", model.RecipeFields.IngredientGroups)
	fm.mapFieldMaskPathToDomainMasks("image_uri", model.RecipeFields.ImageURI)
	fm.mapFieldMaskPathToDomainMasks("visibility", model.RecipeFields.VisibilityLevel)
	fm.mapFieldMaskPathToDomainMasks("recipe_access.name", model.RecipeFields.AccessId)
	fm.mapFieldMaskPathToDomainMasks("recipe_access.permission_level", model.RecipeFields.PermissionLevel)
	fm.mapFieldMaskPathToDomainMasks("recipe_access.state", model.RecipeFields.State)

	return fm
}

func (f *defaultRecipeFieldMasker) mapFieldMaskPathToDomainMasks(fieldMaskPath string, domainMasks ...string) {
	if _, ok := f.fieldMaskFields[fieldMaskPath]; !ok {
		panic(fmt.Sprintf("field mask path %s not found", fieldMaskPath))
	}

	field := f.fieldMaskFields[fieldMaskPath]
	field.domainMasks = append(field.domainMasks, domainMasks...)
	f.fieldMaskFields[fieldMaskPath] = field
}

// GetFieldMaskFromCtx gets the field mask from the context.
func (f *defaultRecipeFieldMasker) GetFieldMaskFromCtx(ctx context.Context) *fieldmaskpb.FieldMask {
	if ctx == nil {
		return nil
	}

	headers, _ := metadata.FromIncomingContext(ctx)

	hMasks, ok := headers[fieldMaskKey]
	if !ok {
		return nil
	}

	masks := []string{}
	for _, m := range hMasks {
		splitMasks := strings.Split(m, ",")
		masks = append(masks, splitMasks...)
	}

	return &fieldmaskpb.FieldMask{
		Paths: masks,
	}
}

// GetReadMask gets the read mask.
func (f *defaultRecipeFieldMasker) GetReadMask(mask *fieldmaskpb.FieldMask) ([]string, error) {
	if len(mask.GetPaths()) == 0 || (len(mask.GetPaths()) == 1 && mask.GetPaths()[0] == "*") {
		mask = &fieldmaskpb.FieldMask{}
		for _, field := range f.fieldMaskFields {
			mask.Paths = append(mask.Paths, field.fieldMaskPath)
		}
	}

	// don't error on invalid paths, just ignore them
	// this is to allow for more flexibility as fields change

	return f.convertToDomainMask(mask)
}

// GetWriteMask gets the write mask.
func (f *defaultRecipeFieldMasker) GetWriteMask(mask *fieldmaskpb.FieldMask) ([]string, error) {
	if len(mask.GetPaths()) == 0 || (len(mask.GetPaths()) == 1 && mask.GetPaths()[0] == "*") {
		mask = &fieldmaskpb.FieldMask{}
		for _, field := range f.fieldMaskFields {
			if !field.readOnly {
				mask.Paths = append(mask.Paths, field.fieldMaskPath)
			}
		}
	}

	// output-only paths are to be ignored
	// as per https://google.aip.dev/161#output-only-fields

	for _, m := range mask.GetPaths() {
		if _, ok := f.fieldMaskFields[m]; !ok {
			return nil, status.Errorf(codes.InvalidArgument, "field mask path %s not found or read-only", m)
		}
	}

	return f.convertToDomainMask(mask)
}

func (f *defaultRecipeFieldMasker) convertToDomainMask(mask *fieldmaskpb.FieldMask) ([]string, error) {
	domainMask := make([]string, 0, len(mask.GetPaths()))
	for _, m := range mask.GetPaths() {
		field, ok := f.fieldMaskFields[m]
		if !ok && len(field.domainMasks) == 0 {
			return nil, status.Errorf(codes.InvalidArgument, "field mask path %s not usable", m)
		} else if ok {
			domainMask = append(domainMask, field.domainMasks...)
		} else {
			return nil, status.Errorf(codes.InvalidArgument, "field mask path %s not found", m)
		}
	}

	return domainMask, nil
}
