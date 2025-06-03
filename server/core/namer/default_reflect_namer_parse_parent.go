package namer

import (
	"fmt"

	"go.einride.tech/aip/resourcename"
)

// MustParseParent calls ParseParent and panics if an error occurs.
// It is a convenience method for callers who expect the parent string to always be valid.
func (n *defaultReflectNamer[T]) MustParseParent(parent string, in *T) int {
	patternIndex, err := n.ParseParent(parent, in)
	if err != nil {
		panic(err)
	}
	return patternIndex
}

// ParseParent parses the parent resource name string and populates the relevant fields in 'in'.
// It returns the pattern index used for parsing, or an error if parsing fails.
func (n *defaultReflectNamer[T]) ParseParent(parent string, in *T) (int, error) {
	patternDetails, patternIndex, err := determineParentPattern(n.patternsDetails, parent)
	if err != nil {
		return 0, err
	}

	err = n._parse(parent, patternDetails, in)
	if err != nil {
		return 0, err
	}

	return patternIndex, nil
}

// determineParentPattern finds the pattern and its index that matches the given parent string.
// If the parent string matches a root pattern (no parent), it returns the root pattern.
// Returns an error if no suitable pattern is found.
func determineParentPattern(patternsDetails map[int]patternDetails, parent string) (patternDetails, int, error) {
	patternIndex := -1

	for i, parentPatternDetails := range patternsDetails {
		// If there are no parent sections, check for a root pattern match
		if len(parentPatternDetails.splitParentPattern) < 1 {
			if resourcename.Match(parentPatternDetails.pattern, parent) || parent == "" {
				return parentPatternDetails, i, nil
			}
			continue
		}

		// Otherwise, check if the parent string matches the parent pattern
		if resourcename.Match(parentPatternDetails.parentPattern, parent) {
			// Save the index and break to return the pattern
			patternIndex = i
			break
		}
	}

	if patternIndex == -1 {
		return patternDetails{}, patternIndex, ErrInvalidParent{
			msg: fmt.Sprintf("invalid parent name: %s", parent),
		}
	}

	return patternsDetails[patternIndex], patternIndex, nil
}
