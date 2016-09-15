{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
// OneP returns a single {{$varNameSingular}} record from the query, and panics on error.
func (q {{$varNameSingular}}Query) OneP() (*{{$tableNameSingular}}) {
	o, err := q.One()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// One returns a single {{$varNameSingular}} record from the query.
func (q {{$varNameSingular}}Query) One() (*{{$tableNameSingular}}, error) {
	o := &{{$tableNameSingular}}{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to execute a one query for {{.Table.Name}}")
	}

	{{if not .NoHooks -}}
	if err := o.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
		return o, err
	}
	{{- end}}

	return o, nil
}

// AllP returns all {{$tableNameSingular}} records from the query, and panics on error.
func (q {{$varNameSingular}}Query) AllP() {{$tableNameSingular}}Slice {
	o, err := q.All()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return o
}

// All returns all {{$tableNameSingular}} records from the query.
func (q {{$varNameSingular}}Query) All() ({{$tableNameSingular}}Slice, error) {
	var o {{$tableNameSingular}}Slice

	err := q.Bind(&o)
	if err != nil {
		return nil, errors.Wrap(err, "{{.PkgName}}: failed to assign all query results to {{$tableNameSingular}} slice")
	}

	{{if not .NoHooks -}}
	if len({{$varNameSingular}}AfterSelectHooks) != 0 {
		for _, obj := range o {
			if err := obj.doAfterSelectHooks(queries.GetExecutor(q.Query)); err != nil {
				return o, err
			}
		}
	}
	{{- end}}

	return o, nil
}

// CountP returns the count of all {{$tableNameSingular}} records in the query, and panics on error.
func (q {{$varNameSingular}}Query) CountP() int64 {
	c, err := q.Count()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return c
}

// Count returns the count of all {{$tableNameSingular}} records in the query.
func (q {{$varNameSingular}}Query) Count() (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "{{.PkgName}}: failed to count {{.Table.Name}} rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table, and panics on error.
func (q {{$varNameSingular}}Query) ExistsP() bool {
	e, err := q.Exists()
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q {{$varNameSingular}}Query) Exists() (bool, error) {
	var count int64

	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRow().Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "{{.PkgName}}: failed to check if {{.Table.Name}} exists")
	}

	return count > 0, nil
}
