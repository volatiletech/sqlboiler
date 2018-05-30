{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $.Tables $.Table . -}}
		{{- $varNameSingular := $.Table.Name | singular | camelCase -}}
		{{- $arg := printf "maybe%s" $txt.LocalTable.NameGo -}}
// Load{{$txt.Function.Name}} allows an eager lookup of values, cached into the
// loaded structs of the objects.
func ({{$varNameSingular}}L) Load{{$txt.Function.Name}}({{if $.NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, singular bool, {{$arg}} interface{}) error {
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
		args[0] = object.{{$txt.LocalTable.ColumnNameGo}}
	} else {
		for i, obj := range slice {
			if obj.R == nil {
				obj.R = &{{$varNameSingular}}R{}
			}
			args[i] = obj.{{$txt.LocalTable.ColumnNameGo}}
		}
	}

	query := fmt.Sprintf(
		"select * from {{.ForeignTable | $.SchemaTable}} where {{.ForeignColumn | $.Quotes}} in (%s)",
		strmangle.Placeholders(dialect.UseIndexPlaceholders, count, 1, 1),
	)

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	{{if $.NoContext -}}
	results, err := e.Query(query, args...)
	{{else -}}
	results, err := e.QueryContext(ctx, query, args...)
	{{end -}}
	if err != nil {
		return errors.Wrap(err, "failed to eager load {{$txt.ForeignTable.NameGo}}")
	}
	defer results.Close()

	var resultSlice []*{{$txt.ForeignTable.NameGo}}
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice {{$txt.ForeignTable.NameGo}}")
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
