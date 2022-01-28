{{- if .Table.IsView -}}
{{- else -}}
{{- $alias := .Aliases.Table .Table.Name -}}
{{- $schemaTable := .Table.Name | .SchemaTable -}}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted -}}
{{- $soft := and .AddSoftDeletes $canSoftDelete }}
{{if .AddGlobal -}}
// DeleteG deletes a single {{$alias.UpSingular}} record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *{{$alias.UpSingular}}) DeleteG({{if not .NoContext}}ctx context.Context{{if $soft}}, hardDelete bool{{end}}{{else}}{{if $soft}}hardDelete bool{{end}}{{end}}) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return o.Delete({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}{{if $soft}}, hardDelete{{end}})
}

{{end -}}

{{if .AddPanic -}}
// DeleteP deletes a single {{$alias.UpSingular}} record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$alias.UpSingular}}) DeleteP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}{{if $soft}}, hardDelete bool{{end}}) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end}}err := o.Delete({{if not .NoContext}}ctx, {{end -}} exec{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// DeleteGP deletes a single {{$alias.UpSingular}} record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$alias.UpSingular}}) DeleteGP({{if not .NoContext}}ctx context.Context{{if $soft}}, hardDelete bool{{end}}{{else}}{{if $soft}}hardDelete bool{{end}}{{end}}) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end}}err := o.Delete({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

// Delete deletes a single {{$alias.UpSingular}} record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$alias.UpSingular}}) Delete({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}{{if $soft}}, hardDelete bool{{end}}) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	if o == nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.New("{{.PkgName}}: no {{$alias.UpSingular}} provided for delete")
	}

	{{if not .NoHooks -}}
	if err := o.doBeforeDeleteHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} err
	}
	{{- end}}

	{{if $soft -}}
	var (
		sql string
		args []interface{}
	)
	if hardDelete {
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), {{$alias.DownSingular}}PrimaryKeyMapping)
		sql = "DELETE FROM {{$schemaTable}} WHERE {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}"
	} else {
		currTime := time.Now().In(boil.GetLocation())
		o.DeletedAt = null.TimeFrom(currTime)
		wl := []string{"{{or $.AutoColumns.Deleted "deleted_at"}}"}
		sql = fmt.Sprintf("UPDATE {{$schemaTable}} SET %s WHERE {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 2 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}",
			strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, wl),
		)
		valueMapping, err := queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, append(wl, {{$alias.DownSingular}}PrimaryKeyColumns...))
		if err != nil {
			return {{if not .NoRowsAffected}}0, {{end -}} err
		}
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), valueMapping)
	}
	{{else -}}
	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), {{$alias.DownSingular}}PrimaryKeyMapping)
	sql := "DELETE FROM {{$schemaTable}} WHERE {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}"
	{{- end}}

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
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to delete from {{.Table.Name}}")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by delete for {{.Table.Name}}")
	}

	{{end -}}

	{{if not .NoHooks -}}
	if err := o.doAfterDeleteHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} err
	}
	{{- end}}

	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
}

{{if .AddGlobal -}}
func (q {{$alias.DownSingular}}Query) DeleteAllG({{if not .NoContext}}ctx context.Context{{end}}{{if $soft}}{{if not .NoContext}}, {{end}}hardDelete bool{{end}}) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return q.DeleteAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}{{if $soft}}, hardDelete{{end}})
}

{{end -}}

