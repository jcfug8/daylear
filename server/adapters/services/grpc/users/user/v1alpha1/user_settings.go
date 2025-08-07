package v1alpha1

import (
	"context"

	convert "github.com/jcfug8/daylear/server/adapters/services/grpc/users/user/v1alpha1/convert"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var userSettingsFieldMap = map[string][]string{
	"name":  {model.UserField_Parent, model.UserField_Id},
	"email": {model.UserField_Email},
}

func (s *UserService) GetUserSettings(ctx context.Context, req *pb.GetUserSettingsRequest) (*pb.UserSettings, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetUserSettings called")

	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	mUser := model.User{}
	_, err = s.userSettingsNamer.Parse(req.GetName(), &mUser)
	if err != nil {
		log.Warn().Err(err).Str("name", req.GetName()).Msg("invalid user settings name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", req.GetName())
	}

	mUser, err = s.domain.GetOwnUser(ctx, authAccount, mUser.Id, s.userSettingsFieldMasker.GetAll())
	if err != nil {
		log.Error().Err(err).Msg("domain.GetUser failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	userSettings, err := convert.UserSettingsToProto(s.userSettingsNamer, mUser)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC GetUserSettings returning successfully")
	return userSettings, nil
}

func (s *UserService) UpdateUserSettings(ctx context.Context, req *pb.UpdateUserSettingsRequest) (*pb.UserSettings, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateUserSettings called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	userSettingsProto := req.GetUserSettings()
	var mUser model.User
	_, err = s.userSettingsNamer.Parse(userSettingsProto.GetName(), &mUser)
	if err != nil {
		log.Warn().Err(err).Str("name", userSettingsProto.GetName()).Msg("invalid user settings name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", userSettingsProto.GetName())
	}

	// No updatable fields, so just return current settings
	// TODO: add update mask and ability to update fields when there come
	mUser, err = s.domain.GetUser(ctx, authAccount, mUser.Parent, mUser.Id, s.userSettingsFieldMasker.GetAll())
	if err != nil {
		log.Error().Err(err).Msg("domain.GetUser failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	userSettings, err := convert.UserSettingsToProto(s.userSettingsNamer, mUser)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	log.Info().Msg("gRPC UpdateUserSettings returning successfully (no fields updated)")
	return userSettings, nil
}
