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

{{- end -}}
