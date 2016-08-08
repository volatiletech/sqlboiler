{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", "}}
// {{$tableNameSingular}}FindG retrieves a single record by ID.
func {{$tableNameSingular}}FindG({{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
  return {{$tableNameSingular}}Find(boil.GetDB(), {{$pkNames | join ", "}}, selectCols...)
}

// {{$tableNameSingular}}FindGP retrieves a single record by ID, and panics on error.
func {{$tableNameSingular}}FindGP({{$pkArgs}}, selectCols ...string) *{{$tableNameSingular}} {
  retobj, err := {{$tableNameSingular}}Find(boil.GetDB(), {{$pkNames | join ", "}}, selectCols...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return retobj
}

// {{$tableNameSingular}}Find retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func {{$tableNameSingular}}Find(exec boil.Executor, {{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
  {{$varNameSingular}} := &{{$tableNameSingular}}{}

  mods := []qm.QueryMod{
    qm.Select(selectCols...),
    qm.From("{{.Table.Name}}"),
    qm.Where(`{{whereClause .Table.PKey.Columns 1}}`, {{$pkNames | join ", "}}),
  }

  q := NewQuery(exec, mods...)

  err := q.Bind({{$varNameSingular}})
  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %v", err)
  }

  return {{$varNameSingular}}, nil
}

// {{$tableNameSingular}}FindP retrieves a single record by ID with an executor, and panics on error.
func {{$tableNameSingular}}FindP(exec boil.Executor, {{$pkArgs}}, selectCols ...string) *{{$tableNameSingular}} {
  retobj, err := {{$tableNameSingular}}Find(exec, {{$pkNames | join ", "}}, selectCols...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return retobj
}
