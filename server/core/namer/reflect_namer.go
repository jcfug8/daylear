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

// NewParentIDNamer creates a namer for the given resource type that can be used to format and parse the name
// of a resource that has a parent and an id.
func NewReflectNamer[T any](
	resource proto.Message,
) (ReflectNamer[T], error) {
	patterns := getPatterns(resource)
	if len(patterns) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", resource)
	}

	// loop over the patterns and get all the keys we need to parse any name
	patternKeys := map[string]patternVar{}
	for _, pattern := range patterns {
		var patternScanner resourcename.Scanner
		patternScanner.Init(pattern)
		for patternScanner.Scan() {
			segment := patternScanner.Segment()
			if segment.IsVariable() {
				patternKeys[segment.Literal().ResourceID()] = patternVar{
					patternKey: segment.Literal().ResourceID(),
					numField:   []int{},
					fieldType:  nil,
				}
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

	patternKeys, err := determineFullPatternVars(t, typ, patternKeys)
	if err != nil {
		return nil, err
	}

	return &defaultReflectNamer[T]{
		patterns:    patterns,
		patternKeys: patternKeys,
	}, nil
}

func determineFullPatternVars[T any](t T, typ reflect.Type, patternKeys map[string]patternVar) (map[string]patternVar, error) {
	return determineFullPatternVarsRec(t, typ, patternKeys, nil)
}

func determineFullPatternVarsRec[T any](t T, typ reflect.Type, patternKeys map[string]patternVar, parentIndex []int) (map[string]patternVar, error) {
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
				patternKeys[parts[1]] = patternVar{
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
			patternKeys, err = determineFullPatternVarsRec(zero, field.Type, patternKeys, currentIndex)
			if err != nil {
				return nil, err
			}
		}
	}
	return patternKeys, nil
}

type defaultReflectNamer[T any] struct {
	patterns    []string
	patternKeys map[string]patternVar
}

type patternVar struct {
	patternKey string
	numField   []int
	fieldType  reflect.Type
}

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
		var patternScanner resourcename.Scanner
		vars := []string{}
		patternScanner.Init(pattern)
		missingData := false
		for patternScanner.Scan() {
			segment := patternScanner.Segment()
			if segment.IsVariable() {
				patternKey := segment.Literal().ResourceID()
				patternVar, ok := n.patternKeys[patternKey]
				if !ok || len(patternVar.numField) == 0 {
					missingData = true
					lastErr = fmt.Errorf("pattern key %s not found in type %T", patternKey, in)
					break
				}

				value := reflect.ValueOf(in).FieldByIndex(patternVar.numField)
				if value.Kind() != patternVar.fieldType.Kind() {
					missingData = true
					lastErr = fmt.Errorf("field %s in type %T is not of type %s", patternKey, in, patternVar.fieldType.Kind())
					break
				}

				formattedValue := ""
				switch value.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if value.Int() == 0 {
						missingData = true
						lastErr = fmt.Errorf("field %s in type %T is zero", patternKey, in)
						break
					}
					formattedValue = strconv.FormatInt(value.Int(), 10)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if value.Uint() == 0 {
						missingData = true
						lastErr = fmt.Errorf("field %s in type %T is zero", patternKey, in)
						break
					}
					formattedValue = strconv.FormatUint(value.Uint(), 10)
				case reflect.String:
					if value.String() == "" {
						missingData = true
						lastErr = fmt.Errorf("field %s in type %T is empty", patternKey, in)
						break
					}
					formattedValue = value.String()
				default:
					missingData = true
					lastErr = fmt.Errorf("unsupported type %s", value.Kind())
					break
				}
				if missingData {
					break
				}
				vars = append(vars, formattedValue)
			}
		}
		if !missingData {
			return resourcename.Sprint(pattern, vars...), nil
		}
	}
	if lastErr == nil {
		lastErr = fmt.Errorf("no suitable pattern found for type %T", in)
	}
	return "", lastErr
}

func (n *defaultReflectNamer[T]) Parse(name string, in T) (T, int, error) {
	pattern, patternIndex, err := determineNamePattern(n.patterns, name)
	if err != nil {
		return in, 0, err
	}

	in, err = n._parse(name, pattern, n.patternKeys, in)
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

	in, err = n._parse(parent, pattern, n.patternKeys, in)
	if err != nil {
		return in, 0, err
	}

	return in, patternIndex, nil
}

func (n *defaultReflectNamer[T]) _parse(name string, pattern string, patternKeys map[string]patternVar, in T) (T, error) {
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
		patternVar, ok := n.patternKeys[patternKey]
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
