package fieldmask

import (
	"reflect"
	"sort"
	"testing"
)

func TestFieldMaskerWithTableFiltering(t *testing.T) {
	// Create a field masker with table information
	mapping := map[string][]Field{
		"title":       {{Name: "title", Table: "calendar"}},
		"description": {{Name: "description", Table: "calendar"}},
		"access": {
			{Name: "access_id", Table: "calendar_access"},
			{Name: "permission_level", Table: "calendar_access"},
			{Name: "state", Table: "calendar_access"},
		},
		"user": {
			{Name: "user_id", Table: "user"},
			{Name: "username", Table: "user"},
			{Name: "user_permission_level", Table: "user_access"},
			{Name: "user_state", Table: "user_access"},
		},
	}

	fmm := NewFieldMasker(mapping)

	// Helper function to sort results for deterministic comparison
	sortResults := func(results []string) []string {
		sorted := make([]string, len(results))
		copy(sorted, results)
		sort.Strings(sorted)
		return sorted
	}

	// Test getting all fields
	allFields := fmm.Get()
	expectedAll := []string{"title", "description", "access_id", "permission_level", "state", "user_id", "username", "user_permission_level", "user_state"}
	if !reflect.DeepEqual(sortResults(allFields), sortResults(expectedAll)) {
		t.Errorf("Expected %v, got %v", expectedAll, allFields)
	}

	// Test excluding a specific table
	fieldsWithoutCalendarAccess := fmm.Get(ExcludeTables("calendar_access"))
	expectedWithoutCalendarAccess := []string{"title", "description", "user_id", "username", "user_permission_level", "user_state"}
	if !reflect.DeepEqual(sortResults(fieldsWithoutCalendarAccess), sortResults(expectedWithoutCalendarAccess)) {
		t.Errorf("Expected %v, got %v", expectedWithoutCalendarAccess, fieldsWithoutCalendarAccess)
	}

	// Test excluding multiple tables
	fieldsWithoutAccessTables := fmm.Get(ExcludeTables("calendar_access", "user_access"))
	expectedWithoutAccessTables := []string{"title", "description", "user_id", "username"}
	if !reflect.DeepEqual(sortResults(fieldsWithoutAccessTables), sortResults(expectedWithoutAccessTables)) {
		t.Errorf("Expected %v, got %v", expectedWithoutAccessTables, fieldsWithoutAccessTables)
	}

	// Test including only a specific table
	fieldsOnlyCalendar := fmm.Get(IncludeTables("calendar"))
	expectedOnlyCalendar := []string{"title", "description"}
	if !reflect.DeepEqual(sortResults(fieldsOnlyCalendar), sortResults(expectedOnlyCalendar)) {
		t.Errorf("Expected %v, got %v", expectedOnlyCalendar, fieldsOnlyCalendar)
	}

	// Test including multiple tables
	fieldsFromMainTables := fmm.Get(IncludeTables("calendar", "user"))
	expectedFromMainTables := []string{"title", "description", "user_id", "username"}
	if !reflect.DeepEqual(sortResults(fieldsFromMainTables), sortResults(expectedFromMainTables)) {
		t.Errorf("Expected %v, got %v", expectedFromMainTables, fieldsFromMainTables)
	}

	// Test converting specific fields with table filtering
	userFields := fmm.Convert([]string{"user"})
	expectedUserFields := []string{"user_id", "username", "user_permission_level", "user_state"}
	if !reflect.DeepEqual(sortResults(userFields), sortResults(expectedUserFields)) {
		t.Errorf("Expected %v, got %v", expectedUserFields, userFields)
	}

	// Test converting specific fields excluding access tables
	userFieldsWithoutAccess := fmm.Convert([]string{"user"}, ExcludeTables("user_access"))
	expectedUserFieldsWithoutAccess := []string{"user_id", "username"}
	if !reflect.DeepEqual(sortResults(userFieldsWithoutAccess), sortResults(expectedUserFieldsWithoutAccess)) {
		t.Errorf("Expected %v, got %v", expectedUserFieldsWithoutAccess, userFieldsWithoutAccess)
	}
}
