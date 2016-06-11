{{if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// Update a single {{$tableNameSingular}} record. It takes a whitelist of
// column_name's that should be updated. The primary key will be used to find
// the record to update.
// WARNING: Update does NOT ignore nil members - only the whitelist can be used
// to control the set of columns that will be saved.
func (o *{{$tableNameSingular}}) Update(whitelist ... string) error {
  return o.UpdateX(boil.GetDB(), whitelist...)
}

func (o *{{$tableNameSingular}}) UpdateX(exec boil.Executor, whitelist ... string) error {
  return o.UpdateAtX(exec, {{titleCaseCommaList "o." .Table.PKey.Columns}}, whitelist...)
}

func (o *{{$tableNameSingular}}) UpdateAt({{primaryKeyFuncSig .Table.Columns .Table.PKey.Columns}}, whitelist ...string) error {
  return o.UpdateAtX(boil.GetDB(), {{camelCaseCommaList "" .Table.PKey.Columns}}, whitelist...)
}

func (o *{{$tableNameSingular}}) UpdateAtX(exec boil.Executor, {{primaryKeyFuncSig .Table.Columns .Table.PKey.Columns}}, whitelist ...string) error {
  if err := o.doBeforeUpdateHooks(); err != nil {
    return err
  }

  if len(whitelist) == 0 {
    cols := {{$varNameSingular}}ColumnsWithoutDefault
    cols = append(boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o), cols...)
    // Subtract primary keys and autoincrement columns
    cols = boil.SetComplement(cols, {{$varNameSingular}}PrimaryKeyColumns)
    cols = boil.SetComplement(cols, {{$varNameSingular}}AutoIncrementColumns)

    whitelist = make([]string, len(cols))
    copy(whitelist, cols)
  }

  var err error
  var query string
  if len(whitelist) != 0 {
    query = fmt.Sprintf(`UPDATE {{.Table.Name}} SET %s WHERE %s`, boil.SetParamNames(whitelist), boil.WherePrimaryKey(len(whitelist)+1, {{commaList .Table.PKey.Columns}}))
    _, err = exec.Exec(query, boil.GetStructValues(o, whitelist...), {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  } else {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}}, could not build a whitelist for row: %s", err)
  }

  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
  }

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}} row: %s", err)
  }

  if err := o.doAfterUpdateHooks(); err != nil {
    return err
  }

  return nil
}

func (q {{$varNameSingular}}Query) UpdateAll(cols M) error {
  boil.SetUpdate(q.Query, cols)

  _, err := boil.ExecQuery(q.Query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update all for {{.Table.Name}}: %s", err)
  }

  return nil
}
{{- end}}
