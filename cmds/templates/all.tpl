{{- $tableName := .TableName -}}
// {{titleCase $tableName}}All retrieves all records.
func {{titleCase $tableName}}All(db boil.DB) ([]*{{titleCase $tableName}}, error) {
  {{$varName := camelCase $tableName -}}
  var {{$varName}} []*{{titleCase $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{selectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
