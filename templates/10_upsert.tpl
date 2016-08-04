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

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) Upsert(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string, whitelist ...string) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for upsert")
  }

  wl, returnColumns := o.generateInsertColumns(whitelist...)

  conflict := make([]string, len(conflictColumns))
  update := make([]string, len(updateColumns))

  copy(conflict, conflictColumns)
  copy(update, updateColumns)

  for i, v := range conflict {
    conflict[i] = strmangle.IdentQuote(v)
  }

  for i, v := range update {
    update[i] = strmangle.IdentQuote(v)
  }

  var err error
  if err := o.doBeforeUpsertHooks(); err != nil {
    return err
  }

  ins := fmt.Sprintf(`INSERT INTO {{.Table.Name}} ("%s") VALUES (%s) ON CONFLICT `, strings.Join(wl, `","`), boil.GenerateParamFlags(len(wl), 1))
  if !update {
    ins := ins + "DO NOTHING"
  } else if len(conflict) != 0 {
    ins := ins + fmt.Sprintf(`("%s") DO UPDATE SET %s`, strings.Join(conflict, `","`))
  } else {
    ins := ins + fmt.Sprintf(`("%s") DO UPDATE SET %s`, strings.Join({{$varNameSingular}}PrimaryKeyColumns, `","`))
  }

  if len(returnColumns) != 0 {
    ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
    err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
  } else {
    _, err = exec.Exec(ins, {{.Table.Columns | columnNames | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
  }

  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins, boil.GetStructValues(o, wl...))
  }

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to upsert for {{.Table.Name}}: %s", err)
  }

  if err := o.doAfterUpsertHooks(); err != nil {
    return err
  }

  return nil
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *{{$tableNameSingular}}) UpsertP(exec boil.Executor, update bool, conflictColumns []string, updateColumns []string,  whitelist ...string) {
  if err := o.Upsert(exec, update, conflictColumns, updateColumns, whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}
