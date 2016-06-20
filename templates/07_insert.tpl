{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// {{$tableNameSingular}}Insert inserts a single record.
func (o *{{$tableNameSingular}}) Insert(whitelist ... string) error {
  return o.InsertX(boil.GetDB(), whitelist...)
}

func (o *{{$tableNameSingular}}) InsertX(exec boil.Executor, whitelist ... string) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
  }

  var wl []string

  wl = append(wl, whitelist...)
  if len(whitelist) == 0 {
    wl = append(wl, {{$varNameSingular}}ColumnsWithoutDefault...)
  }

  wl = append(boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o), wl...)
  wl = boil.SortByKeys({{$varNameSingular}}Columns, wl)

  // Only return the columns with default values that are not in the insert whitelist
  returnColumns := boil.SetComplement({{$varNameSingular}}ColumnsWithDefault, wl)

  var err error
  if err := o.doBeforeCreateHooks(); err != nil {
    return err
  }

  ins := fmt.Sprintf(`INSERT INTO {{.Table.Name}} (%s) VALUES (%s)`, strings.Join(wl, ","), boil.GenerateParamFlags(len(wl), 1))

  {{if driverUsesLastInsertID .DriverName}}
  if len(returnColumns) != 0 {
    result, err := exec.Exec(ins, boil.GetStructValues(o, wl...)...)
    if err != nil {
      return fmt.Errorf("{{.PkgName}}: unable to insert into {{.Table.Name}}: %s", err)
    }

    lastId, err := result.lastInsertId()
    if err != nil || lastId == 0 {
      sel := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}} WHERE %s`, strings.Join(returnColumns, ","), boil.WhereClause(wl))
      rows, err := exec.Query(sel, boil.GetStructValues(o, wl...)...)
      if err != nil {
        return fmt.Errorf("{{.PkgName}}: unable to insert into {{.Table.Name}}: %s", err)
      }
      defer rows.Close()

      i := 0
      ptrs := boil.GetStructPointers(o, returnColumns...)
      for rows.Next() {
        if err := rows.Scan(ptrs[i]); err != nil {
          return fmt.Errorf("{{.PkgName}}: unable to get result of insert, scan failed for column %s index %d: %s\n\n%#v", returnColumns[i], i, err, ptrs)
        }
        i++
      }
    } else if {{$varNameSingular}}AutoIncPrimKey != "" {
      sel := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}} WHERE %s=$1`, strings.Join(returnColumns, ","), {{$varNameSingular}}AutoIncPrimaryKey, lastId)
    }
  } else {
    _, err = exec.Exec(ins, boil.GetStructValues(o, wl...)...)
  }
  {{else}}
  if len(returnColumns) != 0 {
    ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
    err = exec.QueryRow(ins, boil.GetStructValues(o, wl...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
  } else {
    _, err = exec.Exec(ins, {{.Table.Columns | columnNames | prefixStringSlice "o." | join ", "}})
  }
  {{end}}

  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins, boil.GetStructValues(o, wl...))
  }

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to insert into {{.Table.Name}}: %s", err)
  }

  if err := o.doAfterCreateHooks(); err != nil {
    return err
  }

  return nil
}
