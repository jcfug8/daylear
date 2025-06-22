package filter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestConverter() *SQLConverter {
	fieldMapping := map[string]string{
		"user.profile.name":         "user_profile_name",
		"user.profile.email":        "user_profile_email",
		"user.profile.created_at":   "user_profile_created_at",
		"user.profile.hex_value":    "user_profile_hex_value",
		"user.profile.is_active":    "user_profile_is_active",
		"user.profile.score":        "user_profile_score",
		"user.profile.age":          "user_profile_age",
		"user.profile.address.city": "user_profile_address_city",
		"name":                      "user_name",
		"age":                       "user_age",
		"email":                     "user_email",
		"created_at":                "created_at",
		"is_active":                 "is_active",
		"score":                     "user_score",
		"hex_value":                 "hex_value",
	}
	return NewSQLConverter(fieldMapping)
}

func TestSQLConverter_SimpleEquals(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name = 'John'")
	assert.NoError(t, err)
	assert.Equal(t, "user_name = $1", got.WhereClause)
	assert.Equal(t, []interface{}{"John"}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1}, got.UsedColumns)
}

func TestSQLConverter_GreaterThan(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("age > 18")
	assert.NoError(t, err)
	assert.Equal(t, "user_age > $1", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18)}, got.Params)
	assert.Equal(t, map[string]int{"user_age": 1}, got.UsedColumns)
}

func TestSQLConverter_Contains(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("contains(email, '@gmail.com')")
	assert.NoError(t, err)
	assert.Equal(t, "user_email LIKE '%' || $1 || '%'", got.WhereClause)
	assert.Equal(t, []interface{}{"@gmail.com"}, got.Params)
	assert.Equal(t, map[string]int{"user_email": 1}, got.UsedColumns)
}

func TestSQLConverter_LogicalAnd(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("age >= 18 AND is_active = true")
	assert.NoError(t, err)
	assert.Equal(t, "user_age >= $1 AND is_active = $2", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18), true}, got.Params)
	assert.Equal(t, map[string]int{"user_age": 1, "is_active": 1}, got.UsedColumns)
}

func TestSQLConverter_LogicalOr(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name = 'John' OR name = 'Jane'")
	assert.NoError(t, err)
	assert.Equal(t, "user_name = $1 OR user_name = $2", got.WhereClause)
	assert.Equal(t, []interface{}{"John", "Jane"}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 2}, got.UsedColumns)
}

func TestSQLConverter_NotOperator(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("NOT (age < 18)")
	assert.NoError(t, err)
	assert.Equal(t, "user_age >= $1", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18)}, got.Params)
	assert.Equal(t, map[string]int{"user_age": 1}, got.UsedColumns)
}

func TestSQLConverter_ComplexExpression(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("age >= 18 AND (name = 'John' OR contains(email, '@gmail.com'))")
	assert.NoError(t, err)
	assert.Equal(t, "user_age >= $1 AND (user_name = $2 OR user_email LIKE '%' || $3 || '%')", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18), "John", "@gmail.com"}, got.Params)
	assert.Equal(t, map[string]int{"user_age": 1, "user_name": 1, "user_email": 1}, got.UsedColumns)
}

func TestSQLConverter_InOperator(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name = 'John' OR name = 'Jane'")
	assert.NoError(t, err)
	assert.Equal(t, "user_name = $1 OR user_name = $2", got.WhereClause)
	assert.Equal(t, []interface{}{"John", "Jane"}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 2}, got.UsedColumns)
}

func TestSQLConverter_HasOperator(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name:*")
	assert.NoError(t, err)
	assert.Equal(t, "user_name IS NOT NULL", got.WhereClause)
	assert.Equal(t, []interface{}{}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1}, got.UsedColumns)
}

func TestSQLConverter_FloatValue(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("score > 95.5")
	assert.NoError(t, err)
	assert.Equal(t, "user_score > $1", got.WhereClause)
	assert.Equal(t, []interface{}{float64(95.5)}, got.Params)
	assert.Equal(t, map[string]int{"user_score": 1}, got.UsedColumns)
}

func TestSQLConverter_HexValue(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("hex_value = 0xFF")
	assert.NoError(t, err)
	assert.Equal(t, "hex_value = $1", got.WhereClause)
	assert.Equal(t, []interface{}{int64(0xFF)}, got.Params)
	assert.Equal(t, map[string]int{"hex_value": 1}, got.UsedColumns)
}

