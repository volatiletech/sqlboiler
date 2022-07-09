{{- if or .Table.IsJoinTable .Table.IsView -}}
{{- else -}}
	{{- range $rel := .Table.ToOneRelationships -}}
		{{- $ltable := $.Aliases.Table $rel.Table -}}
		{{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
		{{- $relAlias := $ftable.Relationship $rel.Name -}}
		{{- $canSoftDelete := (getTable $.Tables $rel.ForeignTable).CanSoftDelete $.AutoColumns.Deleted }}
// {{$relAlias.Local}} pointed to by the foreign key.
func (o *{{$ltable.UpSingular}}) {{$relAlias.Local}}(mods ...qm.QueryMod) ({{$ftable.DownSingular}}Query) {
	queryMods := []qm.QueryMod{
		qm.Where("{{$rel.ForeignColumn | $.Quotes}} = ?", o.{{$ltable.Column $rel.Column}}),
	}

	queryMods = append(queryMods, mods...)

	return {{$ftable.UpPlural}}(queryMods...)
}
{{- end -}}
{{- end -}}
