{{- if .Table.IsView -}}
{{- else -}}
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

{{- end -}}