func TestSQLConverter_BooleanValue(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("is_active = true")
	assert.NoError(t, err)
	assert.Equal(t, "is_active = $1", got.WhereClause)
	assert.Equal(t, []interface{}{true}, got.Params)
	assert.Equal(t, map[string]int{"is_active": 1}, got.UsedColumns)
}

// durations don't work as expected because the einride/aip-go parser doesn't support them
// func TestSQLConverter_DurationValue(t *testing.T) {
// 	converter := newTestConverter()
// 	got, err := converter.Convert("duration > 20s")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "duration > $1", got.WhereClause)
// 	assert.Equal(t, []interface{}{"20s"}, got.Params)
// }

func TestSQLConverter_TimestampValue(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("created_at > '2012-04-21T11:30:00-04:00'")
	assert.NoError(t, err)
	assert.Equal(t, "created_at > $1", got.WhereClause)
	assert.Equal(t, []interface{}{"2012-04-21T11:30:00-04:00"}, got.Params)
	assert.Equal(t, map[string]int{"created_at": 1}, got.UsedColumns)
}

func TestSQLConverter_WildcardString(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("email = '*.gmail.com'")
	assert.NoError(t, err)
	assert.Equal(t, "user_email LIKE $1", got.WhereClause)
	assert.Equal(t, []interface{}{"%.gmail.com"}, got.Params)
	assert.Equal(t, map[string]int{"user_email": 1}, got.UsedColumns)
}

func TestSQLConverter_NestedField(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("user.profile.name = 'John'")
	assert.NoError(t, err)
	assert.Equal(t, "user_profile_name = $1", got.WhereClause)
	assert.Equal(t, []interface{}{"John"}, got.Params)
	assert.Equal(t, map[string]int{"user_profile_name": 1}, got.UsedColumns)
}

func TestSQLConverter_MultipleLogicalOps(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("age >= 18 AND (name = 'John' OR email = '*.gmail.com') AND is_active = true")
	assert.NoError(t, err)
	assert.Equal(t, "user_age >= $1 AND (user_name = $2 OR user_email LIKE $3) AND is_active = $4", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18), "John", "%.gmail.com", true}, got.Params)
	assert.Equal(t, map[string]int{"user_age": 1, "user_name": 1, "user_email": 1, "is_active": 1}, got.UsedColumns)
}

func TestSQLConverter_NegationWithMinus(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("-age < 18")
	assert.NoError(t, err)
	assert.Equal(t, "user_age >= $1", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18)}, got.Params)
	assert.Equal(t, map[string]int{"user_age": 1}, got.UsedColumns)
}

func TestSQLConverter_ComplexNestedExpression(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("(age >= 18 AND is_active = true) OR (name = 'John' AND email = '*.gmail.com')")
	assert.NoError(t, err)
	assert.Equal(t, "(user_age >= $1 AND is_active = $2) OR (user_name = $3 AND user_email LIKE $4)", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18), true, "John", "%.gmail.com"}, got.Params)
	assert.Equal(t, map[string]int{"user_age": 1, "is_active": 1, "user_name": 1, "user_email": 1}, got.UsedColumns)
}

func TestSQLConverter_StartsWithAndEndsWith(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("starts_with(email, 'test') AND ends_with(email, '.com')")
	assert.NoError(t, err)
	assert.Equal(t, "user_email LIKE $1 || '%' AND user_email LIKE '%' || $2", got.WhereClause)
	assert.Equal(t, []interface{}{"test", ".com"}, got.Params)
	assert.Equal(t, map[string]int{"user_email": 2}, got.UsedColumns)
}

func TestSQLConverter_MultipleHasOperators(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name:* AND email:*")
	assert.NoError(t, err)
	assert.Equal(t, "user_name IS NOT NULL AND user_email IS NOT NULL", got.WhereClause)
	assert.Equal(t, []interface{}{}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1, "user_email": 1}, got.UsedColumns)
}

func TestSQLConverter_NotEqualsWithNull(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name != null")
	assert.NoError(t, err)
	assert.Equal(t, "user_name IS NOT NULL", got.WhereClause)
	assert.Equal(t, []interface{}{}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1}, got.UsedColumns)
}

func TestSQLConverter_EqualsWithNull(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name = null")
	assert.NoError(t, err)
	assert.Equal(t, "user_name IS NULL", got.WhereClause)
	assert.Equal(t, []interface{}{}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1}, got.UsedColumns)
}

