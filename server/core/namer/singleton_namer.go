package namer

import (
	"fmt"
	"strings"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

type SingletonVarGetter[Parent any] func(parent Parent, patternIndex int) []string
type SingletonVarSetter[Parent any] func(vars []string, patternIndex int) (Parent, error)
type SingletonParentGetter[Parent any] func(vars []string, patternIndex int) (Parent, error)

// SingletonNamer is a namer that can be used to format and parse the name
// of a resource that has a parent but no id.
type SingletonNamer[Parent any] interface {
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

// NewSingletonNamer creates a new singleton namer for the given resource type.
func NewSingletonNamer[Parent any](
	resource proto.Message,
	getter SingletonVarGetter[Parent],
	setter SingletonVarSetter[Parent],
	parentGetter SingletonParentGetter[Parent],
) (SingletonNamer[Parent], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultSingletonNamer[Parent]{
		patterns:     patterns,
		getter:       getter,
		setter:       setter,
		parentGetter: parentGetter,
	}, nil
}

type defaultSingletonNamer[Parent any] struct {
	patterns     []string
	getter       SingletonVarGetter[Parent]
	setter       SingletonVarSetter[Parent]
	parentGetter SingletonParentGetter[Parent]
}

func (n *defaultSingletonNamer[Parent]) Format(parent Parent, patternIndex int) (string, error) {
	if patternIndex < 0 || patternIndex >= len(n.patterns) {
		return "", fmt.Errorf("invalid pattern index: %d", patternIndex)
	}

	pattern := n.patterns[patternIndex]
	vars := n.getter(parent, patternIndex)

	return resourcename.Sprint(pattern, vars...), nil
}

func (n *defaultSingletonNamer[Parent]) Parse(name string) (Parent, int, error) {
	var parent Parent

	// find the pattern index
	patternIndex := -1
	for i, pattern := range n.patterns {
		if resourcename.Match(pattern, name) {
			patternIndex = i
			break
		}
	}

	if patternIndex == -1 {
		return parent, 0, fmt.Errorf("invalid name: %s", name)
	}

	pattern := n.patterns[patternIndex]
	var varCount int
	var patternScanner resourcename.Scanner
	patternScanner.Init(pattern)
	for patternScanner.Scan() {
		segment := patternScanner.Segment()
		if segment.IsVariable() {
			varCount++
		}
	}

	vars := make([]string, varCount)
	varsPtr := make([]*string, len(vars))
	for i := range vars {
		varsPtr[i] = &vars[i]
	}

	err := resourcename.Sscan(name, pattern, varsPtr...)
	if err != nil {
		return parent, 0, err
	}

	// Set the values
	parent, err = n.setter(vars, patternIndex)
	if err != nil {
		return parent, 0, err
	}

	return parent, patternIndex, nil
}

func (n *defaultSingletonNamer[Parent]) ParseParent(parent string) (Parent, int, error) {
	var parentParent Parent

	// find the pattern index
	patternIndex := -1
	for i, pattern := range n.patterns {
		splitPattern := strings.Split(pattern, "/")
		if len(splitPattern) < 3 {
			continue
		}
		parentPattern := strings.Join(splitPattern[:len(splitPattern)-1], "/")
		if !resourcename.Match(parentPattern, parent) {
			continue
		}
		patternIndex = i
		break
	}

	if patternIndex == -1 {
		return parentParent, 0, fmt.Errorf("invalid parent: %s", parent)
	}

	pattern := n.patterns[patternIndex]
	var varCount int
	var patternScanner resourcename.Scanner
	patternScanner.Init(pattern)
	for patternScanner.Scan() {
		segment := patternScanner.Segment()
		if segment.IsVariable() {
			varCount++
		}
	}

	vars := make([]string, varCount)
	varsPtr := make([]*string, len(vars))
	for i := range vars {
		varsPtr[i] = &vars[i]
	}

	err := resourcename.Sscan(parent, pattern, varsPtr...)
	if err != nil {
		return parentParent, 0, err
	}

	parentParent, err = n.parentGetter(vars, patternIndex)
	if err != nil {
		return parentParent, 0, err
	}

	return parentParent, patternIndex, nil
}
