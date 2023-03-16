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

{{- if not .AutoColumns.Deleted }}
    // This is a dummy variable to prevent unused regexp import error
    var _ = &regexp.Regexp{}
{{- end }}

{{- if and (.AutoColumns.Deleted) (ne $.AutoColumns.Deleted "deleted_at") }}
    func init() {
        queries.SetRemoveSoftDeleteRgx(regexp.MustCompile("{{$.AutoColumns.Deleted}}[\"'`]? is null"))
    }
{{- end }}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(mods ...qm.QueryMod) *queries.Query {
	q := &queries.Query{}
	queries.SetDialect(q, &dialect)
	qm.Apply(q, mods...)

	return q
}