func TestSQLConverter_MultipleWildcards(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("email = 'test.*.com'")
	assert.NoError(t, err)
	assert.Equal(t, "user_email LIKE $1", got.WhereClause)
	assert.Equal(t, []interface{}{"test.%.com"}, got.Params)
	assert.Equal(t, map[string]int{"user_email": 1}, got.UsedColumns)
}

func TestSQLConverter_ComplexNestedField(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("user.profile.address.city = 'New York'")
	assert.NoError(t, err)
	assert.Equal(t, "user_profile_address_city = $1", got.WhereClause)
	assert.Equal(t, []interface{}{"New York"}, got.Params)
	assert.Equal(t, map[string]int{"user_profile_address_city": 1}, got.UsedColumns)
}

func TestSQLConverter_MultipleContains(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("contains(name, 'John') AND contains(email, '@gmail.com')")
	assert.NoError(t, err)
	assert.Equal(t, "user_name LIKE '%' || $1 || '%' AND user_email LIKE '%' || $2 || '%'", got.WhereClause)
	assert.Equal(t, []interface{}{"John", "@gmail.com"}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1, "user_email": 1}, got.UsedColumns)
}

func TestSQLConverter_ComplexLogicalWithNull(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("(name = null OR email = '*.gmail.com') AND is_active = true")
	assert.NoError(t, err)
	assert.Equal(t, "(user_name IS NULL OR user_email LIKE $1) AND is_active = $2", got.WhereClause)
	assert.Equal(t, []interface{}{"%.gmail.com", true}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1, "user_email": 1, "is_active": 1}, got.UsedColumns)
}

func TestSQLConverter_MultipleStartsWithEndsWith(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("starts_with(name, 'John') AND ends_with(email, '.com') AND contains(hex_value, 'FF')")
	assert.NoError(t, err)
	assert.Equal(t, "user_name LIKE $1 || '%' AND user_email LIKE '%' || $2 AND hex_value LIKE '%' || $3 || '%'", got.WhereClause)
	assert.Equal(t, []interface{}{"John", ".com", "FF"}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1, "user_email": 1, "hex_value": 1}, got.UsedColumns)
}

func TestSQLConverter_ComplexHasOperators(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name:* AND email:* AND (age > 18 OR is_active = true)")
	assert.NoError(t, err)
	assert.Equal(t, "user_name IS NOT NULL AND user_email IS NOT NULL AND (user_age > $1 OR is_active = $2)", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18), true}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1, "user_email": 1, "user_age": 1, "is_active": 1}, got.UsedColumns)
}

func TestSQLConverter_ComplexNotExpressions(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("NOT (name = 'John' AND age > 18) OR NOT (email = '*.gmail.com')")
	assert.NoError(t, err)
	assert.Equal(t, "NOT (user_name = $1 AND user_age > $2) OR NOT (user_email LIKE $3)", got.WhereClause)
	assert.Equal(t, []interface{}{"John", int64(18), "%.gmail.com"}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1, "user_age": 1, "user_email": 1}, got.UsedColumns)
}

func TestSQLConverter_ComplexNestedLogicalWithNull(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("(name = null OR (age > 18 AND is_active = true)) AND (email = '*.gmail.com' OR score > 95.5)")
	assert.NoError(t, err)
	assert.Equal(t, "(user_name IS NULL OR (user_age > $1 AND is_active = $2)) AND (user_email LIKE $3 OR user_score > $4)", got.WhereClause)
	assert.Equal(t, []interface{}{int64(18), true, "%.gmail.com", float64(95.5)}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1, "user_age": 1, "is_active": 1, "user_email": 1, "user_score": 1}, got.UsedColumns)
}

func TestSQLConverter_EmptyFilter(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("")
	assert.NoError(t, err)
	assert.Equal(t, "", got.WhereClause)
	assert.Equal(t, []interface{}{}, got.Params)
	assert.Equal(t, map[string]int{}, got.UsedColumns)
}

func TestSQLConverter_InvalidFilter(t *testing.T) {
	converter := newTestConverter()
	_, err := converter.Convert("invalid filter")
	assert.Error(t, err)
}

func TestSQLConverter_InvalidOperator(t *testing.T) {
	converter := newTestConverter()
	_, err := converter.Convert("name INVALID 'John'")
	assert.Error(t, err)
}

func TestSQLConverter_InvalidFunction(t *testing.T) {
	converter := newTestConverter()
	_, err := converter.Convert("invalid_function(name, 'John')")
	assert.Error(t, err)
}

