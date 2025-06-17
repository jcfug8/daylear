package filter

import (
	"fmt"
	"strconv"
	"strings"

	"go.einride.tech/aip/filtering"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

// SQLConverter converts AIP-160 filter expressions to SQL WHERE clauses
type SQLConverter struct {
	// FieldMapping maps AIP field names to SQL column names
	FieldMapping map[string]string
	// Params holds the parameter values for the SQL query
	Params []interface{}
}

// NewSQLConverter creates a new SQLConverter with the given field mapping
func NewSQLConverter(fieldMapping map[string]string) *SQLConverter {
	return &SQLConverter{
		FieldMapping: fieldMapping,
		Params:       make([]interface{}, 0),
	}
}

// Convert parses an AIP-160 filter expression and converts it to a SQL WHERE clause
func (c *SQLConverter) Convert(filter string) (string, error) {
	if filter == "" {
		return "", nil
	}

	var parser filtering.Parser
	parser.Init(filter)
	parsed, err := parser.Parse()
	if err != nil {
		return "", fmt.Errorf("failed to parse filter: %w", err)
	}

	sql, err := c.convertExpression(parsed.Expr)
	if err != nil {
		return "", fmt.Errorf("failed to convert expression: %w", err)
	}

	return sql, nil
}

func (c *SQLConverter) convertExpression(expr *expr.Expr) (string, error) {
	if expr == nil {
		return "", fmt.Errorf("expression is nil")
	}
	if call := expr.GetCallExpr(); call != nil {
		return c.convertCallExpr(call)
	}
	if cnst := expr.GetConstExpr(); cnst != nil {
		value, err := c.convertConstExpr(cnst)
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("$%d", paramIndex+1), nil
	}
	if ident := expr.GetIdentExpr(); ident != nil {
		value, err := c.convertIdentExpr(ident)
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("$%d", paramIndex+1), nil
	}
	return "", fmt.Errorf("unsupported expression type: %T", expr.ExprKind)
}

func (c *SQLConverter) convertCallExpr(call *expr.Expr_Call) (string, error) {
	function := call.Function
	args := call.Args

	if len(args) < 1 {
		return "", fmt.Errorf("call expression requires at least one argument")
	}

	// Handle logical operations first
	switch function {
	case "AND":
		if len(args) != 2 {
			return "", fmt.Errorf("and operator requires exactly two arguments")
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
			return "", fmt.Errorf("or operator requires exactly two arguments")
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
			return "", fmt.Errorf("not operator requires exactly one argument")
		}
		// Try to optimize NOT expressions
		if notCall, ok := args[0].ExprKind.(*expr.Expr_CallExpr); ok {
			switch notCall.CallExpr.Function {
			case "<":
				if len(notCall.CallExpr.Args) != 2 {
					return "", fmt.Errorf("less than operator requires exactly two arguments")
				}
				field, err := c.convertFieldExpr(notCall.CallExpr.Args[0])
				if err != nil {
					return "", err
				}
				value, err := c.convertValueExpr(notCall.CallExpr.Args[1])
				if err != nil {
					return "", err
				}
				paramIndex := len(c.Params)
				c.Params = append(c.Params, value)
				return fmt.Sprintf("%s >= $%d", field, paramIndex+1), nil
			case ">":
				if len(notCall.CallExpr.Args) != 2 {
					return "", fmt.Errorf("greater than operator requires exactly two arguments")
				}
				field, err := c.convertFieldExpr(notCall.CallExpr.Args[0])
				if err != nil {
					return "", err
				}
				value, err := c.convertValueExpr(notCall.CallExpr.Args[1])
				if err != nil {
					return "", err
				}
				paramIndex := len(c.Params)
				c.Params = append(c.Params, value)
				return fmt.Sprintf("%s <= $%d", field, paramIndex+1), nil
			case "<=":
				if len(notCall.CallExpr.Args) != 2 {
					return "", fmt.Errorf("less than or equal operator requires exactly two arguments")
				}
				field, err := c.convertFieldExpr(notCall.CallExpr.Args[0])
				if err != nil {
					return "", err
				}
				value, err := c.convertValueExpr(notCall.CallExpr.Args[1])
				if err != nil {
					return "", err
				}
				paramIndex := len(c.Params)
				c.Params = append(c.Params, value)
				return fmt.Sprintf("%s > $%d", field, paramIndex+1), nil
			case ">=":
				if len(notCall.CallExpr.Args) != 2 {
					return "", fmt.Errorf("greater than or equal operator requires exactly two arguments")
				}
				field, err := c.convertFieldExpr(notCall.CallExpr.Args[0])
				if err != nil {
					return "", err
				}
				value, err := c.convertValueExpr(notCall.CallExpr.Args[1])
				if err != nil {
					return "", err
				}
				paramIndex := len(c.Params)
				c.Params = append(c.Params, value)
				return fmt.Sprintf("%s < $%d", field, paramIndex+1), nil
			case "=":
				if len(notCall.CallExpr.Args) != 2 {
					return "", fmt.Errorf("equals operator requires exactly two arguments")
				}
				field, err := c.convertFieldExpr(notCall.CallExpr.Args[0])
				if err != nil {
					return "", err
				}
				value, err := c.convertValueExpr(notCall.CallExpr.Args[1])
				if err != nil {
					return "", err
				}
				paramIndex := len(c.Params)
				c.Params = append(c.Params, value)
				// Check if the value is a string containing wildcards
				if strValue, ok := value.(string); ok && strings.Contains(strValue, "*") {
					// Replace * with % for SQL LIKE
					likeValue := strings.ReplaceAll(strValue, "*", "%")
					c.Params[paramIndex] = likeValue
					return fmt.Sprintf("NOT (%s LIKE $%d)", field, paramIndex+1), nil
				}
				return fmt.Sprintf("%s != $%d", field, paramIndex+1), nil
			case "!=":
				if len(notCall.CallExpr.Args) != 2 {
					return "", fmt.Errorf("not equals operator requires exactly two arguments")
				}
				field, err := c.convertFieldExpr(notCall.CallExpr.Args[0])
				if err != nil {
					return "", err
				}
				value, err := c.convertValueExpr(notCall.CallExpr.Args[1])
				if err != nil {
					return "", err
				}
				paramIndex := len(c.Params)
				c.Params = append(c.Params, value)
				// Check if the value is a string containing wildcards
				if strValue, ok := value.(string); ok && strings.Contains(strValue, "*") {
					// Replace * with % for SQL LIKE
					likeValue := strings.ReplaceAll(strValue, "*", "%")
					c.Params[paramIndex] = likeValue
					return fmt.Sprintf("NOT (%s LIKE $%d)", field, paramIndex+1), nil
				}
				return fmt.Sprintf("%s = $%d", field, paramIndex+1), nil
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

	// Handle different function types
	switch function {
	case "=":
		if len(args) != 2 {
			return "", fmt.Errorf("equals operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		if value == nil {
			// Special case: field = null -> IS NULL
			return fmt.Sprintf("%s IS NULL", field), nil
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		// Check if the value is a string containing wildcards
		if strValue, ok := value.(string); ok && strings.Contains(strValue, "*") {
			// Replace * with % for SQL LIKE
			likeValue := strings.ReplaceAll(strValue, "*", "%")
			c.Params[paramIndex] = likeValue
			return fmt.Sprintf("%s LIKE $%d", field, paramIndex+1), nil
		}
		return fmt.Sprintf("%s = $%d", field, paramIndex+1), nil

	case "!=":
		if len(args) != 2 {
			return "", fmt.Errorf("not equals operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		if value == nil {
			// Special case: field != null -> IS NOT NULL
			return fmt.Sprintf("%s IS NOT NULL", field), nil
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		// Check if the value is a string containing wildcards
		if strValue, ok := value.(string); ok && strings.Contains(strValue, "*") {
			// Replace * with % for SQL LIKE
			likeValue := strings.ReplaceAll(strValue, "*", "%")
			c.Params[paramIndex] = likeValue
			return fmt.Sprintf("%s NOT LIKE $%d", field, paramIndex+1), nil
		}
		return fmt.Sprintf("%s != $%d", field, paramIndex+1), nil

	case ">":
		if len(args) != 2 {
			return "", fmt.Errorf("greater than operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("%s > $%d", field, paramIndex+1), nil

	case ">=":
		if len(args) != 2 {
			return "", fmt.Errorf("greater than or equal operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("%s >= $%d", field, paramIndex+1), nil

	case "<":
		if len(args) != 2 {
			return "", fmt.Errorf("less than operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("%s < $%d", field, paramIndex+1), nil

	case "<=":
		if len(args) != 2 {
			return "", fmt.Errorf("less than or equal operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("%s <= $%d", field, paramIndex+1), nil

	case "contains":
		if len(args) != 2 {
			return "", fmt.Errorf("contains operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("%s LIKE '%%' || $%d || '%%'", field, paramIndex+1), nil

	case "starts_with":
		if len(args) != 2 {
			return "", fmt.Errorf("starts with operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("%s LIKE $%d || '%%'", field, paramIndex+1), nil

	case "ends_with":
		if len(args) != 2 {
			return "", fmt.Errorf("ends with operator requires exactly two arguments")
		}
		value, err := c.convertValueExpr(args[1])
		if err != nil {
			return "", err
		}
		paramIndex := len(c.Params)
		c.Params = append(c.Params, value)
		return fmt.Sprintf("%s LIKE '%%' || $%d", field, paramIndex+1), nil

	case ":":
		if len(args) != 2 {
			return "", fmt.Errorf("has operator requires exactly two arguments")
		}
		// For the has operator, we only care about the field name
		return fmt.Sprintf("%s IS NOT NULL", field), nil

	default:
		return "", fmt.Errorf("unsupported function: %s", function)
	}
}

func (c *SQLConverter) convertFieldExpr(expr *expr.Expr) (string, error) {
	if expr == nil {
		return "", fmt.Errorf("field expression is nil")
	}
	if ident := expr.GetIdentExpr(); ident != nil {
		field := ident.Name
		if sqlField, ok := c.FieldMapping[field]; ok {
			return sqlField, nil
		}
		return field, nil
	}
	if selectExpr := expr.GetSelectExpr(); selectExpr != nil {
		// Handle nested field expressions like user.profile.name
		operand, err := c.convertFieldExpr(selectExpr.Operand)
		if err != nil {
			return "", err
		}
		// Convert the nested field path to a single field name by replacing dots with underscores
		field := operand + "_" + selectExpr.Field
		if sqlField, ok := c.FieldMapping[field]; ok {
			return sqlField, nil
		}
		return field, nil
	}
	if call := expr.GetCallExpr(); call != nil {
		// Handle logical operations like AND and OR
		switch call.Function {
		case "and":
			if len(call.Args) != 2 {
				return "", fmt.Errorf("and operator requires exactly two arguments")
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
				return "", fmt.Errorf("or operator requires exactly two arguments")
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
				return "", fmt.Errorf("comparison operator requires exactly two arguments")
			}
			left, err := c.convertFieldExpr(call.Args[0])
			if err != nil {
				return "", err
			}
			right, err := c.convertValueExpr(call.Args[1])
			if err != nil {
				return "", err
			}
			paramIndex := len(c.Params)
			c.Params = append(c.Params, right)
			return fmt.Sprintf("%s %s $%d", left, call.Function, paramIndex+1), nil
		default:
			return "", fmt.Errorf("unsupported function: %s", call.Function)
		}
	}
	return "", fmt.Errorf("unsupported field expression type: %T", expr.ExprKind)
}

func (c *SQLConverter) convertValueExpr(expr *expr.Expr) (interface{}, error) {
	if expr == nil {
		return nil, fmt.Errorf("value expression is nil")
	}
	if cnst := expr.GetConstExpr(); cnst != nil {
		return c.convertConstExpr(cnst)
	}
	if ident := expr.GetIdentExpr(); ident != nil {
		return c.convertIdentExpr(ident)
	}
	return nil, fmt.Errorf("unsupported value expression type: %T", expr.ExprKind)
}

func (c *SQLConverter) convertConstExpr(constExpr *expr.Constant) (interface{}, error) {
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
		return nil, fmt.Errorf("unsupported constant type: %T", constExpr.ConstantKind)
	}
}

func (c *SQLConverter) convertIdentExpr(identExpr *expr.Expr_Ident) (interface{}, error) {
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
