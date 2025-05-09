package convert

import (
	namer "github.com/jcfug8/daylear/server/adapters/services/grpc/circles/circle/v1alpha1/namer"
	model "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// ProtoToCircle converts a protobuf Circle to a model Circle
func ProtoToCircle(CircleNamer namer.CircleNamer, proto *pb.Circle) (model.Circle, error) {
	circle := model.Circle{}
	if proto.Name != "" {
		parent, id, err := CircleNamer.Parse(proto.Name)
		if err != nil {
			return circle, err
		}
		circle.Parent = parent
		circle.Id = id
	}
	circle.Title = proto.Title
	return circle, nil
}

// CircleToProto converts a model Circle to a protobuf Circle
func CircleToProto(CircleNamer namer.CircleNamer, circle model.Circle) (*pb.Circle, error) {
	proto := &pb.Circle{}
	name, err := CircleNamer.Format(circle.Parent, circle.Id)
	if err != nil {
		return proto, err
	}
	proto.Name = name
	proto.Title = circle.Title
	return proto, nil
}

// CircleListToProto converts a slice of model Circles to a slice of protobuf Circles
func CircleListToProto(CircleNamer namer.CircleNamer, circles []model.Circle) ([]*pb.Circle, error) {
	protos := make([]*pb.Circle, len(circles))
	for i, circle := range circles {
		proto, err := CircleToProto(CircleNamer, circle)
		if err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}

// ProtosToCircle converts a slice of protobuf Circles to a slice of model Circles
func ProtosToCircle(CircleNamer namer.CircleNamer, protos []*pb.Circle) ([]model.Circle, error) {
	res := make([]model.Circle, len(protos))
	for i, proto := range protos {
		circle, err := ProtoToCircle(CircleNamer, proto)
		if err != nil {
			return nil, err
		}
		res[i] = circle
	}
	return res, nil
}
