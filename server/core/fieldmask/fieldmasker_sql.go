package fieldmask

import (
	"fmt"
	"reflect"
	"sync"

	"gorm.io/gorm/schema"
)

// NewSQLFieldMasker initializes a FieldMasker and ensures all field values are valid column names for the given table struct.
// It uses GORM's schema reflection to extract column names from the struct and validate the field mapping.
func NewSQLFieldMasker(tableStruct any, mapping map[string][]Field) FieldMasker {
	// Extract column names from the struct using GORM's schema reflection
	validColumns := extractColumnNames(tableStruct)

	// Validate that all field names in the mapping correspond to actual column names
	for key, fields := range mapping {
		for _, field := range fields {
			if !validColumns[field.Name] && !validColumns[field.Alias] {
				// Log or handle invalid column names - for now, we'll just skip them
				// In a production environment, you might want to return an error
				panic(fmt.Sprintf("Warning: field '%s' in key '%s' is not a valid column in the struct\n", field.Name, key))
			}
		}
	}

	return FieldMasker{mapping: mapping}
}

// extractColumnNames uses GORM's schema reflection to extract column names from a struct
func extractColumnNames(tableStruct any) map[string]bool {
	validColumns := make(map[string]bool)

	if tableStruct == nil {
		return validColumns
	}

	// Get the reflect.Type of the struct
	var typ reflect.Type
	switch v := tableStruct.(type) {
	case reflect.Type:
		typ = v
	default:
		typ = reflect.TypeOf(tableStruct)
	}

	// Handle pointer types
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Only process struct types
	if typ.Kind() != reflect.Struct {
		return validColumns
	}

	// Use GORM's schema parser to extract column information
	sch, err := schema.Parse(tableStruct, &sync.Map{}, &schema.NamingStrategy{})
	if err != nil {
		// If parsing fails, return empty map
		return validColumns
	}

	// Extract column names from the parsed schema
	for _, field := range sch.Fields {
		validColumns[field.DBName] = true
	}

	return validColumns
}
