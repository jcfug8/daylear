package namer

import (
	"testing"

	namerv1 "github.com/jcfug8/daylear/server/genapi/api/namer/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type StandardNamedResource struct {
	ParentOne   int64 `aip_pattern:"key=parent_one"`
	ParentTwo   int64 `aip_pattern:"key=parent_two"`
	ParentThree int64 `aip_pattern:"key=parent_three"`
	ID          int64 `aip_pattern:"key=standard_named_resource"`
	OtherField  string
}

func TestReflectNamer_Format0(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, err := namer.Format(0, StandardNamedResource{0, 0, 0, 4, "test"})
	require.NoError(t, err)
	assert.Equal(t, "standardNamedResources/4", got)
}

func TestReflectNamer_Format1(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, err := namer.Format(1, StandardNamedResource{1, 2, 0, 4, "test"})
	require.NoError(t, err)
	assert.Equal(t, "parentOnes/1/parentTwos/2/standardNamedResources/4", got)
}

func TestReflectNamer_Format2(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, err := namer.Format(2, StandardNamedResource{1, 0, 3, 4, "test"})
	require.NoError(t, err)
	assert.Equal(t, "parentOnes/1/parentThrees/3/standardNamedResources/4", got)
}

func TestReflectNamer_Format3(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, err := namer.Format(3, StandardNamedResource{1, 0, 0, 4, "test"})
	require.NoError(t, err)
	assert.Equal(t, "parentOnes/1/standardNamedResources/4", got)
}

func TestReflectNamer_Format4(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, err := namer.Format(4, StandardNamedResource{1, 2, 3, 4, "test"})
	require.NoError(t, err)
	assert.Equal(t, "parentOnes/1/parentTwos/2/parentThrees/3/standardNamedResources/4", got)
}

func TestReflectNamer_Format_NoPattern(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, err := namer.Format(-1, StandardNamedResource{1, 2, 3, 4, "test"})
	require.NoError(t, err)
	assert.Equal(t, "standardNamedResources/4", got)
}

func BenchmarkReflectNamer_Format4(b *testing.B) {
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(b, err)

	for b.Loop() {
		got, err := namer.Format(4, StandardNamedResource{1, 2, 3, 4, "test"})
		require.NoError(b, err)
		assert.Equal(b, "parentOnes/1/parentTwos/2/parentThrees/3/standardNamedResources/4", got)
	}
}

func TestReflectNamer_Parse0(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.Parse("standardNamedResources/4", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 0, patternIndex)
	assert.Equal(t, StandardNamedResource{0, 0, 0, 4, "test"}, got)
}

func TestReflectNamer_Parse1(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.Parse("parentOnes/1/parentTwos/2/standardNamedResources/4", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 1, patternIndex)
	assert.Equal(t, StandardNamedResource{1, 2, 0, 4, "test"}, got)
}

func TestReflectNamer_Parse2(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.Parse("parentOnes/1/parentThrees/3/standardNamedResources/4", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 2, patternIndex)
	assert.Equal(t, StandardNamedResource{1, 0, 3, 4, "test"}, got)
}

func TestReflectNamer_Parse3(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.Parse("parentOnes/1/standardNamedResources/4", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 3, patternIndex)
	assert.Equal(t, StandardNamedResource{1, 0, 0, 4, "test"}, got)
}

func TestReflectNamer_Parse4(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.Parse("parentOnes/1/parentTwos/2/parentThrees/3/standardNamedResources/4", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 4, patternIndex)
	assert.Equal(t, StandardNamedResource{1, 2, 3, 4, "test"}, got)
}

func BenchmarkReflectNamer_Parse4(b *testing.B) {
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(b, err)

	for b.Loop() {
		got, patternIndex, err := namer.Parse("parentOnes/1/parentTwos/2/parentThrees/3/standardNamedResources/4", StandardNamedResource{OtherField: "test"})
		require.NoError(b, err)
		assert.Equal(b, 4, patternIndex)
		assert.Equal(b, StandardNamedResource{1, 2, 3, 4, "test"}, got)
	}
}

func TestReflectNamer_ParseParent0(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.ParseParent("", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 0, patternIndex)
	assert.Equal(t, StandardNamedResource{0, 0, 0, 0, "test"}, got)
}

func TestRelectNamer_ParseParent1(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.ParseParent("parentOnes/1/parentTwos/2", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 1, patternIndex)
	assert.Equal(t, StandardNamedResource{1, 2, 0, 0, "test"}, got)
}

func TestRelectNamer_ParseParent2(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.ParseParent("parentOnes/1/parentThrees/3", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 2, patternIndex)
	assert.Equal(t, StandardNamedResource{1, 0, 3, 0, "test"}, got)
}

func TestRelectNamer_ParseParent3(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.ParseParent("parentOnes/1", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 3, patternIndex)
	assert.Equal(t, StandardNamedResource{1, 0, 0, 0, "test"}, got)
}

func TestRelectNamer_ParseParent4(t *testing.T) {
	t.Parallel()
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(t, err)

	got, patternIndex, err := namer.ParseParent("parentOnes/1/parentTwos/2/parentThrees/3", StandardNamedResource{OtherField: "test"})
	require.NoError(t, err)
	assert.Equal(t, 4, patternIndex)
	assert.Equal(t, StandardNamedResource{1, 2, 3, 0, "test"}, got)
}

func BenchmarkReflectNamer_ParseParent4(b *testing.B) {
	namer, err := NewReflectNamer[StandardNamedResource](
		&namerv1.StandardNamedResource{},
	)
	require.NoError(b, err)

	for b.Loop() {
		got, patternIndex, err := namer.ParseParent("parentOnes/1/parentTwos/2/parentThrees/3", StandardNamedResource{OtherField: "test"})
		require.NoError(b, err)
		assert.Equal(b, 4, patternIndex)
		assert.Equal(b, StandardNamedResource{1, 2, 3, 0, "test"}, got)
	}
}
