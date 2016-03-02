{{- $tableName := .Table -}}
// {{titleCase $tableName}}FieldsFindBy retrieves the specified columns
// for a single record with the specified column values.
// Pass in a pointer to an object with `db` tags that match the column names you wish to retrieve.
// For example: friendName string `db:"friend_name"`
func {{titleCase $tableName}}FieldsFindBy(db boil.DB, columns map[string]interface{}, results interface{}) (*{{titleCase $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("model: no id provided for {{$tableName}} select")
  }
  {{$varName := camelCase $tableName}}
  var {{$varName}} *{{titleCase $tableName}}
  err := db.Select(&{{$varName}}, `SELECT {{selectParamNames $tableName .Columns}} WHERE id=$1`, id)

  if err != nil {
    return nil, fmt.Errorf("models: unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
