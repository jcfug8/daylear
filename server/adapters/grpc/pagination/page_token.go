package pagination

import (
	"github.com/jcfug8/daylear/server/core/model"
	"go.einride.tech/aip/pagination"
)

// EncodePageToken encodes a page token to a string.
func EncodePageToken[T any](t *model.PageToken[T]) string {
	if t == nil {
		return ""
	}

	return pagination.EncodePageTokenStruct(t)
}
