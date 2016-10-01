package qm

import "github.com/vattle/sqlboiler/queries"

// QueryMod to modify the query object
type QueryMod func(q *queries.Query)

// Apply the query mods to the Query object
func Apply(q *queries.Query, mods ...QueryMod) {
	for _, mod := range mods {
		mod(q)
	}
}

// SQL allows you to execute a plain SQL statement
func SQL(sql string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.SetSQL(q, sql, args...)
	}
}

// Load allows you to specify foreign key relationships to eager load
// for your query. Passed in relationships need to be in the format
// MyThing or MyThings.
// Relationship name plurality is important, if your relationship is
// singular, you need to specify the singular form and vice versa.
func Load(relationships ...string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendLoad(q, relationships...)
	}
}

// InnerJoin on another table
func InnerJoin(clause string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.AppendInnerJoin(q, clause, args...)
	}
}

// Select specific columns opposed to all columns
func Select(columns ...string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendSelect(q, columns...)
	}
}

// Where allows you to specify a where clause for your statement
func Where(clause string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.AppendWhere(q, clause, args...)
	}
}

// And allows you to specify a where clause separated by an AND for your statement
// And is a duplicate of the Where function, but allows for more natural looking
// query mod chains, for example: (Where("a=?"), And("b=?"), Or("c=?")))
func And(clause string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.AppendWhere(q, clause, args...)
	}
}

// Or allows you to specify a where clause separated by an OR for your statement
func Or(clause string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.AppendWhere(q, clause, args...)
		queries.SetLastWhereAsOr(q)
	}
}

// WhereIn allows you to specify a "x IN (set)" clause for your where statement
// Example clauses: "column in ?", "(column1,column2) in ?"
func WhereIn(clause string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.AppendIn(q, clause, args...)
	}
}

// AndIn allows you to specify a "x IN (set)" clause separated by an AndIn
// for your where statement. AndIn is a duplicate of the WhereIn function, but
// allows for more natural looking query mod chains, for example:
// (WhereIn("column1 in ?"), AndIn("column2 in ?"), OrIn("column3 in ?"))
func AndIn(clause string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.AppendIn(q, clause, args...)
	}
}

// OrIn allows you to specify an IN clause separated by
// an OR for your where statement
func OrIn(clause string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.AppendIn(q, clause, args...)
		queries.SetLastInAsOr(q)
	}
}

// GroupBy allows you to specify a group by clause for your statement
func GroupBy(clause string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendGroupBy(q, clause)
	}
}

// OrderBy allows you to specify a order by clause for your statement
func OrderBy(clause string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendOrderBy(q, clause)
	}
}

// Having allows you to specify a having clause for your statement
func Having(clause string, args ...interface{}) QueryMod {
	return func(q *queries.Query) {
		queries.AppendHaving(q, clause, args...)
	}
}

// From allows to specify the table for your statement
func From(from string) QueryMod {
	return func(q *queries.Query) {
		queries.AppendFrom(q, from)
	}
}

// Limit the number of returned rows
func Limit(limit int) QueryMod {
	return func(q *queries.Query) {
		queries.SetLimit(q, limit)
	}
}

// Offset into the results
func Offset(offset int) QueryMod {
	return func(q *queries.Query) {
		queries.SetOffset(q, offset)
	}
}

// For inserts a concurrency locking clause at the end of your statement
func For(clause string) QueryMod {
	return func(q *queries.Query) {
		queries.SetFor(q, clause)
	}
}
