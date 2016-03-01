{{- $tableName := .TableName -}}
// {{makeGoName $tableName}}Find retrieves a single record by ID.
func {{makeGoName $tableName}}Find(db boil.DB, id int) (*{{makeGoName $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := makeGoVarName $tableName}}
  var {{$varName}} *{{makeGoName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
