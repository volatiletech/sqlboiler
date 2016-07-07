{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// Insert a single record.
func (o *{{$tableNameSingular}}) Insert(include ... string) error {
  return o.InsertX(boil.GetDB(), include...)
}

// InsertX a single record using an executor.
func (o *{{$tableNameSingular}}) InsertX(exec boil.Executor, include ... string) error {
  if o == nil {
    return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
  }

  var includes []string

  includes = append(includes, include...)
  if len(include) == 0 {
    includes = append(includes, {{$varNameSingular}}ColumnsWithoutDefault...)
  }

  includes = append(boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o), includes...)
  includes = boil.SortByKeys({{$varNameSingular}}Columns, includes)

  // Only return the columns with default values that are not in the insert include
  returnColumns := boil.SetComplement({{$varNameSingular}}ColumnsWithDefault, includes)

  var err error
  if err := o.doBeforeCreateHooks(); err != nil {
    return err
  }

  ins := fmt.Sprintf(`INSERT INTO {{.Table.Name}} ("%s") VALUES (%s)`, strings.Join(includes, `","`), boil.GenerateParamFlags(len(includes), 1))

  {{if driverUsesLastInsertID .DriverName}}
  if len(returnColumns) != 0 {
    result, err := exec.Exec(ins, boil.GetStructValues(o, includes...)...)
    if err != nil {
      return fmt.Errorf("{{.PkgName}}: unable to insert into {{.Table.Name}}: %s", err)
    }

    lastId, err := result.lastInsertId()
    if err != nil || lastId == 0 {
      sel := fmt.Sprintf(`SELECT %s FROM {{.Table.Name}} WHERE %s`, strings.Join(returnColumns, `","`), boil.WhereClause(includes))
      rows, err := exec.Query(sel, boil.GetStructValues(o, includes...)...)
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
    _, err = exec.Exec(ins, boil.GetStructValues(o, includes...)...)
  }
  {{else}}
  if len(returnColumns) != 0 {
    ins = ins + fmt.Sprintf(` RETURNING %s`, strings.Join(returnColumns, ","))
    err = exec.QueryRow(ins, boil.GetStructValues(o, includes...)...).Scan(boil.GetStructPointers(o, returnColumns...)...)
  } else {
    _, err = exec.Exec(ins, {{.Table.Columns | columnNames | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
  }
  {{end}}

  if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, ins, boil.GetStructValues(o, includes...))
  }

  if err != nil {
    return fmt.Errorf("{{.PkgName}}: unable to insert into {{.Table.Name}}: %s", err)
  }

  if err := o.doAfterCreateHooks(); err != nil {
    return err
  }

  return nil
}
