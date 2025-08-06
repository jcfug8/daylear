package fieldmask

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FieldMasker allows mapping keys to a set of values and converting lists of keys to deduplicated lists of values.
type FieldMasker struct {
	mapping map[string][]string
}

// NewFieldMasker initializes a new FieldMasker with the provided mapping.
func NewFieldMasker(mapping map[string][]string) FieldMasker {
	return FieldMasker{mapping: mapping}
}

// NewProtoFieldMasker initializes a FieldMasker and ensures all keys exist as field mask paths in the proto message.
func NewProtoFieldMasker(msg proto.Message, mapping map[string][]string) (FieldMasker, error) {
	validPaths := make(map[string]struct{}, 0)
	collectFieldMasks(msg.ProtoReflect(), "", validPaths)

	for key := range mapping {
		if _, ok := validPaths[key]; !ok {
			return FieldMasker{}, fmt.Errorf("field mask key '%s' does not exist in proto message %T", key, msg)
		}
	}
	return FieldMasker{mapping: mapping}, nil
}

// collectFieldMasks recursively collects field masks from a proto.Message.
func collectFieldMasks(msg protoreflect.Message, prefix string, fields map[string]struct{}) {
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

		fields[fullPath] = struct{}{}

		// Check if the field is a message type and recursively collect its fields
		if fd.Kind() == protoreflect.MessageKind && !fd.IsList() {
			nestedMsg := msg.Get(fd).Message()
			collectFieldMasks(nestedMsg, fullPath, fields)
		}
	}
}

// Convert takes a list of keys and returns a deduplicated list of all mapped values.
func (fmm *FieldMasker) Convert(keys []string) []string {
	if len(keys) == 0 {
		return fmm.GetAll()
	}

	result := make([]string, 0)
	seen := make(map[string]struct{})
	for _, key := range keys {
		if values, ok := fmm.mapping[key]; ok {
			for _, v := range values {
				if _, exists := seen[v]; !exists {
					result = append(result, v)
					seen[v] = struct{}{}
				}
			}
		}
	}
	return result
}

func (fmm *FieldMasker) GetAll() []string {
	result := make([]string, 0)
	for _, values := range fmm.mapping {
		result = append(result, values...)
	}
	return result
}
