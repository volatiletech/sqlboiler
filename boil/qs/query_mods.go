package qs

import "github.com/pobri19/sqlboiler/boil"

type QueryMod func(q *boil.Query)

func (q *boil.Query) Apply(mods ...QueryMod) {
	for _, mod := range mods {
		mod(q)
	}
}

func Limit(limit int) QueryMod {
	return func(q *boil.Query) {
		q.limit = limit
	}
}

func Join(join string) QueryMod {
	return func(q *boil.Query) {
		q.joins = append(q.joins, join)
	}
}

func Select(columns ...string) QueryMod {
	return func(q *boil.Query) {
		q.selectCols = append(q.selectCols, columns...)
	}
}

func Where(clause string, args ...interface{}) QueryMod {
	return func(q *boil.Query) {
		w := where{
			clause: clause,
			args:   args,
		}

		q.where = append(q.where, w)
	}
}

func GroupBy(clause string) QueryMod {
	return func(q *boil.Query) {
		q.groupBy = append(q.groupBy, clause)
	}
}

func OrderBy(clause string) QueryMod {
	return func(q *boil.Query) {
		q.orderBy = append(q.orderBy, clause)
	}
}

func Having(clause string) QueryMod {
	return func(q *boil.Query) {
		q.having = append(q.having, clause)
	}
}

func From(table string) QueryMod {
	return func(q *boil.Query) {
		q.from = table
	}
}
