{{- if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// {{$tableNameSingular}}Find retrieves a single record by ID.
func {{$tableNameSingular}}Find({{primaryKeyFuncSig .Table.Columns .Table.PKey.Columns}}, columns ...string) (*{{$tableNameSingular}}, error) {
  return {{$tableNameSingular}}FindX(boil.GetDB(), {{camelCaseCommaList .Table.PKey.Columns}}, columns...)
}

func {{$tableNameSingular}}FindX(exec boil.Executor, {{primaryKeyFuncSig .Table.Columns .Table.PKey.Columns}}, columns ...string) (*{{$tableNameSingular}}, error) {
  var {{$varNameSingular}} *{{$tableNameSingular}}
  mods := []qs.QueryMod{
    qs.Select(columns...),
    qs.From("{{.Table.Name}}"),
    qs.Where("{{wherePrimaryKey .Table.PKey.Columns 1}}", {{camelCaseCommaList .Table.PKey.Columns}}),
  }

  q := NewQueryX(exec, mods...)

  err := boil.ExecQueryOne(q).Scan(boil.GetStructPointers(&{{$varNameSingular}}, columns...)...)

  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %s", err)
  }

  return {{$varNameSingular}}, nil
}
{{- end -}}
