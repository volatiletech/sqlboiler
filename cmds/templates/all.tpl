{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
type {{$varNameSingular}}Query struct {
  *boil.Query
}

// {{$tableNamePlural}}All retrieves all records.
func {{$tableNamePlural}}(mods ...qs.QueryMod) {{$varNameSingular}}Query {
  return {{$tableNamePlural}}X(boil.GetDB(), mods...)
}

func {{$tableNamePlural}}X(exec boil.Executor, mods ...qs.QueryMod) {{$varNameSingular}}Query {
  mods = append(mods, qs.From("{{.Table.Name}}"))
  return {{$varNameSingular}}Query{NewQueryX(exec, mods...)}
}
