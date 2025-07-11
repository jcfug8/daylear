package namer

import (
	"reflect"

	"go.einride.tech/aip/resourcename"
)

// MustFormatParent calls FormatParent and panics if an error occurs.
// It is a convenience method for callers who expect the parent resource to always be valid for formatting.
func (n *defaultReflectNamer) MustFormatParent(in interface{}, options ...FormatReflectNamerOption) string {
	formatted, err := n.FormatParent(in, options...)
	if err != nil {
		panic(err)
	}
	return formatted
}

// FormatParent formats the parent portion of the resource name using the AIP parent pattern specified by the options.
// By default, it uses the parent pattern at index 0. If patternIndex is -1, it tries all parent patterns
// and returns the first one that matches. Returns an error if no parent pattern matches or if
// required parent fields are missing/invalid.
func (n *defaultReflectNamer) FormatParent(in interface{}, options ...FormatReflectNamerOption) (string, error) {
	if in == nil {
		return "", ErrInvalidField{msg: "input is nil"}
	}
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return "", ErrInvalidField{msg: "input is nil pointer"}
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return "", ErrInvalidField{msg: "input is not a struct or pointer to struct"}
	}

	config := formatReflectNamerConfig{
		patternIndex: 0,
	}
	for _, option := range options {
		config = option(config)
	}

	// Build a list of pattern indices to try
	tryPatterns := []int{}
	if config.patternIndex == -1 {
		for i := range n.patternsDetails {
			tryPatterns = append(tryPatterns, i)
		}
	} else {
		tryPatterns = append(tryPatterns, config.patternIndex)
	}

	// Get type cache entry for the input type
	cacheEntry, err := n.getTypeCacheEntry(v.Type())
	if err != nil {
		return "", err
	}

	var lastErr error
	for _, idx := range tryPatterns {
		patternDetails := n.patternsDetails[idx]
		if patternDetails.parentPattern == "" {
			continue
		}
		// Only use parent patterns where all required parent keys are present and non-empty
		allPresent := true
		for _, patternKey := range patternDetails.parentPatternKeys {
			patternVar, ok := cacheEntry.patternKeyDetails[patternKey]
			if !ok || len(patternVar.fieldIndexPath) == 0 {
				allPresent = false
				break
			}
			field := v.FieldByIndex(patternVar.fieldIndexPath)
			if field.Kind() == reflect.Ptr {
				if field.IsNil() {
					allPresent = false
					break
				}
				field = field.Elem()
			}
			if (field.Kind() == reflect.String && field.String() == "") ||
				(field.Kind() == reflect.Int && field.Int() == 0) ||
				(field.Kind() == reflect.Uint && field.Uint() == 0) {
				allPresent = false
				break
			}
		}
		if !allPresent {
			continue
		}
		vars, err := n.extractPatternValuesForParent(patternDetails, in, cacheEntry)
		if err == nil {
			return resourcename.Sprint(patternDetails.parentPattern, vars...), nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = ErrNoPatternFound{
			msg: "no suitable parent pattern found for type",
		}
	}

	return "", lastErr
}

// extractPatternValuesForParent is like extractPatternValues but for parent patterns.
func (n *defaultReflectNamer) extractPatternValuesForParent(patternDetails patternDetails, in interface{}, cacheEntry *typeCacheEntry) ([]string, error) {
	var vars []string
	for _, patternKey := range patternDetails.parentPatternKeys {
		patternVar, ok := cacheEntry.patternKeyDetails[patternKey]
		if !ok || len(patternVar.fieldIndexPath) == 0 {
			return nil, ErrInvalidField{
				msg: "unable to format: parent pattern key " + patternKey + " not found in type",
			}
		}
		formatted, err := formatFieldValue(patternVar, in)
		if err != nil {
			return nil, err
		}
		vars = append(vars, formatted)
	}
	return vars, nil
}
