package convert

import (
	namer "github.com/jcfug8/daylear/server/adapters/grpc/users/user/v1alpha1/namer"
	model "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// ProtoToUser converts a protobuf User to a model User
func ProtoToUser(UserNamer namer.UserNamer, proto *pb.User) (model.User, error) {
	user := model.User{}
	if proto.Name != "" {
		id, err := UserNamer.Parse(proto.Name)
		if err != nil {
			return user, err
		}
		user.Id = id

	}

	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceProtoToDomain
	user.Email = proto.Email
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return user, nil
}

// UserToProto converts a model User to a protobuf User
func UserToProto(UserNamer namer.UserNamer, user model.User) (*pb.User, error) {
	proto := &pb.User{}
	name, err := UserNamer.Format(user.Id)
	if err != nil {
		return proto, err
	}
	proto.Name = name

	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceDomainToProto
	proto.Email = user.Email
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return proto, nil
}

// UserListToProto converts a slice of model Users to a slice of protobuf OmniUsers
func UserListToProto(UserNamer namer.UserNamer, users []model.User) ([]*pb.User, error) {
	protos := make([]*pb.User, len(users))
	for i, user := range users {
		proto := &pb.User{}
		var err error
		if proto, err = UserToProto(UserNamer, user); err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}

// ProtosToUser converts a slice of protobuf OmniUsers to a slice of model Users
func ProtosToUser(UserNamer namer.UserNamer, protos []*pb.User) ([]model.User, error) {
	res := make([]model.User, len(protos))
	for i, proto := range protos {
		user := model.User{}
		var err error
		if user, err = ProtoToUser(UserNamer, proto); err != nil {
			return nil, err
		}
		res[i] = user
	}
	return res, nil
}
