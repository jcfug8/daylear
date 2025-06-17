package namer

import (
	"testing"

	namerv1 "github.com/jcfug8/daylear/server/genapi/api/namer/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRootResource is a test resource with a root pattern
type TestRootResource struct {
	ID string `aip_pattern:"key=root_named_resource"`
}

// TestRootResource2 is a second test resource with a root pattern
type TestRootResource2 struct {
	RootNamedResource string
}

func TestReflectNamer_Root_Format(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "root resource",
			input: &TestRootResource{
				ID: "resource-123",
			},
			expected: "rootNamedResources/resource-123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.RootNamedResource]()
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

func TestReflectNamer_Root_Format2(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "root resource",
			input: &TestRootResource2{
				RootNamedResource: "resource-123",
			},
			expected: "rootNamedResources/resource-123",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.RootNamedResource]()
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

func TestReflectNamer_Root_Parse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "root resource",
			input: "rootNamedResources/resource-123",
			expected: &TestRootResource{
				ID: "resource-123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.RootNamedResource]()
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

func TestReflectNamer_Root_Parse2(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		wantErr  bool
	}{
		{
			name:  "root resource",
			input: "rootNamedResources/resource-123",
			expected: &TestRootResource2{
				RootNamedResource: "resource-123",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.RootNamedResource]()
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

func TestReflectNamer_Root_Errors(t *testing.T) {
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
			input:   &TestRootResource{},
			wantErr: true,
		},
		{
			name:    "invalid pattern",
			input:   &TestRootResource{ID: "123"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			namer, err := NewReflectNamer[*namerv1.RootNamedResource]()
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
