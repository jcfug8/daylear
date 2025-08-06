package convert

import (
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	// IRIOMO:CUSTOM_CODE_SLOT_START convertResourceImports
	// IRIOMO:CUSTOM_CODE_SLOT_END
)

// ProtoToUser converts a protobuf User to a model User
func ProtoToUser(userNamer namer.ReflectNamer, accessNamer namer.ReflectNamer, proto *pb.User) (model.User, error) {
	user := model.User{}
	if proto.Name != "" {
		_, err := userNamer.Parse(proto.Name, &user)
		if err != nil {
			return user, err
		}
	}

	user.Username = proto.Username
	user.GivenName = proto.GivenName
	user.FamilyName = proto.FamilyName
	user.ImageUri = proto.ImageUri
	user.Bio = proto.Bio

	return user, nil
}

// UserToProto converts a model User to a protobuf User
func UserToProto(userNamer namer.ReflectNamer, accessNamer namer.ReflectNamer, user model.User, options ...namer.FormatReflectNamerOption) (*pb.User, error) {
	proto := &pb.User{}
	name, err := userNamer.Format(user, options...)
	if err != nil {
		return proto, err
	}
	proto.Name = name

	proto.Username = user.Username
	proto.GivenName = user.GivenName
	proto.FamilyName = user.FamilyName
	proto.ImageUri = user.ImageUri
	proto.Bio = user.Bio

	proto.Access = &pb.User_Access{
		PermissionLevel: user.UserAccess.PermissionLevel,
		State:           user.UserAccess.State,
	}
	if user.UserAccess.UserAccessId.UserAccessId != 0 {
		name, err := accessNamer.Format(user.UserAccess)
		if err != nil {
			return proto, err
		}
		proto.Access.Name = name
	}
	if user.UserAccess.Requester.UserId != 0 {
		requesterName, err := userNamer.Format(user.UserAccess.Requester)
		if err != nil {
			return proto, err
		}
		proto.Access.Requester = requesterName
	}

	return proto, nil
}
