{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) UpsertG(update bool, conflictColumns []string, updateColumns []string,  whitelist ...string) error {
  return o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *{{$tableNameSingular}}) UpsertGP(update bool, conflictColumns []string, updateColumns []string,  whitelist ...string) {
  if err := o.Upsert(boil.GetDB(), update, conflictColumns, updateColumns, whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *{{$tableNameSingular}}) UpsertP(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string,  whitelist ...string) {
  if err := o.Upsert(exec, update, conflictColumns, updateColumns, whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) Upsert(exec boil.Executor, updateOnConflict bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for upsert")
  }

  var ret []string
  whitelist, ret = strmangle.InsertColumnSet(
    {{$varNameSingular}}Columns,
    {{$varNameSingular}}ColumnsWithDefault,
    {{$varNameSingular}}ColumnsWithoutDefault,
    boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o),
    whitelist,
  )
  update := strmangle.UpdateColumnSet(
    {{$varNameSingular}}Columns,
    {{$varNameSingular}}PrimaryKeyColumns,
    updateColumns,
  )
  conflict := conflictColumns
  if len(conflict) == 0 {
    conflict = make([]string, len({{$varNameSingular}}PrimaryKeyColumns))
    copy(conflict, {{$varNameSingular}}PrimaryKeyColumns)
  }

  query := generateUpsertQuery("{{.Table.Name}}", updateOnConflict, ret, update, conflict, whitelist)

  var err error
  if err := o.doBeforeUpsertHooks(); err != nil {
    return err
  }

  if boil.DebugMode {
    fmt.Fprintln(boil.DebugWriter, query)
    fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, whitelist...))
  }

  {{- if .UseLastInsertID}}
  return errors.New("don't know how to do this yet")
  {{- else}}
  if len(ret) != 0 {
    err = exec.QueryRow(query, boil.GetStructValues(o, whitelist...)...).Scan(boil.GetStructPointers(o, ret...)...)
  } else {
    _, err = exec.Exec(query, boil.GetStructValues(o, whitelist...)...)
  }
  {{- end}}

  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to upsert for {{.Table.Name}}")
  }

  if err := o.doAfterUpsertHooks(); err != nil {
    return err
  }

  return nil
}
