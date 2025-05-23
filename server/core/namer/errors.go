package namer

import "fmt"

// Errors
var (
	ErrInvalidParent       = fmt.Errorf("invalid parent name")
	ErrInvalidName         = fmt.Errorf("invalid name")
	ErrInvalidPatternIndex = fmt.Errorf("invalid pattern index")
)
