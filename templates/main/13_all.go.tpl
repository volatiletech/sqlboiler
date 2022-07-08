{{- $alias := .Aliases.Table .Table.Name}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted }}
// {{$alias.UpPlural}} retrieves all the records using an executor.
func {{$alias.UpPlural}}(mods ...qm.QueryMod) {{$alias.DownSingular}}Query {
    {{if and .AddSoftDeletes $canSoftDelete -}}
    mods = append(mods, qm.From("{{$schemaTable}}"), qmhelper.WhereIsNull("{{$schemaTable}}.{{or $.AutoColumns.Deleted "deleted_at" | $.Quotes}}"))
    {{else -}}
    mods = append(mods, qm.From("{{$schemaTable}}"))
    {{end -}}

    q := NewQuery(mods...)
    if len(queries.GetSelect(q)) == 0 {
        queries.SetSelect(q, []string{"{{$schemaTable}}.*"})
    }

    return {{$alias.DownSingular}}Query{q}
}
