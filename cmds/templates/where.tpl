{{- $tableName := .Table -}}
{{- $varName := camelCase $tableName -}}
// {{titleCase $tableName}}Where retrieves all records with the specified column values.
func {{titleCase $tableName}}Where(db boil.DB, columns map[string]interface{}) ([]*{{titleCase $tableName}}, error) {
  var {{$varName}} []*{{titleCase $tableName}}
  query := fmt.Sprintf(`SELECT {{selectParamNames $tableName .Columns}} FROM {{$tableName}} WHERE %s`, boil.Where(columns))
  err := db.Select(&{{$varName}}, query, boil.WhereParams(columns)...)

  if err != nil {
    return nil, fmt.Errorf("{{.PkgName}}: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
