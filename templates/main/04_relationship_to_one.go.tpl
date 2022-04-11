{{- if or .Table.IsJoinTable .Table.IsView -}}
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
	}

	queryMods = append(queryMods, mods...)

	return {{$ftable.UpPlural}}(queryMods...)
}
{{- end -}}
{{- end -}}
