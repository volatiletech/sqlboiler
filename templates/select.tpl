{{- $tableName := .TableName -}}
func Select{{makeGoColName $tableName}}(id int, db *sqlx.DB) ({{makeGoColName $tableName}}, error) {
  if id == 0 {
    return nil, errors.New("No ID provided for {{makeGoColName $tableName}} select")
  }
  {{$varName := makeGoVarName $tableName}}
  var {{$varName}} {{makeGoColName $tableName}}
  err := db.Select(&{{$varName}}, `
          SELECT {{makeSelectParamNames $tableName .TableData}}
          WHERE id=$1
        `, id)

  if err != nil {
    return nil, fmt.Errorf("Unable to select from {{$tableName}}: %s", err)
  }

  return {{$varName}}, nil
}
