{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $tableNamePlural := titleCasePlural .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
{{- $varNamePlural := camelCasePlural .Table.Name -}}
type {{$varNameSingular}}Query struct {
  *boil.Query
}

// {{$tableNamePlural}}All retrieves all records.
func {{$tableNamePlural}}(mods ...QueryMod) {{$varNameSingular}}Query {
  return {{$tableNamePlural}}X(boil.GetDB(), mods...)
}

func {{$tableNamePlural}}X(exec boil.Executor, mods ...QueryMod) {{$tableNameSingular}}Query {
  mods = append(mods, boil.From("{{.Table.Name}}"))
  return NewQueryX(exec, mods...)
}
