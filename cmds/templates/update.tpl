{{if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// Update updates a single {{$tableNameSingular}} record.
// whitelist is a list of column_name's that should be updated.
// Update will match against the primary key column to find the record to update.
// WARNING: This Update method will NOT ignore nil members.
// If you pass in nil members, those columnns will be set to null.
func (o *{{$tableNameSingular}}) Update(whitelist ... string) error {
  return o.UpdateX(boil.GetDB(), whitelist...)
}

func (o *{{$tableNameSingular}}) UpdateX(exec boil.Executor, whitelist ... string) error {
  if err := o.doBeforeUpdateHooks(); err != nil {
    return err
  }

  if len(whitelist) == 0 {
    whitelist = {{$varNameSingular}}ColumnsWithoutDefault
    whitelist = append(boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o), whitelist...)
    whitelist = boil.SetComplement(whitelist, {{$varNameSingular}}PrimaryKeyColumns)
    whitelist = boil.SetComplement(whitelist, {{$varNameSingular}}AutoIncrementColumns)
  }

  var err error
  if len(whitelist) != 0 {
    query := fmt.Sprintf(`UPDATE {{.Table.Name}} SET %s WHERE %s`, boil.SetParamNames(whitelist), boil.WherePrimaryKey(len(whitelist)+1, {{commaList .Table.PKey.Columns}}))
    _, err = exec.Exec(query, boil.GetStructValues(o, whitelist...), {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  } else {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}}, could not build a whitelist for row: %s")
  }

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}} row: %s", err)
  }

  if err := o.doAfterUpdateHooks(); err != nil {
    return err
  }

  return nil
}

func (o *{{$tableNameSingular}}) UpdateAt({{primaryKeyFuncSig .Table.Columns .Table.PKey.Columns}}, whitelist ...string) error {
  return o.UpdateAtX(boil.GetDB(), {{camelCaseCommaList .Table.PKey.Columns}}, whitelist...)
}

func (o *{{$tableNameSingular}}) UpdateAtX(exec boil.Executor, {{primaryKeyFuncSig .Table.Columns .Table.PKey.Columns}}, whitelist ...string) error {
  return nil
}

func (v {{$varNameSingular}}Query) UpdateAll(cols M) error {
return nil
}
{{- end}}
