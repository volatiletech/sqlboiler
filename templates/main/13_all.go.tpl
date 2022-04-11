{{- $alias := .Aliases.Table .Table.Name}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted }}
{{- $soft := and .AddSoftDeletes $canSoftDelete }}
// {{$alias.UpPlural}} retrieves all the records using an executor.
func {{$alias.UpPlural}}(mods ...qm.QueryMod) {{$alias.DownSingular}}Query {
    {{if and .AddSoftDeletes $canSoftDelete -}}
    mods = append(mods, qm.From("{{$schemaTable}}"), qmhelper.WhereIsNull("{{$schemaTable}}.{{"deleted_at" | $.Quotes}}"))
    {{else -}}
	mods = append(mods, qm.From("{{$schemaTable}}"))
	{{end -}}
  query := NewQuery(mods...)

  // set all the queries to point to the same value
  tableQuery := {{$alias.DownSingular}}Query{Query: query}
  tableQuery.BaseQuery.Query = query
  tableQuery.SelectQuery.Query = query
  tableQuery.DeleteQuery.Query = query

  return tableQuery
}
