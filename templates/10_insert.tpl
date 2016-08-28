{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// InsertG a single record. See Insert for whitelist behavior description.
func (o *{{$tableNameSingular}}) InsertG(whitelist ... string) error {
  return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *{{$tableNameSingular}}) InsertGP(whitelist ... string) {
  if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *{{$tableNameSingular}}) InsertP(exec boil.Executor, whitelist ... string) {
  if err := o.Insert(exec, whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are inferred (i.e. name, age)
// - All columns with a default, but non-zero are inferred (i.e. health = 75)
func (o *{{$tableNameSingular}}) Insert(exec boil.Executor, whitelist ... string) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
  }

  var err error
  {{- template "timestamp_insert_helper" . }}

  {{if eq .NoHooks false -}}
  if err := o.doBeforeInsertHooks(); err != nil {
    return err
  }
  {{- end}}

  wl, returnColumns := strmangle.InsertColumnSet(
    {{$varNameSingular}}Columns,
    {{$varNameSingular}}ColumnsWithDefault,
    {{$varNameSingular}}ColumnsWithoutDefault,
    boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, {{$varNameSingular}}TitleCases, o),
    whitelist,
  )

  ins := fmt.Sprintf(`INSERT INTO {{.Table.Name}} ("%s") VALUES (%s)`, strings.Join(wl, `","`), strmangle.Placeholders(len(wl), 1, 1))

  {{if .UseLastInsertID}}
  if boil.DebugMode {
    fmt.Fprintln(boil.DebugWriter, ins)
    fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, {{$varNameSingular}}TitleCases, wl...))
  }

  result, err := exec.Exec(ins, boil.GetStructValues(o, {{$varNameSingular}}TitleCases, wl...)...)
  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
  }

  {{if eq .NoHooks false -}}
  if len(returnColumns) == 0 {
      return o.doAfterInsertHooks()
  }
  {{- end}}

  lastID, err := result.LastInsertId()
  if err != nil || lastID == 0 || len({{$varNameSingular}}AutoIncPrimaryKeys) != 1 {
    return ErrSyncFail
  }

  sel := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}} WHERE %s`, strings.Join(returnColumns, `","`), strmangle.WhereClause(1, {{$varNameSingular}}AutoIncPrimaryKeys))
  err = exec.QueryRow(sel, lastID).Scan(boil.GetStructPointers(o, {{$varNameSingular}}TitleCases, returnColumns...))
  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to populate default values for {{.Table.Name}}")
  }
  {{else}}
  if len(returnColumns) != 0 {
    ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
    err = exec.QueryRow(ins, boil.GetStructValues(o, {{$varNameSingular}}TitleCases, wl...)...).Scan(boil.GetStructPointers(o, {{$varNameSingular}}TitleCases, returnColumns...)...)
  } else {
    _, err = exec.Exec(ins, boil.GetStructValues(o, {{$varNameSingular}}TitleCases, wl...)...)
  }

  if boil.DebugMode {
    fmt.Fprintln(boil.DebugWriter, ins)
    fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, {{$varNameSingular}}TitleCases, wl...))
  }

  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
  }
  {{end}}

  {{if eq .NoHooks false -}}
  return o.doAfterInsertHooks()
  {{- else -}}
  return nil
  {{- end}}
}
