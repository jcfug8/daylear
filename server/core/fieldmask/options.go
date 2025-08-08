package fieldmask

import "slices"

type option func(key string, values []Field) []Field

// filterFields returns a new slice containing only the fields that match the condition
func filterFields(fields []Field, condition func(Field) bool) []Field {
	result := make([]Field, 0, len(fields))
	for _, field := range fields {
		if condition(field) {
			result = append(result, field)
		}
	}
	return result
}

func ExcludeKeys(keys ...string) option {
	return func(key string, values []Field) []Field {
		if slices.Contains(keys, key) {
			return []Field{}
		}
		return values
	}
}

func ExcludeValues(valuesToExclude ...string) option {
	return func(key string, values []Field) []Field {
		return filterFields(values, func(value Field) bool {
			return !slices.Contains(valuesToExclude, value.Name) && !slices.Contains(valuesToExclude, value.Alias)
		})
	}
}

// ExcludeTables returns an option that excludes fields from the specified tables.
func ExcludeTables(tables ...string) option {
	return func(key string, fields []Field) []Field {
		return filterFields(fields, func(value Field) bool {
			return !slices.Contains(tables, value.Table)
		})
	}
}

// IncludeTables returns an option that includes only fields from the specified tables.
func IncludeTables(tables ...string) option {
	return func(key string, fields []Field) []Field {
		return filterFields(fields, func(value Field) bool {
			return slices.Contains(tables, value.Table)
		})
	}
}

// OnlyUpdatable returns an option that includes only updatable fields.
func OnlyUpdatable() option {
	return func(key string, fields []Field) []Field {
		return filterFields(fields, func(value Field) bool {
			return value.Updatable
		})
	}
}
