{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// One returns a single {{$varNameSingular}} record from the query.
func (q {{$varNameSingular}}Query) One() (*{{$tableNameSingular}}, error) {
  o := &{{$tableNameSingular}}{}

  boil.SetLimit(q.Query, 1)

  err := q.Bind(o)
  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: failed to execute a one query for {{.Table.Name}}: %s", err)
  }

  return o, nil
}

// OneP returns a single {{$varNameSingular}} record from the query, and panics on error.
func (q {{$varNameSingular}}Query) OneP() (*{{$tableNameSingular}}) {
  o, err := q.One()
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return o
}

// All returns all {{$tableNameSingular}} records from the query.
func (q {{$varNameSingular}}Query) All() ({{$tableNameSingular}}Slice, error) {
  var o {{$tableNameSingular}}Slice

  err := q.Bind(&o)
  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: failed to assign all query results to {{$tableNameSingular}} slice: %s", err)
  }

  return o, nil
}

// AllP returns all {{$tableNameSingular}} records from the query, and panics on error.
func (q {{$varNameSingular}}Query) AllP() {{$tableNameSingular}}Slice {
    o, err := q.All()
    if err != nil {
      panic(boil.WrapErr(err))
    }

    return o
}

// Count returns the count of all {{$tableNameSingular}} records in the query.
func (q {{$varNameSingular}}Query) Count() (int64, error) {
  var count int64

  boil.SetCount(q.Query)

  err := boil.ExecQueryOne(q.Query).Scan(&count)
  if err != nil {
    return 0, fmt.Errorf("{{.PkgName}}: failed to count {{.Table.Name}} rows: %s", err)
  }

  return count, nil
}

// CountP returns the count of all {{$tableNameSingular}} records in the query, and panics on error.
func (q {{$varNameSingular}}Query) CountP() int64 {
  c, err := q.Count()
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return c
}

// Exists checks if the row exists in the table.
func (q {{$varNameSingular}}Query) Exists() (bool, error) {
  var count int64

  boil.SetCount(q.Query)
  boil.SetLimit(q.Query, 1)

  err := boil.ExecQueryOne(q.Query).Scan(&count)
  if err != nil {
    return false, fmt.Errorf("{{.PkgName}}: failed to check if {{.Table.Name}} exists: %s", err)
  }

  return count > 0, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q {{$varNameSingular}}Query) ExistsP() bool {
  e, err := q.Exists()
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return e
}
