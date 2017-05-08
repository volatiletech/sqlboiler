{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
{{- $colDefs := sqlColDefinitions .Table.Columns .Table.PKey.Columns -}}
{{- $pkNames := $colDefs.Names | stringMap .StringFuncs.camelCase | stringMap .StringFuncs.replaceReserved -}}
{{- $pkArgs := joinSlices " " $pkNames $colDefs.Types | join ", "}}
// Find{{$tableNameSingular}}G retrieves a single record by ID.
func Find{{$tableNameSingular}}G({{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
	return Find{{$tableNameSingular}}(boil.GetDB(), {{$pkNames | join ", "}}, selectCols...)
}

// Find{{$tableNameSingular}}GP retrieves a single record by ID, and panics on error.
func Find{{$tableNameSingular}}GP({{$pkArgs}}, selectCols ...string) *{{$tableNameSingular}} {
	retobj, err := Find{{$tableNameSingular}}(boil.GetDB(), {{$pkNames | join ", "}}, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// Find{{$tableNameSingular}} retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func Find{{$tableNameSingular}}(exec boil.Executor, {{$pkArgs}}, selectCols ...string) (*{{$tableNameSingular}}, error) {
	{{$varNameSingular}}Obj := &{{$tableNameSingular}}{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from {{.Table.Name | .SchemaTable}} where {{if .Dialect.IndexPlaceholders}}{{whereClause .LQ .RQ 1 .Table.PKey.Columns}}{{else}}{{whereClause .LQ .RQ 0 .Table.PKey.Columns}}{{end}}", sel,
	)

	q := queries.Raw(exec, query, {{$pkNames | join ", "}})

	err := q.Bind({{$varNameSingular}}Obj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "{{.PkgName}}: unable to select from {{.Table.Name}}")
	}

	return {{$varNameSingular}}Obj, nil
}

// Find{{$tableNameSingular}}P retrieves a single record by ID with an executor, and panics on error.
func Find{{$tableNameSingular}}P(exec boil.Executor, {{$pkArgs}}, selectCols ...string) *{{$tableNameSingular}} {
	retobj, err := Find{{$tableNameSingular}}(exec, {{$pkNames | join ", "}}, selectCols...)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return retobj
}

// FindOne{{$tableNameSingular}} retrieves a single record using filters.
func FindOne{{$tableNameSingular}}(exec boil.Executor, filters {{$tableNameSingular}}Filter) (*{{$tableNameSingular}}, error) {
	{{$varNameSingular}}Obj := &{{$tableNameSingular}}{}

	query := NewQuery(exec, qm.Select("*"), qm.From("{{.Table.Name | .SchemaTable}}"))

	r := reflect.ValueOf(filters)
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		if f.Elem().IsValid() {
			queries.AppendWhere(query, r.Type().Field(i).Tag.Get("boil")+" = ?", f.Elem().Interface())
		}
	}

	queries.SetLimit(query, 1)

	err := query.Bind({{$varNameSingular}}Obj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "{{.PkgName}}: unable to select from {{.Table.Name}}")
	}

	return {{$varNameSingular}}Obj, nil
}

// FindOne{{$tableNameSingular}}G retrieves a single record using filters.
func FindOne{{$tableNameSingular}}G(filters {{$tableNameSingular}}Filter) (*{{$tableNameSingular}}, error) {
	return FindOne{{$tableNameSingular}}(boil.GetDB(), filters)
}

// FindOne{{$tableNameSingular}}OrInit retrieves a single record using filters, or initializes a new record if one is not found.
func FindOne{{$tableNameSingular}}OrInit(exec boil.Executor, filters {{$tableNameSingular}}Filter) (*{{$tableNameSingular}}, error) {
	{{$varNameSingular}}Obj, err := FindOne{{$tableNameSingular}}(exec, filters)
	if err != nil && errors.Cause(err) != sql.ErrNoRows {
		return nil, err
	}

	if {{$varNameSingular}}Obj == nil {
		{{$varNameSingular}}Obj = &{{$tableNameSingular}}{}
		objR := reflect.ValueOf({{$varNameSingular}}Obj).Elem()
		r := reflect.ValueOf(filters)
		for i := 0; i < r.NumField(); i++ {
			f := r.Field(i)
			if f.Elem().IsValid() {
				objR.FieldByName(r.Type().Field(i).Name).Set(f.Elem())
			}
		}
	}

	return {{$varNameSingular}}Obj, nil
}

// FindOne{{$tableNameSingular}}OrInit retrieves a single record using filters, or initializes a new record if one is not found.
func FindOne{{$tableNameSingular}}OrInitG(filters {{$tableNameSingular}}Filter) (*{{$tableNameSingular}}, error) {
	return FindOne{{$tableNameSingular}}OrInit(boil.GetDB(), filters)
}

// FindOne{{$tableNameSingular}}OrInit retrieves a single record using filters, or initializes and inserts a new record if one is not found.
func FindOne{{$tableNameSingular}}OrCreate(exec boil.Executor, filters {{$tableNameSingular}}Filter) (*{{$tableNameSingular}}, error) {
	{{$varNameSingular}}Obj, err := FindOne{{$tableNameSingular}}OrInit(exec, filters)
	if err != nil {
		return nil, err
	}
	if {{$varNameSingular}}Obj.IsNew() {
		err := {{$varNameSingular}}Obj.Insert(exec)
		if err != nil {
			return nil, err
		}
	}
	return {{$varNameSingular}}Obj, nil
}

// FindOne{{$tableNameSingular}}OrInit retrieves a single record using filters, or initializes and inserts a new record if one is not found.
func FindOne{{$tableNameSingular}}OrCreateG(filters {{$tableNameSingular}}Filter) (*{{$tableNameSingular}}, error) {
	return FindOne{{$tableNameSingular}}OrCreate(boil.GetDB(), filters)
}
