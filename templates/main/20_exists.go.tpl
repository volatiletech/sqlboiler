{{- if .Table.IsView -}}
{{- else -}}
{{- $alias := .Aliases.Table .Table.Name -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap (aliasCols $alias) | stringMap .StringFuncs.camelCase | stringMap .StringFuncs.replaceReserved -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", " -}}
{{- $schemaTable := .Table.Name | .SchemaTable -}}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted }}
{{if .AddGlobal -}}
// {{$alias.UpSingular}}ExistsG checks if the {{$alias.UpSingular}} row exists.
func {{$alias.UpSingular}}ExistsG({{if not .NoContext}}ctx context.Context, {{end -}} {{$pkArgs}}) (bool, error) {
	return {{$alias.UpSingular}}Exists({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, {{$pkNames | join ", "}})
}

{{end -}}

{{if .AddPanic -}}
// {{$alias.UpSingular}}ExistsP checks if the {{$alias.UpSingular}} row exists. Panics on error.
func {{$alias.UpSingular}}ExistsP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, {{$pkArgs}}) bool {
	e, err := {{$alias.UpSingular}}Exists({{if not .NoContext}}ctx, {{end -}} exec, {{$pkNames | join ", "}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// {{$alias.UpSingular}}ExistsGP checks if the {{$alias.UpSingular}} row exists. Panics on error.
func {{$alias.UpSingular}}ExistsGP({{if not .NoContext}}ctx context.Context, {{end -}} {{$pkArgs}}) bool {
	e, err := {{$alias.UpSingular}}Exists({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, {{$pkNames | join ", "}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

{{end -}}

// {{$alias.UpSingular}}Exists checks if the {{$alias.UpSingular}} row exists.
func {{$alias.UpSingular}}Exists({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, {{$pkArgs}}) (bool, error) {
	var exists bool
	{{if .Dialect.UseCaseWhenExistsClause -}}
	sql := "select case when exists(select top(1) 1 from {{$schemaTable}} where {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}) then 1 else 0 end"
	{{- else -}}
	sql := "select exists(select 1 from {{$schemaTable}} where {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}{{if and .AddSoftDeletes $canSoftDelete}} and {{or $.AutoColumns.Deleted "deleted_at" | $.Quotes}} is null{{end}} limit 1)"
	{{- end}}

	{{if .NoContext -}}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, {{$pkNames | join ", "}})
	}
	{{else -}}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, {{$pkNames | join ", "}})
	}
	{{end -}}

	{{if .NoContext -}}
	row := exec.QueryRow(sql, {{$pkNames | join ", "}})
	{{else -}}
	row := exec.QueryRowContext(ctx, sql, {{$pkNames | join ", "}})
	{{- end}}

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: unable to check if {{.Table.Name}} exists")
	}

	return exists, nil
}

// Exists checks if the {{$alias.UpSingular}} row exists.
func (o *{{$alias.UpSingular}}) Exists({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (bool, error) {
	return {{$alias.UpSingular}}Exists({{if .NoContext}}exec{{else}}ctx, exec{{end}}, o.{{$.Table.PKey.Columns | stringMap (aliasCols $alias) | join ", o."}})
}

{{- end -}}
