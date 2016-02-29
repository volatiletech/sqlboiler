{{- $tableName := .TableName -}}
// {{makeGoColName $tableName}}FieldsAllBy retrieves the specified columns
// for all records with the specified column values.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{makeGoColName $tableName}}FieldsAllBy(db *sqlx.DB, columns map[string]interface{}, results interface{}) error {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
