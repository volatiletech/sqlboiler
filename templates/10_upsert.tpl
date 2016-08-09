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
func (o *{{$tableNameSingular}}) Upsert(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for upsert")
  }

  columns := o.generateUpsertColumns(conflictColumns, updateColumns, whitelist)
  query := o.generateUpsertQuery(update, columns)

  var err error
  if err := o.doBeforeUpsertHooks(); err != nil {
    return err
  }

  if len(columns.returning) != 0 {
    err = exec.QueryRow(query, boil.GetStructValues(o, columns.whitelist...)...).Scan(boil.GetStructPointers(o, columns.returning...)...)
  } else {
    _, err = exec.Exec(query, {{.Table.Columns | columnNames | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
  }

  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query, boil.GetStructValues(o, columns.whitelist...))
  }

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to upsert for {{.Table.Name}}: %s", err)
  }

  if err := o.doAfterUpsertHooks(); err != nil {
    return err
  }

  return nil
}

// generateUpsertColumns builds an upsertData object, using generated values when necessary.
func (o *{{$tableNameSingular}}) generateUpsertColumns(conflict []string, update []string, whitelist []string) upsertData {
  var upsertCols upsertData

  upsertCols.whitelist, upsertCols.returning = o.generateInsertColumns(whitelist...)

  upsertCols.conflict = make([]string, len(conflict))
  upsertCols.update = make([]string, len(update))

  // generates the ON CONFLICT() columns if none are provided
  upsertCols.conflict = o.generateConflictColumns(conflict...)

  // generate the UPDATE SET columns if none are provided
  upsertCols.update = o.generateUpdateColumns(update...)

  return upsertCols
}

// generateConflictColumns returns the user provided columns.
// If no columns are provided, it returns the primary key columns.
func (o *{{$tableNameSingular}}) generateConflictColumns(columns ...string) []string {
  if len(columns) != 0 {
    return columns
  }

  c := make([]string, len({{$varNameSingular}}PrimaryKeyColumns))
  copy(c, {{$varNameSingular}}PrimaryKeyColumns)

  return c
}

// generateUpsertQuery builds a SQL statement string using the upsertData provided.
func (o *{{$tableNameSingular}}) generateUpsertQuery(update bool, columns upsertData) string {
  var set, query string

  columns.conflict = strmangle.IdentQuoteSlice(columns.conflict)
  columns.whitelist = strmangle.IdentQuoteSlice(columns.whitelist)

  var sets []string
  // Generate the UPDATE SET clause
  for _, v := range columns.update {
    quoted := strmangle.IdentQuote(v)
    sets = append(sets, fmt.Sprintf("%s = EXCLUDED.%s", quoted, quoted))
  }
  set = strings.Join(sets, ", ")

  query = fmt.Sprintf(
    `INSERT INTO {{.Table.Name}} (%s) VALUES (%s) ON CONFLICT`,
    strings.Join(columns.whitelist, `, `),
    strmangle.Placeholders(len(columns.whitelist), 1, 1),
  )

  if !update {
    query = query + " DO NOTHING"
  } else {
    query = fmt.Sprintf(`%s (%s) DO UPDATE SET %s`, query, strings.Join(columns.conflict, `, `), set)
  }

  if len(columns.returning) != 0 {
    query = fmt.Sprintf(`%s RETURNING %s`, query, strings.Join(columns.returning, ","))
  }

  return query
}
