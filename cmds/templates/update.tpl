{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
// {{$tableNameSingular}}Update updates a single record.
func {{$tableNameSingular}}Update(db boil.DB, id int, columns map[string]interface{}) error {
  if id == 0 {
    return errors.New("{{.PkgName}}: no id provided for {{.Table.Name}} update")
  }

  query := fmt.Sprintf(`UPDATE {{.Table.Name}} SET %s WHERE id=$%d`, boil.Update(columns), len(columns))

  _, err := db.Exec(query, id, boil.WhereParams(columns))
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update row with ID %d in {{.Table.Name}}: %s", id, err)
  }

  return nil
}

{{if hasPrimaryKey .Table.PKey -}}
// Update updates a single {{$tableNameSingular}} record.
// Update will match against the primary key column to find the record to update.
// WARNING: This Update method will NOT ignore nil members.
// If you pass in nil members, those columnns will be set to null.
func (o *{{$tableNameSingular}}) Update(db boil.DB) error {
  {{$flagIndex := primaryKeyFlagIndex .Table.Columns .Table.PKey.Columns}}
  _, err := db.Exec("UPDATE {{.Table.Name}} SET {{updateParamNames .Table.Columns .Table.PKey.Columns}} WHERE {{wherePrimaryKey .Table.PKey.Columns $flagIndex}}", {{updateParamVariables "o." .Table.Columns .Table.PKey.Columns}}, {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}} row: %s", err)
  }

  return nil
}
{{- end}}
