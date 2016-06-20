{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// {{$tableNameSingular}}Find retrieves a single record by ID.
func {{$tableNameSingular}}Find({{primaryKeyFuncSig .Table.Columns .Table.PKey.Columns}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
  return {{$tableNameSingular}}FindX(boil.GetDB(), {{.Table.PKey.Columns | stringMap .StringFuncs.camelCase | join ", "}}, selectCols...)
}

func {{$tableNameSingular}}FindX(exec boil.Executor, {{primaryKeyFuncSig .Table.Columns .Table.PKey.Columns}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
  {{$varNameSingular}} := &{{$tableNameSingular}}{}

  mods := []qm.QueryMod{
    qm.Select(selectCols...),
    qm.Table("{{.Table.Name}}"),
    qm.Where("{{wherePrimaryKey .Table.PKey.Columns 1}}", {{.Table.PKey.Columns | stringMap .StringFuncs.camelCase | join ", "}}),
  }

  q := NewQueryX(exec, mods...)

  err := boil.ExecQueryOne(q).Scan(boil.GetStructPointers({{$varNameSingular}}, selectCols...)...)

  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %v", err)
  }

  return {{$varNameSingular}}, nil
}
