package convert

import (
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// ProtoToPublicUser converts a protobuf User to a model User
func ProtoToPublicUser(userNamer namer.ReflectNamer[model.User], proto *pb.PublicUser) (model.User, error) {
	user := model.User{}
	if proto.Name != "" {
		_, err := userNamer.Parse(proto.Name, &user)
		if err != nil {
			return user, err
		}
	}

	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceProtoToDomain
	user.Username = proto.Username
	user.GivenName = proto.GivenName
	user.FamilyName = proto.FamilyName
	// IRIOMO:CUSTOM_CODE_SLOT_END

	return user, nil
}

// PublicUserToProto converts a model User to a protobuf PublicUser
func PublicUserToProto(userNamer namer.ReflectNamer[model.User], user model.User) (*pb.PublicUser, error) {
	proto := &pb.PublicUser{}
	name, err := userNamer.Format(user)
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
func PublicUserListToProto(publicUserNamer namer.ReflectNamer[model.User], users []model.User) ([]*pb.PublicUser, error) {
	protos := make([]*pb.PublicUser, len(users))
	for i, user := range users {
		proto := &pb.PublicUser{}
		var err error
		if proto, err = PublicUserToProto(publicUserNamer, user); err != nil {
			return nil, err
		}
		protos[i] = proto
	}
	return protos, nil
}

// ProtosToPublicUser converts a slice of protobuf PublicUsers to a slice of model Users
func ProtosToPublicUser(publicUserNamer namer.ReflectNamer[model.User], protos []*pb.PublicUser) ([]model.User, error) {
	res := make([]model.User, len(protos))
	for i, proto := range protos {
		user := model.User{}
		var err error
		if user, err = ProtoToPublicUser(publicUserNamer, proto); err != nil {
			return nil, err
		}
		res[i] = user
	}
	return res, nil
}
