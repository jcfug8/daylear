package filtering

import (
	"fmt"
	"strings"

	"go.einride.tech/aip/filtering"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	"gorm.io/gorm/clause"
)

// NewSQLTranspiler creates a new SQL transpiler.
func NewSQLTranspiler(
	fields map[string]Field[clause.Expression],
) Transpiler[clause.Expression] {
	return Transpiler[clause.Expression]{
		And: func(x, y clause.Expression) (clause.Expression, error) {
			if and, ok := x.(*clause.AndConditions); ok {
				and.Exprs = append(and.Exprs, y)
				return and, nil
			}

			and := &clause.AndConditions{Exprs: []clause.Expression{x, y}}

			return and, nil
		},
		Or: func(x, y clause.Expression) (clause.Expression, error) {
			if or, ok := x.(*clause.OrConditions); ok {
				or.Exprs = append(or.Exprs, y)
				return or, nil
			}

			or := &clause.OrConditions{Exprs: []clause.Expression{x, y}}

			return or, nil
		},
		Compare: func(op, path string, value any) (clause.Expression, error) {
			if field, ok := fields[path]; ok {
				return field.Compare(op, path, value)
			}

			return nil, NewPathError("field unsupported")
		},
	}
}

// Transpiler is a generic transpiler for filtering expressions.
type Transpiler[T any] struct {
	And     func(T, T) (T, error)
	Or      func(T, T) (T, error)
	Compare func(string, string, any) (T, error)
	filter  string
	source  *expr.SourceInfo
	info    *TranspileInfo
}

// Transpile transpiles a filter expression into a value of type T.
func (t *Transpiler[T]) Transpile(filter string) (T, *TranspileInfo, error) {
	var value T

	if filter == "" {
		return value, &TranspileInfo{}, nil
	}

	var parser filtering.Parser
	parser.Init(filter)
	parsed, err := parser.Parse()
	if err != nil {
		return value, nil, fmt.Errorf("could not parse: %v", err)
	}

	t = &Transpiler[T]{
		Compare: t.Compare,
		And:     t.And,
		Or:      t.Or,
		filter:  filter,
		source:  parsed.SourceInfo,
		info: &TranspileInfo{
			fields: make(map[string]bool),
		},
	}

	value, err = t.transpileExpr(parsed.Expr)
	if err != nil {
		return value, nil, err
	}

	return value, t.info, nil
}

func (t *Transpiler[T]) transpileExpr(e *expr.Expr) (T, error) {
	var zero T

	if _, ok := e.GetExprKind().(*expr.Expr_CallExpr); ok {
		return t.transpileCallExpr(e)
	}

	return zero, fmt.Errorf("unsupported expr: %v", e)
}

func (t *Transpiler[T]) transpileCallExpr(e *expr.Expr) (T, error) {
	var zero T

	function := e.GetCallExpr().GetFunction()

	switch function {
	case filtering.FunctionEquals:
		return t.transpileCompareCallExpr(e)
	case filtering.FunctionNotEquals:
		return t.transpileCompareCallExpr(e)
	case filtering.FunctionLessThan:
		return t.transpileCompareCallExpr(e)
	case filtering.FunctionLessEquals:
		return t.transpileCompareCallExpr(e)
	case filtering.FunctionGreaterThan:
		return t.transpileCompareCallExpr(e)
	case filtering.FunctionGreaterEquals:
		return t.transpileCompareCallExpr(e)
	case filtering.FunctionNot:
		return t.transpileBinaryLogicalCallExpr(e)
	case filtering.FunctionAnd, filtering.FunctionFuzzyAnd:
		return t.transpileBinaryLogicalCallExpr(e)
	case filtering.FunctionOr:
		return t.transpileBinaryLogicalCallExpr(e)
	default:
		return zero, t.errorf(e, "operator unsupported")
	}
}

func (t *Transpiler[T]) transpileCompareCallExpr(e *expr.Expr) (T, error) {
	var (
		args = e.GetCallExpr().GetArgs()
		op   = e.GetCallExpr().GetFunction()
		zero T
	)

	if err := requireLen(2, args); err != nil {
		return zero, t.errorf(e, "arguments %v", err)
	}

	path, err := t.transpilePathExpr(args[0])
	if err != nil {
		return zero, t.errorf(args[0], "field %v", err)
	}

	value, err := t.transpileValueExpr(args[1])
	if err != nil {
		return zero, t.errorf(args[1], err.Error())
	}

	t.info.addField(path)

	comparison, err := t.Compare(op, path, value)
	if err != nil {
		switch err.(type) {
		case *operatorError:
			return zero, t.errorf(e, err.Error())
		case *pathError:
			return zero, t.errorf(args[0], err.Error())
		case *valueError:
			return zero, t.errorf(args[1], err.Error())
		default:
			return zero, t.errorf(args[0], err.Error())
		}
	}

	return comparison, nil
}

