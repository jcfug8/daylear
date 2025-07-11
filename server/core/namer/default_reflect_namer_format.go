package namer

import (
	"fmt"
	"reflect"
	"strconv"

	"go.einride.tech/aip/resourcename"
)

// MustFormat calls Format and panics if an error occurs.
// It is a convenience method for callers who expect the resource to always be valid for formatting.
func (n *defaultReflectNamer) MustFormat(in interface{}, options ...FormatReflectNamerOption) string {
	formatted, err := n.Format(in, options...)
	if err != nil {
		panic(err)
	}
	return formatted
}

// Format formats the name of the resource using the AIP pattern specified by the options.
// By default, it uses the pattern at index 0. If patternIndex is -1, it tries all patterns
// and returns the first one that matches. Returns an error if no pattern matches or if
// required fields are missing/invalid.
func (n *defaultReflectNamer) Format(in interface{}, options ...FormatReflectNamerOption) (string, error) {
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

	// Get type cache entry for the input type
	cacheEntry, err := n.getTypeCacheEntry(v.Type())
	if err != nil {
		return "", err
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

	var lastErr error
	for _, idx := range tryPatterns {
		patternDetails := n.patternsDetails[idx]
		// Only use patterns where all required keys are present and non-empty
		allPresent := true
		for _, patternKey := range patternDetails.patternKeys {
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
		vars, err := n.extractPatternValues(patternDetails, in, cacheEntry)
		if err == nil {
			return resourcename.Sprint(patternDetails.pattern, vars...), nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = ErrNoPatternFound{
			msg: fmt.Sprintf("no suitable pattern found for type %T", in),
		}
	}

	return "", lastErr
}

// extractPatternValues extracts the values for all variables in the given pattern from the input struct.
// Returns an error if any variable is missing, zero, or empty.
func (n *defaultReflectNamer) extractPatternValues(patternDetails patternDetails, in interface{}, cacheEntry *typeCacheEntry) ([]string, error) {
	var vars []string
	for _, patternKey := range patternDetails.patternKeys {
		patternVar, ok := cacheEntry.patternKeyDetails[patternKey]
		if !ok || len(patternVar.fieldIndexPath) == 0 {
			return nil, ErrInvalidField{
				msg: fmt.Sprintf("unable to format: pattern key %s not found in type %T", patternKey, in),
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

// formatFieldValue formats a struct field value as a string for use in a resource name.
// Returns an error if the value is zero, empty, nil, or of an unsupported type.
func formatFieldValue(pv patternKeyDetails, in interface{}) (string, error) {
	// Get the field from the input struct using the index path
	field := reflect.ValueOf(in)
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	field = field.FieldByIndex(pv.fieldIndexPath)
	if field.Kind() != pv.fieldType.Kind() {
		return "", ErrInvalidFieldType{
			msg: fmt.Sprintf("unable to format: key %s is not of type %s in input struct %T", pv.patternKey, pv.fieldType.Kind(), in),
		}
	}

	// Handle pointer fields
	if field.Kind() == reflect.Ptr {
		if field.IsNil() {
			return "", ErrInvalidFieldValue{
				msg: fmt.Sprintf("unable to format: key %s is nil in input struct %T", pv.patternKey, in),
			}
		}
		field = field.Elem()
	}

	// Format the field depending on its type. Only int, uint, and string are supported.
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if field.Int() == 0 {
			return "", ErrInvalidFieldValue{
				msg: fmt.Sprintf("unable to format: int field is zero for key %s in input struct %T", pv.patternKey, in),
			}
		}
		return strconv.FormatInt(field.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if field.Uint() == 0 {
			return "", ErrInvalidFieldValue{
				msg: fmt.Sprintf("unable to format: uint field is zero for key %s in input struct %T", pv.patternKey, in),
			}
		}
		return strconv.FormatUint(field.Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		if field.Float() == 0 {
			return "", ErrInvalidFieldValue{
				msg: fmt.Sprintf("unable to format: float field is zero for key %s in input struct %T", pv.patternKey, in),
			}
		}
		return strconv.FormatFloat(field.Float(), 'f', -1, 64), nil
	case reflect.String:
		if field.String() == "" {
			return "", ErrInvalidFieldValue{
				msg: fmt.Sprintf("unable to format: string field is empty for key %s in input struct %T", pv.patternKey, in),
			}
		}
		return field.String(), nil
	default:
		return "", ErrInvalidFieldType{
			msg: fmt.Sprintf("unable to format: unsupported type %s for key %s in input struct %T", field.Kind(), pv.patternKey, in),
		}
	}
}
