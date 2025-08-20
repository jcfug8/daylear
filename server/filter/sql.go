package filter

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"go.einride.tech/aip/filtering"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// Error messages
const (
	errNilExpression     = "expression is nil"
	errUnsupportedExpr   = "unsupported expression type: %T"
	errUnsupportedFunc   = "unsupported function: %s"
	errUnsupportedConst  = "unsupported constant type: %T"
	errUnsupportedValue  = "unsupported value expression type: %T"
	errUnsupportedField  = "unsupported field expression type: %T"
	errInvalidArgs       = "%s operator requires exactly %d arguments"
	errInvalidFilter     = "failed to parse filter: %w"
	errInvalidConversion = "failed to convert expression: %w"
)

// SQLConverter converts AIP-160 filter expressions to SQL WHERE clauses
type SQLConverter struct {
	// FieldMapping maps AIP field names to SQL column names
	FieldMapping    map[string]Field
	UseQuestionMark bool // if true, use ? placeholders (GORM style); else use $1, $2... (Postgres style)
}

type Field struct {
	Name            string
	Table           string
	CustomConverter CustomConverter
}

// CustomConverter allows custom logic for field conversions
// Returns (customSQL, used) where used=true means use custom logic, false means fall back to default
type CustomConverter func(ctx *Conversion, field string, operator string, value interface{}) (string, bool)

func (f Field) String() string {
	name := f.Name
	if f.Table != "" {
		name = f.Table + "." + name
	}
	return name
}

type Conversion struct {
	FieldMapping    map[string]Field
	WhereClause     string
	Params          []interface{}
	UsedColumns     map[string]int
	UseQuestionMark bool
}

// NewSQLConverter creates a new SQLConverter with the given field mapping and optional useQuestionMark flag
func NewSQLConverter(fieldMapping map[string]Field, useQuestionMark ...bool) *SQLConverter {
	useQM := false
	if len(useQuestionMark) > 0 {
		useQM = useQuestionMark[0]
	}
	return &SQLConverter{
		FieldMapping:    fieldMapping,
		UseQuestionMark: useQM,
	}
}

// Convert parses an AIP-160 filter expression and converts it to a SQL WHERE clause
func (c *SQLConverter) Convert(filter string) (*Conversion, error) {
	conversion := &Conversion{
		FieldMapping:    c.FieldMapping,
		Params:          make([]interface{}, 0),
		UsedColumns:     make(map[string]int),
		UseQuestionMark: c.UseQuestionMark,
	}

	return conversion.convert(filter)
}

func (c *Conversion) convert(filter string) (*Conversion, error) {
	if filter == "" {
		return c, nil
	}

	var parser filtering.Parser
	parser.Init(filter)
	parsed, err := parser.Parse()
	if err != nil {
		return c, fmt.Errorf(errInvalidFilter, err)
	}

	sql, err := c.convertExpression(parsed.Expr)
	if err != nil {
		return c, fmt.Errorf(errInvalidConversion, err)
	}

	c.WhereClause = sql

	return c, nil
}

// Helper functions for common operations
func (c *Conversion) addParam(value interface{}) string {
	paramIndex := len(c.Params)
	c.Params = append(c.Params, value)
	if c.UseQuestionMark {
		return "?"
	}
	return fmt.Sprintf("$%d", paramIndex+1)
}

func (c *Conversion) convertWildcardString(value string) string {
	return strings.ReplaceAll(value, "*", "%")
}

func (c *Conversion) handleWildcardString(field string, value string) string {
	likeValue := c.convertWildcardString(value)
	return fmt.Sprintf("%s LIKE %s", field, c.addParam(likeValue))
}

func (c *Conversion) handleNullComparison(field string, value interface{}, isEquals bool) string {
	if value == nil {
		if isEquals {
			return fmt.Sprintf("%s IS NULL", field)
		}
		return fmt.Sprintf("%s IS NOT NULL", field)
	}
	// Check if the value is a string containing wildcards
	if strValue, ok := value.(string); ok && strings.Contains(strValue, "*") {
		likeValue := c.convertWildcardString(strValue)
		if isEquals {
			return fmt.Sprintf("%s LIKE %s", field, c.addParam(likeValue))
		}
		return fmt.Sprintf("%s NOT LIKE %s", field, c.addParam(likeValue))
	}
	operator := "="
	if !isEquals {
		operator = "!="
	}
	return fmt.Sprintf("%s %s %s", field, operator, c.addParam(value))
}

func (c *Conversion) convertExpression(expr *expr.Expr) (string, error) {
	if expr == nil {
		return "", errors.New(errNilExpression)
	}
	if call := expr.GetCallExpr(); call != nil {
		return c.convertCallExpr(call)
	}
	if cnst := expr.GetConstExpr(); cnst != nil {
		value, err := c.convertConstExpr(cnst)
		if err != nil {
			return "", err
		}
		return c.addParam(value), nil
	}
	if ident := expr.GetIdentExpr(); ident != nil {
		value, err := c.convertIdentExpr(ident)
		if err != nil {
			return "", err
		}
		return c.addParam(value), nil
	}
	return "", fmt.Errorf(errUnsupportedExpr, expr.ExprKind)
}

