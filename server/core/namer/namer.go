package namer

import (
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type parentType interface {
	GetVars(patternIndex int) []string
	SetVars(patternIndex int, vars []string) parentType
}

type idType interface {
	GetId(patternIndex int) string
	SetId(patternIndex int, id string) idType
}

func getPatterns(resource proto.Message) []string {
	resourceOption := proto.GetExtension(
		resource.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)
	return resourceOption.Pattern
}
