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

  wl, returnColumns := o.generateInsertColumns(whitelist...)

  var err error
  if err := o.doBeforeCreateHooks(); err != nil {
    return err
  }

  ins := fmt.Sprintf(`INSERT INTO {{.Table.Name}} ("%s") VALUES (%s)`, strings.Join(wl, `","`), strmangle.Placeholders(len(wl), 1, 1))

  {{if driverUsesLastInsertID .DriverName}}
  if boil.DebugMode {
    fmt.Fprintln(boil.DebugWriter, ins)
    fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, wl...))
  }

  result, err := exec.Exec(ins, boil.GetStructValues(o, wl...)...)
  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
  }

  if len(returnColumns) == 0 {
      return o.doAfterCreateHooks()
  }

  lastId, err := result.lastInsertId()
  if err != nil || lastId == 0 {
    return ErrSyncFail
  }

  if len({{$varNameSingular}}AutoIncPrimaryKeys) != 1 {
    return ErrSyncFail
  }

  pkey := {{$varNameSingular}}AutoIncPrimaryKeys[0]
  sel := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}} WHERE %s`, strings.Join(returnColumns, `","`), strmangle.WhereClause(1, pkey))
  err := exec.QueryRow(sel, lastId).Scan(boil.GetStructPointers(o, returnColumns...))
  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to populate default values for {{.Table.Name}}")
  }
  {{else}}
  if len(returnColumns) != 0 {
    ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
    err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
  } else {
    _, err = exec.Exec(ins, boil.GetStructValues(o, wl...)...)
  }

  if boil.DebugMode {
    fmt.Fprintln(boil.DebugWriter, ins)
    fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, wl...))
  }

  if err != nil {
    return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
  }
  {{end}}

  return o.doAfterCreateHooks()
}

// generateInsertColumns generates the whitelist columns and return columns for an insert statement
// the return columns are used to get values that are assigned within the database during the
// insert to keep the struct in sync with what's in the db.
// with a whitelist:
// - the whitelist is used for the insert columns
// - the return columns are the result of (columns with default values - the whitelist)
// without a whitelist:
// - start with columns without a default as these always need to be inserted
// - add all columns that have a default in the database but that are non-zero in the struct
// - the return columns are the result of (columns with default values - the previous set)
func (o *{{$tableNameSingular}}) generateInsertColumns(whitelist ...string) ([]string, []string) {
  if len(whitelist) > 0 {
    return whitelist, boil.SetComplement({{$varNameSingular}}ColumnsWithDefault, whitelist)
  }

  var wl []string

  wl = append(wl, {{$varNameSingular}}ColumnsWithoutDefault...)

  wl = boil.SetMerge(boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o), wl)
  wl = boil.SortByKeys({{$varNameSingular}}Columns, wl)

  // Only return the columns with default values that are not in the insert whitelist
  rc := boil.SetComplement({{$varNameSingular}}ColumnsWithDefault, wl)

  return wl, rc
}
