package namer

import (
	"fmt"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

type RootVarGetter[Id any] func(id Id, patternIndex int) string
type RootVarSetter[Id any] func(idVar string, patternIndex int) (Id, error)

// RootNamer is a namer that can be used to format and parse the name
// of a resource that has no parent.
type RootNamer[Id any] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex.
	Format(id Id, patternIndex int) (string, error)

	// Parse parses the name of the resource returning the id
	// and the index of the aip pattern that was used to parse the name.
	Parse(name string) (Id, int, error)
}

// NewRootNamer creates a new root namer for the given resource type.
func NewRootNamer[Id any](
	resource proto.Message,
	getter RootVarGetter[Id],
	setter RootVarSetter[Id],
) (RootNamer[Id], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultRootNamer[Id]{
		patterns: patterns,
		getter:   getter,
		setter:   setter,
	}, nil
}

type defaultRootNamer[Id any] struct {
	patterns []string
	getter   RootVarGetter[Id]
	setter   RootVarSetter[Id]
}

func (n *defaultRootNamer[Id]) Format(id Id, patternIndex int) (string, error) {
	if patternIndex < 0 || patternIndex >= len(n.patterns) {
		return "", fmt.Errorf("invalid pattern index: %d", patternIndex)
	}

	pattern := n.patterns[patternIndex]
	vars := []string{n.getter(id, patternIndex)}

	return resourcename.Sprint(pattern, vars...), nil
}

func (n *defaultRootNamer[Id]) Parse(name string) (Id, int, error) {
	var id Id

	// find the pattern index
	patternIndex := -1
	for i, pattern := range n.patterns {
		if resourcename.Match(pattern, name) {
			patternIndex = i
			break
		}
	}

	if patternIndex == -1 {
		return id, 0, fmt.Errorf("invalid name: %s", name)
	}

	pattern := n.patterns[patternIndex]
	idVar := ""

	err := resourcename.Sscan(name, pattern, &idVar)
	if err != nil {
		return id, 0, err
	}

	id, err = n.setter(idVar, patternIndex)
	if err != nil {
		return id, 0, err
	}

	return id, patternIndex, nil
}
