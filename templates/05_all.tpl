{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// {{$tableNamePlural}}All retrieves all records.
func {{$tableNamePlural}}(mods ...qm.QueryMod) {{$varNameSingular}}Query {
  return {{$tableNamePlural}}X(boil.GetDB(), mods...)
}

// {{$tableNamePlural}}X retrieves all the records using an executor.
func {{$tableNamePlural}}X(exec boil.Executor, mods ...qm.QueryMod) {{$varNameSingular}}Query {
  mods = append(mods, qm.Table("{{.Table.Name}}"))
  return {{$varNameSingular}}Query{NewQueryX(exec, mods...)}
}
