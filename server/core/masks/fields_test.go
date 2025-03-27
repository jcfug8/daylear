package masks_test

import (
	"slices"
	"testing"

	"github.com/jcfug8/daylear/server/core/masks"

	"github.com/stretchr/testify/require"
)

func Test_Equal(t *testing.T) {
	type T = testing.T

	for _, tt := range []struct {
		name     string
		lhs, rhs []string
		equal    bool
	}{{
		name:  "should return true if both nil",
		equal: true,
	}, {
		name: "should return false if lhs is nil",
		rhs:  []string{"a, b, c"},
	}, {
		name: "should return false if nhs is nil",
		lhs:  []string{"a, b, c"},
	}, {
		name: "should return false if different amount of fields",
		lhs:  []string{"a", "b"},
		rhs:  []string{"c"},
	}, {
		name: "should return false if different fields",
		lhs:  []string{"a"},
		rhs:  []string{"b"},
	}, {
		name:  "should return true if equal",
		lhs:   []string{"a, b, c"},
		rhs:   []string{"a, b, c"},
		equal: true,
	}} {
		t.Run(tt.name, func(t *T) {
			require.Equal(t, tt.equal, masks.Equal(tt.lhs, tt.rhs))
		})
	}
}

func Test_Intersection(t *testing.T) {
	type T = testing.T

	t.Run("should return intersection", func(t *T) {
		require.Nil(t, masks.Intersection(nil, nil))
		require.Nil(t, masks.Intersection([]string{"a, b, c"}, nil))
		require.Nil(t, masks.Intersection(nil, []string{"a, b, c"}))

		mask := []string{"a", "d"}
		want := []string{"a"}

		intersect := masks.Intersection([]string{"a", "b", "c"}, mask)
		require.True(t, slices.Contains(intersect, "a"))
		require.False(t, slices.Contains(intersect, "d"))
		require.True(t, masks.Equal(intersect, want))
	})
}

func Test_Map(t *testing.T) {
	type T = testing.T

	t.Run("should return mapped fields", func(t *T) {
		require.Nil(t, masks.Map(nil, nil))
		require.Nil(t, masks.Map([]string{"a", "b", "c"}, nil))

		mappings := masks.FieldMap{
			"a": []string{"1", "2"},
			"c": []string{"3"},
			"e": []string{"4", "5", "2"},
			"f": []string{"6", "7", "8"},
		}

		want := []string{"1", "2", "3", "4", "5"}

		mask := masks.Map([]string{"a", "b", "c", "e"}, mappings)
		require.True(t, masks.Equal(want, mask))
	})
}

func Test_Prefix(t *testing.T) {
	type T = testing.T

	t.Run("should prefix paths", func(t *T) {
		require.Nil(t, masks.Prefix("a", nil))

		mask := []string{"a", "b", "c", "e"}
		want := []string{"a.a", "a.b", "a.c", "a.e"}

		require.True(t, masks.Equal(want, masks.Prefix("a.", mask)))
	})
}

func Test_RemovePaths(t *testing.T) {
	type T = testing.T

	t.Run("should remove paths", func(t *T) {
		require.Nil(t, masks.RemovePaths(nil, "a", "e"))

		mask := []string{"a", "b", "c", "e"}
		want := []string{"b", "c"}

		require.True(t, masks.Equal(want, masks.RemovePaths(mask, "a", "e")))
	})
}
