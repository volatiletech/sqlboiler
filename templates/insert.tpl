{{- $tableName := .TableName -}}
func Insert{{makeGoColName $tableName}}(o *{{makeGoColName $tableName}}, db *sqlx.DB) (int, error) {
  if o == nil {
    return 0, errors.New("No {{makeGoColName $tableName}} provided for insertion")
  }

  var rowID int
  err := db.QueryRow(`
          INSERT INTO {{$tableName}}
          ({{makeGoInsertParamNames .TableData}})
          VALUES({{makeGoInsertParamFlags .TableData}})
          RETURNING id
        `)

  if err != nil {
    return 0, fmt.Errorf("Unable to insert {{$tableName}}: %s", err)
  }

  return rowID, nil
}
