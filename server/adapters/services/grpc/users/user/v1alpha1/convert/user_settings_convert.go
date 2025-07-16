package convert

import (
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
)

// ProtoToUserSettings converts a protobuf UserSettings to a model User
func ProtoToUserSettings(userSettingsNamer namer.ReflectNamer, proto *pb.UserSettings) (model.User, error) {
	user := model.User{}
	if proto.Name != "" {
		_, err := userSettingsNamer.Parse(proto.Name, &user)
		if err != nil {
			return user, err
		}
	}
	user.Email = proto.Email
	return user, nil
}

// UserSettingsToProto converts a model User to a protobuf UserSettings
func UserSettingsToProto(userSettingsNamer namer.ReflectNamer, user model.User) (*pb.UserSettings, error) {
	proto := &pb.UserSettings{}
	name, err := userSettingsNamer.Format(user)
	if err != nil {
		return proto, err
	}
	proto.Name = name
	proto.Email = user.Email
	return proto, nil
}