func TestSQLConverter_InvalidField(t *testing.T) {
	converter := newTestConverter()
	_, err := converter.Convert("invalid_field = 'John'")
	assert.NoError(t, err) // Should not error, just use the field name as is
}

func TestSQLConverter_InvalidValue(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name = invalid_value")
	assert.NoError(t, err) // Should not error, just use the value as is
	assert.Equal(t, "user_name = $1", got.WhereClause)
	assert.Equal(t, []interface{}{"invalid_value"}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1}, got.UsedColumns)
}

func TestSQLConverter_InvalidLogicalOperator(t *testing.T) {
	converter := newTestConverter()
	_, err := converter.Convert("name = 'John' INVALID age > 18")
	assert.Error(t, err)
}

func TestSQLConverter_InvalidNotOperator(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("NOT name = 'John' AND age > 18")
	assert.NoError(t, err)
	assert.Equal(t, "user_name != $1 AND user_age > $2", got.WhereClause)
	assert.Equal(t, []interface{}{"John", int64(18)}, got.Params)
	assert.Equal(t, map[string]int{"user_name": 1, "user_age": 1}, got.UsedColumns)
}

func TestSQLConverter_InvalidFunctionCalls(t *testing.T) {
	tests := []struct {
		name    string
		filter  string
		wantErr bool
	}{
		{
			name:    "Invalid contains",
			filter:  "contains(name)",
			wantErr: true,
		},
		{
			name:    "Invalid starts_with",
			filter:  "starts_with(name)",
			wantErr: true,
		},
		{
			name:    "Invalid ends_with",
			filter:  "ends_with(name)",
			wantErr: true,
		},
		{
			name:    "Invalid has",
			filter:  "name:",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := newTestConverter()
			_, err := converter.Convert(tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSQLConverter_NestedFieldOperations(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		wantSQL  string
		wantArgs []interface{}
		wantErr  bool
		wantUsed map[string]int
	}{
		{
			name:    "Invalid nested field path",
			filter:  "user..profile.name = 'John'",
			wantErr: true,
		},
		{
			name:    "Invalid nested field operator",
			filter:  "user.profile.name INVALID 'John'",
			wantErr: true,
		},
		{
			name:    "Invalid nested field function",
			filter:  "invalid_function(user.profile.name, 'John')",
			wantErr: true,
		},
		{
			name:    "Invalid nested field logical operator",
			filter:  "user.profile.name = 'John' INVALID user.profile.age > 18",
			wantErr: true,
		},
		{
			name:     "Invalid nested field not operator",
			filter:   "NOT user.profile.name = 'John' AND user.profile.age > 18",
			wantSQL:  "user_profile_name != $1 AND user_profile_age > $2",
			wantArgs: []interface{}{"John", int64(18)},
			wantUsed: map[string]int{"user_profile_name": 1, "user_profile_age": 1},
		},
		{
			name:    "Invalid nested field contains operator",
			filter:  "contains(user.profile.name)",
			wantErr: true,
		},
		{
			name:    "Invalid nested field starts_with operator",
			filter:  "starts_with(user.profile.name)",
			wantErr: true,
		},
		{
			name:    "Invalid nested field ends_with operator",
			filter:  "ends_with(user.profile.name)",
			wantErr: true,
		},
		{
			name:    "Invalid nested field has operator",
			filter:  "user.profile.name:",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := newTestConverter()
			got, err := converter.Convert(tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantSQL, got.WhereClause)
				assert.Equal(t, tt.wantArgs, got.Params)
				if tt.wantUsed != nil {
					assert.Equal(t, tt.wantUsed, got.UsedColumns)
				}
			}
		})
	}
}

func TestSQLConverter_NestedFieldValues(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		wantSQL  string
		wantArgs []interface{}
		wantUsed map[string]int
	}{
		{
			name:     "Null value",
			filter:   "user.profile.name = null",
			wantSQL:  "user_profile_name IS NULL",
			wantArgs: []interface{}{},
			wantUsed: map[string]int{"user_profile_name": 1},
		},
		{
			name:     "Not null value",
			filter:   "user.profile.name != null",
			wantSQL:  "user_profile_name IS NOT NULL",
			wantArgs: []interface{}{},
			wantUsed: map[string]int{"user_profile_name": 1},
		},
		{
			name:     "Boolean value",
			filter:   "user.profile.is_active = true",
			wantSQL:  "user_profile_is_active = $1",
			wantArgs: []interface{}{true},
			wantUsed: map[string]int{"user_profile_is_active": 1},
		},
		{
			name:     "Numeric value",
			filter:   "user.profile.age > 18",
			wantSQL:  "user_profile_age > $1",
			wantArgs: []interface{}{int64(18)},
			wantUsed: map[string]int{"user_profile_age": 1},
		},
		{
			name:     "Float value",
			filter:   "user.profile.score > 95.5",
			wantSQL:  "user_profile_score > $1",
			wantArgs: []interface{}{float64(95.5)},
			wantUsed: map[string]int{"user_profile_score": 1},
		},
		{
			name:     "Hex value",
			filter:   "user.profile.hex_value = 0xFF",
			wantSQL:  "user_profile_hex_value = $1",
			wantArgs: []interface{}{int64(0xFF)},
			wantUsed: map[string]int{"user_profile_hex_value": 1},
		},
		{
			name:     "Timestamp value",
			filter:   "user.profile.created_at > '2012-04-21T11:30:00-04:00'",
			wantSQL:  "user_profile_created_at > $1",
			wantArgs: []interface{}{"2012-04-21T11:30:00-04:00"},
			wantUsed: map[string]int{"user_profile_created_at": 1},
		},
		{
			name:     "Wildcard value",
			filter:   "user.profile.email = '*.gmail.com'",
			wantSQL:  "user_profile_email LIKE $1",
			wantArgs: []interface{}{"%.gmail.com"},
			wantUsed: map[string]int{"user_profile_email": 1},
		},
		{
			name:     "Multiple wildcards",
			filter:   "user.profile.email = 'test.*.com'",
			wantSQL:  "user_profile_email LIKE $1",
			wantArgs: []interface{}{"test.%.com"},
			wantUsed: map[string]int{"user_profile_email": 1},
		},
		{
			name:     "Complex nested field",
			filter:   "user.profile.address.city = 'New York'",
			wantSQL:  "user_profile_address_city = $1",
			wantArgs: []interface{}{"New York"},
			wantUsed: map[string]int{"user_profile_address_city": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := newTestConverter()
			got, err := converter.Convert(tt.filter)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantSQL, got.WhereClause)
			assert.Equal(t, tt.wantArgs, got.Params)
			assert.Equal(t, tt.wantUsed, got.UsedColumns)
		})
	}
}

