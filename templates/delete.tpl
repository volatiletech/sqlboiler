{{if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// Delete deletes a single {{$tableNameSingular}} record.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) Delete() error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for deletion")
  }

  return o.DeleteX(boil.GetDB())
}

func (o *{{$tableNameSingular}}) DeleteX(exec boil.Executor) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for deletion")
  }

  var mods []qm.QueryMod

  mods = append(mods,
    qm.Table("{{.Table.Name}}"),
    qm.Where("{{wherePrimaryKey .Table.PKey.Columns 1}}", {{paramsPrimaryKey "o." .Table.PKey.Columns true}}),
  )

  query := NewQueryX(exec, mods...)
  boil.SetDelete(query)

  _, err := boil.ExecQuery(query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete from {{.Table.Name}}: %s", err)
  }

  return nil
}

func (o {{$varNameSingular}}Query) DeleteAll() error {
  if o.Query == nil {
    return errors.New("{{.PkgName}}: no {{$varNameSingular}}Query provided for delete all")
  }

  boil.SetDelete(o.Query)

  _, err := boil.ExecQuery(o.Query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete all from {{.Table.Name}}: %s", err)
  }

  return nil
}

func (o {{$varNameSingular}}Slice) DeleteAll() error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} slice provided for delete all")
  }
  return o.DeleteAllX(boil.GetDB())
}

func (o {{$varNameSingular}}Slice) DeleteAllX(exec boil.Executor) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{$tableNameSingular}} slice provided for delete all")
  }

  var mods []qm.QueryMod

  args := o.inPrimaryKeyArgs()
  in := boil.WherePrimaryKeyIn(len(o), {{commaList .Table.PKey.Columns}})

  mods = append(mods,
    qm.Table("{{.Table.Name}}"),
    qm.Where(in, args...),
  )

  query := NewQueryX(exec, mods...)
  boil.SetDelete(query)

  _, err := boil.ExecQuery(query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete all from {{$varNameSingular}} slice: %s", err)
  }
  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, args)
  }

  return nil
}
{{- end}}
