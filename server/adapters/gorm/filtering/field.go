package filtering

import (
	"cmp"
	"slices"
	"strconv"
	"strings"

	"gorm.io/gorm/clause"
)

// NewSQLField creates a new SQLField with the given column name and operators.
func NewSQLField[T any](col string, ops ...string) *SQLField[T] {
	field := &SQLField[T]{
		ops: make(map[string]bool),
		col: col,
	}

	for _, op := range ops {
		field.ops[op] = true
	}

	return field
}

// SQLField represents a field in a SQL table.
type SQLField[T any] struct {
	col string
	ops map[string]bool
}

// Compare compares the field with the given operator, path, and value.
func (f *SQLField[T]) Compare(
	op, _path string,
	value any,
) (clause.Expression, error) {
	if !f.ops[op] {
		ops := strings.Join(sort(getKeys(f.ops)), ", ")
		return nil, NewOperatorError("operator not in supported [%s]", ops)
	}

	param := parameterize(value)

	value, err := coerceValue[T](value)
	if err != nil {
		return nil, NewValueError(err.Error())
	}

	switch op {
	case "=":
	case "!=":
	case ">=":
	case "<=":
	case ">":
	case "<":
	default:
		err = NewOperatorError("operator unsupported")
	}

	cond := clause.Expr{
		SQL:  f.col + " " + op + " " + param,
		Vars: []interface{}{value},
	}

	return cond, err
}

// Field represents a field in a data structure.
type Field[T any] interface {
	Compare(string, string, any) (T, error)
}

func getKeys[T comparable, U any](m map[T]U) []T {
	keys := make([]T, 0, len(m))

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

func sort[S ~[]E, E cmp.Ordered](xs S) S {
	slices.Sort(xs)
	return xs
}

func quoteString(value any) any {
	if s, ok := value.(string); ok {
		return strconv.Quote(s)
	}

	return value
}

func parameterize(value any) string {
	out := "?"
	if _, ok := value.(interface{ isAny() }); ok {
		out = "ANY(?)"
	}

	return out
}
