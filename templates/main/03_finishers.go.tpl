{{- $alias := .Aliases.Table .Table.Name}}

{{if .AddGlobal -}}
// OneG returns a single {{$alias.DownSingular}} record from the query using the global executor.
func (q {{$alias.DownSingular}}Query) OneG(ctx context.Context) (*{{$alias.UpSingular}}, error) {
	return q.One({ctx, boil.GetContextDB())
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// OneGP returns a single {{$alias.DownSingular}} record from the query using the global executor, and panics on error.
func (q {{$alias.DownSingular}}Query) OneGP(ctx context.Context) *{{$alias.UpSingular}} {
	o, err := q.One(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

{{end -}}

{{if .AddPanic -}}
// OneP returns a single {{$alias.DownSingular}} record from the query, and panics on error.
func (q {{$alias.DownSingular}}Query) OneP(ctx context.Context, exec boil.ContextExecutor) (*{{$alias.UpSingular}}) {
	o, err := q.One(ctx, exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

{{end -}}

// One returns a single {{$alias.DownSingular}} record from the query.
func (q {{$alias.DownSingular}}Query) One(ctx context.Context, exec boil.ContextExecutor) (*{{$alias.UpSingular}}, error) {
	o := &{{$alias.UpSingular}}{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		{{if not .AlwaysWrapErrors -}}
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		{{end -}}
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to execute a one query for {{.Table.Name}}")
	}

	if err := o.doAfterSelectHooks(ctx, exec); err != nil {
		return o, err
	}

	return o, nil
}

{{if .AddGlobal -}}
// AllG returns all {{$alias.UpSingular}} records from the query using the global executor.
func (q {{$alias.DownSingular}}Query) AllG(ctx context.Context) ({{$alias.UpSingular}}Slice, error) {
	return q.All(ctx, boil.GetContextDB())
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// AllGP returns all {{$alias.UpSingular}} records from the query using the global executor, and panics on error.
func (q {{$alias.DownSingular}}Query) AllGP(ctx context.Context) {{$alias.UpSingular}}Slice {
	o, err := q.All(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

{{end -}}

{{if .AddPanic -}}
// AllP returns all {{$alias.UpSingular}} records from the query, and panics on error.
func (q {{$alias.DownSingular}}Query) AllP(ctx context.Context, exec boil.ContextExecutor) {{$alias.UpSingular}}Slice {
	o, err := q.All(ctx,  exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

{{end -}}

// All returns all {{$alias.UpSingular}} records from the query.
func (q {{$alias.DownSingular}}Query) All(ctx context.Context, exec boil.ContextExecutor) ({{$alias.UpSingular}}Slice, error) {
	var o []*{{$alias.UpSingular}}

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to assign all query results to {{$alias.UpSingular}} slice")
	}

	if len({{$alias.DownSingular}}AfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(ctx,  exec); err != nil {
				return o, err
			}
		}
	}

	return o, nil
}

{{if .AddGlobal -}}
// CountG returns the count of all {{$alias.UpSingular}} records in the query using the global executor
func (q {{$alias.DownSingular}}Query) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// CountGP returns the count of all {{$alias.UpSingular}} records in the query using the global executor, and panics on error.
func (q {{$alias.DownSingular}}Query) CountGP(ctx context.Context) int64 {
	c, err := q.Count(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

{{end -}}

{{if .AddPanic -}}
// CountP returns the count of all {{$alias.UpSingular}} records in the query, and panics on error.
func (q {{$alias.DownSingular}}Query) CountP(ctx context.Context, exec boil.ContextExecutor) int64 {
	c, err := q.Count(ctx,  exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

{{end -}}

// Count returns the count of all {{$alias.UpSingular}} records in the query.
func (q {{$alias.DownSingular}}Query) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to count {{.Table.Name}} rows")
	}

	return count, nil
}

{{if .AddGlobal -}}
// ExistsG checks if the row exists in the table using the global executor.
func (q {{$alias.DownSingular}}Query) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

{{end -}}

{{if and .AddGlobal .AddPanic -}}
// ExistsGP checks if the row exists in the table using the global executor, and panics on error.
func (q {{$alias.DownSingular}}Query) ExistsGP(ctx context.Context) bool {
	e, err := q.Exists(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

{{end -}}

{{if .AddPanic -}}
// ExistsP checks if the row exists in the table, and panics on error.
func (q {{$alias.DownSingular}}Query) ExistsP(ctx context.Context, exec boil.ContextExecutor) bool {
	e, err := q.Exists(ctx,  exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

{{end -}}

// Exists checks if the row exists in the table.
func (q {{$alias.DownSingular}}Query) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: failed to check if {{.Table.Name}} exists")
	}

	return count > 0, nil
}
