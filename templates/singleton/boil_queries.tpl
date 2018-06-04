var dialect = drivers.Dialect{
	LQ: 0x{{printf "%x" .Dialect.LQ}},
	RQ: 0x{{printf "%x" .Dialect.RQ}},

	UseIndexPlaceholders:    {{.Dialect.UseIndexPlaceholders}},
	UseLastInsertID:         {{.Dialect.UseLastInsertID}},
	UseSchema:               {{.Dialect.UseSchema}},
	UseDefaultKeyword:       {{.Dialect.UseDefaultKeyword}},
	UseAutoColumns:          {{.Dialect.UseAutoColumns}},
	UseTopClause:            {{.Dialect.UseTopClause}},
	UseOutputClause:         {{.Dialect.UseOutputClause}},
	UseCaseWhenExistsClause: {{.Dialect.UseCaseWhenExistsClause}},
}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(mods ...qm.QueryMod) *queries.Query {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q, mods...)

	return q
}
