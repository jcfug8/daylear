package convert

import (
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// ProtoToUser converts a protobuf User to a model User
func ProtoToUser(userNamer namer.ReflectNamer, proto *pb.User) (model.User, error) {
	user := model.User{}
	if proto.Name != "" {
		_, err := userNamer.Parse(proto.Name, &user)
		if err != nil {
			return user, err
		}
	}

	user.Email = proto.Email
	user.Username = proto.Username
	user.GivenName = proto.GivenName
	user.FamilyName = proto.FamilyName

	return user, nil
}

// UserToProto converts a model User to a protobuf User
func UserToProto(userNamer namer.ReflectNamer, user model.User) (*pb.User, error) {
	proto := &pb.User{}
	name, err := userNamer.Format(user)
	if err != nil {
		return proto, err
	}
	proto.Name = name

	proto.Email = user.Email
	proto.Username = user.Username
	proto.GivenName = user.GivenName
	proto.FamilyName = user.FamilyName

	return proto, nil
}

// UserListToProto converts a slice of model Users to a slice of protobuf OmniUsers
func UserListToProto(userNamer namer.ReflectNamer, users []model.User) ([]*pb.User, error) {
	protos := make([]*pb.User, len(users))
	for i, user := range users {
		proto := &pb.User{}
		var err error
		if proto, err = UserToProto(userNamer, user); err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}

// ProtosToUser converts a slice of protobuf OmniUsers to a slice of model Users
func ProtosToUser(userNamer namer.ReflectNamer, protos []*pb.User) ([]model.User, error) {
	res := make([]model.User, len(protos))
	for i, proto := range protos {
		user := model.User{}
		var err error
		if user, err = ProtoToUser(userNamer, proto); err != nil {
			return nil, err
		}
		res[i] = user
	}
	return res, nil
}
