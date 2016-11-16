{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $txt := txtsFromToMany $dot.Tables $table . -}}
		{{- $schemaForeignTable := .ForeignTable | $dot.SchemaTable}}
// {{$txt.Function.Name}}G retrieves all the {{.ForeignTable | singular}}'s {{$txt.ForeignTable.NameHumanReadable}}
{{- if not (eq $txt.Function.Name $txt.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func (o *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}G(mods ...qm.QueryMod) {{$varNameSingular}}Query {
	return o.{{$txt.Function.Name}}(boil.GetDB(), mods...)
}

// {{$txt.Function.Name}} retrieves all the {{.ForeignTable | singular}}'s {{$txt.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $txt.Function.Name $txt.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func (o *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) {{$varNameSingular}}Query {
	queryMods := []qm.QueryMod{
		qm.Select("{{id 0 | $dot.Quotes}}.*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if .ToJoinTable -}}
	queryMods = append(queryMods,
		qm.InnerJoin("{{.JoinTable | $dot.SchemaTable}} as {{id 1 | $dot.Quotes}} on {{id 0 | $dot.Quotes}}.{{.ForeignColumn | $dot.Quotes}} = {{id 1 | $dot.Quotes}}.{{.JoinForeignColumn | $dot.Quotes}}"),
		qm.Where("{{id 1 | $dot.Quotes}}.{{.JoinLocalColumn | $dot.Quotes}}=?", o.{{$txt.LocalTable.ColumnNameGo}}),
	)
		{{else -}}
	queryMods = append(queryMods,
		qm.Where("{{id 0 | $dot.Quotes}}.{{.ForeignColumn | $dot.Quotes}}=?", o.{{$txt.LocalTable.ColumnNameGo}}),
	)
		{{end}}

	query := {{$txt.ForeignTable.NamePluralGo}}(exec, queryMods...)
	queries.SetFrom(query.Query, "{{$schemaForeignTable}} as {{id 0 | $dot.Quotes}}")
	return query
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinTable */ -}}
