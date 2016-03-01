{{- $tableName := .TableName -}}
// {{makeGoName $tableName}}AllBy retrieves all records with the specified column values.
func {{makeGoName $tableName}}AllBy(db boil.DB, columns map[string]interface{}) ([]*{{makeGoName $tableName}}, error) {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
