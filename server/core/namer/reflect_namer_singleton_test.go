package namer

import (
	"testing"

	namerv1 "github.com/jcfug8/daylear/server/genapi/api/namer/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSingletonResource is a test resource with a singleton pattern
type TestSingletonResource struct {
	ParentOneID   string `aip_pattern:"key=parent_one"`
	ParentTwoID   string `aip_pattern:"key=parent_two"`
	ParentThreeID string `aip_pattern:"key=parent_three"`
	ID            string `aip_pattern:"key=singleton_named_resource"`
}

func TestReflectNamer_Singleton_Format(t *testing.T) {
	tests := []struct {
		name         string
		resource     *TestSingletonResource
		patternIndex int
		want         string
	}{
		{
			name: "root_pattern",
			resource: &TestSingletonResource{
				ID: "singletonNamedResource",
			},
			patternIndex: 0, // singletonNamedResource
			want:         "singletonNamedResource",
		},
		{
			name: "parent_one_two_pattern",
			resource: &TestSingletonResource{
				ParentOneID: "parent-123",
				ParentTwoID: "parent-456",
				ID:          "singletonNamedResource",
			},
			patternIndex: 1, // parentOnes/{parent_one}/parentTwos/{parent_two}/singletonNamedResource
			want:         "parentOnes/parent-123/parentTwos/parent-456/singletonNamedResource",
		},
		{
			name: "parent_one_three_pattern",
			resource: &TestSingletonResource{
				ParentOneID:   "parent-123",
				ParentThreeID: "parent-789",
				ID:            "singletonNamedResource",
			},
			patternIndex: 2, // parentOnes/{parent_one}/parentThrees/{parent_three}/singletonNamedResource
			want:         "parentOnes/parent-123/parentThrees/parent-789/singletonNamedResource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.SingletonNamedResource]()
			require.NoError(t, err)

			got, err := namer.Format(tt.resource, AsPatternIndex(tt.patternIndex))
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReflectNamer_Singleton_FormatParent(t *testing.T) {
	tests := []struct {
		name         string
		resource     *TestSingletonResource
		patternIndex int
		want         string
		expectErr    bool
	}{
		{
			name: "root_pattern",
			resource: &TestSingletonResource{
				ID: "singletonNamedResource",
			},
			patternIndex: 0, // singletonNamedResource
			want:         "",
			expectErr:    true,
		},
		{
			name: "parent_one_two_pattern",
			resource: &TestSingletonResource{
				ParentOneID: "parent-123",
				ParentTwoID: "parent-456",
				ID:          "singletonNamedResource",
			},
			patternIndex: 1, // parentOnes/{parent_one}/parentTwos/{parent_two}/singletonNamedResource
			want:         "parentOnes/parent-123/parentTwos/parent-456",
		},
		{
			name: "parent_one_three_pattern",
			resource: &TestSingletonResource{
				ParentOneID:   "parent-123",
				ParentThreeID: "parent-789",
				ID:            "singletonNamedResource",
			},
			patternIndex: 2, // parentOnes/{parent_one}/parentThrees/{parent_three}/singletonNamedResource
			want:         "parentOnes/parent-123/parentThrees/parent-789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.SingletonNamedResource]()
			require.NoError(t, err)

			got, err := namer.FormatParent(tt.resource, AsPatternIndex(tt.patternIndex))
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestReflectNamer_Singleton_Parse(t *testing.T) {
	tests := []struct {
		name         string
		resource     *TestSingletonResource
		patternIndex int
		want         string
	}{
		{
			name: "root_pattern",
			resource: &TestSingletonResource{
				ID: "singletonNamedResource",
			},
			patternIndex: 0, // singletonNamedResource
			want:         "singletonNamedResource",
		},
		{
			name: "parent_one_two_pattern",
			resource: &TestSingletonResource{
				ParentOneID: "parent-123",
				ParentTwoID: "parent-456",
				ID:          "singletonNamedResource",
			},
			patternIndex: 1, // parentOnes/{parent_one}/parentTwos/{parent_two}/singletonNamedResource
			want:         "parentOnes/parent-123/parentTwos/parent-456/singletonNamedResource",
		},
		{
			name: "parent_one_three_pattern",
			resource: &TestSingletonResource{
				ParentOneID:   "parent-123",
				ParentThreeID: "parent-789",
				ID:            "singletonNamedResource",
			},
			patternIndex: 2, // parentOnes/{parent_one}/parentThrees/{parent_three}/singletonNamedResource
			want:         "parentOnes/parent-123/parentThrees/parent-789/singletonNamedResource",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.SingletonNamedResource]()
			require.NoError(t, err)

			got, err := namer.Parse(tt.want, tt.resource)
			require.NoError(t, err)
			assert.Equal(t, tt.resource, got)
		})
	}
}

func TestReflectNamer_Singleton_ParseParent(t *testing.T) {
	tests := []struct {
		name         string
		resource     *TestSingletonResource
		patternIndex int
		want         string
	}{
		{
			name: "root_pattern",
			resource: &TestSingletonResource{
				ID: "singletonNamedResource",
			},
			patternIndex: 0, // singletonNamedResource
			want:         "",
		},
		{
			name: "parent_one_two_pattern",
			resource: &TestSingletonResource{
				ParentOneID: "parent-123",
				ParentTwoID: "parent-456",
				ID:          "singletonNamedResource",
			},
			patternIndex: 1, // parentOnes/{parent_one}/parentTwos/{parent_two}/singletonNamedResource
			want:         "parentOnes/parent-123/parentTwos/parent-456",
		},
		{
			name: "parent_one_three_pattern",
			resource: &TestSingletonResource{
				ParentOneID:   "parent-123",
				ParentThreeID: "parent-789",
				ID:            "singletonNamedResource",
			},
			patternIndex: 2, // parentOnes/{parent_one}/parentThrees/{parent_three}/singletonNamedResource
			want:         "parentOnes/parent-123/parentThrees/parent-789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.SingletonNamedResource]()
			require.NoError(t, err)

			got, err := namer.ParseParent(tt.want, tt.resource)
			require.NoError(t, err)
			assert.Equal(t, tt.resource, got)
		})
	}
}

func TestReflectNamer_Singleton_Errors(t *testing.T) {
	tests := []struct {
		name         string
		resource     interface{}
		patternIndex int
		wantErr      bool
	}{
		{
			name:         "nil_input",
			resource:     nil,
			patternIndex: 0,
			wantErr:      true,
		},
		{
			name:         "non-struct_input",
			resource:     "not a struct",
			patternIndex: 0,
			wantErr:      true,
		},
		{
			name: "missing_required_field",
			resource: &TestSingletonResource{
				ParentOneID: "parent-123",
				// Missing ParentTwoID for parent_one_two_pattern
				ID: "singletonNamedResource",
			},
			patternIndex: 1, // Try to use parent_one_two_pattern
			wantErr:      true,
		},
		{
			name: "invalid_pattern",
			resource: &TestSingletonResource{
				ParentOneID:   "parent-123",
				ParentTwoID:   "parent-456",
				ParentThreeID: "parent-789", // Invalid: can't have both ParentTwoID and ParentThreeID
				ID:            "singletonNamedResource",
			},
			patternIndex: 1, // Try to use parent_one_two_pattern
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.SingletonNamedResource]()
			require.NoError(t, err)

			_, err = namer.Format(tt.resource, AsPatternIndex(tt.patternIndex))
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
