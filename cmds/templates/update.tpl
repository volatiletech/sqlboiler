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

{{if hasPrimaryKey .Table.Columns -}}
// Update updates a single {{$tableNameSingular}} record.
// Update will match against the primary key column to find the record to update.
// WARNING: This Update method will NOT ignore nil members.
// If you pass in nil members, those columnns will be set to null.
func (o *{{$tableNameSingular}}) Update(db boil.DB) error {
  {{- $pkeyName := getPrimaryKey .Table.Columns -}}
  _, err := db.Exec("UPDATE {{.Table.Name}} SET {{updateParamNames .Table.Columns}} WHERE {{$pkeyName}}=${{len .Table.Columns}}", {{updateParamVariables "o." .Table.Columns}}, o.{{titleCase $pkeyName}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}} row using primary key {{$pkeyName}}: %s", err)
  }

  return nil
}
{{- end}}
