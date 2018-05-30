{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase | stringMap .StringFuncs.replaceReserved -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", "}}
{{if .AddGlobal -}}
// Find{{$tableNameSingular}}G retrieves a single record by ID.
func Find{{$tableNameSingular}}G({{if not .NoContext}}ctx context.Context, {{end -}} {{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
	return Find{{$tableNameSingular}}({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, {{$pkNames | join ", "}}, selectCols...)
}

{{end -}}

{{if .AddPanic -}}
// Find{{$tableNameSingular}}P retrieves a single record by ID with an executor, and panics on error.
func Find{{$tableNameSingular}}P({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, {{$pkArgs}}, selectCols ...string) *{{$tableNameSingular}} {
	retobj, err := Find{{$tableNameSingular}}({{if not .NoContext}}ctx, {{end -}} exec, {{$pkNames | join ", "}}, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// Find{{$tableNameSingular}}GP retrieves a single record by ID, and panics on error.
func Find{{$tableNameSingular}}GP({{if not .NoContext}}ctx context.Context, {{end -}} {{$pkArgs}}, selectCols ...string) *{{$tableNameSingular}} {
	retobj, err := Find{{$tableNameSingular}}({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, {{$pkNames | join ", "}}, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

{{end -}}

// Find{{$tableNameSingular}} retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func Find{{$tableNameSingular}}({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, {{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
	{{$varNameSingular}}Obj := &{{$tableNameSingular}}{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from {{.Table.Name | .SchemaTable}} where {{if .Dialect.UseIndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}", sel,
	)

	q := queries.Raw(query, {{$pkNames | join ", "}})

	err := q.Bind({{if not .NoContext}}ctx{{else}}nil{{end}}, exec, {{$varNameSingular}}Obj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "{{.PkgName}}: unable to select from {{.Table.Name}}")
	}

	return {{$varNameSingular}}Obj, nil
}
