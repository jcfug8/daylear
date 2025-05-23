## Table of Contents
- [Purpose](#purpose)
- [Concepts](#concepts)
- [Usage](#usage)
  - [IDNamer](#1-idnamer)
  - [ParentNamer](#2-parentnamer)
  - [ParentIDNamer](#3-parentidnamer)
- [Implementation Details](#implementation-details)
- [Adding a New Resource Namer](#adding-a-new-resource-namer)
- [Example Test](#example-test)
- [Best Practices](#best-practices)
- [Questions](#questions)

# Namers

The `namer` package provides utilities for formatting and parsing resource names according to [Google AIP resource name patterns](https://google.aip.dev/122). It is designed to work with resources defined in protobufs, supporting both simple and hierarchical (parent/child) resource structures.

## Purpose

- **Reduced Boilerplate:** Reduces the amount of code needed to format and parse resource names. You should only need to create two or three functions that can convert a slice of strings to your resource parent/id and vice versa. It will handle all the pattern matching and parsing for you.
- **Consistent Resource Name Handling:** Helps give a standard way to interact with resource names.

## Concepts

- **Resource Name Patterns:** Defined in protobufs using the `google.api.resource` annotation.
- **IDNamer:** For resources with a single ID and no parent. 
  - A pattern with no parent is referenced as a "root" pattern. 
- **ParentNamer:** For singleton resources with a parent but no ID. 
  - A pattern with a parent but no id is referenced as a "singleton" pattern. 
  - This will also work for singleton resources with a "root" pattern, or no parent, but at that point do you really need a namer?
- **ParentIDNamer:** For resources with both a parent and an ID. 
  - A pattern with a parent and an id is referenced as a "standard" pattern. 
  - This will also work for resources that have both a "root" and "standard" patterns.
- **ReflectNamer:** For resources with a parent and/or id.
  - This will work for any resource pattern, but will require you to define aip-pattern struct tags on the struct passed to parse/format methods.

## Usage

### 1. IDNamer
For resources with a single ID (e.g., `rootNamedResources/{root_named_resource}`) or a "root" pattern:
- Define the ID type
- Define a function that returns the id variable for the given id and pattern index
- Define a function that sets the id variable for the given id variable and pattern index
- Instantiate the IDNamer with the resource, the id variable getter, and the id variable setter

```go
// Define your ID type
type ID struct {
    ID int64
}

func GetIDNamerVars(id ID, patternIndex int) (string, error) { ... }
func SetIDNamerVars(idVar string, patternIndex int) (ID, error) { ... }

namer, err := NewIDNamer(
    &namerv1.RootNamedResource{},
    GetIDNamerVars,
    SetIDNamerVars,
)
```
- **Format:** `namer.Format(ID{123}, 0)` → `"rootNamedResources/123"`
- **Parse:** `namer.Parse("rootNamedResources/123")` → `ID{123}, 0, nil`

### 2. ParentNamer
For singleton resources (e.g., `parentOnes/{parent_one}/singletonNamedResource`):
- Define the parent type
- Define a function that returns the parent variables for the given parent and pattern index
- Define a function that sets the parent variables for the given parent variables and pattern index
- Instantiate the ParentNamer with the resource, the parent variable getter, and the parent variable setter

```go
type Parent struct { 
    One int64
    Two int64
}

func GetSingletonNamedVars(parent Parent, patternIndex int) ([]string, error) { ... }
func SetSingletonNamedParent(vars []string, patternIndex int) (Parent, error) { ... }

namer, err := NewParentNamer(
    &namerv1.SingletonNamedResource{},
    GetSingletonNamedVars,
    SetSingletonNamedParent,
)
```
- **Format:** `namer.Format(Parent{1, 2}, 1)` → `"ones/1/twos/2/singletonNamedResource"`
- **Parse:** `namer.Parse("ones/1/twos/2/singletonNamedResource")` → `Parent{1, 2}, 0, nil`
- **ParseParent:** `namer.ParseParent("ones/1/twos/2")` → `Parent{1, 2}, 0, nil`

### 3. ParentIDNamer
For resources with both a parent and an ID (e.g., `parentOnes/{parent_one}/standardNamedResources/{standard_named_resource}`) or a "standard" pattern:
- Define the parent type
- Define the id type
- Define a function that returns the parent and id variables for the given parent, id, and pattern index
- Define a function that sets the parent variables for the given parent variables and pattern index
- Define a function that sets the id variable for the given id variable and pattern index
- Instantiate the ParentIDNamer with the resource, the parent and id variable getter, the parent variable setter, and the id variable setter

```go
type Parent struct { 
    One int64
    Two int64
}
type ID struct { 
    ID int64
}

func GetStandardNamedVars(parent StandardNamedResourceParent, id StandardNamedResourceId, patternIndex int) ([]string, error) { ... }
func SetStandardNamedID(idVar string, patternIndex int) (StandardNamedResourceId, error) { ... }
func SetStandardNamedParent(vars []string, patternIndex int) (StandardNamedResourceParent, error) { ... }

namer, err := NewParentIDNamer(
    &namerv1.StandardNamedResource{},
    GetStandardNamedVars,
    SetStandardNamedParent,
    SetStandardNamedID,
)
```
- **Format:** `namer.Format(Parent{1, 2}, ID{3}, 0)` → `"ones/1/twos/2/standardNamedResources/3"`
- **Parse:** `namer.Parse("ones/1/twos/2/standardNamedResources/3")` → `Parent{1, 2}, ID{3}, 0, nil`
- **ParseParent:** `namer.ParseParent("ones/1/twos/2")` → `StandardNamedResourceParent{1, 2}, 0, nil`

### 4. ReflectNamer
For resources with a parent and/or id.
- Define the resource type
- Instantiate the ReflectNamer with the resource

```go
type Resource struct {
  ParentOne int64 `aip_pattern:"key=parent_one"`
  ParentTwo int64 `aip_pattern:"key=parent_two"`
  ID int64 `aip_pattern:"key=resource"`
}

namer, err := NewReflectNamer[Resource](
  &namerv1.Resource{},
)
```
- **Format:** `namer.Format(0, Resource{1, 2, 3})` → `"parent_one/1/parent_two/2/resource/3"`
- **Parse:** `namer.Parse("parent_one/1/parent_two/2/resource/3")` → `Resource{1, 2, 3}, 0, nil`
- **ParseParent:** `namer.ParseParent("parent_one/1/parent_two/2")` → `Resource{1, 2, 0}, 0, nil`

## Implementation Details

- **Pattern Extraction:** Patterns are extracted from the protobuf resource descriptor using reflection onces at startup.
- **Pattern Matching:** Uses `go.einride.tech/aip/resourcename` for parsing and formatting.
- **Error Handling:** Returns descriptive errors for invalid patterns, indices, or parse failures.

## Adding a New Resource Namer

1. **Define the resource and its patterns in your proto file.**
2. **Generate the Go code using `buf generate`.**
3. **Implement the ID/Parent structs and getter/setter functions.**
4. **Instantiate the appropriate namer in your adapter or service.**
5. **It recommended to write tests for all expected name formats and parse scenarios.**

## Example Test

See the test files for comprehensive examples:
- `id_namer_root_test.go` : IDNamer for a "root" pattern
- `id_namer_root_singleton_test.go` : IDNamer for a "root" pattern with a singleton
- `parent_namer_test.go` : ParentNamer for a singleton pattern
- `parent_id_namer_test.go` : ParentIDNamer for a standard pattern

---

**Questions?**  
See the code comments and tests for more details, or ask the backend team for guidance on complex resource patterns. 