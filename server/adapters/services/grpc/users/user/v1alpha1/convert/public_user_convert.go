package convert

import (
	namer "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/namer"
	model "github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// ProtoToPublicUser converts a protobuf User to a model User
func ProtoToPublicUser(UserNamer namer.UserNamer, proto *pb.PublicUser) (model.User, error) {
	user := model.User{}
	if proto.Name != "" {
		id, err := UserNamer.Parse(proto.Name)
		if err != nil {
			return user, err
		}
		user.Id = id
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceProtoToDomain
	user.Username = proto.Username
	user.GivenName = proto.GivenName
	user.FamilyName = proto.FamilyName
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return user, nil
}

// PublicUserToProto converts a model User to a protobuf PublicUser
func PublicUserToProto(UserNamer namer.UserNamer, user model.User) (*pb.PublicUser, error) {
	proto := &pb.PublicUser{}
	name, err := UserNamer.Format(user.Id)
	if err != nil {
		return proto, err
	}
	proto.Name = name

	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceDomainToProto
	proto.Username = user.Username
	proto.GivenName = user.GivenName
	proto.FamilyName = user.FamilyName
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return proto, nil
}

// PublicUserListToProto converts a slice of model Users to a slice of protobuf PublicUsers
func PublicUserListToProto(UserNamer namer.UserNamer, users []model.User) ([]*pb.PublicUser, error) {
	protos := make([]*pb.PublicUser, len(users))
	for i, user := range users {
		proto := &pb.PublicUser{}
		var err error
		if proto, err = PublicUserToProto(UserNamer, user); err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}

// ProtosToPublicUser converts a slice of protobuf PublicUsers to a slice of model Users
func ProtosToPublicUser(UserNamer namer.UserNamer, protos []*pb.PublicUser) ([]model.User, error) {
	res := make([]model.User, len(protos))
	for i, proto := range protos {
		user := model.User{}
		var err error
		if user, err = ProtoToPublicUser(UserNamer, proto); err != nil {
			return nil, err
		}
		res[i] = user
	}
	return res, nil
}
