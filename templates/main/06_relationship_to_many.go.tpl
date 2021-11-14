{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range $rel := .Table.ToManyRelationships -}}
		{{- $ltable := $.Aliases.Table $rel.Table -}}
		{{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
		{{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
		{{- $schemaForeignTable := .ForeignTable | $.SchemaTable -}}
		{{- $canSoftDelete := (getTable $.Tables .ForeignTable).CanSoftDelete $.AutoColumns.Deleted}}
// {{$relAlias.Local}} retrieves all the {{.ForeignTable | singular}}'s {{$ftable.UpPlural}} with an executor
{{- if not (eq $relAlias.Local $ftable.UpPlural)}} via {{$rel.ForeignColumn}} column{{- end}}.
func (o *{{$ltable.UpSingular}}) {{$relAlias.Local}}(mods ...qm.QueryMod) {{$ftable.DownSingular}}Query {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if $rel.ToJoinTable -}}
	queryMods = append(queryMods,
		{{$schemaJoinTable := $rel.JoinTable | $.SchemaTable -}}
		qm.InnerJoin("{{$schemaJoinTable}} on {{$schemaForeignTable}}.{{$rel.ForeignColumn | $.Quotes}} = {{$schemaJoinTable}}.{{$rel.JoinForeignColumn | $.Quotes}}"),
		qm.Where("{{$schemaJoinTable}}.{{$rel.JoinLocalColumn | $.Quotes}}=?", o.{{$ltable.Column $rel.Column}}),
	)
		{{else -}}
	queryMods = append(queryMods,
		qm.Where("{{$schemaForeignTable}}.{{$rel.ForeignColumn | $.Quotes}}=?", o.{{$ltable.Column $rel.Column}}),
		{{if and $.AddSoftDeletes $canSoftDelete -}}
		qmhelper.WhereIsNull("{{$schemaForeignTable}}.{{"deleted_at" | $.Quotes}}"),
		{{- end}}
	)
		{{end}}

	query := {{$ftable.UpPlural}}(queryMods...)
	queries.SetFrom(query.Query, "{{$schemaForeignTable}}")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"{{$schemaForeignTable}}.*"})
	}

	return query
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinTable */ -}}
