package namer

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	recipev1alpha1 "github.com/jcfug8/daylear/server/genapi/api/meals/recipe/v1alpha1"
)

// TestParent implements parentType interface for testing
type TestParent struct {
	UserId   int64
	CircleId int64 // Optional circle ID
}

// TestId implements idType interface for testing
type TestId struct {
	RecipeId int64
}

func GetStandardNamerVars(parent TestParent, id TestId, patternIndex int) []string {
	if patternIndex == 0 {
		// Pattern: users/{user}/recipes/{recipe}
		return []string{strconv.FormatInt(parent.UserId, 10), strconv.FormatInt(id.RecipeId, 10)}
	}
	// Pattern: users/{user}/circles/{circle}/recipes/{recipe}
	return []string{strconv.FormatInt(parent.UserId, 10), strconv.FormatInt(parent.CircleId, 10), strconv.FormatInt(id.RecipeId, 10)}
}

func SetStandardNamerVars(vars []string, patternIndex int) (TestParent, TestId, error) {
	parent, err := SetStandardNamerParent(vars[:len(vars)-1], patternIndex)
	if err != nil {
		return TestParent{}, TestId{}, err
	}

	if patternIndex == 0 {
		recipeId, err := strconv.ParseInt(vars[1], 10, 64)
		if err != nil {
			return TestParent{}, TestId{}, err
		}
		return parent, TestId{RecipeId: recipeId}, nil
	}

	recipeId, err := strconv.ParseInt(vars[2], 10, 64)
	if err != nil {
		return TestParent{}, TestId{}, err
	}

	return parent, TestId{RecipeId: recipeId}, nil
}

func SetStandardNamerParent(vars []string, patternIndex int) (TestParent, error) {
	userId, err := strconv.ParseInt(vars[0], 10, 64)
	if err != nil {
		return TestParent{}, err
	}

	if patternIndex == 0 {
		return TestParent{UserId: userId}, nil
	}

	circleId, err := strconv.ParseInt(vars[1], 10, 64)
	if err != nil {
		return TestParent{}, err
	}

	return TestParent{UserId: userId, CircleId: circleId}, nil
}

