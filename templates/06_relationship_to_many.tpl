{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $txt := txtsFromToMany $.Tables $table . -}}
		{{- $schemaForeignTable := .ForeignTable | $.SchemaTable}}
// {{$txt.Function.Name}} retrieves all the {{.ForeignTable | singular}}'s {{$txt.ForeignTable.NameHumanReadable}} with an executor
{{- if not (eq $txt.Function.Name $txt.ForeignTable.NamePluralGo)}} via {{.ForeignColumn}} column{{- end}}.
func (o *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}(mods ...qm.QueryMod) {{$varNameSingular}}Query {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if .ToJoinTable -}}
	queryMods = append(queryMods,
		{{$schemaJoinTable := .JoinTable | $.SchemaTable -}}
		qm.InnerJoin("{{$schemaJoinTable}} on {{$schemaForeignTable}}.{{.ForeignColumn | $.Quotes}} = {{$schemaJoinTable}}.{{.JoinForeignColumn | $.Quotes}}"),
		qm.Where("{{$schemaJoinTable}}.{{.JoinLocalColumn | $.Quotes}}=?", o.{{$txt.LocalTable.ColumnNameGo}}),
	)
		{{else -}}
	queryMods = append(queryMods,
		qm.Where("{{$schemaForeignTable}}.{{.ForeignColumn | $.Quotes}}=?", o.{{$txt.LocalTable.ColumnNameGo}}),
	)
		{{end}}

	query := {{$txt.ForeignTable.NamePluralGo}}(queryMods...)
	queries.SetFrom(query.Query, "{{$schemaForeignTable}}")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"{{$schemaForeignTable}}.*"})
	}

	return query
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinTable */ -}}
