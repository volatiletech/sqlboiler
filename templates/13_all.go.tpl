{{- $alias := .Aliases.Table .Table.Name}}
// {{$alias.UpPlural}} retrieves all the records using an executor.
func {{$alias.UpPlural}}(mods ...qm.QueryMod) {{$alias.DownSingular}}Query {
	mods = append(mods, qm.From("{{.Table.Name | .SchemaTable}}"))
	return {{$alias.DownSingular}}Query{NewQuery(mods...)}
}
