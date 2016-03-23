{{- $tableNameSingular := titleCaseSingular .Table -}}
// {{$tableNameSingular}}Update updates a single record.
func {{$tableNameSingular}}Update(db boil.DB, id int, columns map[string]interface{}) error {
  if id == 0 {
    return errors.New("{{.PkgName}}: no id provided for {{.Table}} update")
  }

  query := fmt.Sprintf(`UPDATE {{.Table}} SET %s WHERE id=$%d`, boil.Update(columns), len(columns))

  _, err := db.Exec(query, id, boil.WhereParams(columns))
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update row with ID %d in {{.Table}}: %s", id, err)
  }

  return nil
}

{{if hasPrimaryKey .Columns -}}
// Update updates a single {{$tableNameSingular}} record.
// Update will match against the primary key column to find the record to update.
// WARNING: This Update method will NOT ignore nil members.
// If you pass in nil members, those columnns will be set to null.
func (o *{{$tableNameSingular}}) Update(db boil.DB) error {
  {{- $pkeyName := getPrimaryKey .Columns -}}
  _, err := db.Exec("UPDATE {{.Table}} SET {{updateParamNames .Columns}} WHERE {{$pkeyName}}=${{len .Columns}}", {{updateParamVariables "o." .Columns}}, o.{{titleCase $pkeyName}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table}} row using primary key {{$pkeyName}}: %s", err)
  }

  return nil
}
{{- end}}
