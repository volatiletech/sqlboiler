{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $tableNameSingular := .ForeignTable | singular | titleCase -}}
		{{- $txt := txtsFromToMany $dot.Tables $table . -}}
		{{- $schemaForeignTable := .ForeignTable | $dot.SchemaTable}}
// {{$txt.Function.Name}}G retrieves all the {{.ForeignTable | singular}}'s {{$txt.ForeignTable.NameHumanReadable}}
{{- if not (eq $txt.Function.Name $txt.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func (o *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}G(mods ...qm.QueryMod) {{$tableNameSingular}}Query {
	return o.{{$txt.Function.Name}}(boil.GetDB(), mods...)
}

// {{$txt.Function.Name}} retrieves all the {{.ForeignTable | singular}}'s {{$txt.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $txt.Function.Name $txt.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func (o *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) {{$tableNameSingular}}Query {
	queryMods := []qm.QueryMod{
		qm.Select("{{$schemaForeignTable}}.*"),
	}

	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if .ToJoinTable -}}
	queryMods = append(queryMods,
		{{$schemaJoinTable := .JoinTable | $.SchemaTable -}}
		qm.InnerJoin("{{$schemaJoinTable}} on {{$schemaForeignTable}}.{{.ForeignColumn | $dot.Quotes}} = {{$schemaJoinTable}}.{{.JoinForeignColumn | $dot.Quotes}}"),
		qm.Where("{{$schemaJoinTable}}.{{.JoinLocalColumn | $dot.Quotes}}=?", o.{{$txt.LocalTable.ColumnNameGo}}),
	)
		{{else -}}
	queryMods = append(queryMods,
		qm.Where("{{$schemaForeignTable}}.{{.ForeignColumn | $dot.Quotes}}=?", o.{{$txt.LocalTable.ColumnNameGo}}),
	)
		{{end}}

	query := {{$txt.ForeignTable.NamePluralGo}}(exec, queryMods...)
	queries.SetFrom(query.Query, "{{$schemaForeignTable}}")
	return query
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinTable */ -}}
