{{- if .Table.IsView -}}
{{- else -}}
{{- $alias := .Aliases.Table .Table.Name -}}
{{- $schemaTable := .Table.Name | .SchemaTable -}}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted -}}
{{- $soft := and .AddSoftDeletes $canSoftDelete }}
{{- $softDelCol := or $.AutoColumns.Deleted "deleted_at"}}
{{if .AddGlobal -}}
// DeleteG deletes a single {{$alias.UpSingular}} record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *{{$alias.UpSingular}}) DeleteG(ctx context.Context{{if $soft}}, hardDelete bool{{end}}) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB(){{if $soft}}, hardDelete{{end}})
}

{{end -}}

{{if .AddPanic -}}
// DeleteP deletes a single {{$alias.UpSingular}} record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$alias.UpSingular}}) DeleteP(ctx context.Context, exec boil.ContextExecutor{{if $soft}}, hardDelete bool{{end}}) int64 {
	rowsAff, err := o.Delete(ctx, exec{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// DeleteGP deletes a single {{$alias.UpSingular}} record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$alias.UpSingular}}) DeleteGP(ctx context.Context{{if $soft}}, hardDelete bool{{end}}) int64 {
	rowsAff, err := o.Delete(ctx, boil.GetContextDB(){{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

// Delete deletes a single {{$alias.UpSingular}} record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$alias.UpSingular}}) Delete(ctx context.Context, exec boil.ContextExecutor{{if $soft}}, hardDelete bool{{end}}) (int64, error) {
	if o == nil {
		return 0, errors.New("{{.PkgName}}: no {{$alias.UpSingular}} provided for delete")
	}

	{{if not .NoHooks -}}
	if err := o.doBeforeDeleteHooks(ctx, exec); err != nil {
		return 0, err
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
		o.{{$alias.Column $softDelCol}} = null.TimeFrom(currTime)
		wl := []string{"{{$softDelCol}}"}
		sql = fmt.Sprintf("UPDATE {{$schemaTable}} SET %s WHERE {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 2 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}",
			strmangle.SetParamNames("{{.LQ}}", "{{.RQ}}", {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, wl),
		)
		valueMapping, err := queries.BindMapping({{$alias.DownSingular}}Type, {{$alias.DownSingular}}Mapping, append(wl, {{$alias.DownSingular}}PrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
		args = queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), valueMapping)
	}
	{{else -}}
	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), {{$alias.DownSingular}}PrimaryKeyMapping)
	sql := "DELETE FROM {{$schemaTable}} WHERE {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}"
	{{- end}}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to delete from {{.Table.Name}}")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by delete for {{.Table.Name}}")
	}

	{{if not .NoHooks -}}
	if err := o.doAfterDeleteHooks(ctx, exec); err != nil {
		return 0, err
	}
	{{- end}}

	return rowsAff, nil
}

{{if .AddGlobal -}}
func (q {{$alias.DownSingular}}Query) DeleteAllG(ctx context.Context{{if $soft}}, hardDelete bool{{end}}) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB(){{if $soft}}, hardDelete{{end}})
}

{{end -}}

{{if .AddPanic -}}
// DeleteAllP deletes all rows, and panics on error.
func (q {{$alias.DownSingular}}Query) DeleteAllP(ctx context.Context, exec boil.ContextExecutor{{if $soft}}, hardDelete bool{{end}}) int64 {
	rowsAff, err := q.DeleteAll(ctx, exec{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// DeleteAllGP deletes all rows, and panics on error.
func (q {{$alias.DownSingular}}Query) DeleteAllGP(ctx context.Context{{if $soft}}, hardDelete bool{{end}}) int64 {
	rowsAff, err := q.DeleteAll(ctx, boil.GetContextDB(){{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

// DeleteAll deletes all matching rows.
func (q {{$alias.DownSingular}}Query) DeleteAll(ctx context.Context, exec boil.ContextExecutor{{if $soft}}, hardDelete bool{{end}}) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("{{.PkgName}}: no {{$alias.DownSingular}}Query provided for delete all")
	}

	{{if $soft -}}
	if hardDelete {
		queries.SetDelete(q.Query)
	} else {
		currTime := time.Now().In(boil.GetLocation())
		queries.SetUpdate(q.Query, M{"{{$softDelCol}}": currTime})
	}
	{{else -}}
	queries.SetDelete(q.Query)
	{{- end}}

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{.Table.Name}}")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by deleteall for {{.Table.Name}}")
	}

	return rowsAff, nil
}

{{if .AddGlobal -}}
// DeleteAllG deletes all rows in the slice.
func (o {{$alias.UpSingular}}Slice) DeleteAllG(ctx context.Context{{if $soft}}, hardDelete bool{{end}}) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB(){{if $soft}}, hardDelete{{end}})
}

{{end -}}

{{if .AddPanic -}}
// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o {{$alias.UpSingular}}Slice) DeleteAllP(ctx context.Context, exec boil.ContextExecutor{{if $soft}}, hardDelete bool{{end}}) int64 {
	rowsAff, err := o.DeleteAll(ctx, exec{{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o {{$alias.UpSingular}}Slice) DeleteAllGP(ctx context.Context{{if $soft}}, hardDelete bool{{end}}) int64 {
	rowsAff, err := o.DeleteAll(ctx, boil.GetContextDB(){{if $soft}}, hardDelete{{end}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return rowsAff
}

{{end -}}

// DeleteAll deletes all rows in the slice, using an executor.
func (o {{$alias.UpSingular}}Slice) DeleteAll(ctx context.Context, exec boil.ContextExecutor{{if $soft}}, hardDelete bool{{end}}) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	{{if not .NoHooks -}}
	if len({{$alias.DownSingular}}BeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(ctx, exec); err != nil {
				return 0, err
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
			obj.{{$alias.Column $softDelCol}} = null.TimeFrom(currTime)
		}
		wl := []string{"{{$softDelCol}}"}
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

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{$alias.DownSingular}} slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by deleteall for {{.Table.Name}}")
	}

	{{if not .NoHooks -}}
	if len({{$alias.DownSingular}}AfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(ctx, exec); err != nil {
				return 0, err
			}
		}
	}
	{{- end}}

	return rowsAff, nil
}

{{- end -}}
