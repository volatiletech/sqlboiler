{{- $tableName := .TableName -}}
// {{makeGoColName $tableName}}All retrieves all records.
func {{makeGoColName $tableName}}All(db *sqlx.DB) ([]*{{makeGoColName $tableName}}, error) {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
