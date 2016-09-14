{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) UpsertG({{if eq .DriverName "postgres"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) error {
	return o.Upsert(boil.GetDB(), {{if eq .DriverName "postgres"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...)
}

// UpsertGP attempts an insert, and does an update or ignore on conflict. Panics on error.
func (o *{{$tableNameSingular}}) UpsertGP({{if eq .DriverName "postgres"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) {
	if err := o.Upsert(boil.GetDB(), {{if eq .DriverName "postgres"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...); err != nil {
	panic(boil.WrapErr(err))
	}
}

// UpsertP attempts an insert using an executor, and does an update or ignore on conflict.
// UpsertP panics on error.
func (o *{{$tableNameSingular}}) UpsertP(exec boil.Executor, {{if eq .DriverName "postgres"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string,	whitelist ...string) {
	if err := o.Upsert(exec, {{if eq .DriverName "postgres"}}updateOnConflict, conflictColumns, {{end}}updateColumns, whitelist...); err != nil {
	panic(boil.WrapErr(err))
	}
}
	
// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
func (o *{{$tableNameSingular}}) Upsert(exec boil.Executor, {{if eq .DriverName "postgres"}}updateOnConflict bool, conflictColumns []string, {{end}}updateColumns []string, whitelist ...string) error {
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

	{{if eq .DriverName "postgres" -}}
	conflict := conflictColumns
	if len(conflict) == 0 {
	conflict = make([]string, len({{$varNameSingular}}PrimaryKeyColumns))
	copy(conflict, {{$varNameSingular}}PrimaryKeyColumns)
	}
	{{- end}}

	{{if eq .DriverName "postgres" -}}
	query := boil.BuildUpsertQueryPostgres(dialect, "{{.Table.Name}}", updateOnConflict, ret, update, conflict, whitelist)
	{{- else if eq .DriverName "mysql" -}}
	query := boil.BuildUpsertQueryMySQL(dialect, "{{.Table.Name}}", update, whitelist)
	{{- end}}

	if boil.DebugMode {
	fmt.Fprintln(boil.DebugWriter, query)
	fmt.Fprintln(boil.DebugWriter, boil.GetStructValues(o, whitelist...))
	}

	{{- if .UseLastInsertID}}
	res, err := exec.Exec(query, boil.GetStructValues(o, whitelist...)...)
	{{- else}}
	if len(ret) != 0 {
	err = exec.QueryRow(query, boil.GetStructValues(o, whitelist...)...).Scan(boil.GetStructPointers(o, ret...)...)
	} else {
	_, err = exec.Exec(query, boil.GetStructValues(o, whitelist...)...)
	}
	{{- end}}

	if err != nil {
	return errors.Wrap(err, "{{.PkgName}}: unable to upsert for {{.Table.Name}}")
	}

	{{if .UseLastInsertID -}}
	if len(ret) != 0 {
	lid, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to get last insert id for {{.Table.Name}}")
	}
	{{$aipk := autoIncPrimaryKey .Table.Columns .Table.PKey}}
	aipk := "{{$aipk.Name}}"
	// if the update did not change anything, lid will be 0
	if lid == 0 && aipk == "" {
		// do a select using all pkeys
	} else if lid != 0 {
		// do a select using all pkeys + lid
	}
	}
	{{- end}}

	{{if not .NoHooks -}}
	if err := o.doAfterUpsertHooks(exec); err != nil {
	return err
	}
	{{- end}}

	return nil
}
