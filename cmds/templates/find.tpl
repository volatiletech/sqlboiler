{{- $tableNameSingular := titleCaseSingular .Table -}}
{{- $dbName := singular .Table -}}
{{- $varNameSingular := camelCaseSingular .Table -}}
// {{$tableNameSingular}}Find retrieves a single record by ID.
func {{$tableNameSingular}}Find(db boil.DB, id int) (*{{$tableNameSingular}}, error) {
  if id == 0 {
    return nil, errors.New("{{.PkgName}}: no id provided for {{.Table}} select")
  }
  var {{$varNameSingular}} *{{$tableNameSingular}}
  err := db.Select(&{{$varNameSingular}}, `SELECT {{selectParamNames $dbName .Columns}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{.Table}}: %s", err)
  }

  return {{$varNameSingular}}, nil
}
