package pagination

import (
	"github.com/jcfug8/daylear/server/core/errz"
	"github.com/jcfug8/daylear/server/core/model"

	"go.einride.tech/aip/pagination"
)

// refer to https://pkg.go.dev/go.einride.tech/aip/pagination

// pageTokenChecksumMask is a random bitmask applied to offset-based page token
// checksums.
//
// Change the bitmask to force checksum failures when changing the page token
// implementation.
const pageTokenChecksumMask uint32 = 0x9acb0442

// ParsePageToken returns a PageToken with the offset (skip) as well as the
// tail of the previous page (if any).
// Refer to https://google.aip.dev/158#skipping-results
func ParsePageToken[T any](req Request) (_ *model.PageToken[T], err error) {
	checksum, err := calculateRequestChecksum(req)
	if err != nil {
		return nil, err
	}

	skip := int32(0)
	if s, ok := req.(skipRequest); ok {
		skip = s.GetSkip()
	}

	checksum ^= pageTokenChecksumMask // apply checksum mask for PageToken
	if req.GetPageToken() == "" {
		return model.NewPageToken[T](req.GetPageSize(), skip, checksum, nil), nil
	}

	pageToken := &model.PageToken[T]{}
	err = pagination.DecodePageTokenStruct(req.GetPageToken(), pageToken)
	if err != nil {
		return nil, err
	}

	if pageToken.RequestChecksum != checksum {
		return nil, errz.NewInvalidArgument(
			"checksum mismatch (got 0x%x but expected 0x%x)",
			pageToken.RequestChecksum, checksum)
	}

	pageToken.Skip = skip

	return pageToken, nil
}
