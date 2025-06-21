package masks

// NewFieldMap will return a new FieldMap.
func NewFieldMap() FieldMap {
	return make(FieldMap)
}

// FieldMap is a map of fields to a list of paths.
type FieldMap map[string][]string

// MapFieldToFields will map the provided field to the provided fields.
func (m FieldMap) MapFieldToFields(field string, fields ...string) FieldMap {
	if m == nil {
		m = NewFieldMap()
	}
	m[field] = fields
	return m
}

// ToStringMap will convert the FieldMap to a map[string]string.
// It will take the first field in the list of fields.
func (m FieldMap) ToStringMap() map[string]string {
	stringMap := make(map[string]string)
	for key, values := range m {
		if len(values) > 0 {
			stringMap[key] = values[0]
		}
	}
	return stringMap
}
