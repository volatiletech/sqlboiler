package qs

import "github.com/nullbio/sqlboiler/boil"

type QueryMod func(q *boil.Query)

func Apply(q *boil.Query, mods ...QueryMod) {
	for _, mod := range mods {
		mod(q)
	}
}

func Or(whereMods ...QueryMod) QueryMod {
	return func(q *boil.Query) {
		if len(whereMods) < 2 {
			// error, needs to be at least 2 for an or
		}
		// add the where mods to query with or seperators
	}
}

func Limit(limit int) QueryMod {
	return func(q *boil.Query) {
		boil.SetLimit(q, limit)
	}
}

func InnerJoin(on string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetInnerJoin(q, on, args...)
	}
}

func OuterJoin(on string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetOuterJoin(q, on, args...)
	}
}

func LeftOuterJoin(on string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetLeftOuterJoin(q, on, args...)
	}
}

func RightOuterJoin(on string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetRightOuterJoin(q, on, args...)
	}
}

func Select(columns ...string) QueryMod {
	return func(q *boil.Query) {
		boil.SetSelect(q, columns...)
	}
}

func Where(clause string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		boil.SetWhere(q, clause, args...)
	}
}

func GroupBy(clause string) QueryMod {
	return func(q *boil.Query) {
		boil.SetGroupBy(q, clause)
	}
}

func OrderBy(clause string) QueryMod {
	return func(q *boil.Query) {
		boil.SetOrderBy(q, clause)
	}
}

func Having(clause string) QueryMod {
	return func(q *boil.Query) {
		boil.SetHaving(q, clause)
	}
}

func Table(table string) QueryMod {
	return func(q *boil.Query) {
		boil.SetTable(q, table)
	}
}
