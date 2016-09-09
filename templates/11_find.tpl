{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $dbName := singular .Table.Name -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", " -}}
// {{$tableNameSingular}}FindG retrieves a single record by ID.
func Find{{$tableNameSingular}}G({{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
  return Find{{$tableNameSingular}}(boil.GetDB(), {{$pkNames | join ", "}}, selectCols...)
}

// {{$tableNameSingular}}FindGP retrieves a single record by ID, and panics on error.
func Find{{$tableNameSingular}}GP({{$pkArgs}}, selectCols ...string) *{{$tableNameSingular}} {
  retobj, err := Find{{$tableNameSingular}}(boil.GetDB(), {{$pkNames | join ", "}}, selectCols...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return retobj
}

// {{$tableNameSingular}}Find retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func Find{{$tableNameSingular}}(exec boil.Executor, {{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
  {{$varNameSingular}}Obj := &{{$tableNameSingular}}{}

  sel := "*"
  if len(selectCols) > 0 {
    sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
  }
  query := fmt.Sprintf(
    `select %s from {{schemaTable .Dialect.LQ .Dialect.RQ .DriverName .Schema .Table.Name}} where {{whereClause 1 .Table.PKey.Columns}}`, sel,
  )

  q := boil.SQL(exec, query, {{$pkNames | join ", "}})

  err := q.Bind({{$varNameSingular}}Obj)
  if err != nil {
    if errors.Cause(err) == sql.ErrNoRows {
      return nil, sql.ErrNoRows
    }
    return nil, errors.Wrap(err, "{{.PkgName}}: unable to select from {{.Table.Name}}")
  }

  return {{$varNameSingular}}Obj, nil
}

// {{$tableNameSingular}}FindP retrieves a single record by ID with an executor, and panics on error.
func Find{{$tableNameSingular}}P(exec boil.Executor, {{$pkArgs}}, selectCols ...string) *{{$tableNameSingular}} {
  retobj, err := Find{{$tableNameSingular}}(exec, {{$pkNames | join ", "}}, selectCols...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return retobj
}
