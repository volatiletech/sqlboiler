{{- if .Table.IsView -}}
{{- else -}}
{{- $alias := .Aliases.Table .Table.Name -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap (aliasCols $alias) | stringMap .StringFuncs.camelCase | stringMap .StringFuncs.replaceReserved -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", " -}}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted }}
{{if .AddGlobal -}}
// Find{{$alias.UpSingular}}G retrieves a single record by ID.
func Find{{$alias.UpSingular}}G({{if not .NoContext}}ctx context.Context, {{end -}} {{$pkArgs}}, selectCols ...string) (*{{$alias.UpSingular}}, error) {
	return Find{{$alias.UpSingular}}({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, {{$pkNames | join ", "}}, selectCols...)
}

{{end -}}

{{if .AddPanic -}}
// Find{{$alias.UpSingular}}P retrieves a single record by ID with an executor, and panics on error.
func Find{{$alias.UpSingular}}P({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, {{$pkArgs}}, selectCols ...string) *{{$alias.UpSingular}} {
	retobj, err := Find{{$alias.UpSingular}}({{if not .NoContext}}ctx, {{end -}} exec, {{$pkNames | join ", "}}, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// Find{{$alias.UpSingular}}GP retrieves a single record by ID, and panics on error.
func Find{{$alias.UpSingular}}GP({{if not .NoContext}}ctx context.Context, {{end -}} {{$pkArgs}}, selectCols ...string) *{{$alias.UpSingular}} {
	retobj, err := Find{{$alias.UpSingular}}({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, {{$pkNames | join ", "}}, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

{{end -}}

// Find{{$alias.UpSingular}} retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func Find{{$alias.UpSingular}}(ctx context.Context, exec boil.ContextExecutor, {{$pkArgs}}, selectCols ...string) (*{{$alias.UpSingular}}, error) {
	{{if eq (len .Table.PKey.Columns) 1 -}}
	return helpers.Find[*{{$alias.UpSingular}}, {{$alias.DownSingular}}Table, {{$alias.DownSingular}}Hooks](ctx, exec, {{$alias.UpSingular}}PK({{index $pkNames 0}}), selectCols...)
	{{- else -}}
		pk := {{$alias.UpSingular}}PK{
		{{range $index, $pkcol := .Table.PKey.Columns -}}
		{{- $column := $.Table.GetColumn $pkcol -}}
		{{$alias.Column $column.Name}}: {{index $pkNames $index}},
		{{end -}}
		}
	return helpers.Find[*{{$alias.UpSingular}}, {{$alias.DownSingular}}Table, {{$alias.DownSingular}}Hooks](ctx, exec, pk, selectCols...)
	{{- end}}
}

{{- end -}}
