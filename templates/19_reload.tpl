{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{if .AddGlobal -}}
// ReloadG refetches the object from the database using the primary keys.
func (o *{{$tableNameSingular}}) ReloadG({{if not .NoContext}}ctx context.Context{{end}}) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{$tableNameSingular}} provided for reload")
	}

	return o.Reload({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}})
}

{{end -}}

{{if .AddPanic -}}
// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *{{$tableNameSingular}}) ReloadP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) {
	if err := o.Reload({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// ReloadGP refetches the object from the database and panics on error.
func (o *{{$tableNameSingular}}) ReloadGP({{if not .NoContext}}ctx context.Context{{end}}) {
	if err := o.Reload({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *{{$tableNameSingular}}) Reload({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) error {
	ret, err := Find{{$tableNameSingular}}({{if not .NoContext}}ctx, {{end -}} exec, {{.Table.PKey.Columns | stringMap .StringFuncs.titleCase | prefixStringSlice "o." | join ", "}})
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

{{if .AddGlobal -}}
// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$tableNameSingular}}Slice) ReloadAllG({{if not .NoContext}}ctx context.Context{{end}}) error {
	if o == nil {
		return errors.New("{{.PkgName}}: empty {{$tableNameSingular}}Slice provided for reload all")
	}

	return o.ReloadAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}})
}

{{end -}}

{{if .AddPanic -}}
// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *{{$tableNameSingular}}Slice) ReloadAllP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) {
	if err := o.ReloadAll({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *{{$tableNameSingular}}Slice) ReloadAllGP({{if not .NoContext}}ctx context.Context{{end}}) {
	if err := o.ReloadAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$tableNameSingular}}Slice) ReloadAll({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	{{$varNamePlural}} := {{$tableNameSingular}}Slice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$varNameSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT {{$schemaTable}}.* FROM {{$schemaTable}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind({{if .NoContext}}nil{{else}}ctx{{end}}, exec, &{{$varNamePlural}})
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to reload all in {{$tableNameSingular}}Slice")
	}

	*o = {{$varNamePlural}}

	return nil
}
