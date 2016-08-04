{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// InsertG a single record. See Insert for whitelist behavior description.
func (o *{{$tableNameSingular}}) InsertG(whitelist ... string) error {
  return o.Insert(boil.GetDB(), whitelist...)
}

// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *{{$tableNameSingular}}) InsertGP(whitelist ... string) {
  if err := o.Insert(boil.GetDB(), whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *{{$tableNameSingular}}) InsertP(exec boil.Executor, whitelist ... string) {
  if err := o.Insert(exec, whitelist...); err != nil {
    panic(boil.WrapErr(err))
  }
}

// Insert a single record using an executor.
// Whitelist behavior: If a whitelist is provided, only those columns supplied are inserted
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns without a default value are inferred (i.e. name, age)
// - All columns with a default, but non-zero are inferred (i.e. health = 75)
func (o *{{$tableNameSingular}}) Insert(exec boil.Executor, whitelist ... string) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
  }

  wl, returnColumns := o.generateInsertColumns(whitelist...)

  var err error
  if err := o.doBeforeCreateHooks(); err != nil {
    return err
  }

  ins := fmt.Sprintf(`INSERT INTO {{.Table.Name}} ("%s") VALUES (%s)`, strings.Join(wl, `","`), boil.GenerateParamFlags(len(wl), 1))

  {{if driverUsesLastInsertID .DriverName}}
  if len(returnColumns) != 0 {
    result, err := exec.Exec(ins, boil.GetStructValues(o, wl...)...)
    if err != nil {
      return fmt.Errorf("{{.PkgName}}: unable to insert into {{.Table.Name}}: %s", err)
    }

    lastId, err := result.lastInsertId()
    if err != nil || lastId == 0 {
      sel := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}} WHERE %s`, strings.Join(returnColumns, `","`), boil.WhereClause(wl))
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
    _, err = exec.Exec(ins, {{.Table.Columns | columnNames | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
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

// generateInsertColumns generates the whitelist columns and return columns for an insert statement
func (o *{{$tableNameSingular}}) generateInsertColumns(whitelist ...string) ([]string, []string) {
  var wl []string

  wl = append(wl, whitelist...)
  if len(whitelist) == 0 {
    wl = append(wl, {{$varNameSingular}}ColumnsWithoutDefault...)
  }

  wl = append(boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o), wl...)
  wl = boil.SortByKeys({{$varNameSingular}}Columns, wl)

  // Only return the columns with default values that are not in the insert whitelist
  rc := boil.SetComplement({{$varNameSingular}}ColumnsWithDefault, wl)

  return wl, rc
}
