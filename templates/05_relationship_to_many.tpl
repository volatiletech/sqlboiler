{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $rel := txtsFromToMany $dot.Tables $table . -}}
		{{- $schemaForeignTable := .ForeignTable | $dot.SchemaTable -}}
// {{$rel.Function.Name}}G retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}}
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}G(mods ...qm.QueryMod) {{$varNameSingular}}Query {
	return {{$rel.Function.Receiver}}.{{$rel.Function.Name}}(boil.GetDB(), mods...)
}

// {{$rel.Function.Name}} retrieves all the {{$rel.LocalTable.NameSingular}}'s {{$rel.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $rel.Function.Name $rel.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) {{$rel.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) {{$varNameSingular}}Query {
	queryMods := []qm.QueryMod{
		qm.Select("{{id 0 | $dot.Quotes}}.*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if .ToJoinTable -}}
	queryMods = append(queryMods,
		qm.InnerJoin("{{.JoinTable | $dot.SchemaTable}} as {{id 1 | $dot.Quotes}} on {{id 0 | $dot.Quotes}}.{{.ForeignColumn | $dot.Quotes}} = {{id 1 | $dot.Quotes}}.{{.JoinForeignColumn | $dot.Quotes}}"),
		qm.Where("{{id 1 | $dot.Quotes}}.{{.JoinLocalColumn | $dot.Quotes}}={{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}", {{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}),
	)
		{{else -}}
	queryMods = append(queryMods,
		qm.Where("{{id 0 | $dot.Quotes}}.{{.ForeignColumn | $dot.Quotes}}={{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}", {{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}),
	)
		{{end}}

	query := {{$rel.ForeignTable.NamePluralGo}}(exec, queryMods...)
	queries.SetFrom(query.Query, "{{$schemaForeignTable}} as {{id 0 | $dot.Quotes}}")
	return query
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinTable */ -}}
