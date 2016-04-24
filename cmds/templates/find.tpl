{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// {{$tableNameSingular}}Find retrieves a single record by ID.
func {{$tableNameSingular}}Find(id int64, columns ...string) (*{{$tableNameSingular}}, error) {
  return {{$tableNameSingular}}FindX(boil.GetDB(), id, columns...)
}

func {{$tableNameSingular}}FindX(exec boil.Executor, id int64, columns ...string) (*{{$tableNameSingular}}, error) {
  if id == 0 {
    return nil, errors.New("{{.PkgName}}: no id provided for {{.Table.Name}} select")
  }

  var {{$varNameSingular}} *{{$tableNameSingular}}
  mods := []qs.QueryMod{
    qs.Select(columns...),
    qs.From("{{.Table.Name}}"),
    qs.Where("id=$1", id),
  }

  q := NewQueryX(exec, mods...)

  err := boil.ExecQueryOne(q).Scan(
  )

    //GetStructPointers({{$varNameSingular}}, columnsthings)
  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %s", err)
  }

  return {{$varNameSingular}}, nil
}
