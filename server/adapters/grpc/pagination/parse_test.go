package pagination

import (
	"testing"

	"github.com/jcfug8/daylear/server/core/model"
	"github.com/stretchr/testify/require"
	"go.einride.tech/aip/pagination"
	freightv1 "go.einride.tech/aip/proto/gen/einride/example/freight/v1"
	"google.golang.org/genproto/googleapis/example/library/v1"
)

func Test_Parse(t *testing.T) {
	t.Parallel()

	type book struct {
		Title string
	}

	t.Run("valid checksums", func(t *testing.T) {
		t.Parallel()
		req1 := &library.ListBooksRequest{
			Parent:   "shelves/1",
			PageSize: 3,
		}
		token1, err := ParsePageToken[book](req1)
		require.NoError(t, err)
		require.Equal(t, int32(3), token1.PageSize)
		require.Zero(t, token1.Skip)
		require.Nil(t, token1.Tail)

		page1 := []book{{
			Title: "foo",
		}, {
			Title: "bar",
		}, {
			Title: "baz",
		}}

		req2 := &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  6,
			PageToken: EncodePageToken(token1.Next(page1)),
		}
		token2, err := ParsePageToken[book](req2)
		require.NoError(t, err)
		require.Equal(t, int32(3), token2.PageSize)
		require.Zero(t, token2.Skip)
		require.Equal(t, page1[2], token2.Tail)

		page2 := []book{{
			Title: "foobar",
		}}

		token3 := token2.Next(page2)
		require.Nil(t, token3)
		require.Zero(t, EncodePageToken(token3))
	})

	t.Run("skip", func(t *testing.T) {
		t.Parallel()

		t.Run("docs example 1", func(t *testing.T) {
			t.Parallel()

			// From https://google.aip.dev/158:
			// A request with no page token and a skip value of 30 returns a
			// single page of results starting with the 31st result.
			token, err := ParsePageToken[book](&freightv1.ListSitesRequest{
				Parent: "shippers/1",
				Skip:   30,
			})
			require.NoError(t, err)
			require.Zero(t, token.PageSize)
			require.Equal(t, int32(30), token.Skip) // 31st result
			require.Nil(t, token.Tail)
		})

		t.Run("docs example 2", func(t *testing.T) {
			t.Parallel()

			// From https://google.aip.dev/158:
			// A request with a page token corresponding to the 51st result
			// (because the first 50 results were returned on the first page)
			// and a skip value of 30 returns a single page of results starting
			// with the 81st result.
			req1 := &freightv1.ListSitesRequest{
				Parent:   "shippers/1",
				PageSize: 50,
			}
			token1, err := ParsePageToken[book](req1)
			require.NoError(t, err)
			require.Equal(t, int32(50), token1.PageSize)
			require.Zero(t, token1.Skip)
			require.Nil(t, token1.Tail)

			page := make([]book, 50)
			page[49] = book{Title: "foo"}

			req2 := &freightv1.ListSitesRequest{
				Parent:    "shippers/1",
				Skip:      30,
				PageSize:  50,
				PageToken: EncodePageToken(token1.Next(page)),
			}
			token2, err := ParsePageToken[book](req2)
			require.NoError(t, err)
			require.Equal(t, int32(50), token2.PageSize)
			require.Equal(t, int32(30), token2.Skip)
			require.Equal(t, page[49], token2.Tail)
		})

		t.Run("handle empty token with skip", func(t *testing.T) {
			t.Parallel()

			req1 := &freightv1.ListSitesRequest{
				Parent:   "shippers/1",
				Skip:     30,
				PageSize: 20,
			}
			token1, err := ParsePageToken[book](req1)
			require.NoError(t, err)
			require.Equal(t, int32(20), token1.PageSize)
			require.Equal(t, int32(30), token1.Skip)
			require.Nil(t, token1.Tail)
		})

		t.Run("handle existing token with another skip", func(t *testing.T) {
			t.Parallel()

			req1 := &freightv1.ListSitesRequest{
				Parent:   "shippers/1",
				Skip:     50,
				PageSize: 20,
			}
			token1, err := ParsePageToken[book](req1)
			require.NoError(t, err)
			require.Equal(t, int32(20), token1.PageSize)
			require.Equal(t, int32(50), token1.Skip)
			require.Nil(t, token1.Tail)

			req2 := &freightv1.ListSitesRequest{
				Parent:    "shippers/1",
				Skip:      30,
				PageToken: EncodePageToken(token1),
			}
			token2, err := ParsePageToken[book](req2)
			require.NoError(t, err)
			require.Equal(t, int32(20), token2.PageSize)
			require.Equal(t, int32(30), token2.Skip)
			require.Nil(t, token2.Tail)

			page := make([]book, 50)
			page[49] = book{Title: "foo"}

			token3 := token2.Next(page)
			require.Equal(t, int32(20), token3.PageSize)
			require.Zero(t, token3.Skip)
			require.Equal(t, page[49], token3.Tail)
		})
	})

	t.Run("invalid format", func(t *testing.T) {
		t.Parallel()

		req1 := &library.ListBooksRequest{
			Parent:    "shelves/1",
			PageSize:  10,
			PageToken: "invalid",
		}
		token1, err := ParsePageToken[book](req1)
		require.ErrorContains(t, err, "decode")
		require.Nil(t, token1)
	})

	t.Run("invalid checksum", func(t *testing.T) {
		t.Parallel()

		req1 := &library.ListBooksRequest{
			Parent:   "shelves/1",
			PageSize: 10,
			PageToken: pagination.EncodePageTokenStruct(model.PageToken[book]{
				Skip:            100,
				RequestChecksum: 1234, // invalid
			}),
		}
		token1, err := ParsePageToken[book](req1)
		require.ErrorContains(t, err, "checksum")
		require.Nil(t, token1)
	})
}
