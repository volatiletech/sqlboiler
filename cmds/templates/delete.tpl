{{if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// Delete deletes a single {{$tableNameSingular}} record.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) Delete() error {
  return o.DeleteX(boil.GetDB())
}

func (o *{{$tableNameSingular}}) DeleteX(exec boil.Executor) error {
  var mods []qs.QueryMod

  mods = append(mods,
    qs.From("{{.Table.Name}}"),
    qs.Where("{{wherePrimaryKey .Table.PKey.Columns 1}}", {{paramsPrimaryKey "o." .Table.PKey.Columns true}}),
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
  boil.SetDelete(o)

  _, err := boil.ExecQuery(query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete all from {{.Table.Name}}: %s", err)
  }

  return nil
}

func (o {{$varNameSingular}}Slice) DeleteAll() error {
  return DeleteAllX(boil.GetDB())
}

func (o {{$varNameSingular}}Slice) DeleteAllX(exec boil.Executor) error {
  var mods []qs.QueryMod

  args := o.inPrimaryKeyArgs()
  in := boil.WherePrimaryKeyIn(len(o), {{primaryKeyStrList .Table.PKey.Columns}})

  mods = append(mods,
    qs.From("{{.Table.Name}}"),
    qs.Where(in, args...),
  )

  query := NewQueryX(exec, mods...)
  boil.SetDelete(query)

  _, err := boil.ExecQuery(query)
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete all from {{$varNameSingular}} slice: %s", err)
  }

  return nil
}
{{- end}}