func TestStandardNamer_Format(t *testing.T) {
	tests := []struct {
		name         string
		parent       TestParent
		id           TestId
		patternIndex int
		want         string
		wantErr      bool
	}{
		{
			name:         "valid format - user recipe",
			parent:       TestParent{UserId: 123},
			id:           TestId{RecipeId: 456},
			patternIndex: 0,
			want:         "users/123/recipes/456",
			wantErr:      false,
		},
		{
			name:         "valid format - circle recipe",
			parent:       TestParent{UserId: 123, CircleId: 789},
			id:           TestId{RecipeId: 456},
			patternIndex: 1,
			want:         "users/123/circles/789/recipes/456",
			wantErr:      false,
		},
		{
			name:         "invalid pattern index",
			parent:       TestParent{UserId: 123},
			id:           TestId{RecipeId: 456},
			patternIndex: 999,
			want:         "",
			wantErr:      true,
		},
		{
			name:         "negative pattern index",
			parent:       TestParent{UserId: 123},
			id:           TestId{RecipeId: 456},
			patternIndex: -1,
			want:         "",
			wantErr:      true,
		},
		{
			name:         "zero values",
			parent:       TestParent{UserId: 0},
			id:           TestId{RecipeId: 0},
			patternIndex: 0,
			want:         "users/0/recipes/0",
			wantErr:      false,
		},
		{
			name:         "negative values",
			parent:       TestParent{UserId: -123, CircleId: -789},
			id:           TestId{RecipeId: -456},
			patternIndex: 1,
			want:         "users/-123/circles/-789/recipes/-456",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("testing", tt.name)
			namer, err := NewStandardNamer(
				&recipev1alpha1.Recipe{},
				GetStandardNamerVars,
				SetStandardNamerVars,
				SetStandardNamerParent,
			)
			require.NoError(t, err)

			got, err := namer.Format(tt.parent, tt.id, tt.patternIndex)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStandardNamer_Parse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  struct {
			parent       TestParent
			id           TestId
			patternIndex int
		}
		wantErr bool
	}{
		{
			name:  "valid parse - user recipe",
			input: "users/123/recipes/456",
			want: struct {
				parent       TestParent
				id           TestId
				patternIndex int
			}{
				parent:       TestParent{UserId: 123},
				id:           TestId{RecipeId: 456},
				patternIndex: 0,
			},
			wantErr: false,
		},
		{
			name:  "valid parse - circle recipe",
			input: "users/123/circles/789/recipes/456",
			want: struct {
				parent       TestParent
				id           TestId
				patternIndex int
			}{
				parent:       TestParent{UserId: 123, CircleId: 789},
				id:           TestId{RecipeId: 456},
				patternIndex: 1,
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid/format",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "missing segments",
			input:   "users/123",
			wantErr: true,
		},
		{
			name:    "extra segments",
			input:   "users/123/recipes/456/extra",
			wantErr: true,
		},
		{
			name:    "invalid user id",
			input:   "users/not-a-number/recipes/456",
			wantErr: true,
		},
		{
			name:    "invalid recipe id",
			input:   "users/123/recipes/not-a-number",
			wantErr: true,
		},
		{
			name:  "zero values",
			input: "users/0/recipes/0",
			want: struct {
				parent       TestParent
				id           TestId
				patternIndex int
			}{
				parent:       TestParent{UserId: 0},
				id:           TestId{RecipeId: 0},
				patternIndex: 0,
			},
			wantErr: false,
		},
		{
			name:  "negative values",
			input: "users/-123/circles/-789/recipes/-456",
			want: struct {
				parent       TestParent
				id           TestId
				patternIndex int
			}{
				parent:       TestParent{UserId: -123, CircleId: -789},
				id:           TestId{RecipeId: -456},
				patternIndex: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("testing", tt.name)
			namer, err := NewStandardNamer(
				&recipev1alpha1.Recipe{},
				GetStandardNamerVars,
				SetStandardNamerVars,
				SetStandardNamerParent,
			)
			require.NoError(t, err)

			gotParent, gotId, gotPatternIndex, err := namer.Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.parent, gotParent)
			assert.Equal(t, tt.want.id, gotId)
			assert.Equal(t, tt.want.patternIndex, gotPatternIndex)
		})
	}
}

func TestStandardNamer_ParseParent(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  struct {
			parent       TestParent
			patternIndex int
		}
		wantErr bool
	}{
		{
			name:  "valid parse - user recipe",
			input: "users/123/recipes",
			want: struct {
				parent       TestParent
				patternIndex int
			}{
				parent:       TestParent{UserId: 123},
				patternIndex: 0,
			},
			wantErr: false,
		},
		{
			name:  "valid parse - circle recipe",
			input: "users/123/circles/789/recipes",
			want: struct {
				parent       TestParent
				patternIndex int
			}{
				parent:       TestParent{UserId: 123, CircleId: 789},
				patternIndex: 1,
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid/format",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "missing segments",
			input:   "users/123",
			wantErr: true,
		},
		{
			name:    "extra segments",
			input:   "users/123/recipes/extra",
			wantErr: true,
		},
		{
			name:    "invalid user id",
			input:   "users/not-a-number/recipes",
			wantErr: true,
		},
		{
			name:    "invalid circle id",
			input:   "users/123/circles/not-a-number/recipes",
			wantErr: true,
		},
		{
			name:  "zero values",
			input: "users/0/recipes",
			want: struct {
				parent       TestParent
				patternIndex int
			}{
				parent:       TestParent{UserId: 0},
				patternIndex: 0,
			},
			wantErr: false,
		},
		{
			name:  "negative values",
			input: "users/-123/circles/-789/recipes",
			want: struct {
				parent       TestParent
				patternIndex int
			}{
				parent:       TestParent{UserId: -123, CircleId: -789},
				patternIndex: 1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("testing", tt.name)
			namer, err := NewStandardNamer(
				&recipev1alpha1.Recipe{},
				GetStandardNamerVars,
				SetStandardNamerVars,
				SetStandardNamerParent,
			)
			require.NoError(t, err)

			gotParent, gotPatternIndex, err := namer.ParseParent(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.parent, gotParent)
			assert.Equal(t, tt.want.patternIndex, gotPatternIndex)
		})
	}
}

func TestStandardNamer_Integration(t *testing.T) {
	namer, err := NewStandardNamer(
		&recipev1alpha1.Recipe{},
		GetStandardNamerVars,
		SetStandardNamerVars,
		SetStandardNamerParent,
	)
	require.NoError(t, err)

	testCases := []struct {
		parent       TestParent
		id           TestId
		patternIndex int
	}{
		{
			parent:       TestParent{UserId: 123},
			id:           TestId{RecipeId: 456},
			patternIndex: 0,
		},
		{
			parent:       TestParent{UserId: 123, CircleId: 789},
			id:           TestId{RecipeId: 456},
			patternIndex: 1,
		},
		{
			parent:       TestParent{UserId: 0},
			id:           TestId{RecipeId: 0},
			patternIndex: 0,
		},
		{
			parent:       TestParent{UserId: -123, CircleId: -789},
			id:           TestId{RecipeId: -456},
			patternIndex: 1,
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			fmt.Println("testing", strconv.Itoa(i))
			// Test Format -> Parse
			formatted, err := namer.Format(tc.parent, tc.id, tc.patternIndex)
			require.NoError(t, err)

			parsedParent, parsedId, patternIndex, err := namer.Parse(formatted)
			require.NoError(t, err)
			assert.Equal(t, tc.parent, parsedParent)
			assert.Equal(t, tc.id, parsedId)
			assert.Equal(t, tc.patternIndex, patternIndex)

			// Test ParseParent
			parentPath := formatted[:strings.LastIndex(formatted, "/")]
			parsedParent, patternIndex, err = namer.ParseParent(parentPath)
			require.NoError(t, err)
			assert.Equal(t, tc.parent, parsedParent)
			assert.Equal(t, tc.patternIndex, patternIndex)
		})
	}
}
