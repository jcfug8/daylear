package namer

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	namerv1 "github.com/jcfug8/daylear/server/genapi/api/namer/v1"
)

// RootNamedResourceID
type RootNamedResourceID struct {
	RootNamedResourceID int64
}

func GetIDNamerVars(id RootNamedResourceID, patternIndex int) (string, error) {
	return strconv.FormatInt(id.RootNamedResourceID, 10), nil
}

func SetIDNamerVars(idVar string, patternIndex int) (RootNamedResourceID, error) {
	rootNamedResourceID, err := strconv.ParseInt(idVar, 10, 64)
	if err != nil {
		return RootNamedResourceID{}, err
	}
	return RootNamedResourceID{RootNamedResourceID: rootNamedResourceID}, nil
}

func NewRootNamedResourceNamer() (IDNamer[RootNamedResourceID], error) {
	return NewIDNamer(
		&namerv1.RootNamedResource{},
		GetIDNamerVars,
		SetIDNamerVars,
	)
}

func TestIDNamer_Format_ValidFormat(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)

	got, err := namer.Format(RootNamedResourceID{RootNamedResourceID: 123}, 0)
	assert.NoError(t, err)
	assert.Equal(t, "rootNamedResources/123", got)
}

func TestIDNamer_Format_InvalidPatternIndex(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)

	got, err := namer.Format(RootNamedResourceID{RootNamedResourceID: 123}, 999)
	assert.Error(t, err)
	assert.Equal(t, "", got)
}

func TestIDNamer_Format_NegativePatternIndex(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)

	got, err := namer.Format(RootNamedResourceID{RootNamedResourceID: 123}, -1)
	assert.Error(t, err)
	assert.Equal(t, "", got)
}

func TestIDNamer_Format_ZeroUserID(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)

	got, err := namer.Format(RootNamedResourceID{RootNamedResourceID: 0}, 0)
	assert.NoError(t, err)
	assert.Equal(t, "rootNamedResources/0", got)
}

func TestIDNamer_Format_NegativeUserID(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)

	got, err := namer.Format(RootNamedResourceID{RootNamedResourceID: -123}, 0)
	assert.NoError(t, err)
	assert.Equal(t, "rootNamedResources/-123", got)
}

func TestIDNamer_Parse_ValidParse(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	gotID, gotPatternIndex, err := namer.Parse("rootNamedResources/123")
	assert.NoError(t, err)
	assert.Equal(t, RootNamedResourceID{RootNamedResourceID: 123}, gotID)
	assert.Equal(t, 0, gotPatternIndex)
}

func TestIDNamer_Parse_InvalidFormat(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	_, _, err = namer.Parse("invalid/format")
	assert.Error(t, err)
}

func TestIDNamer_Parse_InvalidUserID(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	_, _, err = namer.Parse("rootNamedResources/not-a-number")
	assert.Error(t, err)
}

func TestIDNamer_Parse_EmptyString(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	_, _, err = namer.Parse("")
	assert.Error(t, err)
}

func TestIDNamer_Parse_MissingUserID(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	_, _, err = namer.Parse("rootNamedResources/")
	assert.Error(t, err)
}

func TestIDNamer_Parse_ExtraSegments(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	_, _, err = namer.Parse("rootNamedResources/123/extra")
	assert.Error(t, err)
}

func TestIDNamer_Parse_ZeroUserID(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	gotID, gotPatternIndex, err := namer.Parse("rootNamedResources/0")
	assert.NoError(t, err)
	assert.Equal(t, RootNamedResourceID{RootNamedResourceID: 0}, gotID)
	assert.Equal(t, 0, gotPatternIndex)
}

func TestIDNamer_Parse_NegativeUserID(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	gotID, gotPatternIndex, err := namer.Parse("rootNamedResources/-123")
	assert.NoError(t, err)
	assert.Equal(t, RootNamedResourceID{RootNamedResourceID: -123}, gotID)
	assert.Equal(t, 0, gotPatternIndex)
}

func TestIDNamer_Integration_123(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	formatted, err := namer.Format(RootNamedResourceID{RootNamedResourceID: 123}, 0)
	require.NoError(t, err)
	parsedID, patternIndex, err := namer.Parse(formatted)
	require.NoError(t, err)
	assert.Equal(t, RootNamedResourceID{RootNamedResourceID: 123}, parsedID)
	assert.Equal(t, 0, patternIndex)
}

func TestIDNamer_Integration_0(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	formatted, err := namer.Format(RootNamedResourceID{RootNamedResourceID: 0}, 0)
	require.NoError(t, err)
	parsedID, patternIndex, err := namer.Parse(formatted)
	require.NoError(t, err)
	assert.Equal(t, RootNamedResourceID{RootNamedResourceID: 0}, parsedID)
	assert.Equal(t, 0, patternIndex)
}

func TestIDNamer_Integration_Negative123(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	formatted, err := namer.Format(RootNamedResourceID{RootNamedResourceID: -123}, 0)
	require.NoError(t, err)
	parsedID, patternIndex, err := namer.Parse(formatted)
	require.NoError(t, err)
	assert.Equal(t, RootNamedResourceID{RootNamedResourceID: -123}, parsedID)
	assert.Equal(t, 0, patternIndex)
}

func TestIDNamer_Integration_999999999(t *testing.T) {
	t.Parallel()
	namer, err := NewRootNamedResourceNamer()
	require.NoError(t, err)
	formatted, err := namer.Format(RootNamedResourceID{RootNamedResourceID: 999999999}, 0)
	require.NoError(t, err)
	parsedID, patternIndex, err := namer.Parse(formatted)
	require.NoError(t, err)
	assert.Equal(t, RootNamedResourceID{RootNamedResourceID: 999999999}, parsedID)
	assert.Equal(t, 0, patternIndex)
}
