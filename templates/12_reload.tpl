{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
// ReloadGP refetches the object from the database and panics on error.
func (o *{{$tableNameSingular}}) ReloadGP() {
  if err := o.ReloadG(); err != nil {
    panic(boil.WrapErr(err))
  }
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *{{$tableNameSingular}}) ReloadP(exec boil.Executor) {
  if err := o.Reload(exec); err != nil {
    panic(boil.WrapErr(err))
  }
}

// ReloadG refetches the object from the database using the primary keys.
func (o *{{$tableNameSingular}}) ReloadG() error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for reload")
  }

  return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *{{$tableNameSingular}}) Reload(exec boil.Executor) error {
  ret, err := {{$tableNameSingular}}Find(exec, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
  if err != nil {
    return err
  }

  *o = *ret
  return nil
}

func (o *{{$tableNameSingular}}Slice) ReloadAllGP() {
  if err := o.ReloadAllG(); err != nil {
    panic(boil.WrapErr(err))
  }
}

func (o *{{$tableNameSingular}}Slice) ReloadAllP(exec boil.Executor) {
  if err := o.ReloadAll(exec); err != nil {
    panic(boil.WrapErr(err))
  }
}

func (o *{{$tableNameSingular}}Slice) ReloadAllG() error {
  if o == nil {
    return errors.New("{{.PkgName}}: empty {{$tableNameSingular}}Slice provided for reload all")
  }

  return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$tableNameSingular}}Slice) ReloadAll(exec boil.Executor) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} slice provided for reload all")
  }

  if len(*o) == 0 {
    return nil
  }

  {{$varNamePlural}} := {{$tableNameSingular}}Slice{}
  args := o.inPrimaryKeyArgs()

  sql := fmt.Sprintf(
    `SELECT {{.Table.Name}}.* FROM {{.Table.Name}} WHERE (%s) IN (%s)`,
    strings.Join(strmangle.IdentQuoteSlice({{$varNameSingular}}PrimaryKeyColumns), ","),
    strmangle.Placeholders(len(*o) * len({{$varNameSingular}}PrimaryKeyColumns), 1, len({{$varNameSingular}}PrimaryKeyColumns)),
  )

  q := boil.SQL(sql, args...)
  boil.SetExecutor(q, exec)

  err := q.Bind(&{{$varNamePlural}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to reload all in {{$tableNameSingular}}Slice: %v", err)
  }

  *o = {{$varNamePlural}}

  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
    fmt.Fprintln(boil.DebugWriter, args)
  }

  return nil
}
