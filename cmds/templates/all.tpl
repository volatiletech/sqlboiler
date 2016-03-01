{{- $tableName := .TableName -}}
// {{makeGoName $tableName}}All retrieves all records.
func {{makeGoName $tableName}}All(db boil.DB) ([]*{{makeGoName $tableName}}, error) {
  {{$varName := makeGoVarName $tableName -}}
  var {{$varName}} []*{{makeGoName $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{makeSelectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
