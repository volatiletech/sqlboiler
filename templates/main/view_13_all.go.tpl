{{- $alias := .Aliases.View .View.Name}}
{{- $schemaView := .View.Name | .SchemaTable}}
{{- $canSoftDelete := .View.CanSoftDelete $.AutoColumns.Deleted }}
// {{$alias.UpPlural}} retrieves all the records using an executor.
func {{$alias.UpPlural}}(mods ...qm.QueryMod) {{$alias.DownSingular}}Query {
    {{if and .AddSoftDeletes $canSoftDelete -}}
    mods = append(mods, qm.From("{{$schemaView}}"), qmhelper.WhereIsNull("{{$schemaView}}.{{"deleted_at" | $.Quotes}}"))
    {{else -}}
	mods = append(mods, qm.From("{{$schemaView}}"))
	{{end -}}
	return {{$alias.DownSingular}}Query{NewQuery(mods...)}
}
