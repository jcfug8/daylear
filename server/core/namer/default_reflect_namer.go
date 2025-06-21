package namer

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"

	"go.einride.tech/aip/resourcename"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// FieldIndexPath represents the path to a struct field, supporting nested fields.
type FieldIndexPath []int

// patternKeyDetails contains metadata for a pattern variable, including its key name,
// the index path to the field in the struct, and its type.
type patternKeyDetails struct {
	patternKey     string
	fieldIndexPath FieldIndexPath
	fieldType      reflect.Type
}

// patternDetails holds information about a resource name pattern, including its index,
// the pattern string, split pattern segments, variable keys, and parent pattern details.
type patternDetails struct {
	patternIndex       int
	pattern            string
	splitPattern       []string
	patternKeys        []string
	parentPattern      string
	splitParentPattern []string
	parentPatternKeys  []string
}

// typeCacheEntry holds the cached reflection data for a specific type
type typeCacheEntry struct {
	patternKeyDetails map[string]patternKeyDetails
}

// defaultReflectNamer implements ReflectNamer using Go reflection to map struct fields
// to resource name pattern variables.
type defaultReflectNamer struct {
	patternsDetails map[int]patternDetails
	typeCache       sync.Map // map[reflect.Type]*typeCacheEntry
}

// NewReflectNamer creates a ReflectNamer for the given proto message type ProtoType.
// It uses reflection to map struct fields to AIP resource name pattern variables.
//
// If multiple fields are named or tagged with the same key, the last one is used.
// If struct tags are used, they must be of the form `aip_pattern:"key=pattern_key"`.
// Extra patterns can be added (starting at index 100) via options.
// By default, all pattern keys must be present in the struct; this can be disabled with DisableStrictNoMissingStructKeys.
func NewReflectNamer[ProtoType proto.Message](
	options ...newReflectNamerOption,
) (ReflectNamer, error) {
	// Set up configuration from options
	config := newReflectNamerConfig{
		disableStrictNoMissingStructKeys: false,
	}
	for _, option := range options {
		config = option(config)
	}

	// Extract resource name patterns from the proto message and any extra patterns
	var p ProtoType
	patternsDetails, err := getPatternsDetails(p, config.extraPatterns)
	if err != nil {
		return nil, err
	}
	if len(patternsDetails) == 0 {
		return nil, fmt.Errorf("no resource pattern found in %T", p)
	}

	return &defaultReflectNamer{
		patternsDetails: patternsDetails,
	}, nil
}

// getTypeCacheEntry gets or creates a type cache entry for the given type
func (n *defaultReflectNamer) getTypeCacheEntry(typ reflect.Type) (*typeCacheEntry, error) {
	// Handle pointer types
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Check if we already have a cache entry
	if entry, ok := n.typeCache.Load(typ); ok {
		return entry.(*typeCacheEntry), nil
	}

	// Create new cache entry
	patternKeyDetails, err := extractStructKeyDetails(typ)
	if err != nil {
		return nil, err
	}

	entry := &typeCacheEntry{
		patternKeyDetails: patternKeyDetails,
	}

	// Store in cache
	n.typeCache.Store(typ, entry)
	return entry, nil
}

// extractStructKeyDetails uses reflection to map struct fields to pattern keys for the given type.
// Returns a map of pattern key to patternKeyDetails.
func extractStructKeyDetails(typ reflect.Type) (map[string]patternKeyDetails, error) {
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type must be a struct")
	}
	patternVarMap := make(map[string]patternKeyDetails)
	return extractStructKeyDetailsRec(typ, patternVarMap, nil)
}

