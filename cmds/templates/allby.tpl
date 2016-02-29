{{- $tableName := .TableName -}}
// {{makeGoColName $tableName}}AllBy retrieves all records with the specified column values.
func {{makeGoColName $tableName}}AllBy(db *sqlx.DB, columns map[string]interface{}) ([]*{{makeGoColName $tableName}}, error) {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
