{{- $alias := .Aliases.Table .Table.Name -}}
{{- $schemaTable := .Table.Name | .SchemaTable -}}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted }}
{{if .AddGlobal -}}
// ReloadG refetches the object from the database using the primary keys.
func (o *{{$alias.UpSingular}}) ReloadG({{if not .NoContext}}ctx context.Context{{end}}) error {
	if o == nil {
		return errors.New("{{.PkgName}}: no {{$alias.UpSingular}} provided for reload")
	}

	return o.Reload({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}})
}

{{end -}}

{{if .AddPanic -}}
// ReloadP refetches the object from the database with an executor. Panics on error.
func (o *{{$alias.UpSingular}}) ReloadP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) {
	if err := o.Reload({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// ReloadGP refetches the object from the database and panics on error.
func (o *{{$alias.UpSingular}}) ReloadGP({{if not .NoContext}}ctx context.Context{{end}}) {
	if err := o.Reload({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *{{$alias.UpSingular}}) Reload({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) error {
	ret, err := Find{{$alias.UpSingular}}({{if not .NoContext}}ctx, {{end -}} exec, {{.Table.PKey.Columns | stringMap (aliasCols $alias) | prefixStringSlice "o." | join ", "}})
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

{{if .AddGlobal -}}
// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$alias.UpSingular}}Slice) ReloadAllG({{if not .NoContext}}ctx context.Context{{end}}) error {
	if o == nil {
		return errors.New("{{.PkgName}}: empty {{$alias.UpSingular}}Slice provided for reload all")
	}

	return o.ReloadAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}})
}

{{end -}}

{{if .AddPanic -}}
// ReloadAllP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *{{$alias.UpSingular}}Slice) ReloadAllP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) {
	if err := o.ReloadAll({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// ReloadAllGP refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
// Panics on error.
func (o *{{$alias.UpSingular}}Slice) ReloadAllGP({{if not .NoContext}}ctx context.Context{{end}}) {
	if err := o.ReloadAll({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *{{$alias.UpSingular}}Slice) ReloadAll({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := {{$alias.UpSingular}}Slice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), {{$alias.DownSingular}}PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT {{$schemaTable}}.* FROM {{$schemaTable}} WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), {{if .Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, {{$alias.DownSingular}}PrimaryKeyColumns, len(*o)){{if and .AddSoftDeletes $canSoftDelete}} +
		"and {{"deleted_at" | $.Quotes}} is null"
		{{- end}}

	q := queries.Raw(sql, args...)

	err := q.Bind({{if .NoContext}}nil{{else}}ctx{{end}}, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "{{.PkgName}}: unable to reload all in {{$alias.UpSingular}}Slice")
	}

	*o = slice

	return nil
}
