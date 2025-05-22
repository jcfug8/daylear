package namer

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	namerv1 "github.com/jcfug8/daylear/server/genapi/api/namer/v1"
)

// StandardNamedResourceParent
type StandardNamedResourceParent struct {
	ParentOneID   int64
	ParentTwoID   int64
	ParentThreeID int64
}

// StandardNamedResourceId
type StandardNamedResourceId struct {
	ID int64
}

func GetStandardNamedVars(parent StandardNamedResourceParent, id StandardNamedResourceId, patternIndex int) ([]string, error) {
	switch patternIndex {
	case 0:
		return []string{strconv.FormatInt(id.ID, 10)}, nil
	case 1:
		return []string{strconv.FormatInt(parent.ParentOneID, 10), strconv.FormatInt(parent.ParentTwoID, 10), strconv.FormatInt(id.ID, 10)}, nil
	case 2:
		return []string{strconv.FormatInt(parent.ParentOneID, 10), strconv.FormatInt(parent.ParentThreeID, 10), strconv.FormatInt(id.ID, 10)}, nil
	case 3:
		return []string{strconv.FormatInt(parent.ParentOneID, 10), strconv.FormatInt(id.ID, 10)}, nil
	}
	return []string{}, fmt.Errorf("invalid pattern index: %d", patternIndex)
}

func SetStandardNamedVars(vars []string, patternIndex int) (StandardNamedResourceParent, StandardNamedResourceId, error) {
	switch patternIndex {
	case 0, 1, 2, 3:
		parent, err := SetStandardNamedParent(vars[:len(vars)-1], patternIndex)
		if err != nil {
			return StandardNamedResourceParent{}, StandardNamedResourceId{}, err
		}
		id, err := strconv.ParseInt(vars[len(vars)-1], 10, 64)
		if err != nil {
			return StandardNamedResourceParent{}, StandardNamedResourceId{}, err
		}
		return parent, StandardNamedResourceId{ID: id}, nil
	}
	return StandardNamedResourceParent{}, StandardNamedResourceId{}, fmt.Errorf("invalid pattern index: %d", patternIndex)
}

func SetStandardNamedParent(vars []string, patternIndex int) (StandardNamedResourceParent, error) {
	switch patternIndex {
	case 0:
		if len(vars) != 0 {
			return StandardNamedResourceParent{}, fmt.Errorf("expected 0 vars, got %d", len(vars))
		}
		return StandardNamedResourceParent{}, nil
	case 1:

		if len(vars) != 2 {
			return StandardNamedResourceParent{}, fmt.Errorf("expected 2 vars, got %d", len(vars))
		}
		parentOneID, err := strconv.ParseInt(vars[0], 10, 64)
		if err != nil {
			return StandardNamedResourceParent{}, err
		}
		parentTwoID, err := strconv.ParseInt(vars[1], 10, 64)
		if err != nil {
			return StandardNamedResourceParent{}, err
		}
		return StandardNamedResourceParent{ParentOneID: parentOneID, ParentTwoID: parentTwoID}, nil
	case 2:
		if len(vars) != 2 {
			return StandardNamedResourceParent{}, fmt.Errorf("expected 2 vars, got %d", len(vars))
		}
		parentOneID, err := strconv.ParseInt(vars[0], 10, 64)
		if err != nil {
			return StandardNamedResourceParent{}, err
		}
		parentThreeID, err := strconv.ParseInt(vars[1], 10, 64)
		if err != nil {
			return StandardNamedResourceParent{}, err
		}
		return StandardNamedResourceParent{ParentOneID: parentOneID, ParentThreeID: parentThreeID}, nil
	case 3:
		if len(vars) != 1 {
			return StandardNamedResourceParent{}, fmt.Errorf("expected 1 var, got %d", len(vars))
		}
		parentOneID, err := strconv.ParseInt(vars[0], 10, 64)
		if err != nil {
			return StandardNamedResourceParent{}, err
		}
		return StandardNamedResourceParent{ParentOneID: parentOneID}, nil
	}
	return StandardNamedResourceParent{}, fmt.Errorf("invalid pattern index: %d", patternIndex)
}

func newParentIDNamer() (ParentIDNamer[StandardNamedResourceParent, StandardNamedResourceId], error) {
	return NewParentIDNamer(
		&namerv1.StandardNamedResource{},
		GetStandardNamedVars,
		SetStandardNamedVars,
		SetStandardNamedParent,
	)
}

