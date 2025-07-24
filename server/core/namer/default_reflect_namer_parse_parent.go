package namer

import (
	"fmt"
	"reflect"
	"strings"

	"go.einride.tech/aip/resourcename"
)

// MustParseParent calls ParseParent and panics if an error occurs.
// It is a convenience method for callers who expect the parent string to always be valid.
func (n *defaultReflectNamer) MustParseParent(parent string, in interface{}) int {
	patternIndex, err := n.ParseParent(parent, in)
	if err != nil {
		panic(err)
	}
	return patternIndex
}

// ParseParent parses the parent resource name string and populates the relevant fields in 'in'.
// It returns the pattern index used for parsing, or an error if parsing fails.
func (n *defaultReflectNamer) ParseParent(parent string, in interface{}) (int, error) {
	patternDetails, patternIndex, err := determineParentPattern(n.patternsDetails, parent)
	if err != nil {
		return 0, err
	}

	err = n._parseParent(parent, patternDetails, in)
	if err != nil {
		return 0, err
	}

	return patternIndex, nil
}

// _parseParent parses the parent string into the struct pointed to by 'in'.
func (n *defaultReflectNamer) _parseParent(parent string, patternDetails patternDetails, in interface{}) error {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("input must be a non-nil pointer to a struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	// If there are no parent sections, treat as root (empty string)
	if len(patternDetails.splitParentPattern) < 1 {
		if parent == "" {
			return nil
		}
		return ErrInvalidParent{msg: fmt.Sprintf("invalid parent name: %s", parent)}
	}

	splitParent := strings.Split(parent, "/")
	patternSegments := patternDetails.splitParentPattern
	if len(splitParent) != len(patternSegments) {
		return ErrInvalidParent{msg: fmt.Sprintf("invalid parent name: %s", parent)}
	}

	// Extract only the variable segments from the parent string
	var parentVars []string
	for i, seg := range patternSegments {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			parentVars = append(parentVars, splitParent[i])
		}
	}
	if len(parentVars) != len(patternDetails.parentPatternKeys) {
		return ErrInvalidParent{msg: fmt.Sprintf("invalid parent name: %s", parent)}
	}

	cacheEntry, err := n.getTypeCacheEntry(v.Type())
	if err != nil {
		return err
	}

	for i, patternKey := range patternDetails.parentPatternKeys {
		patternVar, ok := cacheEntry.patternKeyDetails[patternKey]
		if !ok {
			return ErrInvalidField{
				msg: fmt.Sprintf("unable to parse: parent pattern key %s not found in type", patternKey),
			}
		}
		if len(patternVar.fieldIndexPath) == 0 {
			return ErrInvalidField{
				msg: fmt.Sprintf("unable to parse: parent pattern key %s found in type but has no field index", patternKey),
			}
		}
		err := setFieldValue(patternVar, parentVars[i], v)
		if err != nil {
			return err
		}
	}

	return nil
}

// determineParentPattern finds the pattern and its index that matches the given parent string.
// If the parent string matches a root pattern (no parent), it returns the root pattern.
// Returns an error if no suitable pattern is found.
func determineParentPattern(patternsDetails map[int]patternDetails, parent string) (patternDetails, int, error) {
	patternIndex := -1

	for i, parentPatternDetails := range patternsDetails {
		// If there are no parent sections, check for a root pattern match
		if len(parentPatternDetails.splitParentPattern) < 1 {
			if resourcename.Match(parentPatternDetails.parentPattern, parent) || parent == "" {
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
