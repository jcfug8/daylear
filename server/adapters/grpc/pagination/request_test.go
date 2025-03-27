package pagination

import (
	"testing"

	"github.com/stretchr/testify/require"
	freightv1 "go.einride.tech/aip/proto/gen/einride/example/freight/v1"
	"google.golang.org/genproto/googleapis/example/library/v1"
)

func Test_CalculateRequestChecksum(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name  string
		req1  Request
		req2  Request
		equal bool
	}{{
		name: "same request",
		req1: &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  100,
			PageToken: "token",
		},
		req2: &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  100,
			PageToken: "token",
		},
		equal: true,
	}, {
		name: "different parents",
		req1: &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  100,
			PageToken: "token",
		},
		req2: &library.ListBooksRequest{
			Parent:    "shelves/2",
			PageSize:  100,
			PageToken: "token",
		},
		equal: false,
	}, {
		name: "different page sizes",
		req1: &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  100,
			PageToken: "token",
		},
		req2: &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  200,
			PageToken: "token",
		},
		equal: true,
	}, {
		name: "different page tokens",
		req1: &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  100,
			PageToken: "token",
		},
		req2: &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  100,
			PageToken: "token2",
		},
		equal: true,
	}, {
		name: "different skips",
		req1: &freightv1.ListSitesRequest{
			Parent:   "shippers/1",
			PageSize: 100,
			Skip:     0,
		},
		req2: &freightv1.ListSitesRequest{
			Parent:   "shippers/1",
			PageSize: 100,
			Skip:     30,
		},
		equal: true,
	}} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sum1, err := calculateRequestChecksum(tt.req1)
			require.NoError(t, err)

			sum2, err := calculateRequestChecksum(tt.req2)
			require.NoError(t, err)

			if tt.equal {
				require.Equal(t, sum1, sum2)
			} else {
				require.NotEqual(t, sum1, sum2)
			}
		})
	}
}
