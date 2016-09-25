{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
// ReloadGP refetches the object from the database and panics on error.
func (o *{{$tableNameSingular}}) ReloadGP() {
	if err := o.ReloadG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *{{$tableNameSingular}}) ReloadP(exec boil.Executor) {
	if err := o.Reload(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadG refetches the object from the database using the primary keys.
func (o *{{$tableNameSingular}}) ReloadG() error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for reload")
	}

	return o.Reload(boil.GetDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *{{$tableNameSingular}}) Reload(exec boil.Executor) error {
	ret, err := Find{{$tableNameSingular}}(exec, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *{{$tableNameSingular}}Slice) ReloadAllGP() {
	if err := o.ReloadAllG(); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *{{$tableNameSingular}}Slice) ReloadAllP(exec boil.Executor) {
	if err := o.ReloadAll(exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$tableNameSingular}}Slice) ReloadAllG() error {
	if o == nil {
		return errors.New("{{.PkgName}}: empty {{$tableNameSingular}}Slice provided for reload all")
	}

	return o.ReloadAll(boil.GetDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$tableNameSingular}}Slice) ReloadAll(exec boil.Executor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	{{$varNamePlural}} := {{$tableNameSingular}}Slice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf(
		"SELECT {{$schemaTable}}.* FROM {{$schemaTable}} WHERE (%s) IN (%s)",
		strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, {{$varNameSingular}}PrimaryKeyColumns), ","),
		strmangle.Placeholders(dialect.IndexPlaceholders, len(*o) * len({{$varNameSingular}}PrimaryKeyColumns), 1, len({{$varNameSingular}}PrimaryKeyColumns)),
	)

	q := queries.Raw(exec, sql, args...)

	err := q.Bind(&{{$varNamePlural}})
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to reload all in {{$tableNameSingular}}Slice")
	}

	*o = {{$varNamePlural}}

	return nil
}
