package namer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.einride.tech/aip/resourcename"
)

// MustParse calls Parse and panics if an error occurs.
// It is a convenience method for callers who expect the resource name to always be valid.
func (n *defaultReflectNamer) MustParse(name string, in interface{}) int {
	patternIndex, err := n.Parse(name, in)
	if err != nil {
		panic(err)
	}
	return patternIndex
}

// Parse parses a resource name string into the struct pointed to by 'in'.
// Returns the index of the pattern used for parsing, or an error if parsing fails.
func (n *defaultReflectNamer) Parse(name string, in interface{}) (int, error) {
	pattern, patternIndex, err := determineNamePattern(n.patternsDetails, name)
	if err != nil {
		return 0, err
	}

	err = n._parse(name, pattern, in)
	if err != nil {
		return 0, err
	}

	return patternIndex, nil
}

// determineNamePattern - determines the pattern index of the given name. It will
// return the pattern index and the pattern if found, otherwise it will return an error.
func determineNamePattern(patternsDetails map[int]patternDetails, name string) (patternDetails, int, error) {
	patternIndex := -1

	for i, patternDetails := range patternsDetails {
		// check if the name pattern matches the name sent in
		if resourcename.Match(patternDetails.pattern, name) {
			patternIndex = i
			break
		}
	}

	if patternIndex == -1 {
		return patternDetails{}, patternIndex, ErrInvalidName{
			msg: fmt.Sprintf("invalid name: %s", name),
		}
	}

	return patternsDetails[patternIndex], patternIndex, nil
}

func (n *defaultReflectNamer) _parse(name string, patternDetails patternDetails, in interface{}) error {
	v := reflect.ValueOf(in)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("input must be a non-nil pointer to a struct")
	}
	v = v.Elem()
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("input must be a pointer to a struct")
	}

	splitName := strings.Split(name, "/")
	for i, segment := range splitName {
		if i%2 == 0 {
			continue
		}
		patternKey := patternDetails.splitPattern[i][1 : len(patternDetails.splitPattern[i])-1]
		cacheEntry, err := n.getTypeCacheEntry(v.Type())
		if err != nil {
			return err
		}
		patternVar, ok := cacheEntry.patternKeyDetails[patternKey]
		if !ok {
			return ErrInvalidField{
				msg: fmt.Sprintf("unable to parse: pattern key %s not found in type %T", patternKey, in),
			}
		}

		if len(patternVar.fieldIndexPath) == 0 {
			return ErrInvalidField{
				msg: fmt.Sprintf("unable to parse: pattern key %s found in type %T but has no field index", patternKey, in),
			}
		}

		err = setFieldValue(patternVar, segment, v)
		if err != nil {
			return err
		}
	}

	return nil
}

// setFieldValue sets the value of the field specified by patternVar in the input struct.
func setFieldValue(pv patternKeyDetails, segment string, v reflect.Value) error {
	// this loop is used to traverse the field index path and set the value of the field.
	// we need to loop instead of just setting the value because we need to handle nested
	// pointer to structs. If the pointer is nil, we need to initialize it.
	for i, idx := range pv.fieldIndexPath {
		field := v.Field(idx)
		// If this is not the last index in the path, we need to traverse deeper
		if i < len(pv.fieldIndexPath)-1 {
			// If the field is a pointer to a struct and is nil, initialize it
			if field.Kind() == reflect.Ptr && field.Type().Elem().Kind() == reflect.Struct {
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				v = field.Elem()
			} else if field.Kind() == reflect.Struct {
				v = field
			} else {
				return fmt.Errorf("unexpected kind %s in field path", field.Kind())
			}
		} else { // Last field: set the value
			// if the segment is a wildcard, we should not set the value
			if segment == "-" {
				return nil
			}
			kind := field.Kind()
			// handle if the field is a pointer
			if kind == reflect.Ptr {
				elemKind := field.Type().Elem().Kind()
				if field.IsNil() {
					field.Set(reflect.New(field.Type().Elem()))
				}
				field = field.Elem()
				kind = elemKind

			}
			switch kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				v, err := strconv.ParseInt(segment, 10, 64)
				if err != nil {
					return ErrInvalidFieldValue{msg: fmt.Sprintf("unable to parse: invalid int value %s for key %s in input struct", segment, pv.patternKey)}
				}
				field.SetInt(v)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				v, err := strconv.ParseUint(segment, 10, 64)
				if err != nil {
					return ErrInvalidFieldValue{msg: fmt.Sprintf("unable to parse: invalid uint value %s for key %s in input struct", segment, pv.patternKey)}
				}
				field.SetUint(v)
			case reflect.Float32, reflect.Float64:
				v, err := strconv.ParseFloat(segment, 64)
				if err != nil {
					return ErrInvalidFieldValue{msg: fmt.Sprintf("unable to parse: invalid float value %s for key %s in input struct", segment, pv.patternKey)}
				}
				field.SetFloat(v)
			case reflect.String:
				field.SetString(segment)
			default:
				return ErrInvalidFieldType{msg: fmt.Sprintf("unable to parse: unsupported type %s for key %s in input struct", field.Kind(), pv.patternKey)}
			}
		}
	}
	return nil
}
