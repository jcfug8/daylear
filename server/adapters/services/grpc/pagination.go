package grpc

import (
	"go.einride.tech/aip/pagination"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

type PaginationConfig struct {
	DefaultPageSize int32
	MaxPageSize     int32
}

type PaginationRequest interface {
	proto.Message
	GetPageSize() int32
	GetPageToken() string
}

func SetupPagination(request PaginationRequest, config PaginationConfig) (pagination.PageToken, int32, error) {
	pageToken, err := pagination.ParsePageToken(request)
	if err != nil {
		return pagination.PageToken{}, 0, status.Errorf(codes.InvalidArgument, "invalid page token")
	}

	pageSize := request.GetPageSize()
	// Set default page size if not specified
	if pageSize == 0 {
		pageSize = config.DefaultPageSize
	}

	// Cap at max page size
	if pageSize > config.MaxPageSize {
		pageSize = config.MaxPageSize
	}

	return pageToken, pageSize, nil
}