func (t *Transpiler[T]) transpilePathExpr(e *expr.Expr) (string, error) {
	switch kind := e.ExprKind.(type) {
	case *expr.Expr_ConstExpr:
		switch kind := kind.ConstExpr.ConstantKind.(type) {
		case *expr.Constant_StringValue:
			return kind.StringValue, nil
		default:
			return "", fmt.Errorf("path must be string")
		}
	case *expr.Expr_IdentExpr:
		return kind.IdentExpr.Name, nil
	case *expr.Expr_SelectExpr:
		path, err := t.transpilePathExpr(kind.SelectExpr.Operand)
		if err != nil {
			return "", err
		}

		path += "." + kind.SelectExpr.Field

		return path, nil
	default:
		return "", fmt.Errorf("unsupported path expression")
	}
}

func (t *Transpiler[T]) transpileValueExpr(e *expr.Expr) (any, error) {
	switch e.ExprKind.(type) {
	case *expr.Expr_CallExpr:
		var res any
		var err error

		switch e.GetCallExpr().GetFunction() {
		case "any":
			res, err = t.transpileAnyCallExpr(e)
		default:
			err = fmt.Errorf("unsupported function")
		}

		if err != nil {
			return nil, err
		}

		return res, err
	case *expr.Expr_ConstExpr:
		return t.transpileConstExpr(e)
	case *expr.Expr_IdentExpr:
		return t.transpileIdentExpr(e)
	default:
		return "", fmt.Errorf("unsupported value expression")
	}
}

func (*Transpiler[T]) transpileIdentExpr(e *expr.Expr) (any, error) {
	value := e.GetIdentExpr().GetName()

	switch value {
	case "false":
		return false, nil
	case "true":
		return true, nil
	default:
		return value, nil
	}
}

func (t *Transpiler[T]) transpileAnyCallExpr(e *expr.Expr) (any, error) {
	var ok bool
	var xs any

	for _, e := range e.GetCallExpr().GetArgs() {
		x, err := t.transpileValueExpr(e)
		if err != nil {
			return nil, err
		}

		switch x := x.(type) {
		case float64:
			xs, ok = appendAny(xs, x)
		case uint64:
			xs, ok = appendAny(xs, x)
		case string:
			xs, ok = appendAny(xs, x)
		case int64:
			xs, ok = appendAny(xs, x)
		case bool:
			xs, ok = appendAny(xs, x)
		default:
			supported := "float64, uint64, int64, string, bool"
			return nil, fmt.Errorf("any values must be [%s]", supported)
		}

		if !ok {
			return nil, fmt.Errorf("any values must be same type")
		}
	}

	return xs, nil
}

func (*Transpiler[T]) transpileConstExpr(e *expr.Expr) (any, error) {
	switch kind := e.GetConstExpr().GetConstantKind().(type) {
	case *expr.Constant_Int64Value:
		return kind.Int64Value, nil
	case *expr.Constant_NullValue:
		return nil, nil
	case *expr.Constant_BoolValue:
		return kind.BoolValue, nil
	case *expr.Constant_Uint64Value:
		return kind.Uint64Value, nil
	case *expr.Constant_DoubleValue:
		return kind.DoubleValue, nil
	case *expr.Constant_StringValue:
		return kind.StringValue, nil
	case *expr.Constant_BytesValue:
		return kind.BytesValue, nil
	default:
		return nil, fmt.Errorf("unsupported constant type")
	}
}

func (t *Transpiler[T]) transpileBinaryLogicalCallExpr(
	e *expr.Expr,
) (T, error) {
	var (
		args = e.GetCallExpr().GetArgs()
		op   = e.GetCallExpr().GetFunction()
		zero T
	)

	if err := requireLen(2, args); err != nil {
		return zero, t.errorf(e, "arguments: %v", err)
	}

	x, err := t.transpileExpr(args[0])
	if err != nil {
		return zero, err
	}

	y, err := t.transpileExpr(args[1])
	if err != nil {
		return zero, err
	}

	switch {
	case t.And != nil && (op == "FUZZY" || op == "AND"):
		return t.And(x, y)
	case t.Or != nil && op == "OR":
		return t.Or(x, y)
	default:
		return zero, t.errorf(e, "operator unsupported")
	}
}

func (t *Transpiler[T]) errorf(e *expr.Expr, format string, args ...any) error {
	pos := t.source.Positions[e.Id]

	if function := e.GetCallExpr().GetFunction(); function != "" {
		// The position offset for binary operators appears to always be
		// identical to the position offset of the left value in the binary
		// expression. Possibly a bug in the filtering module? To work
		// around the position can be shifted forward when the expression
		// defines a function and it's not found at the defined position.

		if n := strings.Index(t.filter[pos:], function); n > 0 {
			pos += int32(n)
		}
	}

	var line int32
	var tail int32

	for _, offset := range t.source.LineOffsets {
		if pos <= offset {
			tail = offset
			break
		}

		line = offset + 1
	}

	if tail == 0 {
		tail = int32(len(t.filter))
	}

	return fmt.Errorf("could not parse:\n%s\n%s^ %v",
		t.filter[:tail],
		strings.Repeat(" ", int(pos-line)),
		fmt.Errorf(format, args...))
}
