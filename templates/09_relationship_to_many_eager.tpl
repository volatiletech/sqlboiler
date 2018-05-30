{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $varNameSingular := $.Table.Name | singular | camelCase -}}
		{{- $txt := txtsFromToMany $.Tables $.Table . -}}
		{{- $arg := printf "maybe%s" $txt.LocalTable.NameGo -}}
		{{- $schemaForeignTable := .ForeignTable | $.SchemaTable}}
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
			{{- $schemaJoinTable := .JoinTable | $.SchemaTable -}}
	query := fmt.Sprintf(
		"select {{id 0 | $.Quotes}}.*, {{id 1 | $.Quotes}}.{{.JoinLocalColumn | $.Quotes}} from {{$schemaForeignTable}} as {{id 0 | $.Quotes}} inner join {{$schemaJoinTable}} as {{id 1 | $.Quotes}} on {{id 0 | $.Quotes}}.{{.ForeignColumn | $.Quotes}} = {{id 1 | $.Quotes}}.{{.JoinForeignColumn | $.Quotes}} where {{id 1 | $.Quotes}}.{{.JoinLocalColumn | $.Quotes}} in (%s)",
		strmangle.Placeholders(dialect.UseIndexPlaceholders, count, 1, 1),
	)
		{{else -}}
	query := fmt.Sprintf(
		"select * from {{$schemaForeignTable}} where {{.ForeignColumn | $.Quotes}} in (%s)",
		strmangle.Placeholders(dialect.UseIndexPlaceholders, count, 1, 1),
	)
		{{end -}}

	if boil.DebugMode {
		fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
	}

	{{if $.NoContext -}}
	results, err := e.Query(query, args...)
	{{else -}}
	results, err := e.QueryContext(ctx, query, args...)
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

	if err = results.Err(); err != nil {
		return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignTable}}")
	}
	{{else -}}
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice {{.ForeignTable}}")
	}
	{{end}}

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