// extractStructKeyDetailsRec recursively traverses struct fields to build a map of pattern keys to their details.
func extractStructKeyDetailsRec(typ reflect.Type, patternVarMap map[string]patternKeyDetails, parentIndex []int) (map[string]patternKeyDetails, error) {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("aip_pattern")
		currentIndex := append(parentIndex, i)
		var patternKeys []string
		if tag != "" {
			parts := strings.Split(tag, "=")
			if len(parts) == 2 && parts[0] == "key" {
				patternKeys = strings.Split(parts[1], ",")
			}
		}
		if len(patternKeys) == 0 {
			patternKeys = []string{SnakeCase(field.Name)}
		}

		// Recurse into struct or pointer-to-struct fields
		if field.Type.Kind() == reflect.Struct || (field.Type.Kind() == reflect.Ptr && field.Type.Elem().Kind() == reflect.Struct) {
			if field.Type.Kind() == reflect.Ptr {
				field.Type = field.Type.Elem()
			}
			var err error
			patternVarMap, err = extractStructKeyDetailsRec(field.Type, patternVarMap, currentIndex)
			if err != nil {
				return nil, err
			}
		} else {
			kind := field.Type.Kind()
			if kind == reflect.Ptr {
				kind = field.Type.Elem().Kind()
			}
			switch kind {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
				reflect.Float32, reflect.Float64,
				reflect.String:
			default:
				continue
			}
			for _, patternKey := range patternKeys {
				patternVarMap[patternKey] = patternKeyDetails{
					patternKey:     patternKey,
					fieldIndexPath: currentIndex,
					fieldType:      field.Type,
				}
			}
		}
	}
	return patternVarMap, nil
}

// getPatternsDetails parses the proto resource descriptor and returns all resource name patterns,
// including any extra patterns provided. Patterns are indexed by their order (extra patterns start at 100).
func getPatternsDetails(resource proto.Message, extraPatterns []string) (map[int]patternDetails, error) {
	resourceOption := proto.GetExtension(
		resource.ProtoReflect().Descriptor().Options().(*descriptorpb.MessageOptions),
		annotations.E_Resource,
	).(*annotations.ResourceDescriptor)

	// Start the extra patterns at index 100
	patterns := make([]string, 100)
	copy(patterns, resourceOption.Pattern)
	for _, pattern := range extraPatterns {
		if err := resourcename.Validate(pattern); err != nil {
			return nil, fmt.Errorf("extra pattern %s is invalid: %w", pattern, err)
		}
	}
	patterns = append(patterns, extraPatterns...)

	patternsDetails := make(map[int]patternDetails)
	for i, pattern := range patterns {
		if pattern == "" {
			continue
		}

		splitPattern := strings.Split(pattern, "/")

		splitParentPattern := []string{}
		// If there are more than 3 sections, derive the parent pattern
		if len(splitPattern) > 3 {
			// Get the parent pattern sections as if it's a singleton pattern
			splitParentPattern = splitPattern[:len(splitPattern)-1]
			// If it now has an odd number of sections, it was a collection name pattern; remove another section
			if len(splitParentPattern)%2 == 1 {
				splitParentPattern = splitParentPattern[:len(splitParentPattern)-1]
			}
		}

		// Extract pattern keys from the pattern
		patternKeys := extractPatternKeys(pattern)
		parentPatternKeys := extractPatternKeys(strings.Join(splitParentPattern, "/"))

		patternsDetails[i] = patternDetails{
			patternIndex:       i,
			pattern:            pattern,
			splitPattern:       splitPattern,
			patternKeys:        patternKeys,
			parentPattern:      strings.Join(splitParentPattern, "/"),
			splitParentPattern: splitParentPattern,
			parentPatternKeys:  parentPatternKeys,
		}
	}
	return patternsDetails, nil
}

// extractPatternKeys extracts the variable keys from a pattern string.
// For example, "users/{user}/posts/{post}" returns ["user", "post"].
func extractPatternKeys(pattern string) []string {
	var keys []string
	parts := strings.Split(pattern, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			key := part[1 : len(part)-1]
			keys = append(keys, key)
		}
	}
	return keys
}

// SnakeCase converts a string to snake_case, used for mapping struct field names to pattern keys.
func SnakeCase(s string) string {
	var out []rune
	for i, r := range s {
		if r == ' ' || r == '-' {
			out = append(out, '_')
			continue
		}
		if i > 0 && r >= 'A' && r <= 'Z' && (s[i-1] >= 'a' && s[i-1] <= 'z') {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(r))
	}
	return string(out)
}
