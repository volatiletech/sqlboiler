{{- if .Table.IsView -}}
{{- else -}}
{{- $alias := .Aliases.Table .Table.Name -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{if .AddGlobal -}}
// UpdateG a single {{$alias.UpSingular}} record using the global executor.
// See Update for more documentation.
func (o *{{$alias.UpSingular}}) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

{{end -}}

{{if .AddPanic -}}
// UpdateP uses an executor to update the {{$alias.UpSingular}}, and panics on error.
// See Update for more documentation.
func (o *{{$alias.UpSingular}}) UpdateP(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) int64 {
	rowsAff, err := o.Update(ctx, exec, columns)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// UpdateGP a single {{$alias.UpSingular}} record using the global executor. Panics on error.
// See Update for more documentation.
func (o *{{$alias.UpSingular}}) UpdateGP(ctx context.Context, columns boil.Columns) int64 {
	rowsAff, err := o.Update(ctx, boil.GetContextDB(), columns)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

// Update uses an executor to update the {{$alias.UpSingular}}.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *{{$alias.UpSingular}}) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	{{- template "timestamp_update_helper" . -}}

	var err error
	{{if not .NoHooks -}}
	if err = o.doBeforeUpdateHooks(ctx, exec); err != nil {
		return 0, err
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
			return 0, errors.New("{{.PkgName}}: unable to update {{.Table.Name}}, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE {{$schemaTable}} SET %s WHERE %s",
			strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, wl),
			strmangle.WhereClause("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}len(wl)+1{{else}}0{{end}}, {{$alias.DownSingular}}PrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, append(wl, {{$alias.DownSingular}}PrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}

	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to update {{.Table.Name}} row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by update for {{.Table.Name}}")
	}

	if !cached {
		{{$alias.DownSingular}}UpdateCacheMut.Lock()
		{{$alias.DownSingular}}UpdateCache[key] = cache
		{{$alias.DownSingular}}UpdateCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return rowsAff, o.doAfterUpdateHooks(ctx, exec)
	{{- else -}}
	return rowsAff, nil
	{{- end}}
}

{{if .AddPanic -}}
// UpdateAllP updates all rows with matching column names, and panics on error.
func (q {{$alias.DownSingular}}Query) UpdateAllP(ctx context.Context, exec boil.ContextExecutor, cols M) int64 {
	rowsAff, err := q.UpdateAll(ctx, exec, cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}


{{if .AddGlobal -}}
// UpdateAllG updates all rows with the specified column values.
func (q {{$alias.DownSingular}}Query) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

{{end -}}


{{if and .AddGlobal .AddPanic -}}
// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (q {{$alias.DownSingular}}Query) UpdateAllGP(ctx context.Context, cols M) int64 {
	rowsAff, err := q.UpdateAll(ctx, boil.GetContextDB(), cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}


// UpdateAll updates all rows with the specified column values.
func (q {{$alias.DownSingular}}Query) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to update all for {{.Table.Name}}")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to retrieve rows affected for {{.Table.Name}}")
	}

	return rowsAff, nil
}

{{if .AddGlobal -}}
// UpdateAllG updates all rows with the specified column values.
func (o {{$alias.UpSingular}}Slice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// UpdateAllGP updates all rows with the specified column values, and panics on error.
func (o {{$alias.UpSingular}}Slice) UpdateAllGP(ctx context.Context, cols M) int64 {
	rowsAff, err := o.UpdateAll(ctx, boil.GetContextDB(), cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

{{if .AddPanic -}}
// UpdateAllP updates all rows with the specified column values, and panics on error.
func (o {{$alias.UpSingular}}Slice) UpdateAllP(ctx context.Context, exec boil.ContextExecutor, cols M) int64 {
	rowsAff, err := o.UpdateAll(ctx, exec, cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o {{$alias.UpSingular}}Slice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("{{.PkgName}}: update all requires at least one column argument")
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

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to update all in {{$alias.DownSingular}} slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to retrieve rows affected all in update all {{$alias.DownSingular}}")
	}

	return rowsAff, nil
}

{{- end -}}
