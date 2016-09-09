{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// {{$tableNamePlural}}G retrieves all records.
func {{$tableNamePlural}}G(mods ...qm.QueryMod) {{$varNameSingular}}Query {
  return {{$tableNamePlural}}(boil.GetDB(), mods...)
}

// {{$tableNamePlural}} retrieves all the records using an executor.
func {{$tableNamePlural}}(exec boil.Executor, mods ...qm.QueryMod) {{$varNameSingular}}Query {
  mods = append(mods, qm.From(`{{schemaTable .Dialect.LQ .Dialect.RQ .DriverName .Schema .Table.Name}}`))
  return {{$varNameSingular}}Query{NewQuery(exec, mods...)}
}
