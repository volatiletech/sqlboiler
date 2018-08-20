{{- $alias := .Aliases.Table .Table.Name}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{if .AddGlobal -}}
// InsertG a single record. See Insert for whitelist behavior description.
func (o *{{$alias.UpSingular}}) InsertG({{if not .NoContext}}ctx context.Context, {{end -}} columns boil.Columns) error {
	return o.Insert({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, columns)
}

{{end -}}

{{if .AddPanic -}}
// InsertP a single record using an executor, and panics on error. See Insert
// for whitelist behavior description.
func (o *{{$alias.UpSingular}}) InsertP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, columns boil.Columns) {
	if err := o.Insert({{if not .NoContext}}ctx, {{end -}} exec, columns); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// InsertGP a single record, and panics on error. See Insert for whitelist
// behavior description.
func (o *{{$alias.UpSingular}}) InsertGP({{if not .NoContext}}ctx context.Context, {{end -}} columns boil.Columns) {
	if err := o.Insert({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, columns); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *{{$alias.UpSingular}}) Insert({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, columns boil.Columns) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
	}

	var err error
	{{- template "timestamp_insert_helper" . }}

	{{if not .NoHooks -}}
	if err := o.doBeforeInsertHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		return err
	}
	{{- end}}

	nzDefaults := queries.NonZeroDefaultSet({{$alias.DownSingular}}ColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	{{$alias.DownSingular}}InsertCacheMut.RLock()
	cache, cached := {{$alias.DownSingular}}InsertCache[key]
	{{$alias.DownSingular}}InsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			{{$alias.DownSingular}}Columns,
			{{$alias.DownSingular}}ColumnsWithDefault,
			{{$alias.DownSingular}}ColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO {{$schemaTable}} ({{.LQ}}%s{{.RQ}}) %%sVALUES (%s)%%s", strings.Join(wl, "{{.RQ}},{{.LQ}}"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			{{if .Dialect.UseDefaultKeyword -}}
			cache.query = "INSERT INTO {{$schemaTable}} %sDEFAULT VALUES%s"
			{{else -}}
			cache.query = "INSERT INTO {{$schemaTable}} () VALUES ()%s%s"
			{{end -}}
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			{{if .Dialect.UseLastInsertID -}}
			cache.retQuery = fmt.Sprintf("SELECT {{.LQ}}%s{{.RQ}} FROM {{$schemaTable}} WHERE %s", strings.Join(returnColumns, "{{.RQ}},{{.LQ}}"), strmangle.WhereClause("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, {{$alias.DownSingular}}PrimaryKeyColumns))
			{{else -}}
				{{if .Dialect.UseOutputClause -}}
			queryOutput = fmt.Sprintf("OUTPUT INSERTED.{{.LQ}}%s{{.RQ}} ", strings.Join(returnColumns, "{{.RQ}},INSERTED.{{.LQ}}"))
				{{else -}}
			queryReturning = fmt.Sprintf(" RETURNING {{.LQ}}%s{{.RQ}}", strings.Join(returnColumns, "{{.RQ}},{{.LQ}}"))
				{{end -}}
			{{end -}}
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	{{if .Dialect.UseLastInsertID -}}
	{{- $canLastInsertID := .Table.CanLastInsertID -}}
	{{if $canLastInsertID -}}
		{{if .NoContext -}}
	result, err := exec.Exec(cache.query, vals...)
		{{else -}}
	result, err := exec.ExecContext(ctx, cache.query, vals...)
		{{end -}}
	{{else -}}
		{{if .NoContext -}}
	_, err = exec.Exec(cache.query, vals...)
		{{else -}}
	_, err = exec.ExecContext(ctx, cache.query, vals...)
		{{end -}}
	{{- end}}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
	}
	
	{{if $canLastInsertID -}}
	var lastID int64
	{{- end}}
	var identifierCols []interface{}

	if len(cache.retMapping) == 0 {
		goto CacheNoHooks
	}

	{{if $canLastInsertID -}}
	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	{{$colName := index .Table.PKey.Columns 0 -}}
	{{- $col := .Table.GetColumn $colName -}}
	{{- $colTitled := $colName | titleCase}}
	o.{{$colTitled}} = {{$col.Type}}(lastID)
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == {{$alias.DownSingular}}Mapping["{{$colTitled}}"] {
		goto CacheNoHooks
	}
	{{- end}}

	identifierCols = []interface{}{
		{{range .Table.PKey.Columns -}}
		o.{{$alias.Column .}},
		{{end -}}
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	{{if .NoContext -}}
	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	{{else -}}
	err = exec.QueryRowContext(ctx, cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	{{end -}}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to populate default values for {{.Table.Name}}")
	}
	{{else}}
	if len(cache.retMapping) != 0 {
		{{if .NoContext -}}
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
		{{else -}}
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
		{{end -}}
	} else {
		{{if .NoContext -}}
		_, err = exec.Exec(cache.query, vals...)
		{{else -}}
		_, err = exec.ExecContext(ctx, cache.query, vals...)
		{{end -}}
	}

	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
	}
	{{end}}

{{if .Dialect.UseLastInsertID -}}
CacheNoHooks:
{{- end}}
	if !cached {
		{{$alias.DownSingular}}InsertCacheMut.Lock()
		{{$alias.DownSingular}}InsertCache[key] = cache
		{{$alias.DownSingular}}InsertCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return o.doAfterInsertHooks({{if not .NoContext}}ctx, {{end -}} exec)
	{{- else -}}
	return nil
	{{- end}}
}

// InsertAll inserts all rows with the specified column values, using an executor.
func (o {{$alias.UpSingular}}Slice) InsertAll({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, columns boil.Columns) error {
	ln := int64(len(o))
	if ln == 0 {
		return nil
	}

	{{if .Dialect.UseLastInsertID -}}
	var sql string
	{{else -}}
	var sql, retQuery string
	{{end -}}
	vals := []interface{}{}
	retMaps := []interface{}{}
	for i, row := range o {
		if err := row.doBeforeInsertHooks(ctx, exec); err != nil {
			return err
		}

		{{- if not .NoAutoTimestamps -}}
		{{- $colNames := .Table.Columns | columnNames -}}
		{{if containsAny $colNames "created_at" "updated_at"}}
		currTime := time.Now().In(boil.GetLocation())
		{{range $ind, $col := .Table.Columns}}
		{{- if eq $col.Name "created_at" -}}
		{{- if eq $col.Type "time.Time" }}
		if row.CreatedAt.IsZero() {
			row.CreatedAt = currTime
		}
		{{- else}}
		if queries.MustTime(row.CreatedAt).IsZero() {
			queries.SetScanner(&row.CreatedAt, currTime)
		}
		{{- end -}}
		{{- end -}}
		{{- if eq $col.Name "updated_at" -}}
		{{- if eq $col.Type "time.Time"}}
		if row.UpdatedAt.IsZero() {
			row.UpdatedAt = currTime
		}
		{{- else}}
		if queries.MustTime(row.UpdatedAt).IsZero() {
			queries.SetScanner(&row.UpdatedAt, currTime)
		}
		{{- end -}}
		{{- end -}}
		{{end}}
		{{end}}
		{{- end}}

		nzDefaults := queries.NonZeroDefaultSet({{$alias.DownSingular}}ColumnsWithDefault, row)

		wl, returnColumns := columns.InsertColumnSet(
			{{$alias.DownSingular}}Columns,
			{{$alias.DownSingular}}ColumnsWithDefault,
			{{$alias.DownSingular}}ColumnsWithoutDefault,
			nzDefaults,
		)

		if i == 0 {
			sql = "INSERT INTO {{$schemaTable}} " + "({{.LQ}}" + strings.Join(wl, "{{.RQ}},{{.LQ}}") + "{{.RQ}})" + " VALUES "
			{{if not .Dialect.UseLastInsertID -}}
			retQuery = " RETURNING {{.LQ}}" + strings.Join(returnColumns, "{{.RQ}},{{.LQ}}") + "{{.RQ}} "
			{{end -}}
		}

		sql += strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), len(vals)+1, len(wl))
		if i != len(o)-1 {
			sql += ","
		}

		valMapping, err := queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, wl)
		if err != nil {
			return err
		}
		retMapping, err := queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, returnColumns)
		if err != nil {
			return err
		}
		value := reflect.Indirect(reflect.ValueOf(row))
		vals = append(vals, queries.ValuesFromMapping(value, valMapping)...)
		retMaps = append(retMaps, queries.PtrsFromMapping(value, retMapping)...)
	}

	{{if not .Dialect.UseLastInsertID -}}
	sql += retQuery
	{{end -}}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, vals...)
	}

	{{if .Dialect.UseLastInsertID -}}
	{{- $canLastInsertID := .Table.CanLastInsertID -}}
	{{if $canLastInsertID -}}
		{{if .NoContext -}}
	result, err := exec.Exec(ctx, sql, vals...)
		{{else -}}
	result, err := exec.ExecContext(ctx, sql, vals...)
		{{end -}}
	{{else -}}
		{{if .NoContext -}}
	_, err := exec.Exec(ctx, sql, vals...)
		{{else -}}
	_, err := exec.ExecContext(ctx, sql, vals...)
		{{end -}}
	{{- end}}
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
	}
	{{if $canLastInsertID -}}
	var lastID int64
	{{- end}}

	{{if $canLastInsertID -}}
	lastID, err = result.LastInsertId()
	if err != nil {
		return ErrSyncFail
	}

	{{$colName := index .Table.PKey.Columns 0 -}}
	{{- $col := .Table.GetColumn $colName -}}
	{{- $colTitled := $colName | titleCase -}}
	for i := 0; i < len(o); i++ {
		o[i].{{$colTitled}} = {{$col.Type}}(lastID)+{{$col.Type}}(i)
	}
	{{- end}}

	{{else -}}
	if len(retMaps) == 0 {
		{{if .NoContext -}}
		_, err := exec.Exec(ctx, sql, vals...)
		{{else -}}
		_, err := exec.ExecContext(ctx, sql, vals...)
		{{end -}}
		if err != nil {
			return errors.Wrap(err, "models: unable to insert all in user slice")
		}
		return nil
	}

	{{if .NoContext -}}
	rows, err := exec.Query(ctx, sql, vals...)
	{{else -}}
	rows, err := exec.QueryContext(ctx, sql, vals...)
	{{end -}}
	if err != nil {
		return errors.Wrap(err, "models: unable to insert all in user slice")
	}
	defer rows.Close()

	for i := 0; rows.Next() || i < len(retMaps); i++ {
		if err := rows.Scan(retMaps[i]); err != nil {
			return err
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	{{end -}}

	return nil
}
