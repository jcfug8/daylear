package namer

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// ReflectNamer is a namer that can be used to format and parse the name
// of a resource that has a parent and an id.
type ReflectNamer[T any] interface {
	// Format formats the name of the resource using the aip pattern
	// specified by the patternIndex. If patternIndex is -1, it will
	// format using the first pattern that is possible.
	Format(patternIndex int, in T) (string, error)

	// Parse parses the name of the resource returning the parent structure,
	// the id structure and the index of the aip pattern that was used to
	// parse the name.
	Parse(name string, in T) (T, int, error)

	// ParseParent parses the parent of the resource returning the parent
	// structure and the index of the aip pattern that was used to parse
	// the name.
	ParseParent(parent string, in T) (T, int, error)
}

// FieldIndexPath is a slice of ints representing the path to a struct field (for nested fields).
type FieldIndexPath []int

// patternVar holds metadata for a pattern variable.
type patternVar struct {
	patternKey string
	numField   FieldIndexPath
	fieldType  reflect.Type
}

// defaultReflectNamer implements ReflectNamer using reflection.
type defaultReflectNamer[T any] struct {
	patterns      []string
	patternVarMap map[string]patternVar
}

// NewReflectNamer creates a namer for the given resource type that can be used to format and parse the name
// of a resource that has a parent and an id.
func NewReflectNamer[T any](
	resource proto.Message,
) (ReflectNamer[T], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	// loop over the patterns and get all the keys we need to parse any name
	patternVarMap := map[string]patternVar{}
	for _, pattern := range patterns {
		for _, key := range extractPatternVars(pattern) {
			patternVarMap[key] = patternVar{
				patternKey: key,
				numField:   FieldIndexPath{},
				fieldType:  nil,
			}
		}
	}

	var t T
	typ := reflect.TypeOf(t)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type T must be a struct")
	}

	patternVarMap, err := determineFullPatternVars(t, typ, patternVarMap)
	if err != nil {
		return nil, err
	}

	return &defaultReflectNamer[T]{
		patterns:      patterns,
		patternVarMap: patternVarMap,
	}, nil
}

// Format formats the name of the resource using the aip pattern specified by the patternIndex.
// If patternIndex is -1, it will format using the first pattern that is possible.
func (n *defaultReflectNamer[T]) Format(patternIndex int, in T) (string, error) {
	tryPatterns := []int{}
	if patternIndex == -1 {
		for i := range n.patterns {
			tryPatterns = append(tryPatterns, i)
		}
	} else {
		tryPatterns = append(tryPatterns, patternIndex)
	}

	var lastErr error
	for _, idx := range tryPatterns {
		pattern := n.patterns[idx]
		vars, err := n.extractPatternValues(pattern, in)
		if err == nil {
			return resourcename.Sprint(pattern, vars...), nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no suitable pattern found for type %T", in)
	}
	return "", lastErr
}

// extractPatternValues returns the values for the variables in the pattern from the input struct.
// Returns error if any variable is missing or zero/empty.
func (n *defaultReflectNamer[T]) extractPatternValues(pattern string, in T) ([]string, error) {
	var vars []string
	for _, patternKey := range extractPatternVars(pattern) {
		patternVar, err := n.getPatternVar(patternKey)
		if err != nil {
			return nil, err
		}
		value, err := getFieldValue(in, patternVar)
		if err != nil {
			return nil, err
		}
		formatted, err := formatReflectValue(value)
		if err != nil {
			return nil, err
		}
		vars = append(vars, formatted)
	}
	return vars, nil
}

// getPatternVar returns the patternVar for a given key, or an error if not found.
func (n *defaultReflectNamer[T]) getPatternVar(key string) (patternVar, error) {
	pv, ok := n.patternVarMap[key]
	if !ok || len(pv.numField) == 0 {
		return patternVar{}, fmt.Errorf("pattern key %s not found in type", key)
	}
	return pv, nil
}

// getFieldValue returns the reflect.Value for the field specified by patternVar in the input struct.
func getFieldValue[T any](in T, pv patternVar) (reflect.Value, error) {
	val := reflect.ValueOf(in)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if !val.IsValid() {
		return reflect.Value{}, fmt.Errorf("invalid value for input struct")
	}
	field := val.FieldByIndex(pv.numField)
	if field.Kind() != pv.fieldType.Kind() {
		return reflect.Value{}, fmt.Errorf("field %s is not of type %s", pv.patternKey, pv.fieldType.Kind())
	}
	return field, nil
}

// formatReflectValue formats a reflect.Value as a string, or returns an error if zero/empty/unsupported.
func formatReflectValue(value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value.Int() == 0 {
			return "", fmt.Errorf("int field is zero")
		}
		return strconv.FormatInt(value.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value.Uint() == 0 {
			return "", fmt.Errorf("uint field is zero")
		}
		return strconv.FormatUint(value.Uint(), 10), nil
	case reflect.String:
		if value.String() == "" {
			return "", fmt.Errorf("string field is empty")
		}
		return value.String(), nil
	default:
		return "", fmt.Errorf("unsupported type %s", value.Kind())
	}
}

// extractPatternVars returns the variable keys from a pattern string.
func extractPatternVars(pattern string) []string {
	var keys []string
	var scanner resourcename.Scanner
	scanner.Init(pattern)
	for scanner.Scan() {
		segment := scanner.Segment()
		if segment.IsVariable() {
			keys = append(keys, segment.Literal().ResourceID())
		}
	}
	return keys
}

func determineFullPatternVars[T any](t T, typ reflect.Type, patternVarMap map[string]patternVar) (map[string]patternVar, error) {
	return determineFullPatternVarsRec(t, typ, patternVarMap, nil)
}

func determineFullPatternVarsRec[T any](t T, typ reflect.Type, patternVarMap map[string]patternVar, parentIndex []int) (map[string]patternVar, error) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("aip_pattern")
		currentIndex := append(parentIndex, i)
		if tag != "" {
			parts := strings.Split(tag, "=")
			if len(parts) == 2 && parts[0] == "key" {
				if field.Type.Kind() == reflect.Struct {
					return nil, fmt.Errorf("field %s in type %T is a struct (tagged field should not be a struct)", field.Name, t)
				} else if field.Type.Kind() == reflect.Slice {
					return nil, fmt.Errorf("field %s in type %T is a slice", field.Name, t)
				} else if field.Type.Kind() == reflect.Map {
					return nil, fmt.Errorf("field %s in type %T is a map", field.Name, t)
				} else if field.Type.Kind() == reflect.Chan {
					return nil, fmt.Errorf("field %s in type %T is a channel", field.Name, t)
				} else if field.Type.Kind() == reflect.Func {
					return nil, fmt.Errorf("field %s in type %T is a function", field.Name, t)
				}
				patternVarMap[parts[1]] = patternVar{
					patternKey: parts[1],
					numField:   currentIndex,
					fieldType:  field.Type,
				}
			}
		}
		// If the field is a struct (and not time.Time or similar), recurse
		if field.Type.Kind() == reflect.Struct && field.Type.PkgPath() != "time" {
			var zero T
			var err error
			patternVarMap, err = determineFullPatternVarsRec(zero, field.Type, patternVarMap, currentIndex)
			if err != nil {
				return nil, err
			}
		}
	}
	return patternVarMap, nil
}

