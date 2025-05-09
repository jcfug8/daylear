package fieldmasker

import (
	"go.uber.org/fx"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	fieldMaskKey = "tcn-field-mask"
)

// Module -
var Module = fx.Module(
	"fieldmasker",
	fx.Provide(
		NewCircleFieldMasker,
	),
)

type fieldMaskField struct {
	fieldMaskPath string
	domainMasks   []string
	readOnly      bool
}

// collectFieldMasks recursively collects field masks from a proto.Message.
func collectFieldMasks(msg protoreflect.Message, prefix string, fields map[string]fieldMaskField) {
	msgDescriptor := msg.Descriptor()

	// Iterate over all field descriptors
	for i := 0; i < msgDescriptor.Fields().Len(); i++ {
		fd := msgDescriptor.Fields().Get(i)
		fieldName := string(fd.Name())
		fullPath := fieldName
		if prefix != "" {
			fullPath = prefix + "." + fieldName
		}

		// Check for field behavior annotations
		readOnly := false
		if options := fd.Options().(*descriptorpb.FieldOptions); options != nil {
			behaviors := proto.GetExtension(options, annotations.E_FieldBehavior).([]annotations.FieldBehavior)
			for _, behavior := range behaviors {
				if behavior == annotations.FieldBehavior_OUTPUT_ONLY || behavior == annotations.FieldBehavior_IMMUTABLE || behavior == annotations.FieldBehavior_IDENTIFIER {
					readOnly = true
					break
				}
			}
		}

		fields[fullPath] = fieldMaskField{
			fieldMaskPath: fullPath,
			domainMasks:   []string{},
			readOnly:      readOnly,
		}

		// Check if the field is a message type and recursively collect its fields
		if fd.Kind() == protoreflect.MessageKind && !fd.IsList() {
			nestedMsg := msg.Get(fd).Message()
			collectFieldMasks(nestedMsg, fullPath, fields)
		}
	}
}
