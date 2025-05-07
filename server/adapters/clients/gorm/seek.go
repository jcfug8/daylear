package gorm

import "gorm.io/gorm/clause"

// Seek builds a GORM clause.Expression that seeks to the given values.
// refer to https://vladmihalcea.com/sql-seek-keyset-pagination/
func Seek(
	orders []clause.OrderByColumn,
	values map[string]any,
) clause.Expression {
	var asc, desc bool
	columns := make([]clause.Column, len(orders))
	args := make([]any, 0, len(orders))

	if values == nil {
		values = make(map[string]any)
	}

	for i, order := range orders {
		columns[i] = order.Column
		args = append(args, values[order.Column.Name])

		desc = desc || order.Desc
		asc = asc || !order.Desc
	}

	if !(asc && desc) {
		if desc {
			return clause.Lt{Column: columns, Value: args}
		}

		return clause.Gt{Column: columns, Value: args}
	}

	var exprs []clause.Expression
	for len(orders) > 0 {
		tail := len(orders) - 1

		ops := make([]clause.Expression, 0, len(orders))
		for pos, column := range orders {
			if pos != tail {
				ops = append(ops, clause.Eq{
					Column: column.Column,
					Value:  values[column.Column.Name],
				})
				continue
			}

			if column.Desc {
				ops = append(ops, clause.Lt{
					Column: column.Column,
					Value:  values[column.Column.Name],
				})
			} else {
				ops = append(ops, clause.Gt{
					Column: column.Column,
					Value:  values[column.Column.Name],
				})
			}
		}

		exprs = append(exprs, clause.And(ops...))

		orders = orders[:tail]
	}

	return clause.Or(exprs...)
}
