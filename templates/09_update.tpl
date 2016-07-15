{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", "}}
// Update a single {{$tableNameSingular}} record.
// Update takes a whitelist of column names that should be updated.
// The primary key will be used to find the record to update.
func (o *{{$tableNameSingular}}) Update(whitelist ...string) error {
  return o.UpdateX(boil.GetDB(), whitelist...)
}

// UpdateX uses an executor to update the {{$tableNameSingular}}.
func (o *{{$tableNameSingular}}) UpdateX(exec boil.Executor, whitelist ... string) error {
  return o.UpdateAtX(exec, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}}, whitelist...)
}

// UpdateAt updates the {{$tableNameSingular}} using the primary key to find the row to update.
func (o *{{$tableNameSingular}}) UpdateAt({{$pkArgs}}, whitelist ...string) error {
  return o.UpdateAtX(boil.GetDB(), {{$pkNames | join ", "}}, whitelist...)
}

// UpdateAtX uses an executor to update the {{$tableNameSingular}} using the primary key to find the row to update.
func (o *{{$tableNameSingular}}) UpdateAtX(exec boil.Executor, {{$pkArgs}}, whitelist ...string) error {
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

func (q {{$varNameSingular}}Query) UpdateAll(cols M) error {
  boil.SetUpdate(q.Query, cols)

  _, err := boil.ExecQuery(q.Query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to update all for {{.Table.Name}}: %s", err)
  }

  return nil
}

// generateUpdateColumns generates the whitelist columns for an update statement
func (o *{{$tableNameSingular}}) generateUpdateColumns(whitelist ...string) []string {
  if len(whitelist) != 0 {
    return whitelist
  }

  var wl []string
  cols := {{$varNameSingular}}ColumnsWithoutDefault
  cols = append(boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o), cols...)
  // Subtract primary keys and autoincrement columns
  cols = boil.SetComplement(cols, {{$varNameSingular}}PrimaryKeyColumns)
  cols = boil.SetComplement(cols, {{$varNameSingular}}AutoIncrementColumns)

  wl = make([]string, len(cols))
  copy(wl, cols)

  return wl
}
