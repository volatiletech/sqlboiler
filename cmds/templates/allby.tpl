{{- $tableName := .TableName -}}
// {{titleCase $tableName}}AllBy retrieves all records with the specified column values.
func {{titleCase $tableName}}AllBy(db boil.DB, columns map[string]interface{}) ([]*{{titleCase $tableName}}, error) {
  {{$varName := camelCase $tableName -}}
  var {{$varName}} []*{{titleCase $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{selectParamNames $tableName .TableData}}`)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