func (n *defaultReflectNamer[T]) Parse(name string, in T) (T, int, error) {
	pattern, patternIndex, err := determineNamePattern(n.patterns, name)
	if err != nil {
		return in, 0, err
	}

	in, err = n._parse(name, pattern, in)
	if err != nil {
		return in, 0, err
	}

	return in, patternIndex, nil
}

func (n *defaultReflectNamer[T]) ParseParent(parent string, in T) (T, int, error) {
	pattern, patternIndex, err := determineParentPattern(n.patterns, parent)
	if err != nil {
		return in, 0, err
	}

	in, err = n._parse(parent, pattern, in)
	if err != nil {
		return in, 0, err
	}

	return in, patternIndex, nil
}

func (n *defaultReflectNamer[T]) _parse(name string, pattern string, in T) (T, error) {
	var ptrIn reflect.Value
	val := reflect.ValueOf(in)
	typ := val.Type()
	if typ.Kind() == reflect.Ptr {
		ptrIn = val
	} else {
		// make a pointer to the type
		ptrIn = reflect.New(typ)
		ptrIn.Elem().Set(val)
	}

	splitName := strings.Split(name, "/")
	splitPattern := strings.Split(pattern, "/")
	for i, segment := range splitName {
		if i%2 == 0 {
			continue
		}
		patternKey := splitPattern[i][1 : len(splitPattern[i])-1]
		patternVar, ok := n.patternVarMap[patternKey]
		if !ok {
			return in, fmt.Errorf("pattern key %s not found in type %T", patternKey, in)
		}

		if len(patternVar.numField) == 0 {
			return in, fmt.Errorf("pattern key %s not found in type %T", patternKey, in)
		}

		value := ptrIn.Elem().FieldByIndex(patternVar.numField)
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.ParseInt(segment, 10, 64)
			if err != nil {
				return in, err
			}
			value.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.ParseUint(segment, 10, 64)
			if err != nil {
				return in, err
			}
			value.SetUint(v)
		case reflect.String:
			value.SetString(segment)
		default:
			return in, fmt.Errorf("unsupported type %s", value.Kind())
		}
	}

	if ptrIn.Kind() == reflect.Ptr {
		in = ptrIn.Elem().Interface().(T)
	}

	return in, nil
}

// getPatterns - parse the given proto resource and returns the google aip name patterns.
func getPatterns(resource proto.Message) []string {
	resourceOption := proto.GetExtension(
		resource.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)
	return resourceOption.Pattern
}

// determineNamePattern - determines the pattern index of the given name. It will
// return the pattern index and the pattern if found, otherwise it will return an error.
func determineNamePattern(patterns []string, name string) (string, int, error) {
	patternIndex := -1

	for i, pattern := range patterns {
		// check if the name pattern matches the name sent in
		if resourcename.Match(pattern, name) {
			patternIndex = i
			break
		}
	}

	if patternIndex == -1 {
		return "", patternIndex, ErrInvalidName
	}

	return patterns[patternIndex], patternIndex, nil
}

// determineParentPattern - determines the pattern index of the given parent. It will
// return the pattern index and the pattern if found, otherwise it will return an error.
// If the parent passed in is a root pattern, a pattern with no parent, it will return the
// root pattern if found, otherwise it will return an error.
func determineParentPattern(patterns []string, parent string) (string, int, error) {
	patternIndex := -1
	parentPattern := ""

	for i, pattern := range patterns {
		splitNamePattern := strings.Split(pattern, "/")

		// if there are less than 3 sections, it cannot have a parent
		if len(splitNamePattern) < 3 {
			if resourcename.Match(pattern, parent) || parent == "" {
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
		return parentPattern, patternIndex, ErrInvalidParent
	}

	return parentPattern, patternIndex, nil
}
