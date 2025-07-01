package fieldmasker

import (
	"context"
	"fmt"
	"strings"

	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"github.com/jcfug8/daylear/server/ports/fileretriever"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	// IRIOMO:CUSTOM_CODE_SLOT_START importFieldMasker
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

var _ CircleFieldMasker = &defaultCircleFieldMasker{}

// CircleFieldMasker is an interface for handling field masks for circle.
type CircleFieldMasker interface {
	GetFieldMaskFromCtx(ctx context.Context) *fieldmaskpb.FieldMask
	GetReadMask(*fieldmaskpb.FieldMask) ([]string, error)
	GetWriteMask(*fieldmaskpb.FieldMask) ([]string, error)
}

type defaultCircleFieldMasker struct {
	fieldMaskFields map[string]fieldMaskField
}

// NewCircleFieldMasker creates a new CircleFieldMasker.
func NewCircleFieldMasker() CircleFieldMasker {
	t := new(pb.Circle)
	fm := &defaultCircleFieldMasker{
		fieldMaskFields: make(map[string]fieldMaskField),
	}

	// Recursively collect field masks
	collectFieldMasks(t.ProtoReflect(), "", fm.fieldMaskFields)

	// IRIOMO:CUSTOM_CODE_SLOT_START resourceNamerMapFields
	fm.mapFieldMaskPathToDomainMasks("name", model.CircleFields.Id)
	fm.mapFieldMaskPathToDomainMasks("title", model.CircleFields.Title)
	fm.mapFieldMaskPathToDomainMasks("visibility", model.CircleFields.Visibility)
	fm.mapFieldMaskPathToDomainMasks("permission", model.CircleFields.Permission)
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return fm
}

func (f *defaultCircleFieldMasker) mapFieldMaskPathToDomainMasks(fieldMaskPath string, domainMasks ...string) {
	if _, ok := f.fieldMaskFields[fieldMaskPath]; !ok {
		panic(fmt.Sprintf("field mask path %s not found", fieldMaskPath))
	}

	field := f.fieldMaskFields[fieldMaskPath]
	field.domainMasks = append(field.domainMasks, domainMasks...)
	f.fieldMaskFields[fieldMaskPath] = field
}

// GetFieldMaskFromCtx gets the field mask from the context.
func (f *defaultCircleFieldMasker) GetFieldMaskFromCtx(ctx context.Context) *fieldmaskpb.FieldMask {
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
func (f *defaultCircleFieldMasker) GetReadMask(mask *fieldmaskpb.FieldMask) ([]string, error) {
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
func (f *defaultCircleFieldMasker) GetWriteMask(mask *fieldmaskpb.FieldMask) ([]string, error) {
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
			return nil, fileretriever.ErrInvalidArgument{Msg: fmt.Sprintf("field mask path %s not found or read-only", m)}
		}
	}

	return f.convertToDomainMask(mask)
}

func (f *defaultCircleFieldMasker) convertToDomainMask(mask *fieldmaskpb.FieldMask) ([]string, error) {
	domainMask := make([]string, 0, len(mask.GetPaths()))
	for _, m := range mask.GetPaths() {
		field, ok := f.fieldMaskFields[m]
		if !ok && len(field.domainMasks) == 0 {
			return nil, fileretriever.ErrInvalidArgument{Msg: fmt.Sprintf("field mask path %s not usable", m)}
		} else if ok {
			domainMask = append(domainMask, field.domainMasks...)
		} else {
			return nil, fileretriever.ErrInvalidArgument{Msg: fmt.Sprintf("field mask path %s not found", m)}
		}
	}

	return domainMask, nil
}
