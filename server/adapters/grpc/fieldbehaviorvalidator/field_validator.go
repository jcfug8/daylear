package fieldbehaviorvalidator

import (
	"github.com/jcfug8/daylear/server/core/errz"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type FieldBehaviorValidator interface {
	Validate(t proto.Message) error
}

type defaultFieldValidator struct{}

func NewCreateFieldValidator() FieldBehaviorValidator {
	return &defaultFieldValidator{}
}

// This method loops through all the fields in the proto message and
// verifies them against the google field behavior annotations
func (f *defaultFieldValidator) Validate(t proto.Message) error {
	return validate(t.ProtoReflect())
}

// validate recursively validates fields in a proto.Message.
func validate(msg protoreflect.Message) error {
	msgDescriptor := msg.Descriptor()

	for i := 0; i < msgDescriptor.Fields().Len(); i++ {
		fd := msgDescriptor.Fields().Get(i)

		// Check for field behavior annotations
		if options := fd.Options().(*descriptorpb.FieldOptions); options != nil {
			behaviors := proto.GetExtension(options, annotations.E_FieldBehavior).([]annotations.FieldBehavior)
			for _, behavior := range behaviors {
				switch behavior {
				case annotations.FieldBehavior_OUTPUT_ONLY:
					// Check if the field is set and not empty
					msg.Clear(fd)
				case annotations.FieldBehavior_REQUIRED:
					// Check if the field is set and not empty
					if !msg.Has(fd) || fd.HasDefault() {
						return errz.NewInvalidArgument("field %s is required", fd.Name())
					}
				}
			}
		}

		// Recursively validate nested message fields
		if fd.Kind() == protoreflect.MessageKind && !fd.IsMap() && !fd.IsList() {
			nestedMsg := msg.Get(fd).Message()
			if nestedMsg.IsValid() {
				if err := validate(nestedMsg); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
