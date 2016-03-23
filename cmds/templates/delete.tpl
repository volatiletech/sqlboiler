{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
// {{$tableNameSingular}}Delete deletes a single record.
func {{$tableNameSingular}}Delete(db boil.DB, id int) error {
  if id == 0 {
    return errors.New("{{.PkgName}}: no id provided for {{.Table.Name}} delete")
  }

  _, err := db.Exec("DELETE FROM {{.Table.Name}} WHERE id=$1", id)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete from {{.Table.Name}}: %s", err)
  }

  return nil
}

{{if hasPrimaryKey .Table.Columns -}}
// Delete deletes a single {{$tableNameSingular}} record.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) Delete(db boil.DB) error {
  {{- $pkeyName := getPrimaryKey .Table.Columns -}}
  _, err := db.Exec("DELETE FROM {{.Table.Name}} WHERE {{$pkeyName}}=$1", o.{{titleCase $pkeyName}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete from {{.Table.Name}}: %s", err)
  }

  return nil
}
{{- end}}
