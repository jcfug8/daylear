package v1alpha1

import (
	"context"

	"github.com/jcfug8/daylear/server/adapters/services/http/libs/headers"
	"github.com/jcfug8/daylear/server/core/model"
	pb "github.com/jcfug8/daylear/server/genapi/api/circles/circle/v1alpha1"
	"go.einride.tech/aip/fieldbehavior"
	"go.einride.tech/aip/pagination"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

const (
	accessMaxPageSize     int32 = 1000
	accessDefaultPageSize int32 = 100
)

func (s *CircleService) CreateAccess(ctx context.Context, request *pb.CreateAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// check field behavior
	fieldbehavior.ClearFields(request, annotations.FieldBehavior_OUTPUT_ONLY)
	err = fieldbehavior.ValidateRequiredFields(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request data: %v", err)
	}

	// convert proto to model
	pbAccess := request.GetAccess()
	pbAccess.Name = ""
	modelAccess, err := s.convertAccessProtoToModel(pbAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid access: %v", err)
	}

	// parse parent
	_, err = s.accessNamer.ParseParent(request.Parent, &modelAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// create access
	createdAccess, err := s.domain.CreateCircleAccess(ctx, authAccount, modelAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access: %v", err)
	}

	// convert model to proto
	pbAccess, err = s.convertAccessModelToProto(createdAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
	}

	// check field behavior
	fieldbehavior.ClearFields(pbAccess, annotations.FieldBehavior_INPUT_ONLY)

	return pbAccess, nil
}

func (s *CircleService) DeleteAccess(ctx context.Context, request *pb.DeleteAccessRequest) (*emptypb.Empty, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// parse name
	circleAccess := &model.CircleAccess{}
	_, err = s.accessNamer.Parse(request.Name, circleAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// delete access
	err = s.domain.DeleteCircleAccess(ctx, authAccount, circleAccess.CircleAccessParent, circleAccess.CircleAccessId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete access: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *CircleService) GetAccess(ctx context.Context, request *pb.GetAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// parse name
	circleAccess := &model.CircleAccess{}
	_, err = s.accessNamer.Parse(request.Name, circleAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// get access
	access, err := s.domain.GetCircleAccess(ctx, authAccount, circleAccess.CircleAccessParent, circleAccess.CircleAccessId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get access: %v", err)
	}

	// convert model to proto
	pbAccess, err := s.convertAccessModelToProto(access)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
	}

	// check field behavior
	fieldbehavior.ClearFields(pbAccess, annotations.FieldBehavior_INPUT_ONLY)

	return pbAccess, nil
}

func (s *CircleService) ListAccesses(ctx context.Context, request *pb.ListAccessesRequest) (*pb.ListAccessesResponse, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// set requester
	var circleAccessParent model.CircleAccessParent

	// parse parent
	_, err = s.accessNamer.ParseParent(request.Parent, &circleAccessParent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
	}

	// pagination
	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	if request.GetPageSize() == 0 {
		request.PageSize = accessDefaultPageSize
	}
	request.PageSize = min(request.PageSize, accessMaxPageSize)

	// list accesses
	accesses, err := s.domain.ListCircleAccesses(ctx, authAccount, circleAccessParent, request.GetPageSize(), pageToken.Offset, request.GetFilter())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list accesses: %v", err)
	}

	// convert models to protos
	pbAccesses := make([]*pb.Access, len(accesses))
	for i, access := range accesses {
		pbAccess, err := s.convertAccessModelToProto(access)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
		}
		pbAccesses[i] = pbAccess
	}

	// return response
	return &pb.ListAccessesResponse{
		NextPageToken: pageToken.Next(request).String(),
		Accesses:      pbAccesses,
	}, nil
}

func (s *CircleService) UpdateAccess(ctx context.Context, request *pb.UpdateAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// convert proto to model
	pbAccess := request.GetAccess()
	modelAccess, err := s.convertAccessProtoToModel(pbAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid access: %v", err)
	}

	// update access
	updatedAccess, err := s.domain.UpdateCircleAccess(ctx, authAccount, modelAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update access: %v", err)
	}

	// convert model to proto
	pbAccess, err = s.convertAccessModelToProto(updatedAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
	}

	return pbAccess, nil
}

func (s *CircleService) AcceptAccess(ctx context.Context, request *pb.AcceptAccessRequest) (*pb.Access, error) {
	authAccount, err := headers.ParseAuthData(ctx, s.circleNamer)
	if err != nil {
		return nil, err
	}

	// parse name
	circleAccess := &model.CircleAccess{}
	_, err = s.accessNamer.Parse(request.Name, circleAccess)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
	}

	// accept access
	acceptedAccess, err := s.domain.AcceptCircleAccess(ctx, authAccount, circleAccess.CircleAccessParent, circleAccess.CircleAccessId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to accept access: %v", err)
	}

	// convert model to proto
	pbAccess, err := s.convertAccessModelToProto(acceptedAccess)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to convert model to proto: %v", err)
	}

	return pbAccess, nil
}

// Helper conversion functions (to be implemented)
func (s *CircleService) convertAccessProtoToModel(pbAccess *pb.Access) (model.CircleAccess, error) {
	// TODO: Implement conversion logic
	return model.CircleAccess{}, nil
}

func (s *CircleService) convertAccessModelToProto(modelAccess model.CircleAccess) (*pb.Access, error) {
	// TODO: Implement conversion logic
	return &pb.Access{}, nil
}

// TODO: Add accessNamer field to CircleService
