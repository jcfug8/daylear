package namer

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	namerv1 "github.com/jcfug8/daylear/server/genapi/api/namer/v1"
)

// --- SingletonNamedResource tests ---

type SingletonNamedResourceParent struct {
	ParentOneID   int64
	ParentTwoID   int64
	ParentThreeID int64
}

// Pattern 0: singletonNamedResource (no parent)
// Pattern 1: parentOnes/{parent_one}/parentTwos/{parent_two}/singletonNamedResource
// Pattern 2: parentOnes/{parent_one}/parentThrees/{parent_three}/singletonNamedResource

func GetSingletonNamedVars(parent SingletonNamedResourceParent, patternIndex int) ([]string, error) {
	switch patternIndex {
	case 0:
		return []string{}, nil
	case 1:
		return []string{strconv.FormatInt(parent.ParentOneID, 10), strconv.FormatInt(parent.ParentTwoID, 10)}, nil
	case 2:
		return []string{strconv.FormatInt(parent.ParentOneID, 10), strconv.FormatInt(parent.ParentThreeID, 10)}, nil
	}
	return nil, fmt.Errorf("invalid pattern index: %d", patternIndex)
}

func SetSingletonNamedVars(vars []string, patternIndex int) (SingletonNamedResourceParent, error) {
	switch patternIndex {
	case 0:
		if len(vars) != 0 {
			return SingletonNamedResourceParent{}, fmt.Errorf("expected 0 vars, got %d", len(vars))
		}
		return SingletonNamedResourceParent{}, nil
	case 1:
		if len(vars) != 2 {
			return SingletonNamedResourceParent{}, fmt.Errorf("expected 2 vars, got %d", len(vars))
		}
		parentOneID, err := strconv.ParseInt(vars[0], 10, 64)
		if err != nil {
			return SingletonNamedResourceParent{}, err
		}
		parentTwoID, err := strconv.ParseInt(vars[1], 10, 64)
		if err != nil {
			return SingletonNamedResourceParent{}, err
		}
		return SingletonNamedResourceParent{ParentOneID: parentOneID, ParentTwoID: parentTwoID}, nil
	case 2:
		if len(vars) != 2 {
			return SingletonNamedResourceParent{}, fmt.Errorf("expected 2 vars, got %d", len(vars))
		}
		parentOneID, err := strconv.ParseInt(vars[0], 10, 64)
		if err != nil {
			return SingletonNamedResourceParent{}, err
		}
		parentThreeID, err := strconv.ParseInt(vars[1], 10, 64)
		if err != nil {
			return SingletonNamedResourceParent{}, err
		}
		return SingletonNamedResourceParent{ParentOneID: parentOneID, ParentThreeID: parentThreeID}, nil
	}
	return SingletonNamedResourceParent{}, fmt.Errorf("invalid pattern index: %d", patternIndex)
}

func SetSingletonNamedParent(vars []string, patternIndex int) (SingletonNamedResourceParent, error) {
	return SetSingletonNamedVars(vars, patternIndex)
}

func newParentNamer() (ParentNamer[SingletonNamedResourceParent], error) {
	return NewParentNamer(
		&namerv1.SingletonNamedResource{},
		GetSingletonNamedVars,
		SetSingletonNamedVars,
		SetSingletonNamedParent,
	)
}

func TestParentNamer_Format_Pattern0(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parent := SingletonNamedResourceParent{}
	expected := "singletonNamedResource"
	got, err := namer.Format(parent, 0)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParentNamer_Format_Pattern1(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parent := SingletonNamedResourceParent{1, 2, 0}
	expected := "parentOnes/1/parentTwos/2/singletonNamedResource"
	got, err := namer.Format(parent, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParentNamer_Format_Pattern2(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parent := SingletonNamedResourceParent{3, 0, 4}
	expected := "parentOnes/3/parentThrees/4/singletonNamedResource"
	got, err := namer.Format(parent, 2)
	assert.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestParentNamer_Format_InvalidPatternIndex(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parent := SingletonNamedResourceParent{1, 2, 3}
	_, err = namer.Format(parent, 99)
	assert.Error(t, err)
}

func TestParentNamer_Parse_Pattern0(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	name := "singletonNamedResource"
	parent, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, SingletonNamedResourceParent{}, parent)
	assert.Equal(t, 0, pattern)
}

func TestParentNamer_Parse_Pattern1(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	name := "parentOnes/1/parentTwos/2/singletonNamedResource"
	parent, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, SingletonNamedResourceParent{1, 2, 0}, parent)
	assert.Equal(t, 1, pattern)
}

func TestParentNamer_Parse_Pattern2(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	name := "parentOnes/3/parentThrees/4/singletonNamedResource"
	parent, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, SingletonNamedResourceParent{3, 0, 4}, parent)
	assert.Equal(t, 2, pattern)
}

func TestParentNamer_Parse_InvalidName(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	_, _, err = namer.Parse("invalid/format")
	assert.Error(t, err)
}

func TestParentNamer_ParseParent_Pattern0(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parentStr := "singletonNamedResource"
	parent, pattern, err := namer.ParseParent(parentStr)
	assert.NoError(t, err)
	assert.Equal(t, SingletonNamedResourceParent{}, parent)
	assert.Equal(t, 0, pattern)
}

func TestParentNamer_ParseParent_Pattern1(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parentStr := "parentOnes/1/parentTwos/2"
	parent, pattern, err := namer.ParseParent(parentStr)
	assert.NoError(t, err)
	assert.Equal(t, SingletonNamedResourceParent{1, 2, 0}, parent)
	assert.Equal(t, 1, pattern)
}

func TestParentNamer_ParseParent_Pattern2(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parentStr := "parentOnes/3/parentThrees/4"
	parent, pattern, err := namer.ParseParent(parentStr)
	assert.NoError(t, err)
	assert.Equal(t, SingletonNamedResourceParent{3, 0, 4}, parent)
	assert.Equal(t, 2, pattern)
}

func TestParentNamer_ParseParent_InvalidParent(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	_, _, err = namer.ParseParent("invalid/parent")
	assert.Error(t, err)
}

func TestParentNamer_Integration_Pattern0(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parent := SingletonNamedResourceParent{}
	name, err := namer.Format(parent, 0)
	require.NoError(t, err)
	parsedParent, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, parent, parsedParent)
	assert.Equal(t, 0, pattern)
}

func TestParentNamer_Integration_Pattern1(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parent := SingletonNamedResourceParent{1, 2, 0}
	name, err := namer.Format(parent, 1)
	require.NoError(t, err)
	parsedParent, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, parent, parsedParent)
	assert.Equal(t, 1, pattern)
}

func TestParentNamer_Integration_Pattern2(t *testing.T) {
	t.Parallel()
	namer, err := newParentNamer()
	require.NoError(t, err)
	parent := SingletonNamedResourceParent{3, 0, 4}
	name, err := namer.Format(parent, 2)
	require.NoError(t, err)
	parsedParent, pattern, err := namer.Parse(name)
	assert.NoError(t, err)
	assert.Equal(t, parent, parsedParent)
	assert.Equal(t, 2, pattern)
}
