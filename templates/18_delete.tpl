{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// DeleteG deletes a single {{$tableNameSingular}} record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) DeleteG() {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return o.Delete(boil.GetDB())
}

// DeleteP deletes a single {{$tableNameSingular}} record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$tableNameSingular}}) DeleteP(exec boil.Executor) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end}}err := o.Delete(exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

// DeleteGP deletes a single {{$tableNameSingular}} record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$tableNameSingular}}) DeleteGP() {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end}}err := o.Delete(boil.GetDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

// Delete deletes a single {{$tableNameSingular}} record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) Delete(exec boil.Executor) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	if o == nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for delete")
	}

	{{if not .NoHooks -}}
	if err := o.doBeforeDeleteHooks(exec); err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} err
	}
	{{- end}}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), {{$varNameSingular}}PrimaryKeyMapping)
	sql := "DELETE FROM {{$schemaTable}} WHERE {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}"

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
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to delete from {{.Table.Name}}")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by delete for {{.Table.Name}}")
	}

	{{end -}}

	{{if not .NoHooks -}}
	if err := o.doAfterDeleteHooks(exec); err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} err
	}
	{{- end}}

	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q {{$varNameSingular}}Query) DeleteAllP() {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := q.DeleteAll()
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

// DeleteAll deletes all matching rows.
func (q {{$varNameSingular}}Query) DeleteAll() {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	if q.Query == nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.New("{{.PkgName}}: no {{$varNameSingular}}Query provided for delete all")
	}

	queries.SetDelete(q.Query)

	{{if .NoRowsAffected -}}
	_, err := q.Query.Exec()
	{{else -}}
	result, err := q.Query.Exec()
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

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o {{$tableNameSingular}}Slice) DeleteAllGP() {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := o.DeleteAll(boil.GetDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

// DeleteAllG deletes all rows in the slice.
func (o {{$tableNameSingular}}Slice) DeleteAllG() {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o {{$tableNameSingular}}Slice) DeleteAllP(exec boil.Executor) {{if not .NoRowsAffected}}int64{{end -}} {
	{{if not .NoRowsAffected}}rowsAff, {{end -}} err := o.DeleteAll(exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}
	{{- if not .NoRowsAffected}}

	return rowsAff
	{{end -}}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o {{$tableNameSingular}}Slice) DeleteAll(exec boil.Executor) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
	if o == nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.New("{{.PkgName}}: no {{$tableNameSingular}} slice provided for delete all")
	}

	if len(o) == 0 {
		return {{if not .NoRowsAffected}}0, {{end -}} nil
	}

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}BeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return {{if not .NoRowsAffected}}0, {{end -}} err
			}
		}
	}
	{{- end}}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM {{$schemaTable}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	{{if .NoRowsAffected -}}
	_, err := exec.Exec(sql, args...)
	{{else -}}
	result, err := exec.Exec(sql, args...)
	{{end -}}
	if err != nil {
		return {{if not .NoRowsAffected}}0, {{end -}} errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{$varNameSingular}} slice")
	}

	{{if not .NoRowsAffected -}}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to get rows affected by deleteall for {{.Table.Name}}")
	}

	{{end -}}

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}AfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return {{if not .NoRowsAffected}}0, {{end -}} err
			}
		}
	}
	{{- end}}

	return {{if not .NoRowsAffected}}rowsAff, {{end -}} nil
}