func (c *Conversion) convertCallExpr(call *expr.Expr_Call) (string, error) {
	function := call.Function
	args := call.Args

	if len(args) < 1 {
		return "", fmt.Errorf(errInvalidArgs, function, 1)
	}

	// Handle logical operations first
	switch function {
	case "AND":
		if len(args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, function, 2)
		}
		left, err := c.convertExpression(args[0])
		if err != nil {
			return "", err
		}
		right, err := c.convertExpression(args[1])
		if err != nil {
			return "", err
		}
		// Only add parens if left or right contains OR at the top level
		if needsParens(args[0], "OR") {
			left = "(" + left + ")"
		}
		if needsParens(args[1], "OR") {
			right = "(" + right + ")"
		}
		return fmt.Sprintf("%s AND %s", left, right), nil

	case "OR":
		if len(args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, function, 2)
		}
		left, err := c.convertExpression(args[0])
		if err != nil {
			return "", err
		}
		right, err := c.convertExpression(args[1])
		if err != nil {
			return "", err
		}
		// Always add parens if left or right is AND at the top level
		if needsParens(args[0], "AND") {
			left = "(" + left + ")"
		}
		if needsParens(args[1], "AND") {
			right = "(" + right + ")"
		}
		return fmt.Sprintf("%s OR %s", left, right), nil

	case "NOT":
		if len(args) != 1 {
			return "", fmt.Errorf(errInvalidArgs, function, 1)
		}
		// Try to optimize NOT expressions
		if notCall, ok := args[0].ExprKind.(*expr.Expr_CallExpr); ok {
			switch notCall.CallExpr.Function {
			case "<", ">", "<=", ">=", "=", "!=":
				if len(notCall.CallExpr.Args) != 2 {
					return "", fmt.Errorf(errInvalidArgs, notCall.CallExpr.Function, 2)
				}
				field, err := c.convertFieldExpr(notCall.CallExpr.Args[0])
				if err != nil {
					return "", err
				}
				value, err := c.convertValueExpr(notCall.CallExpr.Args[1])
				if err != nil {
					return "", err
				}
				// Handle wildcard strings
				if strValue, ok := value.(string); ok && strings.Contains(strValue, "*") {
					return fmt.Sprintf("NOT (%s)", c.handleWildcardString(field, strValue)), nil
				}
				// Handle null values
				if value == nil {
					return fmt.Sprintf("%s IS NOT NULL", field), nil
				}
				// Handle regular comparisons
				operator := notCall.CallExpr.Function
				switch operator {
				case "<":
					return fmt.Sprintf("%s >= %s", field, c.addParam(value)), nil
				case ">":
					return fmt.Sprintf("%s <= %s", field, c.addParam(value)), nil
				case "<=":
					return fmt.Sprintf("%s > %s", field, c.addParam(value)), nil
				case ">=":
					return fmt.Sprintf("%s < %s", field, c.addParam(value)), nil
				case "=":
					return fmt.Sprintf("%s != %s", field, c.addParam(value)), nil
				case "!=":
					return fmt.Sprintf("%s = %s", field, c.addParam(value)), nil
				}
			}
		}
		// If we can't optimize, fall back to the standard NOT expression
		expr, err := c.convertExpression(args[0])
		if err != nil {
			return "", err
		}
		// Don't add extra parentheses if the expression already has them
		if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
			return fmt.Sprintf("NOT %s", expr), nil
		}
		return fmt.Sprintf("NOT (%s)", expr), nil
	}

	// For other operations, get the field name from the first argument
	field, err := c.convertFieldExpr(args[0])
	if err != nil {
		return "", err
	}

	simpleField, err := c.convertFieldExprRecursive(args[0])
	if err != nil {
		return "", err
	}

	// Check if this field has a custom converter before handling default logic
	if fieldInfo, exists := c.FieldMapping[simpleField]; exists && fieldInfo.CustomConverter != nil {
		// Try custom converter first
		if len(args) == 2 {
			value, err := c.convertValueExpr(args[1])
			if err != nil {
				return "", err
			}

			if customSQL, used := fieldInfo.CustomConverter(c, field, function, value); used {
				return customSQL, nil
			}
			// If custom converter returns false, fall through to default behavior
		}
	}

	// Handle different function types
	switch function {
	case "=", "!=":
		if len(args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, function, 2)
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		isEquals := function == "="
		return c.handleNullComparison(field, value, isEquals), nil

	case ">", ">=", "<", "<=":
		if len(args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, function, 2)
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s %s", field, function, c.addParam(value)), nil

	case "contains":
		if len(args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, function, 2)
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s LIKE '%%' || %s || '%%'", field, c.addParam(value)), nil

	case "starts_with":
		if len(args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, function, 2)
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s LIKE %s || '%%'", field, c.addParam(value)), nil

	case "ends_with":
		if len(args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, function, 2)
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s LIKE '%%' || %s", field, c.addParam(value)), nil

	case ":":
		if len(args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, function, 2)
		}
		// For the has operator, we only care about the field name
		return fmt.Sprintf("%s IS NOT NULL", field), nil

	default:
		return "", fmt.Errorf(errUnsupportedFunc, function)
	}
}

