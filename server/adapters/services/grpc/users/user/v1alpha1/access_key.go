package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/grpc"
	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/logutil"
	"github.com/jcfug8/daylear/server/core/model"
	"github.com/jcfug8/daylear/server/core/namer"
	pb "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	accessKeyMaxPageSize     int32 = 100
	accessKeyDefaultPageSize int32 = 20
)

var accessKeyFieldMap = map[string][]string{
	"name":        {model.AccessKeyField_Parent, model.AccessKeyField_AccessKeyId},
	"title":       {model.AccessKeyField_Title},
	"description": {model.AccessKeyField_Description},
}

// CreateAccessKey creates a new access key
func (s *UserService) CreateAccessKey(ctx context.Context, request *pb.CreateAccessKeyRequest) (response *pb.AccessKey, err error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC CreateAccessKey called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessRequestFieldBehavior(request)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, err
	}

	// convert proto to model
	accessKeyProto := request.GetAccessKey()
	accessKeyProto.Name = ""
	_, mAccessKey, err := s.ProtoToAccessKey(accessKeyProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	_, err = s.accessKeyNamer.ParseParent(request.GetParent(), &mAccessKey.Parent)
	if err != nil {
		log.Warn().Err(err).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	// create access key
	mAccessKey, err = s.domain.CreateAccessKey(ctx, authAccount, mAccessKey)
	if err != nil {
		log.Error().Err(err).Msg("domain.CreateAccessKey failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	accessKeyProto, err = s.AccessKeyToProto(mAccessKey)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(accessKeyProto)
	log.Info().Msg("gRPC CreateAccessKey returning successfully")
	return accessKeyProto, nil
}

// DeleteAccessKey deletes an access key
func (s *UserService) DeleteAccessKey(ctx context.Context, request *pb.DeleteAccessKeyRequest) (*emptypb.Empty, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC DeleteAccessKey called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mAccessKey model.AccessKey
	_, err = s.accessKeyNamer.Parse(request.GetName(), &mAccessKey)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	_, err = s.domain.DeleteAccessKey(ctx, authAccount, mAccessKey.AccessKeyId)
	if err != nil {
		log.Error().Err(err).Msg("domain.DeleteAccessKey failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Info().Msg("gRPC DeleteAccessKey returning successfully")
	return &emptypb.Empty{}, nil
}

// GetAccessKey retrieves an access key
func (s *UserService) GetAccessKey(ctx context.Context, request *pb.GetAccessKeyRequest) (*pb.AccessKey, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC GetAccessKey called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	var mAccessKey model.AccessKey
	_, err = s.accessKeyNamer.Parse(request.GetName(), &mAccessKey)
	if err != nil {
		log.Warn().Err(err).Str("name", request.GetName()).Msg("invalid name")
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", request.GetName())
	}

	mAccessKey, err = s.domain.GetAccessKey(ctx, authAccount, mAccessKey.Parent, mAccessKey.AccessKeyId, nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.GetAccessKey failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	accessKeyProto, err := s.AccessKeyToProto(mAccessKey)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(accessKeyProto)
	log.Info().Msg("gRPC GetAccessKey returning successfully")
	return accessKeyProto, nil
}

// ListAccessKeys lists access keys for a user
func (s *UserService) ListAccessKeys(ctx context.Context, request *pb.ListAccessKeysRequest) (*pb.ListAccessKeysResponse, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC ListAccessKeys called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// parse parent
	var mAccessKey model.AccessKey
	_, err = s.accessKeyNamer.ParseParent(request.GetParent(), &mAccessKey)
	if err != nil {
		log.Warn().Err(err).Str("parent", request.GetParent()).Msg("invalid parent")
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", request.GetParent())
	}

	pageToken, pageSize, err := grpc.SetupPagination(request, grpc.PaginationConfig{
		DefaultPageSize: accessKeyDefaultPageSize,
		MaxPageSize:     accessKeyMaxPageSize,
	})
	if err != nil {
		log.Warn().Err(err).Msg("pagination setup failed")
		return nil, err
	}

	// list access keys
	mAccessKeys, err := s.domain.ListAccessKeys(ctx, authAccount, mAccessKey.Parent.UserId, pageSize, pageToken.Offset, request.GetFilter(), nil)
	if err != nil {
		log.Error().Err(err).Msg("domain.ListAccessKeys failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert models to protos
	accessKeyProtos := make([]*pb.AccessKey, len(mAccessKeys))
	for i, mAccessKey := range mAccessKeys {
		accessKeyProto, err := s.AccessKeyToProto(mAccessKey)
		if err != nil {
			log.Error().Err(err).Msg("unable to prepare response")
			return nil, status.Error(codes.Internal, "unable to prepare response")
		}
		accessKeyProtos[i] = accessKeyProto
	}

	// check field behavior
	for _, accessKeyProto := range accessKeyProtos {
		grpc.ProcessResponseFieldBehavior(accessKeyProto)
	}

	// create response
	response := &pb.ListAccessKeysResponse{
		AccessKeys: accessKeyProtos,
	}

	// add next page token if there are more results
	if len(mAccessKeys) == int(pageSize) {
		response.NextPageToken = pageToken.Next(request).String()
	}

	log.Info().Msg("gRPC ListAccessKeys returning successfully")
	return response, nil
}

// UpdateAccessKey updates an access key
func (s *UserService) UpdateAccessKey(ctx context.Context, request *pb.UpdateAccessKeyRequest) (*pb.AccessKey, error) {
	log := logutil.EnrichLoggerWithContext(s.log, ctx)
	log.Info().Msg("gRPC UpdateAccessKey called")
	authAccount, err := headers.ParseAuthData(ctx)
	if err != nil {
		log.Warn().Err(err).Msg("failed to parse auth data")
		return nil, err
	}

	// check field behavior
	err = grpc.ProcessUpdateRequestFieldBehavior(request)
	if err != nil {
		log.Warn().Err(err).Msg("invalid request data")
		return nil, err
	}

	// convert proto to model
	accessKeyProto := request.GetAccessKey()
	_, mAccessKey, err := s.ProtoToAccessKey(accessKeyProto)
	if err != nil {
		log.Warn().Err(err).Msg("unable to convert proto to model")
		return nil, status.Error(codes.InvalidArgument, "invalid request data")
	}

	// get update mask
	fieldMask := request.GetUpdateMask()
	updateMask := s.accessKeyFieldMasker.Convert(fieldMask.GetPaths())

	// update access key
	mAccessKey, err = s.domain.UpdateAccessKey(ctx, authAccount, mAccessKey, updateMask)
	if err != nil {
		log.Error().Err(err).Msg("domain.UpdateAccessKey failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	// convert model to proto
	accessKeyProto, err = s.AccessKeyToProto(mAccessKey)
	if err != nil {
		log.Error().Err(err).Msg("unable to prepare response")
		return nil, status.Error(codes.Internal, "unable to prepare response")
	}

	// check field behavior
	grpc.ProcessResponseFieldBehavior(accessKeyProto)
	log.Info().Msg("gRPC UpdateAccessKey returning successfully")
	return accessKeyProto, nil
}

// ProtoToAccessKey converts a proto AccessKey to a model AccessKey
func (s *UserService) ProtoToAccessKey(proto *pb.AccessKey) (nameIndex int, accessKey model.AccessKey, err error) {
	accessKey = model.AccessKey{
		Title:       proto.GetTitle(),
		Description: proto.GetDescription(),
	}

	// Parse parent from name if provided
	if proto.GetName() != "" {
		nameIndex, err = s.accessKeyNamer.Parse(proto.GetName(), &accessKey)
		if err != nil {
			return 0, model.AccessKey{}, err
		}
	}

	return nameIndex, accessKey, nil
}

// AccessKeyToProto converts a model AccessKey to a proto AccessKey
func (s *UserService) AccessKeyToProto(accessKey model.AccessKey, options ...namer.FormatReflectNamerOption) (*pb.AccessKey, error) {
	proto := &pb.AccessKey{
		Title:                accessKey.Title,
		Description:          accessKey.Description,
		UnencryptedAccessKey: accessKey.UnencryptedAccessKey,
	}

	// Generate name
	if accessKey.AccessKeyId.AccessKeyId != 0 {
		name, err := s.accessKeyNamer.Format(accessKey, options...)
		if err != nil {
			return nil, err
		}
		proto.Name = name
	}

	return proto, nil
}
