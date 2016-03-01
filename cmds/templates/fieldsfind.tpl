{{- $tableName := .TableName -}}
// {{titleCase $tableName}}FieldsFind retrieves the specified columns for a single record by ID.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{titleCase $tableName}}FieldsFind(db boil.DB, id int, results interface{}) (*{{titleCase $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := camelCase $tableName}}
  var {{$varName}} *{{titleCase $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{selectParamNames $tableName .TableData}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
