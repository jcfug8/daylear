package namer

import (
	"testing"

	namerv1 "github.com/jcfug8/daylear/server/genapi/api/namer/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStandardResource is a test resource with a parent pattern
type TestStandardResource struct {
	ParentOneID string `aip_pattern:"key=parent_one"`
	ID          string `aip_pattern:"key=standard_named_resource"`
}

func TestReflectNamer_Standard_Format(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "standard resource",
			input: &TestStandardResource{
				ParentOneID: "parent-123",
				ID:          "resource-456",
			},
			expected: "parentOnes/parent-123/standardNamedResources/resource-456",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.StandardNamedResource]()
			require.NoError(t, err)

			got, err := namer.Format(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestReflectNamer_Standard_FormatParent(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "standard resource parent",
			input: &TestStandardResource{
				ParentOneID: "parent-123",
				ID:          "resource-456",
			},
			expected: "parentOnes/parent-123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.StandardNamedResource]()
			require.NoError(t, err)

			got, err := namer.FormatParent(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestReflectNamer_Standard_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "standard resource",
			input: "parentOnes/parent-123/standardNamedResources/resource-456",
			expected: &TestStandardResource{
				ParentOneID: "parent-123",
				ID:          "resource-456",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.StandardNamedResource]()
			require.NoError(t, err)

			// Create a new instance of the expected type
			result := tt.expected

			_, err = namer.Parse(tt.input, result)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReflectNamer_Standard_ParseParent(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "standard resource parent",
			input: "parents/parent-123",
			expected: &TestStandardResource{
				ParentOneID: "parent-123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.StandardNamedResource]()
			require.NoError(t, err)

			// Create a new instance of the expected type
			result := tt.expected

			_, err = namer.ParseParent(tt.input, result)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestReflectNamer_Standard_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "non-struct input",
			input:   "not a struct",
			wantErr: true,
		},
		{
			name:    "missing required field",
			input:   &TestStandardResource{},
			wantErr: true,
		},
		{
			name:    "invalid pattern",
			input:   &TestStandardResource{ParentOneID: "123", ID: "456"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.StandardNamedResource]()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			_, err = namer.Format(tt.input)
			assert.Error(t, err)
		})
	}
}
