{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{if .AddGlobal -}}
// OneG returns a single {{$varNameSingular}} record from the query using the global executor.
func (q {{$varNameSingular}}Query) OneG({{if not .NoContext}}ctx context.Context{{end}}) (*{{$tableNameSingular}}, error) {
	return q.One({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end -}})
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// OneGP returns a single {{$varNameSingular}} record from the query using the global executor, and panics on error.
func (q {{$varNameSingular}}Query) OneGP({{if not .NoContext}}ctx context.Context{{end}}) *{{$tableNameSingular}} {
	o, err := q.One({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end -}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

{{end -}}

{{if .AddPanic -}}
// OneP returns a single {{$varNameSingular}} record from the query, and panics on error.
func (q {{$varNameSingular}}Query) OneP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (*{{$tableNameSingular}}) {
	o, err := q.One({{if not .NoContext}}ctx, {{end -}} exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

{{end -}}

// One returns a single {{$varNameSingular}} record from the query.
func (q {{$varNameSingular}}Query) One({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (*{{$tableNameSingular}}, error) {
	o := &{{$tableNameSingular}}{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind({{if .NoContext}}nil{{else}}ctx{{end}}, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to execute a one query for {{.Table.Name}}")
	}

	{{if not .NoHooks -}}
	if err := o.doAfterSelectHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
		return o, err
	}
	{{- end}}

	return o, nil
}

{{if .AddGlobal -}}
// AllG returns all {{$tableNameSingular}} records from the query using the global executor.
func (q {{$varNameSingular}}Query) AllG({{if not .NoContext}}ctx context.Context{{end}}) ({{$tableNameSingular}}Slice, error) {
	return q.All({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end -}})
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// AllGP returns all {{$tableNameSingular}} records from the query using the global executor, and panics on error.
func (q {{$varNameSingular}}Query) AllGP({{if not .NoContext}}ctx context.Context{{end}}) {{$tableNameSingular}}Slice {
	o, err := q.All({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end -}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

{{end -}}

{{if .AddPanic -}}
// AllP returns all {{$tableNameSingular}} records from the query, and panics on error.
func (q {{$varNameSingular}}Query) AllP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) {{$tableNameSingular}}Slice {
	o, err := q.All({{if not .NoContext}}ctx, {{end -}} exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

{{end -}}

// All returns all {{$tableNameSingular}} records from the query.
func (q {{$varNameSingular}}Query) All({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) ({{$tableNameSingular}}Slice, error) {
	var o []*{{$tableNameSingular}}

	err := q.Bind({{if .NoContext}}nil{{else}}ctx{{end}}, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to assign all query results to {{$tableNameSingular}} slice")
	}

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}AfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks({{if not .NoContext}}ctx, {{end -}} exec); err != nil {
				return o, err
			}
		}
	}
	{{- end}}

	return o, nil
}

{{if .AddGlobal -}}
// CountG returns the count of all {{$tableNameSingular}} records in the query, and panics on error.
func (q {{$varNameSingular}}Query) CountG({{if not .NoContext}}ctx context.Context{{end}}) (int64, error) {
	return q.Count({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end -}})
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// CountGP returns the count of all {{$tableNameSingular}} records in the query using the global executor, and panics on error.
func (q {{$varNameSingular}}Query) CountGP({{if not .NoContext}}ctx context.Context{{end}}) int64 {
	c, err := q.Count({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end -}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

{{end -}}

{{if .AddPanic -}}
// CountP returns the count of all {{$tableNameSingular}} records in the query, and panics on error.
func (q {{$varNameSingular}}Query) CountP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) int64 {
	c, err := q.Count({{if not .NoContext}}ctx, {{end -}} exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

{{end -}}

// Count returns the count of all {{$tableNameSingular}} records in the query.
func (q {{$varNameSingular}}Query) Count({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	{{if .NoContext -}}
	err := q.Query.QueryRow(exec).Scan(&count)
	{{else -}}
	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	{{end -}}
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to count {{.Table.Name}} rows")
	}

	return count, nil
}

{{if .AddGlobal -}}
// ExistsG checks if the row exists in the table, and panics on error.
func (q {{$varNameSingular}}Query) ExistsG({{if not .NoContext}}ctx context.Context{{end}}) (bool, error) {
	return q.Exists({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end -}})
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// ExistsGP checks if the row exists in the table using the global executor, and panics on error.
func (q {{$varNameSingular}}Query) ExistsGP({{if not .NoContext}}ctx context.Context{{end}}) bool {
	e, err := q.Exists({{if .NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end -}})
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

{{end -}}

{{if .AddPanic -}}
// ExistsP checks if the row exists in the table, and panics on error.
func (q {{$varNameSingular}}Query) ExistsP({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) bool {
	e, err := q.Exists({{if not .NoContext}}ctx, {{end -}} exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

{{end -}}

// Exists checks if the row exists in the table.
func (q {{$varNameSingular}}Query) Exists({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	{{if .NoContext -}}
	err := q.Query.QueryRow(exec).Scan(&count)
	{{else -}}
	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	{{end -}}
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: failed to check if {{.Table.Name}} exists")
	}

	return count > 0, nil
}
