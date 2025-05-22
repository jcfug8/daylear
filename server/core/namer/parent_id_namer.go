package namer

import (
	"fmt"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

type StandardVarGetter[Parent any, ID any] func(parent Parent, id ID, patternIndex int) ([]string, error)
type StandardVarSetter[Parent any, ID any] func(vars []string, patternIndex int) (Parent, ID, error)
type StandardParentSetter[Parent any, ID any] func(vars []string, patternIndex int) (Parent, error)

// ParentIDNamer is a namer that can be used to format and parse the name
// of a resource that has a parent and an id.
type ParentIDNamer[Parent any, ID any] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex.
	Format(parent Parent, id ID, patternIndex int) (string, error)

	// Parse parses the name of the resource returning the parent structure,
	// the id structure and the index of the aip pattern that was used to
	// parse the name.
	Parse(name string) (Parent, ID, int, error)

	// ParseParent parses the parent of the resource returning the parent
	// structure and the index of the aip pattern that was used to parse
	// the name.
	ParseParent(parent string) (Parent, int, error)
}

// NewParentIDNamer creates a new standard namer for the given resource type.
func NewParentIDNamer[Parent any, ID any](
	resource proto.Message,
	getter StandardVarGetter[Parent, ID],
	setter StandardVarSetter[Parent, ID],
	parentSetter StandardParentSetter[Parent, ID],
) (ParentIDNamer[Parent, ID], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultParentIDNamer[Parent, ID]{
		patterns:     patterns,
		getter:       getter,
		setter:       setter,
		parentSetter: parentSetter,
	}, nil
}

type defaultParentIDNamer[Parent any, ID any] struct {
	patterns     []string
	getter       StandardVarGetter[Parent, ID]
	setter       StandardVarSetter[Parent, ID]
	parentSetter StandardParentSetter[Parent, ID]
}

func (n *defaultParentIDNamer[Parent, ID]) Format(parent Parent, id ID, patternIndex int) (string, error) {
	if patternIndex < 0 || patternIndex >= len(n.patterns) {
		return "", fmt.Errorf("invalid pattern index: %d", patternIndex)
	}

	pattern := n.patterns[patternIndex]
	vars, err := n.getter(parent, id, patternIndex)
	if err != nil {
		return "", err
	}

	return resourcename.Sprint(pattern, vars...), nil
}

func (n *defaultParentIDNamer[Parent, ID]) Parse(name string) (Parent, ID, int, error) {
	var parent Parent
	var id ID

	pattern, patternIndex, err := determineNamePatternIndex(n.patterns, name)
	if err != nil {
		return parent, id, 0, err
	}

	vars, err := scan(pattern, name)
	if err != nil {
		return parent, id, patternIndex, err
	}

	// Set the values
	parent, id, err = n.setter(vars, patternIndex)
	if err != nil {
		return parent, id, 0, err
	}

	return parent, id, patternIndex, nil
}

func (n *defaultParentIDNamer[Parent, ID]) ParseParent(parent string) (Parent, int, error) {
	var parentStruct Parent

	pattern, patternIndex, err := determineParentPatternIndex(n.patterns, parent)
	if err != nil {
		return parentStruct, 0, err
	}

	if pattern == "" {
		return parentStruct, patternIndex, nil
	}

	vars, err := scan(pattern, parent)
	if err != nil {
		return parentStruct, patternIndex, err
	}

	parentStruct, err = n.parentSetter(vars, patternIndex)
	if err != nil {
		return parentStruct, 0, err
	}

	return parentStruct, patternIndex, nil
}
