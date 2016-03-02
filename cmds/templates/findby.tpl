{{- $tableName := .Table -}}
// {{titleCase $tableName}}FindBy retrieves a single record with the specified column values.
func {{titleCase $tableName}}FindBy(db boil.DB, columns map[string]interface{}) (*{{titleCase $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := camelCase $tableName}}
  var {{$varName}} *{{titleCase $tableName}}
  err := db.Select(&{{$varName}}, fmt.Sprintf(`SELECT {{selectParamNames $tableName .Columns}} WHERE %s=$1`, column), value)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
