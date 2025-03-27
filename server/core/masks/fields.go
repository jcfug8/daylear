package masks

import "slices"

// Equal will return true if the two FieldMasks are equal.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for _, path := range a {
		if !slices.Contains(b, path) {
			return false
		}
	}

	return true
}

// Intersection will return a new FieldMask that is the intersection of the two
// FieldMasks.
func Intersection(a, b []string) []string {
	if len(a) == 0 || len(b) == 0 {
		return nil
	}

	max := min(len(a), len(b))
	m := make(map[string]struct{}, max)
	paths := make([]string, 0, max)

	for _, path := range a {
		if _, ok := m[path]; ok || !slices.Contains(b, path) {
			continue
		}

		paths = append(paths, path)
		m[path] = struct{}{}
	}

	return paths
}

// Map will return a new FieldMask populated from the values of the provided
// FieldMap using the provived FieldMask as the keys. Meaning, if the mappings
// contain a path that is present in the provided FieldMask, the list of fields
// for that path (from the mappings) will be included in the new FieldMask.
func Map(mask []string, mappings FieldMap) []string {
	if len(mask) == 0 || len(mappings) == 0 {
		return nil
	}

	m := make(map[string]struct{}, len(mask))
	paths := make([]string, 0, len(mask))
	for path, fields := range mappings {
		if !slices.Contains(mask, path) {
			continue
		}

		for _, field := range fields {
			if _, ok := m[field]; ok {
				continue
			}

			paths = append(paths, field)
			m[field] = struct{}{}
		}
	}

	return paths
}

// Prefix will return a new FieldMask where all paths are prefixed with the
// provided prefix.
func Prefix(prefix string, paths []string) []string {
	if prefix == "" || len(paths) == 0 {
		return paths
	}

	prefixed := make([]string, len(paths))
	for i, path := range paths {
		prefixed[i] = prefix + path
	}

	return prefixed
}

// RemovePaths will return a new FieldMask with the provided paths removed.
func RemovePaths(mask []string, paths ...string) []string {
	if len(mask) == 0 {
		return nil
	}

	m := make(map[string]struct{}, len(mask))
	for _, path := range mask {
		m[path] = struct{}{}
	}

	for _, path := range paths {
		delete(m, path)
	}

	paths = make([]string, 0, len(m))
	for path := range m {
		paths = append(paths, path)
	}

	return paths
}
