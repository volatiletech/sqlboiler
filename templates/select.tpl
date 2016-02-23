func Insert{{makeGoColName $tableName}}(o *{{makeGoColName $tableName}}, db *sqlx.DB) (int, error) {
  if o == nil {
    return 0, errors.New("No {{objName}} provided for insertion")
  }

  var rowID int
  err := db.QueryRow(
    `INSERT INTO {{tableName}}
    ({{makeGoInsertParamNames tableData}})
    VALUES({{makeGoInsertParamFlags tableData}})
    RETURNING id`
  )

  if err != nil {
    return 0, fmt.Errorf("Unable to insert {{objName}}: %s", err)
  }

  return rowID, nil
}
