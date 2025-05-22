package namer

import (
	"fmt"
	"strings"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

type IDVarGetter[ID any] func(id ID, patternIndex int) ([]string, error)
type ParentVarGetter[Parent any] func(parent Parent, patternIndex int) ([]string, error)
type ParentIDVarGetter[Parent any, ID any] func(parent Parent, id ID, patternIndex int) ([]string, error)

type IDVarSetter[ID any] func(vars []string, patternIndex int) (ID, error)
type ParentVarSetter[Parent any] func(vars []string, patternIndex int) (Parent, error)
type ParentIDVarSetter[Parent any, ID any] func(vars []string, patternIndex int) (Parent, ID, error)

func getPatterns(resource proto.Message) []string {
	resourceOption := proto.GetExtension(
		resource.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)
	return resourceOption.Pattern
}

func determineNamePatternIndex(patterns []string, name string) (string, int, error) {
	patternIndex := -1

	for i, pattern := range patterns {
		// check if the name pattern matches the name sent in
		if resourcename.Match(pattern, name) {
			patternIndex = i
			break
		}
	}

	if patternIndex == -1 {
		return "", patternIndex, fmt.Errorf("invalid name: %s", name)
	}

	return patterns[patternIndex], patternIndex, nil
}

func determineParentPatternIndex(patterns []string, parent string) (string, int, error) {
	patternIndex := -1
	parentPattern := ""

	for i, pattern := range patterns {
		splitNamePattern := strings.Split(pattern, "/")

		// if there are less than 3 sections, it cannot have a parent
		if len(splitNamePattern) < 3 {
			if resourcename.Match(pattern, parent) {
				return "", i, nil
			}
			continue
		}

		// get the parent pattern sections as if its a singleton pattern
		splitParentPattern := splitNamePattern[:len(splitNamePattern)-1]
		// if there is an odd number of sections, it was a collection name pattern
		// so take another section out
		if len(splitParentPattern)%2 == 1 {
			splitParentPattern = splitParentPattern[:len(splitParentPattern)-1]
		}

		// check if the parent pattern matches the parent sent in
		parentPattern = strings.Join(splitParentPattern, "/")
		if resourcename.Match(parentPattern, parent) {
			// save the index and break out so we can return the pattern
			patternIndex = i
			break
		}
	}

	if patternIndex == -1 {
		return parentPattern, patternIndex, fmt.Errorf("invalid parent: %s", parent)
	}

	return parentPattern, patternIndex, nil
}

func scan(pattern string, nameOrParent string) ([]string, error) {
	var varCount int

	// count the number of variables in the pattern
	var patternScanner resourcename.Scanner
	patternScanner.Init(pattern)
	for patternScanner.Scan() {
		segment := patternScanner.Segment()
		if segment.IsVariable() {
			varCount++
		}
	}

	// create a slices of the correct size
	vars := make([]string, varCount)
	varsPtr := make([]*string, len(vars))
	for i := range vars {
		varsPtr[i] = &vars[i]
	}

	// scan the name or parent into the pointer slice
	err := resourcename.Sscan(nameOrParent, pattern, varsPtr...)
	if err != nil {
		return nil, err
	}

	return vars, nil
}
