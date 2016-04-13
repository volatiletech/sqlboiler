{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{if hasPrimaryKey .Table.PKey -}}
// Update updates a single {{$tableNameSingular}} record.
// whitelist is a list of column_name's that should be updated.
// Update will match against the primary key column to find the record to update.
// WARNING: This Update method will NOT ignore nil members.
// If you pass in nil members, those columnns will be set to null.
func (o *{{$tableNameSingular}}) UpdateX(db Executor, whitelist ... string) error {
  if err := o.doBeforeUpdateHooks(); err != nil {
    return err
  }

  {{$flagIndex := primaryKeyFlagIndex .Table.Columns .Table.PKey.Columns}}
  var err error
  if len(whitelist) == 0 {
    _, err = db.Exec("UPDATE {{.Table.Name}} SET {{updateParamNames .Table.Columns .Table.PKey.Columns}} WHERE {{wherePrimaryKey .Table.PKey.Columns $flagIndex}}", {{updateParamVariables "o." .Table.Columns .Table.PKey.Columns}}, {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  } else {
    query := fmt.Sprintf(`UPDATE {{.Table.Name}} SET %s WHERE %s`, boil.SetParamNames(whitelist), WherePrimaryKey(len(whitelist)+1, {{commaList .Table.PKey.Columns}}))
    _, err = db.Exec(query, {{updateParamVariables "o." .Table.Columns .Table.PKey.Columns}}, {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  }

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}} row: %s", err)
  }

  if err := o.doAfterUpdateHooks(); err != nil {
    return err
  }

  return nil
}

func (o *{{$tableNameSingular}}) Update(whitelist ... string) error {
  return o.UpdateX(boil.GetDB(), whitelist...)
}
{{- end}}
