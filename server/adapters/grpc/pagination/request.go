package pagination

import (
	"hash/crc32"

	"github.com/jcfug8/daylear/server/core/errz"

	"google.golang.org/protobuf/proto"
)

// Request defines the interface for all pagination requests.
type Request interface {
	proto.Message
	GetPageSize() int32
	GetPageToken() string
}

type skipRequest interface {
	proto.Message
	// GetSkip returns the skip of the request.
	// See: https://google.aip.dev/158#skipping-results
	GetSkip() int32
}

func calculateRequestChecksum(req Request) (uint32, error) {
	cloned := proto.Clone(req)

	r := cloned.ProtoReflect()
	r.Clear(r.Descriptor().Fields().ByName("page_token"))
	r.Clear(r.Descriptor().Fields().ByName("page_size"))

	if _, ok := req.(skipRequest); ok {
		r.Clear(r.Descriptor().Fields().ByName("skip"))
	}

	data, err := proto.Marshal(cloned)
	if err != nil {
		return 0, errz.NewInternal("request checksum error: %v", err)
	}

	return crc32.ChecksumIEEE(data), nil
}
