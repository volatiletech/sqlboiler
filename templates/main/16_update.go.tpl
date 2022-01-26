{{- if .Table.IsView -}}
{{- else -}}
{{- $alias := .Aliases.Table .Table.Name -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{if .AddGlobal -}}
// UpdateG a single {{$alias.UpSingular}} record using the global executor.
// See Update for more documentation.
func (o *{{$alias.UpSingular}}) UpdateG({{if not .NoContext}}ctx context.Context, {{end -}} columns boil.Columns) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return o.Update({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, columns)
}

{{end -}}

{{if .AddPanic -}}
// UpdateP uses an executor to update the {{$alias.UpSingular}}, and panics on error.
// See Update for more documentation.
func (o *{{$alias.UpSingular}}) UpdateP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, columns boil.Columns) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end}}err := o.Update({{if not .NoContext}}ctx, {{end -}} exec, columns)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// UpdateGP a single {{$alias.UpSingular}} record using the global executor. Panics on error.
// See Update for more documentation.
func (o *{{$alias.UpSingular}}) UpdateGP({{if not .NoContext}}ctx context.Context, {{end -}} columns boil.Columns) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end}}err := o.Update({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, columns)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

// Update uses an executor to update the {{$alias.UpSingular}}.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *{{$alias.UpSingular}}) Update({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, columns boil.Columns) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	{{- template "timestamp_update_helper" . -}}

	var err error
	{{if not .NoHooks -}}
	if err = o.doBeforeUpdateHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} err
	}
	{{end -}}

	key := makeCacheKey(columns, nil)
	{{$alias.DownSingular}}UpdateCacheMut.RLock()
	cache, cached := {{$alias.DownSingular}}UpdateCache[key]
	{{$alias.DownSingular}}UpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			{{$alias.DownSingular}}AllColumns,
			{{$alias.DownSingular}}PrimaryKeyColumns,
		)
		{{- if filterColumnsByAuto true .Table.Columns }}
		wl = strmangle.SetComplement(wl, {{$alias.DownSingular}}GeneratedColumns)
		{{end}}
		{{if not .NoAutoTimestamps}}
		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		{{end -}}
		if len(wl) == 0 {
			return {{if not .NoRowsAffected}}0, {{end -}} errors.New("{{.PkgName}}: unable to update {{.Table.Name}}, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE {{$schemaTable}} SET %s WHERE %s",
			strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, wl),
			strmangle.WhereClause("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}len(wl)+1{{else}}0{{end}}, {{$alias.DownSingular}}PrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, append(wl, {{$alias.DownSingular}}PrimaryKeyColumns...))
		if err != nil {
			return {{if not .NoRowsAffected}}0, {{end -}} err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	{{if .NoContext -}}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}
	{{else -}}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	{{end -}}

	{{if .NoRowsAffected -}}
		{{if .NoContext -}}
	_, err = exec.Exec(cache.query, values...)
		{{else -}}
	_, err = exec.ExecContext(ctx, cache.query, values...)
		{{end -}}
	{{else -}}
	var result sql.Result
		{{if .NoContext -}}
	result, err = exec.Exec(cache.query, values...)
		{{else -}}
	result, err = exec.ExecContext(ctx, cache.query, values...)
		{{end -}}
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
		{{$alias.DownSingular}}UpdateCacheMut.Lock()
		{{$alias.DownSingular}}UpdateCache[key] = cache
		{{$alias.DownSingular}}UpdateCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return {{if not .NoRowsAffected}}rowsAff, {{end -}} o.doAfterUpdateHooks({{if not .NoContext}}ctx, {{end -}} exec)
	{{- else -}}
	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
	{{- end}}
}

{{if .AddPanic -}}
// UpdateAllP updates all rows with matching column names, and panics on error.
func (q {{$alias.DownSingular}}Query) UpdateAllP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, cols M) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := q.UpdateAll({{if not .NoContext}}ctx, {{end -}} exec, cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}


{{if .AddGlobal -}}
// UpdateAllG updates all rows with the specified column values.
func (q {{$alias.DownSingular}}Query) UpdateAllG({{if not .NoContext}}ctx context.Context, {{end -}} cols M) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return q.UpdateAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, cols)
}

{{end -}}


{{if and .AddGlobal .AddPanic -}}
// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (q {{$alias.DownSingular}}Query) UpdateAllGP({{if not .NoContext}}ctx context.Context, {{end -}} cols M) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := q.UpdateAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}


// UpdateAll updates all rows with the specified column values.
func (q {{$alias.DownSingular}}Query) UpdateAll({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, cols M) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	queries.SetUpdate(q.Query, cols)

	{{if .NoRowsAffected -}}
		{{if .NoContext -}}
	_, err := q.Query.Exec(exec)
		{{else -}}
	_, err := q.Query.ExecContext(ctx, exec)
		{{end -}}
	{{else -}}
		{{if .NoContext -}}
	result, err := q.Query.Exec(exec)
		{{else -}}
	result, err := q.Query.ExecContext(ctx, exec)
		{{end -}}
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
func (o {{$alias.UpSingular}}Slice) UpdateAllG({{if not .NoContext}}ctx context.Context, {{end -}} cols M) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return o.UpdateAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, cols)
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o {{$alias.UpSingular}}Slice) UpdateAllGP({{if not .NoContext}}ctx context.Context, {{end -}} cols M) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := o.UpdateAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, cols)
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
func (o {{$alias.UpSingular}}Slice) UpdateAllP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, cols M) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := o.UpdateAll({{if not .NoContext}}ctx, {{end -}} exec, cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o {{$alias.UpSingular}}Slice) UpdateAll({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, cols M) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$alias.DownSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE {{$schemaTable}} SET %s WHERE %s",
		strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.UseIndexPlaceholders}}len(colNames)+1{{else}}0{{end}}, {{$alias.DownSingular}}PrimaryKeyColumns, len(o)))

	{{if .NoContext -}}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}
	{{else -}}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	{{end -}}

	{{if .NoRowsAffected -}}
		{{if .NoContext -}}
	_, err := exec.Exec(sql, args...)
		{{else -}}
	_, err := exec.ExecContext(ctx, sql, args...)
		{{end -}}
	{{else -}}
		{{if .NoContext -}}
	result, err := exec.Exec(sql, args...)
		{{else -}}
	result, err := exec.ExecContext(ctx, sql, args...)
		{{end -}}
	{{end -}}
	if err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to update all in {{$alias.DownSingular}} slice")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to retrieve rows affected all in update all {{$alias.DownSingular}}")
	}
	{{end -}}

	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
}

{{- end -}}
