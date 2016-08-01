package qm

import "github.com/nullbio/sqlboiler/boil"

// QueryMod to modify the query object
type QueryMod func(q *boil.Query)

// Apply the query mods to the Query object
func Apply(q *boil.Query, mods ...QueryMod) {
	for _, mod := range mods {
		mod(q)
	}
}

// SQL allows you to execute a plain SQL statement
func SQL(sql string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetSQL(q, sql, args...)
	}
}

// Or surrounds where clauses to join them with OR as opposed to AND
func Or(whereMods ...QueryMod) QueryMod {
	return func(q *boil.Query) {
		if len(whereMods) < 2 {
			panic("Or requires at least two arguments")
		}

		for _, w := range whereMods {
			w(q)
			boil.SetLastWhereAsOr(q)
		}
	}
}

// Limit the number of returned rows
func Limit(limit int) QueryMod {
	return func(q *boil.Query) {
		boil.SetLimit(q, limit)
	}
}

// Offset into the results
func Offset(offset int) QueryMod {
	return func(q *boil.Query) {
		boil.SetOffset(q, offset)
	}
}

// InnerJoin on another table
func InnerJoin(stmt string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetInnerJoin(q, stmt, args...)
	}
}

// OuterJoin on another table
func OuterJoin(stmt string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetOuterJoin(q, stmt, args...)
	}
}

// LeftOuterJoin on another table
func LeftOuterJoin(stmt string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetLeftOuterJoin(q, stmt, args...)
	}
}

// RightOuterJoin on another table
func RightOuterJoin(stmt string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetRightOuterJoin(q, stmt, args...)
	}
}

// Select specific columns opposed to all columns
func Select(columns ...string) QueryMod {
	return func(q *boil.Query) {
		boil.SetSelect(q, columns...)
	}
}

// Where allows you to specify a where clause for your statement
func Where(clause string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetWhere(q, clause, args...)
	}
}

// GroupBy allows you to specify a group by clause for your statement
func GroupBy(clause string) QueryMod {
	return func(q *boil.Query) {
		boil.SetGroupBy(q, clause)
	}
}

// OrderBy allows you to specify a order by clause for your statement
func OrderBy(clause string) QueryMod {
	return func(q *boil.Query) {
		boil.SetOrderBy(q, clause)
	}
}

// Having allows you to specify a having clause for your statement
func Having(clause string) QueryMod {
	return func(q *boil.Query) {
		boil.SetHaving(q, clause)
	}
}

// From allows to specify the table for your statement
func From(from string) QueryMod {
	return func(q *boil.Query) {
		boil.SetFrom(q, from)
	}
}
