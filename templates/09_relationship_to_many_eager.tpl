{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $varNameSingular := $dot.Table.Name | singular | camelCase -}}
		{{- $txt := txtsFromToMany $dot.Tables $dot.Table . -}}
		{{- $arg := printf "maybe%s" $txt.LocalTable.NameGo -}}
		{{- $schemaForeignTable := .ForeignTable | $dot.SchemaTable}}
// Load{{$txt.Function.Name}} allows an eager lookup of values, cached into the
// loaded structs of the objects.
func ({{$varNameSingular}}L) Load{{$txt.Function.Name}}(e boil.Executor, singular bool, {{$arg}} interface{}) error {
	var slice []*{{$txt.LocalTable.NameGo}}
	var object *{{$txt.LocalTable.NameGo}}

	count := 1
	if singular {
		object = {{$arg}}.(*{{$txt.LocalTable.NameGo}})
	} else {
		slice = *{{$arg}}.(*[]*{{$txt.LocalTable.NameGo}})
		count = len(slice)
	}

	args := make([]interface{}, count)
	if singular {
		if object.R == nil {
			object.R = &{{$varNameSingular}}R{}
		}
		args[0] = object.{{.Column | titleCase}}
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &{{$varNameSingular}}R{}
			}
			args[i] = obj.{{.Column | titleCase}}
		}
	}

		{{if .ToJoinTable -}}
			{{- $schemaJoinTable := .JoinTable | $dot.SchemaTable -}}
	q := {{$txt.ForeignTable.NamePluralGo}}(e,
		qm.Select("{{id 0 | $dot.Quotes}}.*", "{{id 1 | $dot.Quotes}}.{{.JoinLocalColumn | $dot.Quotes}}"),
		qm.From("{{$schemaForeignTable}} as {{id 0 | $dot.Quotes}}"),
		qm.InnerJoin("{{$schemaJoinTable}} as {{id 1 | $dot.Quotes}} on {{id 0 | $dot.Quotes}}.{{.ForeignColumn | $dot.Quotes}} = {{id 1 | $dot.Quotes}}.{{.JoinForeignColumn | $dot.Quotes}}"),
		qm.WhereIn(
			"{{id 1 | $dot.Quotes}}.{{.JoinLocalColumn | $dot.Quotes}} in ?",
			args...,
		),
	)
		{{else -}}
	q := {{$txt.ForeignTable.NamePluralGo}}(e,
		qm.WhereIn(
			"{{.ForeignColumn | $dot.Quotes}} in ?",
			args...,
		),
	)
		{{end -}}

	{{if not $dot.NoHooks -}}
	if err := q.doSelectHooks(queries.GetExecutor(q.Query)); nil != err {
		return err
	}
	{{- end}}

	results, err := q.Query.Query()
	if err != nil {
		return errors.Wrap(err, "failed to eager load {{.ForeignTable}}")
	}
	defer results.Close()

	var resultSlice []*{{$txt.ForeignTable.NameGo}}
	{{if .ToJoinTable -}}
	{{- $foreignTable := getTable $dot.Tables .ForeignTable -}}
	{{- $joinTable := getTable $dot.Tables .JoinTable -}}
	{{- $localCol := $joinTable.GetColumn .JoinLocalColumn}}
	var localJoinCols []{{$localCol.Type}}
	for results.Next() {
		one := new({{$txt.ForeignTable.NameGo}})
		var localJoinCol {{$localCol.Type}}

		err = results.Scan({{$foreignTable.Columns | columnNames | stringMap $dot.StringFuncs.titleCase | prefixStringSlice "&one." | join ", "}}, &localJoinCol)
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignTable}}")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}

	if err = results.Err(); err != nil {
		return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignTable}}")
	}
	{{else -}}
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice {{.ForeignTable}}")
	}
	{{end}}

	{{if not $dot.NoHooks -}}
	if len({{.ForeignTable | singular | camelCase}}AfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks(e); err != nil {
				return err
			}
		}
	}

	{{- end}}
	if singular {
		object.R.{{$txt.Function.Name}} = resultSlice
		return nil
	}

	{{if .ToJoinTable -}}
	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			{{if $txt.Function.UsesBytes -}}
			if 0 == bytes.Compare(local.{{$txt.Function.LocalAssignment}}, localJoinCol) {
			{{else -}}
			if local.{{$txt.Function.LocalAssignment}} == localJoinCol {
			{{end -}}
				local.R.{{$txt.Function.Name}} = append(local.R.{{$txt.Function.Name}}, foreign)
				break
			}
		}
	}
	{{else -}}
	for _, foreign := range resultSlice {
		for _, local := range slice {
			{{if $txt.Function.UsesBytes -}}
			if 0 == bytes.Compare(local.{{$txt.Function.LocalAssignment}}, foreign.{{$txt.Function.ForeignAssignment}}) {
			{{else -}}
			if local.{{$txt.Function.LocalAssignment}} == foreign.{{$txt.Function.ForeignAssignment}} {
			{{end -}}
				local.R.{{$txt.Function.Name}} = append(local.R.{{$txt.Function.Name}}, foreign)
				break
			}
		}
	}
	{{end}}

	return nil
}

{{end -}}{{/* range tomany */}}
{{- end -}}{{/* if IsJoinTable */}}
