{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase}}
// {{$tableNamePlural}}G retrieves all records.
func {{$tableNamePlural}}G(mods ...qm.QueryMod) {{$tableNameSingular}}Query {
	return {{$tableNamePlural}}(boil.GetDB(), mods...)
}

// {{$tableNamePlural}} retrieves all the records using an executor.
func {{$tableNamePlural}}(exec boil.Executor, mods ...qm.QueryMod) {{$tableNameSingular}}Query {
	mods = append(mods, qm.From("{{.Table.Name | .SchemaTable}}"))
	return {{$tableNameSingular}}Query{NewQuery(exec, mods...)}
}
