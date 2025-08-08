# FieldMask Package

The fieldmask package provides functionality for mapping field names to database columns with table information and filtering field lists based on table sources.

## Features

- Map field names to database columns with table information
- Filter fields by specific tables
- Include or exclude fields from multiple tables
- Type-safe table-based filtering

## Basic Usage

```go
// Create a field masker with table information
mapping := map[string][]fieldmask.Field{
    "title": {{Name: "title", Table: "calendar"}},
    "description": {{Name: "description", Table: "calendar"}},
    "access": {
        {Name: "access_id", Table: "calendar_access"},
        {Name: "permission_level", Table: "calendar_access"},
        {Name: "state", Table: "calendar_access"},
    },
}

fmm := fieldmask.NewFieldMasker(mapping)

// Get all fields
allFields := fmm.Get()

// Convert specific fields
fields := fmm.Convert([]string{"title", "access"})
```

## Table-Based Filtering

The new Field-based approach allows you to filter fields based on which table they come from:

### Exclude Fields from Specific Tables

```go
// Exclude fields from a single table
fieldsWithoutAccess := fmm.Get(fieldmask.ExcludeTable(&fmm, "calendar_access"))

// Exclude fields from multiple tables
fieldsWithoutAccessTables := fmm.Get(fieldmask.ExcludeTables(&fmm, "calendar_access", "user_access"))
```

### Include Only Fields from Specific Tables

```go
// Include only fields from a single table
calendarFields := fmm.Get(fieldmask.IncludeTable(&fmm, "calendar"))

// Include only fields from multiple tables
mainTableFields := fmm.Get(fieldmask.IncludeTables(&fmm, "calendar", "user"))
```

### Combine with Field Conversion

```go
// Convert specific fields while excluding certain tables
userFieldsWithoutAccess := fmm.Convert([]string{"user"}, 
    fieldmask.ExcludeTables(&fmm, "user_access"))
```

## Available Options

- `ExcludeKeys(keys ...string)`: Exclude specific field keys
- `ExcludeValues(valuesToExclude ...string)`: Exclude specific column values
- `ExcludeTable(fmm *FieldMasker, table string)`: Exclude fields from a specific table
- `ExcludeTables(fmm *FieldMasker, tables ...string)`: Exclude fields from multiple tables
- `IncludeTable(fmm *FieldMasker, table string)`: Include only fields from a specific table
- `IncludeTables(fmm *FieldMasker, tables ...string)`: Include only fields from multiple tables

## Example with Calendar Model

```go
var CalendarFieldMasker = fieldmask.NewFieldMasker(map[string][]fieldmask.Field{
    cmodel.CalendarField_Title: {{Name: CalendarColumn_Title, Table: "calendar"}},
    cmodel.CalendarField_Description: {{Name: CalendarColumn_Description, Table: "calendar"}},
    cmodel.CalendarField_CalendarAccess: {
        {Name: CalendarAccessColumn_CalendarAccessId, Table: "calendar_access"},
        {Name: CalendarAccessColumn_PermissionLevel, Table: "calendar_access"},
        {Name: CalendarAccessColumn_State, Table: "calendar_access"},
    },
})

// Get only main table fields (exclude join table fields)
mainFields := CalendarFieldMasker.Get(fieldmask.ExcludeTable(&CalendarFieldMasker, "calendar_access"))

// Get only access table fields
accessFields := CalendarFieldMasker.Get(fieldmask.IncludeTable(&CalendarFieldMasker, "calendar_access"))

// Get fields from multiple tables
mainAndAccessFields := CalendarFieldMasker.Get(fieldmask.IncludeTables(&CalendarFieldMasker, "calendar", "calendar_access"))
```

## Benefits

1. **Explicit Table Information**: Each field knows which table it belongs to
2. **Flexible Filtering**: Filter by specific tables, not just "join" vs "main"
3. **Multiple Table Support**: A field can have columns from multiple tables
4. **Type Safety**: Compile-time checking of table names
5. **Granular Control**: Include or exclude specific tables as needed
