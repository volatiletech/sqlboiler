{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase | stringMap .StringFuncs.replaceReserved -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", " -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// {{$tableNameSingular}}Exists checks if the {{$tableNameSingular}} row exists.
func {{$tableNameSingular}}Exists(exec boil.Executor, {{$pkArgs}}) (bool, error) {
	var exists bool
	{{if eq .DriverName "mssql" -}}
	sql := "select case when exists(select top(1) 1 from {{$schemaTable}} where {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}) then 1 else 0 end"
	{{- else -}}
	sql := "select exists(select 1 from {{$schemaTable}} where {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}} limit 1)"
	{{- end}}

	if boil.DebugMode {
	  qStr, err := interpolateParams(sql, {{$pkNames | join ", "}})
	  if err != nil {
	    return false, err
	  }
    fmt.Fprintln(boil.DebugWriter, qStr)
	}

	row := exec.QueryRow(sql, {{$pkNames | join ", "}})

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: unable to check if {{.Table.Name}} exists")
	}

	return exists, nil
}

// {{$tableNameSingular}}ExistsG checks if the {{$tableNameSingular}} row exists.
func {{$tableNameSingular}}ExistsG({{$pkArgs}}) (bool, error) {
	return {{$tableNameSingular}}Exists(boil.GetDB(), {{$pkNames | join ", "}})
}

// {{$tableNameSingular}}ExistsGP checks if the {{$tableNameSingular}} row exists. Panics on error.
func {{$tableNameSingular}}ExistsGP({{$pkArgs}}) bool {
	e, err := {{$tableNameSingular}}Exists(boil.GetDB(), {{$pkNames | join ", "}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// {{$tableNameSingular}}ExistsP checks if the {{$tableNameSingular}} row exists. Panics on error.
func {{$tableNameSingular}}ExistsP(exec boil.Executor, {{$pkArgs}}) bool {
	e, err := {{$tableNameSingular}}Exists(exec, {{$pkNames | join ", "}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// IsNew() checks if record exists in db (aka if its primary key is set).
func (o *{{$tableNameSingular}}) IsNew() bool {
	r := reflect.ValueOf(o).Elem()
	for i := 0; i < r.NumField(); i++ {
		column := r.Type().Field(i).Tag.Get("boil")
		for _, pkColumn := range {{$varNameSingular}}PrimaryKeyColumns {
			if column == pkColumn {
				field := r.Field(i)
				if field.Interface() != reflect.Zero(field.Type()).Interface() {
					return false
				}
			}
		}
	}
	return true
}

// Save() inserts the record if it does not exist, or updates it if it does.
func (o *{{$tableNameSingular}}) Save(exec boil.Executor, whitelist ...string) error {
  if o.IsNew() {
    return o.Insert(exec, whitelist...)
  } else {
    return o.Update(exec, whitelist...)
  }
}

// SaveG() inserts the record if it does not exist, or updates it if it does.
func (o *{{$tableNameSingular}}) SaveG(whitelist ...string) error {
  if o.IsNew() {
    return o.InsertG(whitelist...)
  } else {
    return o.UpdateG(whitelist...)
  }
}
