package namer

import (
	"fmt"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/protobuf/proto"
)

// StandardNamer is a namer that can be used to format and parse the name
// of a resource that has both a parent and an id.
type StandardNamer[Parent parentType, Id idType] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex.
	Format(parent Parent, id Id, patternIndex int) (string, error)

	// Parse parses the name of the resource returning the parent and id
	// and the index of the aip pattern that was used to parse the name.
	Parse(name string) (Parent, Id, int, error)

	// ParseParent parses the parent of the resource returning the parent
	// and the index of the aip pattern that was used to parse the parent.
	ParseParent(name string) (Parent, int, error)
}

func NewStandardNamer[Parent parentType, Id idType](resource proto.Message) (StandardNamer[Parent, Id], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	return &defaultStandardNamer[Parent, Id]{
		patterns: patterns,
	}, nil
}

type defaultStandardNamer[Parent parentType, Id idType] struct {
	patterns []string
}

func (n *defaultStandardNamer[Parent, Id]) Format(parent Parent, id Id, patternIndex int) (string, error) {
	if patternIndex < 0 || patternIndex >= len(n.patterns) {
		return "", fmt.Errorf("invalid pattern index: %d", patternIndex)
	}

	pattern := n.patterns[patternIndex]
	vars := append(parent.GetVars(patternIndex), id.GetId(patternIndex))

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
	vars := append(parent.GetVars(patternIndex), id.GetId(patternIndex))
	varsPtr := make([]*string, len(vars))
	for i := range vars {
		varsPtr[i] = &vars[i]
	}

	err := resourcename.Sscan(name, pattern, varsPtr...)
	if err != nil {
		return parent, id, 0, err
	}

	// Set the values and handle pointer types
	parent = parent.SetVars(patternIndex, vars[:len(vars)-1]).(Parent)
	id = id.SetId(patternIndex, vars[len(vars)-1]).(Id)

	return parent, id, patternIndex, nil
}

func (n *defaultStandardNamer[Parent, Id]) ParseParent(name string) (Parent, int, error) {
	var parent Parent
	parsedParent, _, patternIndex, err := n.Parse(name)
	if err != nil {
		return parent, 0, err
	}
	return parsedParent, patternIndex, nil
}