func (c *Conversion) convertFieldExpr(expr *expr.Expr) (string, error) {
	field, err := c.convertFieldExprRecursive(expr)
	if err != nil {
		return "", err
	}
	return c.useField(field)
}

func (c *Conversion) convertFieldExprRecursive(expr *expr.Expr) (string, error) {
	var field string
	var err error

	if expr == nil {
		return "", errors.New(errNilExpression)
	} else if ident := expr.GetIdentExpr(); ident != nil {
		field = ident.Name
	} else if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		// Handle nested field expressions like user.profile.name
		operand, err := c.convertFieldExprRecursive(selectExpr.Operand)
		if err != nil {
			return "", err
		}
		// Convert the nested field path to a single field name by replacing dots with underscores
		field = operand + "." + selectExpr.Field
	} else if call := expr.GetCallExpr(); call != nil {
		field, err = c.convertFieldCallExpr(call)
		if err != nil {
			return "", err
		}
	} else {
		return "", fmt.Errorf(errUnsupportedField, expr.ExprKind)
	}

	return field, nil
}

func (c *Conversion) useField(field string) (string, error) {
	sqlField, ok := c.FieldMapping[field]
	if !ok {
		return "", fmt.Errorf("field not found in field mapping: %s", field)
	}

	if _, ok := c.UsedColumns[sqlField.String()]; !ok {
		c.UsedColumns[sqlField.String()] = 0
	}
	c.UsedColumns[sqlField.String()]++
	return sqlField.String(), nil
}

func (c *Conversion) convertFieldCallExpr(call *expr.Expr_Call) (string, error) {
	switch call.Function {
	case "and":
		if len(call.Args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, call.Function, 2)
		}
		left, err := c.convertFieldExpr(call.Args[0])
		if err != nil {
			return "", err
		}
		right, err := c.convertFieldExpr(call.Args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s AND %s)", left, right), nil

	case "or":
		if len(call.Args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, call.Function, 2)
		}
		left, err := c.convertFieldExpr(call.Args[0])
		if err != nil {
			return "", err
		}
		right, err := c.convertFieldExpr(call.Args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("(%s OR %s)", left, right), nil

	case ">=", "<=", ">", "<", "=", "!=":
		if len(call.Args) != 2 {
			return "", fmt.Errorf(errInvalidArgs, call.Function, 2)
		}
		left, err := c.convertFieldExpr(call.Args[0])
		if err != nil {
			return "", err
		}
		right, err := c.convertValueExpr(call.Args[1])
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s %s", left, call.Function, c.addParam(right)), nil

	default:
		return "", fmt.Errorf(errUnsupportedFunc, call.Function)
	}
}

func (c *Conversion) convertValueExpr(expr *expr.Expr) (interface{}, error) {
	if expr == nil {
		return nil, errors.New(errNilExpression)
	}
	if cnst := expr.GetConstExpr(); cnst != nil {
		return c.convertConstExpr(cnst)
	}
	if ident := expr.GetIdentExpr(); ident != nil {
		return c.convertIdentExpr(ident)
	}
	return nil, fmt.Errorf(errUnsupportedValue, expr.ExprKind)
}

func (c *Conversion) convertConstExpr(constExpr *expr.Constant) (interface{}, error) {
	switch k := constExpr.ConstantKind.(type) {
	case *expr.Constant_StringValue:
		// Special case: treat "null" string as nil
		if k.StringValue == "null" {
			return nil, nil
		}
		// Always treat as string, do not parse as duration
		return k.StringValue, nil
	case *expr.Constant_Int64Value:
		return k.Int64Value, nil
	case *expr.Constant_Uint64Value:
		return k.Uint64Value, nil
	case *expr.Constant_DoubleValue:
		return k.DoubleValue, nil
	case *expr.Constant_BoolValue:
		return k.BoolValue, nil
	case *expr.Constant_NullValue:
		return nil, nil
	default:
		return nil, fmt.Errorf(errUnsupportedConst, constExpr.ConstantKind)
	}
}

func (c *Conversion) convertIdentExpr(identExpr *expr.Expr_Ident) (interface{}, error) {
	name := identExpr.Name
	switch name {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	}
	// Check for duration pattern: digits followed by 's', e.g., 20s, 1.5s
	if strings.HasSuffix(name, "s") {
		numeric := name[:len(name)-1]
		if _, err := strconv.ParseFloat(numeric, 64); err == nil {
			return name, nil // treat as duration string, e.g., "20s"
		}
	}
	return name, nil
}

// Helper function to determine if a node is a logical op needing parens
func needsParens(expr *expr.Expr, op string) bool {
	if call := expr.GetCallExpr(); call != nil {
		if call.Function == op {
			return true
		}
	}
	return false
}
