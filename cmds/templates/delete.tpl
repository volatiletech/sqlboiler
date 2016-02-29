{{- $tableName := .TableName -}}
// {{makeGoColName $tableName}}Delete deletes a single record.
func {{makeGoColName $tableName}}Delete(db *sqlx.DB, id int) error {
  if id == nil {
    return nil, errors.New("model: no id provided for {{$tableName}} delete")
  }

  err := db.Exec("DELETE FROM {{$tableName}} WHERE id=$1", id)
  if err != nil {
    return errors.New("model: unable to delete from {{$tableName}}: %s", err)
  }

  return nil
}
