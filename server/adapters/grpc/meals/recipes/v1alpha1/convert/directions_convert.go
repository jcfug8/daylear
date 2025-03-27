package convert

import (
	model "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// ProtoToDirection converts a proto Directions to a domain Directions.
func ProtoToDirection(proto *pb.Recipe_Direction) model.RecipeDirection {
	direction := model.RecipeDirection{}
	direction.Steps = proto.Steps
	direction.Title = proto.Title
	return direction
}

// DirectionToProto converts a domain Directions to a proto Directions.
func DirectionToProto(direction model.RecipeDirection) *pb.Recipe_Direction {
	return &pb.Recipe_Direction{
		Steps: direction.Steps,
		Title: direction.Title,
	}
}

// ProtosToDirections converts a slice of proto Directions to a slice of domain Directions.
func ProtosToDirections(protos []*pb.Recipe_Direction) []model.RecipeDirection {
	directions := make([]model.RecipeDirection, len(protos))
	for i, proto := range protos {
		direction := ProtoToDirection(proto)
		directions[i] = direction
	}
	return directions
}

// DirectionsToProtos converts a slice of domain Directions to a slice of proto Directions.
func DirectionsToProtos(directions []model.RecipeDirection) []*pb.Recipe_Direction {
	protos := make([]*pb.Recipe_Direction, len(directions))
	for i, direction := range directions {
		proto := DirectionToProto(direction)
		protos[i] = proto
	}
	return protos
}
