package convert

import (
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// ProtoToCircle converts a protobuf Circle to a model Circle
func ProtoToCircle(circleNamer namer.ReflectNamer, proto *pb.Circle) (model.Circle, error) {
	circle := model.Circle{}
	if proto.Name != "" && circleNamer != nil {
		_, err := circleNamer.Parse(proto.Name, &circle)
		if err != nil {
			return circle, err
		}
	}
	circle.Title = proto.Title
	circle.Visibility = proto.Visibility
	return circle, nil
}

// CircleToProto converts a model Circle to a protobuf Circle
func CircleToProto(circleNamer namer.ReflectNamer, circle model.Circle) (*pb.Circle, error) {
	proto := &pb.Circle{}
	name, err := circleNamer.Format(circle)
	if err != nil {
		return proto, err
	}
	proto.Name = name

	proto.Title = circle.Title
	proto.Visibility = circle.Visibility
	return proto, nil
}

// CircleListToProto converts a slice of model Circles to a slice of protobuf Circles
func CircleListToProto(circleNamer namer.ReflectNamer, circles []model.Circle) ([]*pb.Circle, error) {
	protos := make([]*pb.Circle, len(circles))
	for i, circle := range circles {
		proto, err := CircleToProto(circleNamer, circle)
		if err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}

// ProtosToCircle converts a slice of protobuf Circles to a slice of model Circles
func ProtosToCircle(circleNamer namer.ReflectNamer, protos []*pb.Circle) ([]model.Circle, error) {
	res := make([]model.Circle, len(protos))
	for i, proto := range protos {
		circle, err := ProtoToCircle(circleNamer, proto)
		if err != nil {
			return nil, err
		}
		res[i] = circle
	}
	return res, nil
}
