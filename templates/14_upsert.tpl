{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable -}}
// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) UpsertG({{if ne .DriverName "mysql"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) error {
	return o.Upsert(boil.GetDB(), {{if ne .DriverName "mysql"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *{{$tableNameSingular}}) UpsertGP({{if ne .DriverName "mysql"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), {{if ne .DriverName "mysql"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *{{$tableNameSingular}}) UpsertP(exec boil.Executor, {{if ne .DriverName "mysql"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) {
	if err := o.Upsert(exec, {{if ne .DriverName "mysql"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...); err != nil {
		panic(boil.WrapErr(err))
	}
}
	
// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) Upsert(exec boil.Executor, {{if ne .DriverName "mysql"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string, whitelist ...string) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for upsert")
	}

	{{- template "timestamp_upsert_helper" . }}

	{{if not .NoHooks -}}
	if err := o.doBeforeUpsertHooks(exec); err != nil {
		return err
	}
	{{- end}}

	var err error
	var ret []string
	whitelist, ret = strmangle.InsertColumnSet(
		{{$varNameSingular}}Columns,
		{{$varNameSingular}}ColumnsWithDefault,
		{{$varNameSingular}}ColumnsWithoutDefault,
		boil.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o),
		whitelist,
	)
	update := strmangle.UpdateColumnSet(
		{{$varNameSingular}}Columns,
		{{$varNameSingular}}PrimaryKeyColumns,
		updateColumns,
	)

	{{if ne .DriverName "mysql" -}}
	conflict := conflictColumns
	if len(conflict) == 0 {
		conflict = make([]string, len({{$varNameSingular}}PrimaryKeyColumns))
		copy(conflict, {{$varNameSingular}}PrimaryKeyColumns)
	}
	query := boil.BuildUpsertQueryPostgres(dialect, "{{$schemaTable}}", updateOnConflict, ret, update, conflict, whitelist)
	{{- else -}}
	query := boil.BuildUpsertQueryMySQL(dialect, "{{.Table.Name}}", update, whitelist)
	{{- end}}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, whitelist...))
	}

	{{- if .UseLastInsertID}}
	result, err := exec.Exec(query, boil.GetStructValues(o, whitelist...)...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to upsert for {{.Table.Name}}")
	}

	if len(ret) == 0 {
	{{if not .NoHooks -}}
		return o.doAfterUpsertHooks(exec)
	{{else -}}
		return nil
	{{end -}}
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	var identifierCols []interface{}
	if lastID != 0 {
		{{- $colName := index .Table.PKey.Columns 0 -}}
		{{- $col := .Table.GetColumn $colName -}}
		o.{{$colName | singular | titleCase}} = {{$col.Type}}(lastID)
		identifierCols = []interface{}{lastID}
	} else {
		identifierCols = []interface{}{
			{{range .Table.PKey.Columns -}}
			o.{{. | singular | titleCase}},
			{{end -}}
		}
	}

	if lastID != 0 && len(ret) == 1 {
		retQuery := fmt.Sprintf(
			"SELECT %s FROM {{.LQ}}{{.Table.Name}}{{.RQ}} WHERE {{whereClause .LQ .RQ 0 .Table.PKey.Columns}}",
			strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, ret), ","),
		)

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, ret)
			fmt.Fprintln(boil.DebugWriter, identifierCols...)
		}

		err = exec.QueryRow(retQuery, identifierCols...).Scan(boil.GetStructPointers(o, ret...)...)
		if err != nil {
			return errors.Wrap(err, "{{.PkgName}}: unable to populate default values for {{.Table.Name}}")
		}
	}
	{{- else}}
	if len(ret) != 0 {
		err = exec.QueryRow(query, boil.GetStructValues(o, whitelist...)...).Scan(boil.GetStructPointers(o, ret...)...)
	} else {
		_, err = exec.Exec(query, boil.GetStructValues(o, whitelist...)...)
	}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to upsert for {{.Table.Name}}")
	}
	{{- end}}

	{{if not .NoHooks -}}
	if err := o.doAfterUpsertHooks(exec); err != nil {
		return err
	}
	{{- end}}

	return nil
}
