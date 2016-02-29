{{- $tableName := .TableName -}}
// {{makeGoColName $tableName}}FieldsAll retrieves the specified columns for all records.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{makeGoColName $tableName}}FieldsAll(db *sqlx.DB, results interface{}) error {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
