{{if hasPrimaryKey .Table.PKey -}}
{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
{{- $varNameSingular := camelCaseSingular .Table.Name -}}
// Delete deletes a single {{$tableNameSingular}} record.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) Delete(mods ...QueryMod) error {
  _, err := db.Exec("DELETE FROM {{.Table.Name}} WHERE {{wherePrimaryKey .Table.PKey.Columns 1}}", {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete from {{.Table.Name}}: %s", err)
  }

  return nil
}
{{- end}}

func (o {{$varNameSingular}}Query) DeleteAll() error {
  _, err := db.Exec("DELETE FROM {{.Table.Name}} WHERE {{wherePrimaryKey .Table.PKey.Columns 1}}", {{paramsPrimaryKey "o." .Table.PKey.Columns true}})
  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to delete from {{.Table.Name}}: %s", err)
  }

  return nil
}
