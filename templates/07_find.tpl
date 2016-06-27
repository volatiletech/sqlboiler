{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", "}}
// {{$tableNameSingular}}Find retrieves a single record by ID.
func {{$tableNameSingular}}Find({{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
  return {{$tableNameSingular}}FindX(boil.GetDB(), {{$pkNames | join ", "}}, selectCols...)
}

func {{$tableNameSingular}}FindX(exec boil.Executor, {{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
  {{$varNameSingular}} := &{{$tableNameSingular}}{}

  mods := []qm.QueryMod{
    qm.Select(selectCols...),
    qm.Table("{{.Table.Name}}"),
    qm.Where(`{{whereClause .Table.PKey.Columns 1}}`, {{$pkNames | join ", "}}),
  }

  q := NewQueryX(exec, mods...)

  err := boil.ExecQueryOne(q).Scan(boil.GetStructPointers({{$varNameSingular}}, selectCols...)...)

  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %v", err)
  }

  return {{$varNameSingular}}, nil
}
