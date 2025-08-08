package fieldmask

import "strings"

// Field represents a database column with its table information
type Field struct {
	Name      string
	Table     string
	Alias     string
	Updatable bool
}

func (f Field) String() string {
	name := f.Name
	if f.Table != "" {
		name = f.Table + "." + name
	}
	if f.Alias != "" {
		name = name + " as " + f.Alias
	}
	return name
}

// FieldMasker allows mapping keys to a set of fields and converting lists of keys to deduplicated lists of field names.
type FieldMasker struct {
	mapping map[string][]Field
}

// NewFieldMasker initializes a new FieldMasker with the provided mapping.
func NewFieldMasker(mapping map[string][]Field) FieldMasker {
	return FieldMasker{mapping: mapping}
}

// Convert takes a list of keys and returns a deduplicated list of all mapped values.
func (fmm *FieldMasker) Convert(keys []string, opts ...option) []string {
	if len(keys) == 0 {
		return fmm.Get(opts...)
	}

	result := make([]string, 0)
	seen := make(map[string]struct{})
	for _, key := range keys {
		if fields, ok := fmm.mapping[key]; ok {
			for _, opt := range opts {
				fields = opt(key, fields)
			}
			for _, field := range fields {
				if _, exists := seen[field.String()]; !exists {
					result = append(result, field.String())
					seen[field.String()] = struct{}{}
				}
			}
		}
	}
	return result
}

func (fmm *FieldMasker) Get(opts ...option) []string {
	result := make([]string, 0)
	for key, fields := range fmm.mapping {
		for _, opt := range opts {
			fields = opt(key, fields)
		}
		for _, field := range fields {
			result = append(result, field.String())
		}
	}
	return result
}

func ContainsAny(fields []string, fieldsToCheck ...string) bool {
	for _, field := range fields {
		for _, fieldToCheck := range fieldsToCheck {
			if field == fieldToCheck {
				return true
			} else {
				fieldParts := strings.Split(field, ".")
				if len(fieldParts) == 2 && strings.HasPrefix(fieldParts[1], fieldToCheck) {
					return true
				}
			}
		}
	}
	return false
}
