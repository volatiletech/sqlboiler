{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
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
// - All columns without a default value are included (i.e. name, age)
// - All columns with a default, but non-zero are included (i.e. health = 75)
func (o *{{$tableNameSingular}}) Insert(exec boil.Executor, whitelist ... string) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{.Table.Name}} provided for insertion")
	}

	var err error
	{{- template "timestamp_insert_helper" . }}

	{{if not .NoHooks -}}
	if err := o.doBeforeInsertHooks(exec); err != nil {
		return err
	}
	{{- end}}

	nzDefaults := queries.NonZeroDefaultSet({{$varNameSingular}}ColumnsWithDefault, o)

	key := makeCacheKey(whitelist, nzDefaults)
	{{$varNameSingular}}InsertCacheMut.RLock()
	cache, cached := {{$varNameSingular}}InsertCache[key]
	{{$varNameSingular}}InsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := strmangle.InsertColumnSet(
			{{$varNameSingular}}Columns,
			{{$varNameSingular}}ColumnsWithDefault,
			{{$varNameSingular}}ColumnsWithoutDefault,
			nzDefaults,
			whitelist,
		)

		cache.valueMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping({{$varNameSingular}}Type, {{$varNameSingular}}Mapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO {{$schemaTable}} ({{.LQ}}%s{{.RQ}}) %%sVALUES (%s)%%s", strings.Join(wl, "{{.RQ}},{{.LQ}}"), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			{{if eq .DriverName "mysql" -}}
			cache.query = "INSERT INTO {{$schemaTable}} () VALUES ()"
			{{else -}}
			cache.query = "INSERT INTO {{$schemaTable}} DEFAULT VALUES"
			{{end -}}
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			{{if .UseLastInsertID -}}
			cache.retQuery = fmt.Sprintf("SELECT {{.LQ}}%s{{.RQ}} FROM {{$schemaTable}} WHERE %s", strings.Join(returnColumns, "{{.RQ}},{{.LQ}}"), strmangle.WhereClause("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns))
			{{else -}}
				{{if ne .DriverName "mssql" -}}
			queryReturning = fmt.Sprintf(" RETURNING {{.LQ}}%s{{.RQ}}", strings.Join(returnColumns, "{{.RQ}},{{.LQ}}"))
				{{else -}}
			queryOutput = fmt.Sprintf("OUTPUT INSERTED.{{.LQ}}%s{{.RQ}} ", strings.Join(returnColumns, "{{.RQ}},INSERTED.{{.LQ}}"))
				{{end -}}
			{{end -}}
		}

		if len(wl) != 0 {
			cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	{{if .UseLastInsertID -}}
	{{- $canLastInsertID := .Table.CanLastInsertID -}}
	{{if $canLastInsertID -}}
	result, err := exec.Exec(cache.query, vals...)
	{{else -}}
	_, err = exec.Exec(cache.query, vals...)
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
	if lastID != 0 && len(cache.retMapping) == 1 && cache.retMapping[0] == {{$varNameSingular}}Mapping["{{$colTitled}}"] {
		goto CacheNoHooks
	}
	{{- end}}

	identifierCols = []interface{}{
		{{range .Table.PKey.Columns -}}
		o.{{. | titleCase}},
		{{end -}}
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.retQuery)
		fmt.Fprintln(boil.DebugWriter, identifierCols...)
	}

	err = exec.QueryRow(cache.retQuery, identifierCols...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to populate default values for {{.Table.Name}}")
	}
	{{else}}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRow(cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.Exec(cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to insert into {{.Table.Name}}")
	}
	{{end}}

{{if .UseLastInsertID -}}
CacheNoHooks:
{{- end}}
	if !cached {
		{{$varNameSingular}}InsertCacheMut.Lock()
		{{$varNameSingular}}InsertCache[key] = cache
		{{$varNameSingular}}InsertCacheMut.Unlock()
	}

	{{if not .NoHooks -}}
	return o.doAfterInsertHooks(exec)
	{{- else -}}
	return nil
	{{- end}}
}
