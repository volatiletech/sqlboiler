{{- $tableNameSingular := titleCaseSingular .Table.Name -}}
// {{$tableNameSingular}}Insert inserts a single record.
func (o *{{$tableNameSingular}}) Insert(mods ...QueryMod) error {
  if o == nil {
    return 0, errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
  }

  if err := o.doBeforeCreateHooks(); err != nil {
    return 0, err
  }

  var rowID int
  err := boil.GetDB().QueryRow(`INSERT INTO {{.Table.Name}} ({{insertParamNames .Table.Columns}}) VALUES({{insertParamFlags .Table.Columns}}) RETURNING id`, {{insertParamVariables "o." .Table.Columns}}).Scan(&rowID)

  if err != nil {
    return 0, fmt.Errorf("{{.PkgName}}: unable to insert {{.Table.Name}}: %s", err)
  }

  if err := o.doAfterCreateHooks(); err != nil {
    return 0, err
  }

  return rowID, nil
}

func (o *{{$tableNameSingular}}) InsertX(exec boil.Executor, mods ...QueryMod) error {
  if o == nil {
    return 0, errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
  }

  if err := o.doBeforeCreateHooks(); err != nil {
    return 0, err
  }

  var rowID int
  err := boil.GetDB().QueryRow(`INSERT INTO {{.Table.Name}} ({{insertParamNames .Table.Columns}}) VALUES({{insertParamFlags .Table.Columns}}) RETURNING id`, {{insertParamVariables "o." .Table.Columns}}).Scan(&rowID)

  if err != nil {
    return 0, fmt.Errorf("{{.PkgName}}: unable to insert {{.Table.Name}}: %s", err)
  }

  if err := o.doAfterCreateHooks(); err != nil {
    return 0, err
  }

  return rowID, nil
}