func TestSQLConverter_ComplexNestedExpressions(t *testing.T) {
	tests := []struct {
		name     string
		filter   string
		wantSQL  string
		wantArgs []interface{}
		wantUsed map[string]int
	}{
		{
			name:     "Multiple contains",
			filter:   "contains(user.profile.name, 'John') AND contains(user.profile.email, '@gmail.com')",
			wantSQL:  "user_profile_name LIKE '%' || $1 || '%' AND user_profile_email LIKE '%' || $2 || '%'",
			wantArgs: []interface{}{"John", "@gmail.com"},
			wantUsed: map[string]int{
				"user_profile_name":  1,
				"user_profile_email": 1,
			},
		},
		{
			name:     "Complex logical with null",
			filter:   "(user.profile.name = null OR (user.profile.age > 18 AND user.profile.is_active = true)) AND (user.profile.email = '*.gmail.com' OR user.profile.score > 95.5)",
			wantSQL:  "(user_profile_name IS NULL OR (user_profile_age > $1 AND user_profile_is_active = $2)) AND (user_profile_email LIKE $3 OR user_profile_score > $4)",
			wantArgs: []interface{}{int64(18), true, "%.gmail.com", float64(95.5)},
			wantUsed: map[string]int{
				"user_profile_name":      1,
				"user_profile_age":       1,
				"user_profile_is_active": 1,
				"user_profile_email":     1,
				"user_profile_score":     1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			converter := newTestConverter()
			got, err := converter.Convert(tt.filter)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantSQL, got.WhereClause)
			assert.Equal(t, tt.wantArgs, got.Params)
			assert.Equal(t, tt.wantUsed, got.UsedColumns)
		})
	}
}

func TestSQLConverter_UsedColumns(t *testing.T) {
	converter := newTestConverter()
	got, err := converter.Convert("name = 'John' AND user.profile.name = 'Test'")
	assert.NoError(t, err)
	assert.Equal(t, "user_name = $1 AND user_profile_name = $2", got.WhereClause)
	assert.Equal(t, []interface{}{"John", "Test"}, got.Params)
	expectedUsedColumns := map[string]int{
		"user_name":         1,
		"user_profile_name": 1,
	}
	assert.Equal(t, expectedUsedColumns, got.UsedColumns)
}
