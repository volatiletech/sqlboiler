{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", " -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// {{$tableNameSingular}}Exists checks if the {{$tableNameSingular}} row exists.
func {{$tableNameSingular}}Exists(exec boil.Executor, {{$pkArgs}}) (bool, error) {
	var exists bool

	sql := "select exists(select 1 from {{$schemaTable}} where {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}} limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, {{$pkNames | join ", "}})
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
