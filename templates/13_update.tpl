{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", " -}}
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
  err := o.Update(exec, whitelist...)
  if err != nil {
    panic(boil.WrapErr(err))
  }
}

// Update uses an executor to update the {{$tableNameSingular}}.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *{{$tableNameSingular}}) Update(exec boil.Executor, whitelist ... string) error {
  {{- template "timestamp_update_helper" . -}}

  var err error
  {{if not .NoHooks -}}
  if err = o.doBeforeUpdateHooks(exec); err != nil {
    return err
  }
  {{end -}}

  key := makeCacheKey(whitelist, nil)
  {{$varNameSingular}}UpdateCacheMut.RLock()
  cache, cached := {{$varNameSingular}}UpdateCache[key]
  {{$varNameSingular}}UpdateCacheMut.RUnlock()

  if !cached {
    wl := strmangle.UpdateColumnSet({{$varNameSingular}}Columns, {{$varNameSingular}}PrimaryKeyColumns, whitelist)

    cache.query = fmt.Sprintf(`UPDATE {{schemaTable .DriverName .Schema .Table.Name}} SET %s WHERE %s`, strmangle.SetParamNames(wl), strmangle.WhereClause(len(wl)+1, {{$varNameSingular}}PrimaryKeyColumns))
    cache.valueMapping, err = boil.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, append(wl, {{$varNameSingular}}PrimaryKeyColumns...))
    if err != nil {
      return err
    }
  }

  if len(cache.valueMapping) == 0 {
    return errors.New("{{.PkgName}}: unable to update {{.Table.Name}}, could not build whitelist")
  }

  values := boil.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

  if boil.DebugMode {
    fmt.Fprintln(boil.DebugWriter, cache.query)
    fmt.Fprintln(boil.DebugWriter, values)
  }

  result, err := exec.Exec(cache.query, values...)
  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to update {{.Table.Name}} row")
  }

  if r, err := result.RowsAffected(); err == nil && r != 1 {
    return errors.Errorf("failed to update single row, updated %d rows", r)
  }

  if !cached {
    {{$varNameSingular}}UpdateCacheMut.Lock()
    {{$varNameSingular}}UpdateCache[key] = cache
    {{$varNameSingular}}UpdateCacheMut.Unlock()
  }

  {{if not .NoHooks -}}
  return o.doAfterUpdateHooks(exec)
  {{- else -}}
  return nil
  {{- end}}
}

// UpdateAllP updates all rows with matching column names, and panics on error.
func (q {{$varNameSingular}}Query) UpdateAllP(cols M) {
  if err := q.UpdateAll(cols); err != nil {
    panic(boil.WrapErr(err))
  }
}

// UpdateAll updates all rows with the specified column values.
func (q {{$varNameSingular}}Query) UpdateAll(cols M) error {
  boil.SetUpdate(q.Query, cols)

  _, err := boil.ExecQuery(q.Query)
  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to update all for {{.Table.Name}}")
  }

  return nil
}

// UpdateAllG updates all rows with the specified column values.
func (o {{$tableNameSingular}}Slice) UpdateAllG(cols M) error {
  return o.UpdateAll(boil.GetDB(), cols)
}

// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o {{$tableNameSingular}}Slice) UpdateAllGP(cols M) {
  if err := o.UpdateAll(boil.GetDB(), cols); err != nil {
    panic(boil.WrapErr(err))
  }
}

// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o {{$tableNameSingular}}Slice) UpdateAllP(exec boil.Executor, cols M) {
  if err := o.UpdateAll(exec, cols); err != nil {
    panic(boil.WrapErr(err))
  }
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o {{$tableNameSingular}}Slice) UpdateAll(exec boil.Executor, cols M) error {
  ln := int64(len(o))
  if ln == 0 {
    return nil
  }

  if len(cols) == 0 {
    return errors.New("{{.PkgName}}: update all requires at least one column argument")
  }

  colNames := make([]string, len(cols))
  args := make([]interface{}, len(cols))

  i := 0
  for name, value := range cols {
    colNames[i] = strmangle.IdentQuote(name)
    args[i] = value
    i++
  }

  // Append all of the primary key values for each column
  args = append(args, o.inPrimaryKeyArgs()...)

  sql := fmt.Sprintf(
    `UPDATE {{schemaTable .DriverName .Schema .Table.Name}} SET (%s) = (%s) WHERE (%s) IN (%s)`,
    strings.Join(colNames, ", "),
    strmangle.Placeholders(len(colNames), 1, 1),
    strings.Join(strmangle.IdentQuoteSlice({{$varNameSingular}}PrimaryKeyColumns), ","),
    strmangle.Placeholders(len(o) * len({{$varNameSingular}}PrimaryKeyColumns), len(colNames)+1, len({{$varNameSingular}}PrimaryKeyColumns)),
  )

  if boil.DebugMode {
    fmt.Fprintln(boil.DebugWriter, sql)
    fmt.Fprintln(boil.DebugWriter, args...)
  }

  result, err := exec.Exec(sql, args...)
  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to update all in {{$varNameSingular}} slice")
  }

  if r, err := result.RowsAffected(); err == nil && r != ln {
    return errors.Errorf("failed to update %d rows, only affected %d", ln, r)
  }

  return nil
}
