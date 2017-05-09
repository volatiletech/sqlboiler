{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table . -}}
		{{- $tableNameSingular := .ForeignTable | singular | titleCase}}
// {{$txt.Function.Name}}G pointed to by the foreign key.
func (o *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}G(mods ...qm.QueryMod) {{$tableNameSingular}}Query {
	return o.{{$txt.Function.Name}}(boil.GetDB(), mods...)
}

// {{$txt.Function.Name}} pointed to by the foreign key.
func (o *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) ({{$tableNameSingular}}Query) {
	queryMods := []qm.QueryMod{
		qm.Where("{{$txt.ForeignTable.ColumnName}}=?", o.{{$txt.LocalTable.ColumnNameGo}}),
	}

	queryMods = append(queryMods, mods...)

	query := {{$txt.ForeignTable.NamePluralGo}}(exec, queryMods...)
	queries.SetFrom(query.Query, "{{.ForeignTable | $dot.SchemaTable}}")

	return query
}
{{- end -}}
{{- end -}}
