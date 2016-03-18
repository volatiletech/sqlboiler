{{- $tableNameSingular := titleCaseSingular .Table -}}
// {{$tableNameSingular}}Delete deletes a single record.
func {{$tableNameSingular}}Delete(db boil.DB, id int) error {
  if id == nil {
    return nil, errors.New("{{.PkgName}}: no id provided for {{.Table}} delete")
  }

  err := db.Exec("DELETE FROM {{.Table}} WHERE id=$1", id)
  if err != nil {
    return errors.New("{{.PkgName}}: unable to delete from {{.Table}}: %s", err)
  }

  return nil
}
