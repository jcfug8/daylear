package namer

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

// SingletonNamer is a namer that can be used to format and parse the name
// of a resource that has only parent (a singleton).
type SingletonNamer[Parent any] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex.
	Format(parent Parent, patternIndex int) (string, error)

	// Parse parses the name of the resource returning the parent
	// and the index of the aip pattern that was used to parse the name.
	Parse(name string) (Parent, int, error)

	// ParseParent parses the parent of the resource returning the parent
	// and the index of the aip pattern that was used to parse the parent.
	ParseParent(name string) (Parent, int, error)
}

func NewSingletonNamer[Parent any](resource proto.Message) (SingletonNamer[Parent], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultSingletonNamer[Parent]{
		patterns: patterns,
	}, nil
}

type defaultSingletonNamer[Parent any] struct {
	patterns []string
}

func (n *defaultSingletonNamer[Parent]) Format(parent Parent, patternIndex int) (string, error) {
	return "", nil
}

func (n *defaultSingletonNamer[Parent]) Parse(name string) (Parent, int, error) {
	var parent Parent
	return parent, 0, nil
}

func (n *defaultSingletonNamer[Parent]) ParseParent(name string) (Parent, int, error) {
	var parent Parent
	return parent, 0, nil
}
