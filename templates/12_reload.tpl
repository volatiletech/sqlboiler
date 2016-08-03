{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
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
