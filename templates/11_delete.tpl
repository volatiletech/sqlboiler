{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// DeleteP deletes a single {{$tableNameSingular}} record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$tableNameSingular}}) DeleteP(exec boil.Executor) {
  if err := o.Delete(exec); err != nil {
    panic(boil.WrapErr(err))
  }
}

// DeleteG deletes a single {{$tableNameSingular}} record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) DeleteG() error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for deletion")
  }

  return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single {{$tableNameSingular}} record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$tableNameSingular}}) DeleteGP() {
  if err := o.DeleteG(); err != nil {
    panic(boil.WrapErr(err))
  }
}

// Delete deletes a single {{$tableNameSingular}} record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) Delete(exec boil.Executor) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for deletion")
  }

  var mods []qm.QueryMod

  mods = append(mods,
    qm.From("{{.Table.Name}}"),
    qm.Where(`{{whereClause 1 .Table.PKey.Columns}}`, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}}),
  )

  query := NewQuery(exec, mods...)
  boil.SetDelete(query)

  _, err := boil.ExecQuery(query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete from {{.Table.Name}}: %s", err)
  }

  return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q {{$varNameSingular}}Query) DeleteAllP() {
    if err := q.DeleteAll(); err != nil {
      panic(boil.WrapErr(err))
    }
}

// DeleteAll deletes all matching rows.
func (q {{$varNameSingular}}Query) DeleteAll() error {
  if q.Query == nil {
    return errors.New("{{.PkgName}}: no {{$varNameSingular}}Query provided for delete all")
  }

  boil.SetDelete(q.Query)

  _, err := boil.ExecQuery(q.Query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete all from {{.Table.Name}}: %s", err)
  }

  return nil
}

// DeleteAll deletes all rows in the slice, and panics on error.
func (o {{$tableNameSingular}}Slice) DeleteAllGP() {
  if err := o.DeleteAllG(); err != nil {
    panic(boil.WrapErr(err))
  }
}

// DeleteAllG deletes all rows in the slice.
func (o {{$tableNameSingular}}Slice) DeleteAllG() error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} slice provided for delete all")
  }
  return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o {{$tableNameSingular}}Slice) DeleteAllP(exec boil.Executor) {
  if err := o.DeleteAll(exec); err != nil {
    panic(boil.WrapErr(err))
  }
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o {{$tableNameSingular}}Slice) DeleteAll(exec boil.Executor) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} slice provided for delete all")
  }

  if len(o) == 0 {
    return nil
  }

  args := o.inPrimaryKeyArgs()

  sql := fmt.Sprintf(
    `DELETE FROM {{.Table.Name}} WHERE (%s) IN (%s)`,
    strings.Join(strmangle.IdentQuoteSlice({{$varNameSingular}}PrimaryKeyColumns), ","),
    strmangle.Placeholders(len(o) * len({{$varNameSingular}}PrimaryKeyColumns), 1, len({{$varNameSingular}}PrimaryKeyColumns)),
  )

  q := boil.SQL(sql, args...)
  boil.SetExecutor(q, exec)

  _, err := boil.ExecQuery(q)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete all from {{$varNameSingular}} slice: %s", err)
  }

  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
  }

  return nil
}
