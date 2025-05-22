package namer

import (
	"fmt"
	"strings"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

type StandardVarGetter[Parent any, Id any] func(parent Parent, id Id, patternIndex int) []string
type StandardVarSetter[Parent any, Id any] func(vars []string, patternIndex int) (Parent, Id, error)
type StandardParentSetter[Parent any, Id any] func(vars []string, patternIndex int) (Parent, error)

// StandardNamer is a namer that can be used to format and parse the name
// of a resource that has a parent and an id.
type StandardNamer[Parent any, Id any] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex.
	Format(parent Parent, id Id, patternIndex int) (string, error)

	// Parse parses the name of the resource returning the parent structure,
	// the id structure and the index of the aip pattern that was used to
	// parse the name.
	Parse(name string) (Parent, Id, int, error)

	// ParseParent parses the parent of the resource returning the parent
	// structure and the index of the aip pattern that was used to parse
	// the name.
	ParseParent(parent string) (Parent, int, error)
}

// NewStandardNamer creates a new standard namer for the given resource type.
func NewStandardNamer[Parent any, Id any](
	resource proto.Message,
	getter StandardVarGetter[Parent, Id],
	setter StandardVarSetter[Parent, Id],
	parentSetter StandardParentSetter[Parent, Id],
) (StandardNamer[Parent, Id], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultStandardNamer[Parent, Id]{
		patterns:     patterns,
		getter:       getter,
		setter:       setter,
		parentSetter: parentSetter,
	}, nil
}

type defaultStandardNamer[Parent any, Id any] struct {
	patterns     []string
	getter       StandardVarGetter[Parent, Id]
	setter       StandardVarSetter[Parent, Id]
	parentSetter StandardParentSetter[Parent, Id]
}

func (n *defaultStandardNamer[Parent, Id]) Format(parent Parent, id Id, patternIndex int) (string, error) {
	if patternIndex < 0 || patternIndex >= len(n.patterns) {
		return "", fmt.Errorf("invalid pattern index: %d", patternIndex)
	}

	pattern := n.patterns[patternIndex]
	vars := n.getter(parent, id, patternIndex)

	return resourcename.Sprint(pattern, vars...), nil
}

func (n *defaultStandardNamer[Parent, Id]) Parse(name string) (Parent, Id, int, error) {
	var parent Parent
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
		return parent, id, 0, fmt.Errorf("invalid name: %s", name)
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
		return parent, id, 0, err
	}

	// Set the values
	parent, id, err = n.setter(vars, patternIndex)
	if err != nil {
		return parent, id, 0, err
	}

	return parent, id, patternIndex, nil
}

func (n *defaultStandardNamer[Parent, Id]) ParseParent(parent string) (Parent, int, error) {
	var parentParent Parent

	// find the pattern index
	patternIndex := -1
	for i, pattern := range n.patterns {
		splitPattern := strings.Split(pattern, "/")
		if len(splitPattern) < 4 {
			continue
		}
		parentPattern := strings.Join(splitPattern[:len(splitPattern)-2], "/")
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

	parentParent, err = n.parentSetter(vars, patternIndex)
	if err != nil {
		return parentParent, 0, err
	}

	return parentParent, patternIndex, nil
}
