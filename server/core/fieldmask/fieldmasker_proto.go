package fieldmask

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// NewProtoFieldMasker initializes a FieldMasker and ensures all keys exist as field mask paths in the proto message.
func NewProtoFieldMasker(msg proto.Message, mapping map[string][]string) (FieldMasker, error) {
	m := make(map[string][]Field, 0)
	for key, fields := range mapping {
		mFields := m[key]
		for _, field := range fields {
			mFields = append(mFields, Field{Name: field})
		}
		m[key] = mFields
	}

	validPaths := make(map[string]bool, 0)
	collectProtoFieldNames(msg.ProtoReflect(), "", validPaths)

	for key := range mapping {
		if _, ok := validPaths[key]; !ok {
			return FieldMasker{}, fmt.Errorf("field mask key '%s' does not exist in proto message %T", key, msg)
		}
	}

	return FieldMasker{mapping: m}, nil
}

// collectProtoFieldNames recursively collects field masks from a proto.Message.
func collectProtoFieldNames(msg protoreflect.Message, prefix string, fields map[string]bool) {
	msgDescriptor := msg.Descriptor()

	// Iterate over all field descriptors
	for i := 0; i < msgDescriptor.Fields().Len(); i++ {
		fd := msgDescriptor.Fields().Get(i)
		fieldName := string(fd.Name())
		fullPath := fieldName
		if prefix != "" {
			fullPath = prefix + "." + fieldName
		}

		// // Check for field behavior annotations
		// readOnly := false
		// if options := fd.Options().(*descriptorpb.FieldOptions); options != nil {
		// 	behaviors := proto.GetExtension(options, annotations.E_FieldBehavior).([]annotations.FieldBehavior)
		// 	for _, behavior := range behaviors {
		// 		if behavior == annotations.FieldBehavior_OUTPUT_ONLY || behavior == annotations.FieldBehavior_IMMUTABLE || behavior == annotations.FieldBehavior_IDENTIFIER {
		// 			readOnly = true
		// 			break
		// 		}
		// 	}
		// }

		fields[fullPath] = true

		// Check if the field is a message type and recursively collect its fields
		if fd.Kind() == protoreflect.MessageKind && !fd.IsList() {
			nestedMsg := msg.Get(fd).Message()
			collectProtoFieldNames(nestedMsg, fullPath, fields)
		}
	}
}
