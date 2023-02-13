{{- if or .Table.IsJoinTable .Table.IsView -}}
{{- else -}}
	{{- range $rel := .Table.ToManyRelationships -}}
		{{- $ltable := $.Aliases.Table $rel.Table -}}
		{{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
		{{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
		{{- $schemaForeignTable := .ForeignTable | $.SchemaTable -}}
		{{- $foreignTable := .ForeignTable }}
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
	)
		{{end}}

	return {{$ftable.UpPlural}}(queryMods...)
}

// {{$relAlias.Local}}WithSchema retrieves all the {{.ForeignTable | singular}}'s {{$ftable.UpPlural}} with an executor
{{- if not (eq $relAlias.Local $ftable.UpPlural)}} via {{$rel.ForeignColumn}} column{{- end}}.
func (o *{{$ltable.UpSingular}}) {{$relAlias.Local}}WithSchema(schema string, mods ...qm.QueryMod) {{$ftable.DownSingular}}Query {
	schemaForeignTable := fmt.Sprintf("%s.{{$foreignTable}}", schema)
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if $rel.ToJoinTable -}}
        schemaJoinTable := fmt.Sprintf("%s.{{$rel.JoinTable}}", schema)
	queryMods = append(queryMods,
		qm.InnerJoin(fmt.Sprintf("%s on %s.{{$rel.ForeignColumn}} = %s.{{$rel.JoinForeignColumn}}", schemaJoinTable, schemaForeignTable, schemaJoinTable)),
		qm.Where(fmt.Sprintf("%s.{{$rel.JoinLocalColumn}}=?", schemaJoinTable), o.{{$ltable.Column $rel.Column}}),
	)
		{{else -}}
	queryMods = append(queryMods,
		qm.Where(fmt.Sprintf("%s.{{$rel.ForeignColumn}}=?", schemaForeignTable), o.{{$ltable.Column $rel.Column}}),
	)
		{{end}}

	return {{$ftable.UpPlural}}WithSchema(schema, queryMods...)
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinTable */ -}}
