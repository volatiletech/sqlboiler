{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $dbName := singular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// {{$tableNameSingular}}Find retrieves a single record by ID.
func {{$tableNameSingular}}Find(id int) (*{{$tableNameSingular}}, error) {
  if id == 0 {
    return nil, errors.New("{{.PkgName}}: no id provided for {{.Table.Name}} select")
  }
  var {{$varNameSingular}} *{{$tableNameSingular}}
  err := boil.GetDB().Select(&{{$varNameSingular}}, `SELECT {{selectParamNames $dbName .Table.Columns}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table.Name}}: %s", err)
  }

  return {{$varNameSingular}}, nil
}
