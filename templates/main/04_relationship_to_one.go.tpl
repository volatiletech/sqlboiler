{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range $fkey := .Table.FKeys -}}
		{{- $ltable := $.Aliases.Table $fkey.Table -}}
		{{- $ftable := $.Aliases.Table $fkey.ForeignTable -}}
		{{- $rel := $ltable.Relationship $fkey.Name -}}
		{{- $canSoftDelete := (getTable $.Tables $fkey.ForeignTable).CanSoftDelete $.AutoColumns.Deleted }}
// {{$rel.Foreign}} pointed to by the foreign key.
func (o *{{$ltable.UpSingular}}) {{$rel.Foreign}}(mods ...qm.QueryMod) ({{$ftable.DownSingular}}Query) {
	queryMods := []qm.QueryMod{
		qm.Where("{{$fkey.ForeignColumn | $.Quotes}} = ?", o.{{$ltable.Column $fkey.Column}}),
		{{if and $.AddSoftDeletes $canSoftDelete -}}
		qmhelper.WhereIsNull("deleted_at"),
		{{- end}}
	}

	queryMods = append(queryMods, mods...)

	query := {{$ftable.UpPlural}}(queryMods...)
	queries.SetFrom(query.Query, "{{.ForeignTable | $.SchemaTable}}")

	return query
}
{{- end -}}
{{- end -}}
