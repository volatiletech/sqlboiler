{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// DeleteP deletes a single {{$tableNameSingular}} record with an executor.
// DeleteP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$tableNameSingular}}) DeleteP(exec boil.Executor) {
	if err := o.Delete(exec); err != nil {
	panic(boil.WrapErr(err))
	}
}

// DeleteG deletes a single {{$tableNameSingular}} record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) DeleteG() error {
	if o == nil {
	return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for deletion")
	}

	return o.Delete(boil.GetDB())
}

// DeleteGP deletes a single {{$tableNameSingular}} record.
// DeleteGP will match against the primary key column to find the record to delete.
// Panics on error.
func (o *{{$tableNameSingular}}) DeleteGP() {
	if err := o.DeleteG(); err != nil {
	panic(boil.WrapErr(err))
	}
}

// Delete deletes a single {{$tableNameSingular}} record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *{{$tableNameSingular}}) Delete(exec boil.Executor) error {
	if o == nil {
	return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for delete")
	}

	{{if not .NoHooks -}}
	if err := o.doBeforeDeleteHooks(exec); err != nil {
	return err
	}
	{{- end}}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), {{$varNameSingular}}PrimaryKeyMapping)
	sql := "DELETE FROM {{$schemaTable}} WHERE {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}"

	if boil.DebugMode {
	fmt.Fprintln(boil.DebugWriter, sql)
	fmt.Fprintln(boil.DebugWriter, args...)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
	return errors.Wrap(err, "{{.PkgName}}: unable to delete from {{.Table.Name}}")
	}

	{{if not .NoHooks -}}
	if err := o.doAfterDeleteHooks(exec); err != nil {
	return err
	}
	{{- end}}

	return nil
}

// DeleteAllP deletes all rows, and panics on error.
func (q {{$varNameSingular}}Query) DeleteAllP() {
	if err := q.DeleteAll(); err != nil {
	panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all matching rows.
func (q {{$varNameSingular}}Query) DeleteAll() error {
	if q.Query == nil {
	return errors.New("{{.PkgName}}: no {{$varNameSingular}}Query provided for delete all")
	}

	queries.SetDelete(q.Query)

	_, err := q.Query.Exec()
	if err != nil {
	return errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{.Table.Name}}")
	}

	return nil
}

// DeleteAllGP deletes all rows in the slice, and panics on error.
func (o {{$tableNameSingular}}Slice) DeleteAllGP() {
	if err := o.DeleteAllG(); err != nil {
	panic(boil.WrapErr(err))
	}
}

// DeleteAllG deletes all rows in the slice.
func (o {{$tableNameSingular}}Slice) DeleteAllG() error {
	if o == nil {
	return errors.New("{{.PkgName}}: no {{$tableNameSingular}} slice provided for delete all")
	}
	return o.DeleteAll(boil.GetDB())
}

// DeleteAllP deletes all rows in the slice, using an executor, and panics on error.
func (o {{$tableNameSingular}}Slice) DeleteAllP(exec boil.Executor) {
	if err := o.DeleteAll(exec); err != nil {
	panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o {{$tableNameSingular}}Slice) DeleteAll(exec boil.Executor) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{$tableNameSingular}} slice provided for delete all")
	}

	if len(o) == 0 {
		return nil
	}

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}BeforeDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doBeforeDeleteHooks(exec); err != nil {
				return err
			}
		}
	}
	{{- end}}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"DELETE FROM {{$schemaTable}} WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, {{$varNameSingular}}PrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(o) * len({{$varNameSingular}}PrimaryKeyColumns), 1, len({{$varNameSingular}}PrimaryKeyColumns)),
	)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	_, err := exec.Exec(sql, args...)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to delete all from {{$varNameSingular}} slice")
	}

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}AfterDeleteHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterDeleteHooks(exec); err != nil {
				return err
			}
		}
	}
	{{- end}}

	return nil
}
