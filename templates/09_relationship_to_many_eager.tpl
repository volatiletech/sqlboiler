{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $varNameSingular := .Table | singular | camelCase -}}
		{{- $foreignNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $txt := txtsFromToMany $.Tables $.Table . -}}
		{{- $arg := printf "maybe%s" $txt.LocalTable.NameGo -}}
		{{- $schemaForeignTable := .ForeignTable | $.SchemaTable}}
// Load{{$txt.Function.Name}} allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func ({{$varNameSingular}}L) Load{{$txt.Function.Name}}({{if $.NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, singular bool, {{$arg}} interface{}, mods queries.Applicator) error {
	var slice []*{{$txt.LocalTable.NameGo}}
	var object *{{$txt.LocalTable.NameGo}}

	if singular {
		object = {{$arg}}.(*{{$txt.LocalTable.NameGo}})
	} else {
		slice = *{{$arg}}.(*[]*{{$txt.LocalTable.NameGo}})
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &{{$varNameSingular}}R{}
		}
		args = append(args, object.{{.Column | titleCase}})
	} else {
		Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &{{$varNameSingular}}R{}
			}

			for _, a := range args {
				{{if $txt.Function.UsesPrimitives -}}
				if a == obj.{{$txt.LocalTable.ColumnNameGo}} {
				{{else -}}
				if queries.Equal(a, obj.{{$txt.LocalTable.ColumnNameGo}}) {
				{{end -}}
					continue Outer
				}
			}

			args = append(args, obj.{{.Column | titleCase}})
		}
	}

		{{if .ToJoinTable -}}
			{{- $schemaJoinTable := .JoinTable | $.SchemaTable -}}
	query := NewQuery(
		qm.Select("{{$schemaForeignTable}}.*, {{id 0 | $.Quotes}}.{{.JoinLocalColumn | $.Quotes}}"),
		qm.From("{{$schemaForeignTable}}"),
		qm.InnerJoin("{{$schemaJoinTable}} as {{id 0 | $.Quotes}} on {{$schemaForeignTable}}.{{.ForeignColumn | $.Quotes}} = {{id 0 | $.Quotes}}.{{.JoinForeignColumn | $.Quotes}}"),
		qm.WhereIn("{{id 0 | $.Quotes}}.{{.JoinLocalColumn | $.Quotes}} in ?", args...),
	)
		{{else -}}
	query := NewQuery(qm.From(`{{if $.Dialect.UseSchema}}{{$.Schema}}.{{end}}{{.ForeignTable}}`), qm.WhereIn(`{{.ForeignColumn}} in ?`, args...))
		{{end -}}
	if mods != nil {
		mods.Apply(query)
	}

	{{if $.NoContext -}}
	results, err := query.Query(query, e)
	{{else -}}
	results, err := query.QueryContext(ctx, e)
	{{end -}}
	if err != nil {
		return errors.Wrap(err, "failed to eager load {{.ForeignTable}}")
	}
	defer results.Close()

	var resultSlice []*{{$txt.ForeignTable.NameGo}}
	{{if .ToJoinTable -}}
	{{- $foreignTable := getTable $.Tables .ForeignTable -}}
	{{- $joinTable := getTable $.Tables .JoinTable -}}
	{{- $localCol := $joinTable.GetColumn .JoinLocalColumn}}
	var localJoinCols []{{$localCol.Type}}
	for results.Next() {
		one := new({{$txt.ForeignTable.NameGo}})
		var localJoinCol {{$localCol.Type}}

		err = results.Scan({{$foreignTable.Columns | columnNames | stringMap $.StringFuncs.titleCase | prefixStringSlice "&one." | join ", "}}, &localJoinCol)
		if err = results.Err(); err != nil {
			return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignTable}}")
		}

		resultSlice = append(resultSlice, one)
		localJoinCols = append(localJoinCols, localJoinCol)
	}
	{{- else -}}
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice {{.ForeignTable}}")
	}
	{{- end}}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on {{.ForeignTable}}")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for {{.ForeignTable}}")
	}

	{{if not $.NoHooks -}}
	if len({{.ForeignTable | singular | camelCase}}AfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks({{if $.NoContext}}e{{else}}ctx, e{{end -}}); err != nil {
				return err
			}
		}
	}

	{{- end}}
	if singular {
		object.R.{{$txt.Function.Name}} = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &{{$foreignNameSingular}}R{}
			}
			{{if .ToJoinTable -}}
			foreign.R.{{$txt.Function.ForeignName}} = append(foreign.R.{{$txt.Function.ForeignName}}, object)
			{{else -}}
			foreign.R.{{$txt.Function.ForeignName}} = object
			{{end -}}
		}
		return nil
	}

	{{if .ToJoinTable -}}
	for i, foreign := range resultSlice {
		localJoinCol := localJoinCols[i]
		for _, local := range slice {
			{{if $txt.Function.UsesPrimitives -}}
			if local.{{$txt.LocalTable.ColumnNameGo}} == localJoinCol {
			{{else -}}
			if queries.Equal(local.{{$txt.LocalTable.ColumnNameGo}}, localJoinCol) {
			{{end -}}
				local.R.{{$txt.Function.Name}} = append(local.R.{{$txt.Function.Name}}, foreign)
				if foreign.R == nil {
					foreign.R = &{{$foreignNameSingular}}R{}
				}
				foreign.R.{{$txt.Function.ForeignName}} = append(foreign.R.{{$txt.Function.ForeignName}}, local)
				break
			}
		}
	}
	{{else -}}
	for _, foreign := range resultSlice {
		for _, local := range slice {
			{{if $txt.Function.UsesPrimitives -}}
			if local.{{$txt.LocalTable.ColumnNameGo}} == foreign.{{$txt.ForeignTable.ColumnNameGo}} {
			{{else -}}
			if queries.Equal(local.{{$txt.LocalTable.ColumnNameGo}}, foreign.{{$txt.ForeignTable.ColumnNameGo}}) {
			{{end -}}
				local.R.{{$txt.Function.Name}} = append(local.R.{{$txt.Function.Name}}, foreign)
				if foreign.R == nil {
					foreign.R = &{{$foreignNameSingular}}R{}
				}
				foreign.R.{{$txt.Function.ForeignName}} = local
				break
			}
		}
	}
	{{end}}

	return nil
}

{{end -}}{{/* range tomany */}}
{{- end -}}{{/* if IsJoinTable */}}
