package namer

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	userv1alpha1 "github.com/jcfug8/daylear/server/genapi/api/users/user/v1alpha1"
)

// TestUserId implements idType interface for testing
type TestUserId struct {
	UserId int64
}

func GetRootNamerVars(id TestUserId, patternIndex int) string {
	return strconv.FormatInt(id.UserId, 10)
}

func SetRootNamerVars(idVar string, patternIndex int) (TestUserId, error) {
	userId, err := strconv.ParseInt(idVar, 10, 64)
	if err != nil {
		return TestUserId{}, err
	}
	return TestUserId{UserId: userId}, nil
}

func TestRootNamer_Format(t *testing.T) {
	tests := []struct {
		name         string
		id           TestUserId
		patternIndex int
		want         string
		wantErr      bool
	}{
		{
			name:         "valid format",
			id:           TestUserId{UserId: 123},
			patternIndex: 0,
			want:         "users/123",
			wantErr:      false,
		},
		{
			name:         "invalid pattern index",
			id:           TestUserId{UserId: 123},
			patternIndex: 999,
			want:         "",
			wantErr:      true,
		},
		{
			name:         "negative pattern index",
			id:           TestUserId{UserId: 123},
			patternIndex: -1,
			want:         "",
			wantErr:      true,
		},
		{
			name:         "zero user id",
			id:           TestUserId{UserId: 0},
			patternIndex: 0,
			want:         "users/0",
			wantErr:      false,
		},
		{
			name:         "negative user id",
			id:           TestUserId{UserId: -123},
			patternIndex: 0,
			want:         "users/-123",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("testing", tt.name)
			namer, err := NewRootNamer(
				&userv1alpha1.User{},
				GetRootNamerVars,
				SetRootNamerVars,
			)
			require.NoError(t, err)

			got, err := namer.Format(tt.id, tt.patternIndex)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRootNamer_Parse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  struct {
			id           TestUserId
			patternIndex int
		}
		wantErr bool
	}{
		{
			name:  "valid parse",
			input: "users/123",
			want: struct {
				id           TestUserId
				patternIndex int
			}{
				id:           TestUserId{UserId: 123},
				patternIndex: 0,
			},
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "invalid/format",
			wantErr: true,
		},
		{
			name:    "invalid user id",
			input:   "users/not-a-number",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "missing user id",
			input:   "users/",
			wantErr: true,
		},
		{
			name:    "extra segments",
			input:   "users/123/extra",
			wantErr: true,
		},
		{
			name:  "zero user id",
			input: "users/0",
			want: struct {
				id           TestUserId
				patternIndex int
			}{
				id:           TestUserId{UserId: 0},
				patternIndex: 0,
			},
			wantErr: false,
		},
		{
			name:  "negative user id",
			input: "users/-123",
			want: struct {
				id           TestUserId
				patternIndex int
			}{
				id:           TestUserId{UserId: -123},
				patternIndex: 0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("testing", tt.name)
			namer, err := NewRootNamer(
				&userv1alpha1.User{},
				GetRootNamerVars,
				SetRootNamerVars,
			)
			require.NoError(t, err)

			gotId, gotPatternIndex, err := namer.Parse(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want.id, gotId)
			assert.Equal(t, tt.want.patternIndex, gotPatternIndex)
		})
	}
}

func TestRootNamer_Integration(t *testing.T) {
	// Test that we can format and then parse back to the same values
	namer, err := NewRootNamer(
		&userv1alpha1.User{},
		GetRootNamerVars,
		SetRootNamerVars,
	)
	require.NoError(t, err)

	testCases := []TestUserId{
		{UserId: 123},
		{UserId: 0},
		{UserId: -123},
		{UserId: 999999999},
	}

	for _, originalId := range testCases {
		t.Run(strconv.FormatInt(originalId.UserId, 10), func(t *testing.T) {
			formatted, err := namer.Format(originalId, 0)
			require.NoError(t, err)

			parsedId, patternIndex, err := namer.Parse(formatted)
			require.NoError(t, err)
			assert.Equal(t, originalId, parsedId)
			assert.Equal(t, 0, patternIndex)
		})
	}
}
