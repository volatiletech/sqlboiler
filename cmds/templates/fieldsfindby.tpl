{{- $tableName := .TableName -}}
// {{makeGoName $tableName}}FieldsFindBy retrieves the specified columns
// for a single record with the specified column values.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{makeGoName $tableName}}FieldsFindBy(db boil.DB, columns map[string]interface{}, results interface{}) (*{{makeGoName $tableName}}, error) {
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
