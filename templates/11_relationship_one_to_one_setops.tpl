{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table . -}}
		{{- $foreignNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $varNameSingular := .Table | singular | camelCase}}
		{{- $schemaTable := .Table | $dot.SchemaTable -}}
// Set{{$txt.Function.Name}} of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.Name}} to related.
// Adds o to related.R.{{$txt.Function.ForeignName}}.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}(exec boil.Executor, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table {{$txt.ForeignTable}}")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE {{$schemaTable}} SET %s WHERE %s",
		strmangle.SetParamNames("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.Column}}"{{"}"}}),
		strmangle.WhereClause("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}2{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns),
	)
	values := []interface{}{related.{{$txt.ForeignTable.ColumnNameGo}}, o.{{$dot.Table.PKey.Columns | stringMap $dot.StringFuncs.titleCase | join ", o."}}{{"}"}}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table {{.Table}}")
	}

	o.{{$txt.Function.LocalAssignment}} = related.{{$txt.Function.ForeignAssignment}}
	{{if .Nullable -}}
	o.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
	{{- end}}

	if o.R == nil {
		o.R = &{{$varNameSingular}}R{
			{{$txt.Function.Name}}: related,
		}
	} else {
		o.R.{{$txt.Function.Name}} = related
	}

	{{if .Unique -}}
	if related.R == nil {
		related.R = &{{$foreignNameSingular}}R{
			{{$txt.Function.ForeignName}}: o,
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = o
	}
	{{else -}}
	if related.R == nil {
		related.R = &{{$foreignNameSingular}}R{
			{{$txt.Function.ForeignName}}: {{$txt.LocalTable.NameGo}}Slice{{"{"}}o{{"}"}},
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = append(related.R.{{$txt.Function.ForeignName}}, o)
	}
	{{- end}}

	return nil
}

		{{- if .ForeignColumnNullable}}
// Remove{{$txt.Function.Name}} relationship.
// Sets o.R.{{$txt.Function.Name}} to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}(exec boil.Executor, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
	if err = related.Update(exec, "{{.ForeignColumn}}"); err != nil {
		related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.{{$txt.Function.Name}} = nil
	if related == nil || related.R == nil {
		return nil
	}

	related.R.{{$txt.Function.ForeignName}} = nil
	return nil
}
{{end -}}{{/* if foreignkey nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
