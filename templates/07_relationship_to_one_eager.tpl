{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range .Table.FKeys -}}
		{{- $varNameSingular := $.Table.Name | singular | camelCase -}}
		{{- $txt := txtsFromFKey $.Tables $.Table . -}}
		{{- $arg := printf "maybe%s" $txt.LocalTable.NameGo}}
// Load{{$txt.Function.Name}} allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
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
		args = append(args, object.{{$txt.LocalTable.ColumnNameGo}})
	} else {
		Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &{{$varNameSingular}}R{}
			}

			for _, a := range args {
				{{if $txt.Function.UsesBytes -}}
				if 0 == bytes.Compare(a.([]byte), obj.{{$txt.Function.LocalAssignment}}) {
				{{else -}}
				if a == obj.{{$txt.Function.LocalAssignment}} {
				{{end -}}
					continue Outer
				}
			}

			args = append(args, obj.{{$txt.LocalTable.ColumnNameGo}})
		}
	}

	query := NewQuery(qm.From(`{{if $.Dialect.UseSchema}}{{$.Schema}}.{{end}}{{.ForeignTable}}`), qm.WhereIn(`{{.ForeignColumn}} in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	{{if $.NoContext -}}
	results, err := query.Query(e)
	{{else -}}
	results, err := query.QueryContext(ctx, e)
	{{end -}}
	if err != nil {
		return errors.Wrap(err, "failed to eager load {{$txt.ForeignTable.NameGo}}")
	}

	var resultSlice []*{{$txt.ForeignTable.NameGo}}
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice {{$txt.ForeignTable.NameGo}}")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for {{.ForeignTable}}")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for {{.ForeignTable}}")
	}

	{{if not $.NoHooks -}}
	if len({{$varNameSingular}}AfterSelectHooks) != 0 {
		for _, obj := range resultSlice {
			if err := obj.doAfterSelectHooks({{if $.NoContext}}e{{else}}ctx, e{{end}}); err != nil {
				return err
			}
		}
	}
	{{- end}}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		object.R.{{$txt.Function.Name}} = resultSlice[0]
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			{{if $txt.Function.UsesBytes -}}
			if 0 == bytes.Compare(local.{{$txt.Function.LocalAssignment}}, foreign.{{$txt.Function.ForeignAssignment}}) {
			{{else -}}
			if local.{{$txt.Function.LocalAssignment}} == foreign.{{$txt.Function.ForeignAssignment}} {
			{{end -}}
				local.R.{{$txt.Function.Name}} = foreign
				break
			}
		}
	}

	return nil
}
{{end -}}{{/* range */}}
{{end}}{{/* join table */}}
