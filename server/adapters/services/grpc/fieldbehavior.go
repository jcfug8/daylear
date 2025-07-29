package grpc

import (
	"go.einride.tech/aip/fieldbehavior"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func ProcessRequestFieldBehavior(request proto.Message) error {
	fieldbehavior.ClearFields(request, annotations.FieldBehavior_OUTPUT_ONLY)
	err := fieldbehavior.ValidateRequiredFields(request)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}
	return nil
}

func ProcessUpdateRequestFieldBehavior(request proto.Message) error {
	fieldbehavior.ClearFields(request, annotations.FieldBehavior_OUTPUT_ONLY)
	return nil
}

func ProcessResponseFieldBehavior(response proto.Message) {
	fieldbehavior.ClearFields(response, annotations.FieldBehavior_INPUT_ONLY)
}
