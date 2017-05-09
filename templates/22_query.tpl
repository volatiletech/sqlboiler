{{- $tableNameSingular := .Table.Name | singular | titleCase -}}

// {{$tableNameSingular}}NewQuery filters query results
func {{$tableNameSingular}}NewQuery(exec boil.Executor) *{{$tableNameSingular}}Query {
	return &{{$tableNameSingular}}Query{NewQuery(exec, qm.Select("*"), qm.From("{{.Table.Name | .SchemaTable}}"))}
}

// {{$tableNameSingular}}NewQuery filters query results
func {{$tableNameSingular}}NewQueryG() *{{$tableNameSingular}}Query {
	return {{$tableNameSingular}}NewQuery(boil.GetDB())
}

// Where filters query results
func (q *{{$tableNameSingular}}Query) Where(filters {{$tableNameSingular}}Filter) *{{$tableNameSingular}}Query {
	r := reflect.ValueOf(filters)
	for i := 0; i < r.NumField(); i++ {
		f := r.Field(i)
		if f.Elem().IsValid() {
			if nullable, ok := f.Elem().Interface().(Nullable); ok && nullable.IsZero() {
				queries.AppendWhere(q.Query, r.Type().Field(i).Tag.Get("boil")+" IS NULL")
			} else {
				queries.AppendWhere(q.Query, r.Type().Field(i).Tag.Get("boil")+" = ?", f.Elem().Interface())
			}
		}
	}
	return q
}

// Limit limits query results
func (q *{{$tableNameSingular}}Query) Limit(limit int) *{{$tableNameSingular}}Query {
	queries.SetLimit(q.Query, limit)
	return q
}