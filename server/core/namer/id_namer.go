package namer

import (
	"fmt"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

type RootVarGetter[ID any] func(id ID, patternIndex int) (string, error)
type RootVarSetter[ID any] func(idVar string, patternIndex int) (ID, error)

// IDNamer is a namer that can be used to format and parse the name
// of a resource that has no parent.
type IDNamer[ID any] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex.
	Format(id ID, patternIndex int) (string, error)

	// Parse parses the name of the resource returning the id
	// and the index of the aip pattern that was used to parse the name.
	Parse(name string) (ID, int, error)
}

// NewIDNamer creates a new root namer for the given resource type.
func NewIDNamer[ID any](
	resource proto.Message,
	getter RootVarGetter[ID],
	setter RootVarSetter[ID],
) (IDNamer[ID], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultIDNamer[ID]{
		patterns: patterns,
		getter:   getter,
		setter:   setter,
	}, nil
}

type defaultIDNamer[ID any] struct {
	patterns []string
	getter   RootVarGetter[ID]
	setter   RootVarSetter[ID]
}

func (n *defaultIDNamer[ID]) Format(id ID, patternIndex int) (string, error) {
	if patternIndex < 0 || patternIndex >= len(n.patterns) {
		return "", fmt.Errorf("invalid pattern index: %d", patternIndex)
	}

	pattern := n.patterns[patternIndex]
	idVar, err := n.getter(id, patternIndex)
	if err != nil {
		return "", err
	}

	return resourcename.Sprint(pattern, idVar), nil
}

func (n *defaultIDNamer[ID]) Parse(name string) (ID, int, error) {
	var id ID

	pattern, patternIndex, err := determineNamePatternIndex(n.patterns, name)
	if err != nil {
		return id, 0, err
	}

	vars, err := scan(pattern, name)
	if err != nil {
		return id, patternIndex, err
	}

	var idVar string
	if len(vars) == 1 {
		idVar = vars[0]
	} else if len(vars) > 1 {
		return id, 0, fmt.Errorf("expected 0 or 1 variables in the name, got %d", len(vars))
	}

	id, err = n.setter(idVar, patternIndex)
	if err != nil {
		return id, 0, err
	}

	return id, patternIndex, nil
}
