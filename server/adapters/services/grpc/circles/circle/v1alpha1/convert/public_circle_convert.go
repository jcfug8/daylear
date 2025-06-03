package convert

import (
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
)

// PublicCircleToProto converts a model.Circle to a protobuf PublicCircle
func PublicCircleToProto(publicCircleNamer namer.ReflectNamer[model.Circle], circle model.Circle) (*pb.PublicCircle, error) {
	proto := &pb.PublicCircle{}
	name, err := publicCircleNamer.Format(circle)
	if err != nil {
		return proto, err
	}
	proto.Name = name
	proto.Title = circle.Title
	return proto, nil
}

// PublicCircleListToProto converts a slice of model.Circle to a slice of protobuf PublicCircle
func PublicCircleListToProto(publicCircleNamer namer.ReflectNamer[model.Circle], circles []model.Circle) ([]*pb.PublicCircle, error) {
	protos := make([]*pb.PublicCircle, len(circles))
	for i, circle := range circles {
		proto, err := PublicCircleToProto(publicCircleNamer, circle)
		if err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}
