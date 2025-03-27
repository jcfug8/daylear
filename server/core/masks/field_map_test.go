package masks_test

import (
	"testing"

	"github.com/jcfug8/daylear/server/core/masks"

	"github.com/stretchr/testify/require"
)

func TestFieldMap(t *testing.T) {
	type T = testing.T

	t.Run("should create a new FieldMap", func(t *T) {
		m := masks.NewFieldMap().
			MapFieldToFields("field", "1", "2", "3").
			MapFieldToFields("field2", "4", "5")

		require.Equal(t, masks.FieldMap(map[string][]string{
			"field":  {"1", "2", "3"},
			"field2": {"4", "5"},
		}), m)
	})
}
