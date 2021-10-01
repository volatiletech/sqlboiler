{{- $alias := .Aliases.Table .Table.Name}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted }}
// {{$alias.UpPlural}} retrieves all the records using an executor.
func {{$alias.UpPlural}}(mods ...qm.QueryMod) {{$alias.DownSingular}}Query {
    {{if and .AddSoftDeletes $canSoftDelete -}}
    mods = append(mods, qm.From("{{$schemaTable}}"), qmhelper.WhereIsNull("{{$schemaTable}}.{{"deleted_at" | $.Quotes}}"))
    {{else -}}
	mods = append(mods, qm.From("{{$schemaTable}}"))
	{{end -}}
	return {{$alias.DownSingular}}Query{NewQuery(mods...)}
}
