package fieldmask

import "slices"

type option func(key string, values []string) []string

func ExcludeKeys(keys ...string) option {
	return func(key string, values []string) []string {
		if slices.Contains(keys, key) {
			return []string{}
		}
		return values
	}
}

func ExcludeValues(valuesToExclude ...string) option {
	return func(key string, values []string) []string {
		return slices.DeleteFunc(values, func(value string) bool {
			return slices.Contains(valuesToExclude, value)
		})
	}
}
