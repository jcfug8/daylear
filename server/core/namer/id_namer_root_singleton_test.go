package namer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	namerv1 "github.com/jcfug8/daylear/server/genapi/api/namer/v1"
)

// RootSingletonNamedResourceID is empty because the singleton pattern has no variables
// (pattern: "rootSingletonNamedResource")
type RootSingletonNamedResourceID struct{}

func GetRootSingletonNamedVars(id RootSingletonNamedResourceID, patternIndex int) (string, error) {
	return "", nil
}

func SetRootSingletonNamedVars(idVar string, patternIndex int) (RootSingletonNamedResourceID, error) {
	// No variables to parse for singleton
	return RootSingletonNamedResourceID{}, nil
}

func NewRootSingletonNamedNamer() (IDNamer[RootSingletonNamedResourceID], error) {
	return NewIDNamer(
		&namerv1.RootSingletonNamedResource{},
		GetRootSingletonNamedVars,
		SetRootSingletonNamedVars,
	)
}

func TestIDNamer_RootSingleton_Format(t *testing.T) {
	t.Parallel()
	namer, err := NewRootSingletonNamedNamer()
	require.NoError(t, err)
	got, err := namer.Format(RootSingletonNamedResourceID{}, 0)
	assert.NoError(t, err)
	assert.Equal(t, "rootSingletonNamedResource", got)
}

func TestIDNamer_RootSingleton_Format_InvalidPatternIndex(t *testing.T) {
	t.Parallel()
	namer, err := NewRootSingletonNamedNamer()
	require.NoError(t, err)
	_, err = namer.Format(RootSingletonNamedResourceID{}, 999)
	assert.Error(t, err)
}

func TestIDNamer_RootSingleton_Parse(t *testing.T) {
	t.Parallel()
	namer, err := NewRootSingletonNamedNamer()
	require.NoError(t, err)
	gotID, gotPatternIndex, err := namer.Parse("rootSingletonNamedResource")
	assert.NoError(t, err)
	assert.Equal(t, RootSingletonNamedResourceID{}, gotID)
	assert.Equal(t, 0, gotPatternIndex)
}

func TestIDNamer_RootSingleton_Parse_ExtraVariables(t *testing.T) {
	t.Parallel()
	namer, err := NewRootSingletonNamedNamer()
	require.NoError(t, err)
	_, _, err = namer.Parse("rootSingletonNamedResource/extra")
	assert.Error(t, err)
	_, _, err = namer.Parse("")
	assert.Error(t, err)
}

func TestIDNamer_RootSingleton_Parse_InvalidName(t *testing.T) {
	t.Parallel()
	namer, err := NewRootSingletonNamedNamer()
	require.NoError(t, err)
	_, _, err = namer.Parse("cats")
	assert.Error(t, err)
	_, _, err = namer.Parse("")
	assert.Error(t, err)
}
