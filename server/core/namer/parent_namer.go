package namer

import (
	"fmt"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

type SingletonVarGetter[Parent any] func(parent Parent, patternIndex int) ([]string, error)
type SingletonVarSetter[Parent any] func(vars []string, patternIndex int) (Parent, error)
type SingletonParentSetter[Parent any] func(vars []string, patternIndex int) (Parent, error)

// ParentNamer is a namer that can be used to format and parse the name
// of a resource that has a parent but no id.
type ParentNamer[Parent any] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex.
	Format(parent Parent, patternIndex int) (string, error)

	// Parse parses the name of the resource returning the parent
	// and the index of the aip pattern that was used to parse the name.
	Parse(name string) (Parent, int, error)

	// ParseParent parses the parent of the resource returning the parent
	// and the index of the aip pattern that was used to parse the name.
	ParseParent(name string) (Parent, int, error)
}

// NewParentNamer creates a new singleton namer for the given resource type.
func NewParentNamer[Parent any](
	resource proto.Message,
	getter SingletonVarGetter[Parent],
	setter SingletonVarSetter[Parent],
	parentSetter SingletonParentSetter[Parent],
) (ParentNamer[Parent], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultParentNamer[Parent]{
		patterns:     patterns,
		getter:       getter,
		setter:       setter,
		parentSetter: parentSetter,
	}, nil
}

type defaultParentNamer[Parent any] struct {
	patterns     []string
	getter       SingletonVarGetter[Parent]
	setter       SingletonVarSetter[Parent]
	parentSetter SingletonParentSetter[Parent]
}

func (n *defaultParentNamer[Parent]) Format(parent Parent, patternIndex int) (string, error) {
	if patternIndex < 0 || patternIndex >= len(n.patterns) {
		return "", fmt.Errorf("invalid pattern index: %d", patternIndex)
	}

	pattern := n.patterns[patternIndex]
	vars, err := n.getter(parent, patternIndex)
	if err != nil {
		return "", err
	}

	return resourcename.Sprint(pattern, vars...), nil
}

func (n *defaultParentNamer[Parent]) Parse(name string) (Parent, int, error) {
	var parent Parent

	pattern, patternIndex, err := determineNamePatternIndex(n.patterns, name)
	if err != nil {
		return parent, 0, err
	}

	vars, err := scan(pattern, name)
	if err != nil {
		return parent, patternIndex, err
	}

	// Set the values
	parent, err = n.setter(vars, patternIndex)
	if err != nil {
		return parent, 0, err
	}

	return parent, patternIndex, nil
}

func (n *defaultParentNamer[Parent]) ParseParent(parent string) (Parent, int, error) {
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
