{{- $tableName := .Table -}}
// {{titleCase $tableName}}Find retrieves a single record by ID.
func {{titleCase $tableName}}Find(db boil.DB, id int) (*{{titleCase $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := camelCase $tableName}}
  var {{$varName}} *{{titleCase $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{selectParamNames $tableName .Columns}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
