{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", "}}
// UpdateG a single {{$tableNameSingular}} record. See Update for
// whitelist behavior description.
func (o *{{$tableNameSingular}}) UpdateG(whitelist ...string) error {
  return o.Update(boil.GetDB(), whitelist...)
}

// UpdateGP a single {{$tableNameSingular}} record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *{{$tableNameSingular}}) UpdateGP(whitelist ...string) {
  if err := o.Update(boil.GetDB(), whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// UpdateP uses an executor to update the {{$tableNameSingular}}, and panics on error.
// See Update for whitelist behavior description.
func (o *{{$tableNameSingular}}) UpdateP(exec boil.Executor, whitelist ... string) {
  err := o.UpdateAt(exec, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}}, whitelist...)
  if err != nil {
    panic(boil.WrapErr(err))
  }
}

// Update uses an executor to update the {{$tableNameSingular}}.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
func (o *{{$tableNameSingular}}) Update(exec boil.Executor, whitelist ... string) error {
  return o.UpdateAt(exec, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}}, whitelist...)
}

// UpdateAtG updates the {{$tableNameSingular}} using the primary key to find the row to update.
func (o *{{$tableNameSingular}}) UpdateAtG({{$pkArgs}}, whitelist ...string) error {
  return o.UpdateAt(boil.GetDB(), {{$pkNames | join ", "}}, whitelist...)
}

// UpdateAtGP updates the {{$tableNameSingular}} using the primary key to find the row to update. Panics on error.
func (o *{{$tableNameSingular}}) UpdateAtGP({{$pkArgs}}, whitelist ...string) {
  if err := o.UpdateAt(boil.GetDB(), {{$pkNames | join ", "}}, whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// UpdateAt uses an executor to update the {{$tableNameSingular}} using the primary key to find the row to update.
func (o *{{$tableNameSingular}}) UpdateAt(exec boil.Executor, {{$pkArgs}}, whitelist ...string) error {
  if err := o.doBeforeUpdateHooks(); err != nil {
    return err
  }

  var err error
  var query string
  var values []interface{}

  wl := o.generateUpdateColumns(whitelist...)

  if len(wl) != 0 {
    query = fmt.Sprintf(`UPDATE {{.Table.Name}} SET %s WHERE %s`, boil.SetParamNames(wl), boil.WherePrimaryKey(len(wl)+1, {{.Table.PKey.Columns | stringMap .StringFuncs.quoteWrap | join ", "}}))
    values = boil.GetStructValues(o, wl...)
    values = append(values, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
    _, err = exec.Exec(query, values...)
  } else {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}}, could not build whitelist")
  }

  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
    fmt.Fprintln(boil.DebugWriter, values)
  }

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update {{.Table.Name}} row: %s", err)
  }

  if err := o.doAfterUpdateHooks(); err != nil {
    return err
  }

  return nil
}

// UpdateAtP uses an executor to update the {{$tableNameSingular}} using the primary key to find the row to update.
// Panics on error.
func (o *{{$tableNameSingular}}) UpdateAtP(exec boil.Executor, {{$pkArgs}}, whitelist ...string) {
  if err := o.UpdateAt(exec, {{$pkNames | join ", "}}, whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// UpdateAll updates all rows with matching column names.
func (q {{$varNameSingular}}Query) UpdateAll(cols M) error {
  boil.SetUpdate(q.Query, cols)

  _, err := boil.ExecQuery(q.Query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update all for {{.Table.Name}}: %s", err)
  }

  return nil
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q {{$varNameSingular}}Query) UpdateAllP(cols M) {
  if err := q.UpdateAll(cols); err != nil {
    panic(boil.WrapErr(err))
  }
}

// generateUpdateColumns generates the whitelist columns for an update statement
// if a whitelist is supplied, it's returned
// if a whitelist is missing then we begin with all columns
// then we remove the primary key columns
func (o *{{$tableNameSingular}}) generateUpdateColumns(whitelist ...string) []string {
  if len(whitelist) != 0 {
    return whitelist
  }

  return boil.SetComplement({{$varNameSingular}}Columns, {{$varNameSingular}}PrimaryKeyColumns)
}