{{if .AddPanic -}}
// DeleteAllP deletes all rows, and panics on error.
func (q {{$alias.DownSingular}}Query) DeleteAllP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}{{if $soft}}, hardDelete bool{{end}}) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := q.DeleteAll({{if not .NoContext}}ctx, {{end -}} exec{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// DeleteAllGP deletes all rows, and panics on error.
func (q {{$alias.DownSingular}}Query) DeleteAllGP({{if not .NoContext}}ctx context.Context, {{end}}{{if $soft}}hardDelete bool{{end}}) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := q.DeleteAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

// DeleteAll deletes all matching rows.
func (q {{$alias.DownSingular}}Query) DeleteAll({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}{{if $soft}}, hardDelete bool{{end}}) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	if q.Query == nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.New("{{.PkgName}}: no {{$alias.DownSingular}}Query provided for delete all")
	}

	{{if $soft -}}
	if hardDelete {
		queries.SetDelete(q.Query)
	} else {
		currTime := time.Now().In(boil.GetLocation())
		queries.SetUpdate(q.Query, M{"{{or $.AutoColumns.Deleted "deleted_at"}}": currTime})
	}
	{{else -}}
	queries.SetDelete(q.Query)
	{{- end}}

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
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{.Table.Name}}")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by deleteall for {{.Table.Name}}")
	}

	{{end -}}

	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
}

{{if .AddGlobal -}}
// DeleteAllG deletes all rows in the slice.
func (o {{$alias.UpSingular}}Slice) DeleteAllG({{if not .NoContext}}ctx context.Context{{if $soft}}, hardDelete bool{{end}}{{else}}{{if $soft}}hardDelete bool{{end}}{{end}}) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return o.DeleteAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}{{if $soft}}, hardDelete{{end}})
}

{{end -}}

{{if .AddPanic -}}
// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o {{$alias.UpSingular}}Slice) DeleteAllP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}{{if $soft}}, hardDelete bool{{end}}) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := o.DeleteAll({{if not .NoContext}}ctx, {{end -}} exec{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o {{$alias.UpSingular}}Slice) DeleteAllGP({{if not .NoContext}}ctx context.Context{{if $soft}}, hardDelete bool{{end}}{{else}}{{if $soft}}hardDelete bool{{end}}{{end}}) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := o.DeleteAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

{{end -}}

// DeleteAll deletes all rows in the slice, using an executor.
func (o {{$alias.UpSingular}}Slice) DeleteAll({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}{{if $soft}}, hardDelete bool{{end}}) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	if len(o) == 0 {
		return {{if not .NoRowsAffected}}0, {{end -}} nil
	}

	{{if not .NoHooks -}}
	if len({{$alias.DownSingular}}BeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
				return {{if not .NoRowsAffected}}0, {{end -}} err
			}
		}
	}
	{{- end}}

	{{if $soft -}}
	var (
		sql string
		args []interface{}
	)
	if hardDelete {
		for _, obj := range o {
    		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$alias.DownSingular}}PrimaryKeyMapping)
    		args = append(args, pkeyArgs...)
    	}
		sql = "DELETE FROM {{$schemaTable}} WHERE " +
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, {{$alias.DownSingular}}PrimaryKeyColumns, len(o))
	} else {
		currTime := time.Now().In(boil.GetLocation())
		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$alias.DownSingular}}PrimaryKeyMapping)
			args = append(args, pkeyArgs...)
			obj.DeletedAt = null.TimeFrom(currTime)
		}
		wl := []string{"{{or $.AutoColumns.Deleted "deleted_at"}}"}
		sql = fmt.Sprintf("UPDATE {{$schemaTable}} SET %s WHERE " +
			strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.UseIndexPlaceholders}}2{{else}}0{{end}}, {{$alias.DownSingular}}PrimaryKeyColumns, len(o)),
			strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, wl),
		)
		args = append([]interface{}{currTime}, args...)
	}
	{{else -}}
	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$alias.DownSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM {{$schemaTable}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, {{$alias.DownSingular}}PrimaryKeyColumns, len(o))
	{{- end}}

	{{if .NoContext -}}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}
	{{else -}}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
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
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{$alias.DownSingular}} slice")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by deleteall for {{.Table.Name}}")
	}

	{{end -}}

	{{if not .NoHooks -}}
	if len({{$alias.DownSingular}}AfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
				return {{if not .NoRowsAffected}}0, {{end -}} err
			}
		}
	}
	{{- end}}

	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
}

{{- end -}}
