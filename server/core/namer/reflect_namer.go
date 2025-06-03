package namer

// ReflectNamer defines an interface for formatting and parsing resource names
// according to AIP (Google API Improvement Proposals) patterns using reflection.
// It supports resources with parent and ID fields, and can handle multiple patterns.
type ReflectNamer[T any] interface {
	// Format returns the formatted resource name string using the specified AIP pattern.
	// The pattern is selected by the provided options (e.g., AsPatternIndex).
	// If no pattern index is provided, index 0 is used. If -1 is provided, all patterns
	// are tried and the first that matches is used.
	//
	// Returns an error if the resource cannot be formatted due to missing/invalid fields
	// or if the pattern index is invalid.
	Format(in T, options ...formatReflectNamerOption) (string, error)

	// MustFormat is like Format but panics if an error occurs.
	// Use when you expect formatting to always succeed.
	MustFormat(in T, options ...formatReflectNamerOption) string

	// FormatParent returns the formatted parent resource name string using the specified AIP pattern.
	// The pattern is selected by the provided options (e.g., AsPatternIndex).
	// If no pattern index is provided, index 0 is used. If -1 is provided, all patterns
	// are tried and the first that matches is used.
	//
	// Returns an error if the resource cannot be formatted due to missing/invalid fields
	// or if the pattern index is invalid.
	FormatParent(in T, options ...formatReflectNamerOption) (string, error)

	// MustFormatParent is like FormatParent but panics if an error occurs.
	// Use when you expect formatting to always succeed.
	MustFormatParent(in T, options ...formatReflectNamerOption) string

	// Parse parses a resource name string into the struct pointed to by 'in'.
	// Returns the index of the pattern used for parsing, or an error if parsing fails.
	//
	// Returns an error if the resource cannot be parsed due to missing/invalid fields or
	// if the resource name passed in is not a valid resource name for the resource.
	Parse(name string, in *T) (int, error)

	// MustParse is like Parse but panics if an error occurs.
	// Use when you expect parsing to always succeed.
	MustParse(name string, in *T) int

	// ParseParent parses a parent resource name string into the struct pointed to by 'in'.
	// Returns the index of the pattern used for parsing, or an error if parsing fails.
	//
	// Returns an error if the parent cannot be parsed due to missing/invalid fields or
	// if the parent passed in is not a valid parent for the resource.
	ParseParent(parent string, in *T) (int, error)

	// MustParseParent is like ParseParent but panics if an error occurs.
	// Use when you expect parsing to always succeed.
	MustParseParent(parent string, in *T) int
}
