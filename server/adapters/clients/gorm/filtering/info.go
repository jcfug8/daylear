package filtering

import (
	"fmt"
	"strings"
)

// TranspileInfo holds information about the transpilation process.
type TranspileInfo struct {
	fields map[string]bool
	params []any
}

func (info *TranspileInfo) addField(path string) {
	parts := strings.Split(path, ".")

	for i := 0; i < len(parts); i++ {
		path := strings.Join(parts[:i+1], ".")
		info.fields[path] = true
	}
}

// HasField returns true if the field is present in the transpile info.
func (info *TranspileInfo) HasField(path string) bool {
	return info.fields[path]
}

func requireLen[T any](want int, list []T) error {
	if got := len(list); got != want {
		return fmt.Errorf("want %d, got %d", want, got)
	}

	return nil
}
