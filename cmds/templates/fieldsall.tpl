{{- $tableName := .Table -}}
// {{titleCase $tableName}}FieldsAll retrieves the specified columns for all records.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{titleCase $tableName}}FieldsAll(db boil.DB, results interface{}) error {
  {{$varName := camelCase $tableName -}}
  var {{$varName}} []*{{titleCase $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{selectParamNames $tableName .Columns}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
