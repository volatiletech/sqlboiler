{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{if .AddGlobal -}}
// UpdateG a single {{$tableNameSingular}} record. See Update for
// whitelist behavior description.
func (o *{{$tableNameSingular}}) UpdateG(whitelist ...string) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return o.Update(boil.GetDB(), whitelist...)
}

{{end -}}

{{if .AddPanic -}}
// UpdateP uses an executor to update the {{$tableNameSingular}}, and panics on error.
// See Update for whitelist behavior description.
func (o *{{$tableNameSingular}}) UpdateP(exec boil.Executor, whitelist ... string) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end}}err := o.Update(exec, whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// UpdateGP a single {{$tableNameSingular}} record.
// UpdateGP takes a whitelist of column names that should be updated.
// Panics on error. See Update for whitelist behavior description.
func (o *{{$tableNameSingular}}) UpdateGP(whitelist ...string) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end}}err := o.Update(boil.GetDB(), whitelist...)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

// Update uses an executor to update the {{$tableNameSingular}}.
// Whitelist behavior: If a whitelist is provided, only the columns given are updated.
// No whitelist behavior: Without a whitelist, columns are inferred by the following rules:
// - All columns are inferred to start with
// - All primary keys are subtracted from this set
// Update does not automatically update the record in case of default values. Use .Reload()
// to refresh the records.
func (o *{{$tableNameSingular}}) Update(exec boil.Executor, whitelist ... string) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	{{- template "timestamp_update_helper" . -}}

	var err error
	{{if not .NoHooks -}}
	if err = o.doBeforeUpdateHooks(exec); err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} err
	}
	{{end -}}

	key := makeCacheKey(whitelist, nil)
	{{$varNameSingular}}UpdateCacheMut.RLock()
	cache, cached := {{$varNameSingular}}UpdateCache[key]
	{{$varNameSingular}}UpdateCacheMut.RUnlock()

	if !cached {
		wl := strmangle.UpdateColumnSet(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}PrimaryKeyColumns,
			whitelist,
		)
		{{if .Dialect.UseAutoColumns -}}
		wl = strmangle.SetComplement(wl, {{$varNameSingular}}ColumnsWithAuto)
		{{end}}
		{{if not .NoAutoTimestamps}}
		if len(whitelist) == 0 {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		{{end -}}
		if len(wl) == 0 {
			return {{if not .NoRowsAffected}}0, {{end -}} errors.New("{{.PkgName}}: unable to update {{.Table.Name}}, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE {{$schemaTable}} SET %s WHERE %s",
			strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, wl),
			strmangle.WhereClause("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}len(wl)+1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, append(wl, {{$varNameSingular}}PrimaryKeyColumns...))
		if err != nil {
			return {{if not .NoRowsAffected}}0, {{end -}} err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	{{if .NoRowsAffected -}}
	_, err = exec.Exec(cache.query, values...)
	{{else -}}
	var result sql.Result
	result, err = exec.Exec(cache.query, values...)
	{{end -}}
	if err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to update {{.Table.Name}} row")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by update for {{.Table.Name}}")
	}

	{{end -}}

	if !cached {
		{{$varNameSingular}}UpdateCacheMut.Lock()
		{{$varNameSingular}}UpdateCache[key] = cache
		{{$varNameSingular}}UpdateCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return {{if not .NoRowsAffected}}rowsAff, {{end -}} o.doAfterUpdateHooks(exec)
	{{- else -}}
	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
	{{- end}}
}

{{if .AddPanic -}}
// UpdateAllP updates all rows with matching column names, and panics on error.
func (q {{$varNameSingular}}Query) UpdateAllP(cols M) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := q.UpdateAll(cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

// UpdateAll updates all rows with the specified column values.
func (q {{$varNameSingular}}Query) UpdateAll(cols M) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	queries.SetUpdate(q.Query, cols)

	{{if .NoRowsAffected -}}
	_, err := q.Query.Exec()
	{{else -}}
	result, err := q.Query.Exec()
	{{end -}}
	if err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to update all for {{.Table.Name}}")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to retrieve rows affected for {{.Table.Name}}")
	}

	{{end -}}

	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
}

{{if .AddGlobal -}}
// UpdateAllG updates all rows with the specified column values.
func (o {{$tableNameSingular}}Slice) UpdateAllG(cols M) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return o.UpdateAll(boil.GetDB(), cols)
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o {{$tableNameSingular}}Slice) UpdateAllGP(cols M) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := o.UpdateAll(boil.GetDB(), cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

{{if .AddPanic -}}
// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o {{$tableNameSingular}}Slice) UpdateAllP(exec boil.Executor, cols M) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := o.UpdateAll(exec, cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o {{$tableNameSingular}}Slice) UpdateAll(exec boil.Executor, cols M) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	ln := int64(len(o))
	if ln == 0 {
		return {{if not .NoRowsAffected}}0, {{end -}} nil
	}

	if len(cols) == 0 {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.New("{{.PkgName}}: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}
	
	sql := fmt.Sprintf("UPDATE {{$schemaTable}} SET %s WHERE %s",
		strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.UseIndexPlaceholders}}len(colNames)+1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	{{if .NoRowsAffected -}}
	_, err := exec.Exec(sql, args...)
	{{else -}}
	result, err := exec.Exec(sql, args...)
	{{end -}}
	if err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to update all in {{$varNameSingular}} slice")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to retrieve rows affected all in update all {{$varNameSingular}}")
	}
	{{end -}}

	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
}
