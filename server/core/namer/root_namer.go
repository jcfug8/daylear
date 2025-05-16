package namer

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// RootNamer is a namer that can be used to format and parse the name
// of a resource that has only an id.
type RootNamer[Id any] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex.
	Format(id Id, patternIndex int) (string, error)

	// Parse parses the name of the resource returning the id
	// and the index of the aip pattern that was used to parse the name.
	Parse(name string) (Id, int, error)
}

func NewRootNamer[Id any](resource proto.Message) (RootNamer[Id], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultRootNamer[Id]{
		patterns: patterns,
	}, nil
}

type defaultRootNamer[Id any] struct {
	patterns []string
}

func (n *defaultRootNamer[Id]) Format(id Id, patternIndex int) (string, error) {
	return "", nil
}

func (n *defaultRootNamer[Id]) Parse(name string) (Id, int, error) {
	var id Id
	return id, 0, nil
}
