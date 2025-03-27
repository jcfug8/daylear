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