func TestParentIDNamer_Format_Pattern0(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parent := StandardNamedResourceParent{}
	id := StandardNamedResourceId{55}
	expected := "standardNamedResources/55"
	got, err := namer.Format(parent, id, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParentIDNamer_Format_Pattern1(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parent := StandardNamedResourceParent{1, 2, 0}
	id := StandardNamedResourceId{3}
	expected := "parentOnes/1/parentTwos/2/standardNamedResources/3"
	got, err := namer.Format(parent, id, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParentIDNamer_Format_Pattern2(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parent := StandardNamedResourceParent{4, 0, 5}
	id := StandardNamedResourceId{6}
	expected := "parentOnes/4/parentThrees/5/standardNamedResources/6"
	got, err := namer.Format(parent, id, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParentIDNamer_Format_Pattern3(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parent := StandardNamedResourceParent{7, 0, 0}
	id := StandardNamedResourceId{8}
	expected := "parentOnes/7/standardNamedResources/8"
	got, err := namer.Format(parent, id, 3)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParentIDNamer_Format_InvalidPatternIndex(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	_, err = namer.Format(StandardNamedResourceParent{1, 2, 3}, StandardNamedResourceId{4}, 99)
	assert.Error(t, err)
}

func TestParentIDNamer_Parse_Pattern0(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	name := "standardNamedResources/55"
	parent, id, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, StandardNamedResourceParent{}, parent)
	assert.Equal(t, StandardNamedResourceId{55}, id)
	assert.Equal(t, 0, pattern)
}

func TestParentIDNamer_Parse_Pattern1(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	name := "parentOnes/1/parentTwos/2/standardNamedResources/3"
	parent, id, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, StandardNamedResourceParent{1, 2, 0}, parent)
	assert.Equal(t, StandardNamedResourceId{3}, id)
	assert.Equal(t, 1, pattern)
}

func TestParentIDNamer_Parse_Pattern2(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	name := "parentOnes/4/parentThrees/5/standardNamedResources/6"
	parent, id, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, StandardNamedResourceParent{4, 0, 5}, parent)
	assert.Equal(t, StandardNamedResourceId{6}, id)
	assert.Equal(t, 2, pattern)
}

func TestParentIDNamer_Parse_Pattern3(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	name := "parentOnes/7/standardNamedResources/8"
	parent, id, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, StandardNamedResourceParent{7, 0, 0}, parent)
	assert.Equal(t, StandardNamedResourceId{8}, id)
	assert.Equal(t, 3, pattern)
}

func TestParentIDNamer_Parse_InvalidName(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	_, _, _, err = namer.Parse("invalid/format")
	assert.Error(t, err)
}

func TestParentIDNamer_ParseParent_Pattern0(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parentStr := "standardNamedResources/3"
	parent, pattern, err := namer.ParseParent(parentStr)
	assert.NoError(t, err)
	assert.Equal(t, StandardNamedResourceParent{}, parent)
	assert.Equal(t, 0, pattern)
}

func TestParentIDNamer_ParseParent_Pattern1(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parentStr := "parentOnes/1/parentTwos/2"
	parent, pattern, err := namer.ParseParent(parentStr)
	assert.NoError(t, err)
	assert.Equal(t, StandardNamedResourceParent{1, 2, 0}, parent)
	assert.Equal(t, 1, pattern)
}

func TestParentIDNamer_ParseParent_Pattern2(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parentStr := "parentOnes/4/parentThrees/5"
	parent, pattern, err := namer.ParseParent(parentStr)
	assert.NoError(t, err)
	assert.Equal(t, StandardNamedResourceParent{4, 0, 5}, parent)
	assert.Equal(t, 2, pattern)
}

func TestParentIDNamer_ParseParent_Pattern3(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parentStr := "parentOnes/7"
	parent, pattern, err := namer.ParseParent(parentStr)
	assert.NoError(t, err)
	assert.Equal(t, StandardNamedResourceParent{7, 0, 0}, parent)
	assert.Equal(t, 3, pattern)
}

func TestParentIDNamer_ParseParent_InvalidParent(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	_, _, err = namer.ParseParent("invalid/parent")
	assert.Error(t, err)
}

func TestParentIDNamer_Integration0(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parent := StandardNamedResourceParent{0, 0, 0}
	id := StandardNamedResourceId{30}
	name, err := namer.Format(parent, id, 0)
	require.NoError(t, err)
	parsedParent, parsedId, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, parent, parsedParent)
	assert.Equal(t, id, parsedId)
	assert.Equal(t, 0, pattern)
}

func TestParentIDNamer_Integration1(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parent := StandardNamedResourceParent{1, 2, 0}
	id := StandardNamedResourceId{40}
	name, err := namer.Format(parent, id, 1)
	require.NoError(t, err)
	parsedParent, parsedId, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, parent, parsedParent)
	assert.Equal(t, id, parsedId)
	assert.Equal(t, 1, pattern)
}

func TestParentIDNamer_Integration2(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parent := StandardNamedResourceParent{1, 0, 3}
	id := StandardNamedResourceId{20}
	name, err := namer.Format(parent, id, 2)
	require.NoError(t, err)
	parsedParent, parsedId, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, parent, parsedParent)
	assert.Equal(t, id, parsedId)
	assert.Equal(t, 2, pattern)
}

func TestParentIDNamer_Integration3(t *testing.T) {
	t.Parallel()
	namer, err := newParentIDNamer()
	require.NoError(t, err)
	parent := StandardNamedResourceParent{1, 0, 0}
	id := StandardNamedResourceId{55}
	name, err := namer.Format(parent, id, 3)
	require.NoError(t, err)
	parsedParent, parsedId, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, parent, parsedParent)
	assert.Equal(t, id, parsedId)
	assert.Equal(t, 3, pattern)
}
