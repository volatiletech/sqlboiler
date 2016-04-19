{{if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// Delete deletes a single {{$tableNameSingular}} record.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) Delete(mods ...QueryMod) error {
  return o.DeleteX(boil.GetDB(), mods...)
}

func (o *{{$tableNameSingular}}) DeleteX(exec boil.Executor, mods ...QueryMod) error {
  mods = append(mods,
    boil.From("{{.table.Name}}"),
    boil.Where("{{wherePrimaryKey .Table.Pkey.Columns 1}}", {{paramsPrimaryKey "o." .Table.PKey.Columns true}}),
  )

  query := NewQueryX(exec, mods...)

  _, err := exec.Exec("DELETE FROM {{.Table.Name}} WHERE {{wherePrimaryKey .Table.PKey.Columns 1}}", {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete from {{.Table.Name}}: %s", err)
  }

  return nil
}

func (o {{$varNameSingular}}Query) DeleteAll() error {
  _, err := db.Exec("DELETE FROM {{.Table.Name}} WHERE {{wherePrimaryKey .Table.PKey.Columns 1}}", {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete from {{.Table.Name}}: %s", err)
  }

  return nil
}
{{- end}}
