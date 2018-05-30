{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range .Table.FKeys -}}
		{{- $txt := txtsFromFKey $.Tables $.Table . -}}
		{{- $foreignNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $varNameSingular := .Table | singular | camelCase}}
		{{- $schemaTable := .Table | $.SchemaTable}}
{{if $.AddGlobal -}}
// Set{{$txt.Function.Name}}G of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.Name}} to related.
// Adds o to related.R.{{$txt.Function.ForeignName}}.
// Uses the global database handle.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}G({{if not $.NoContext}}ctx context.Context, {{end -}} insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	return o.Set{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, insert, related)
}

{{end -}}

{{if $.AddPanic -}}
// Set{{$txt.Function.Name}}P of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.Name}} to related.
// Adds o to related.R.{{$txt.Function.ForeignName}}.
// Panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}P({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, insert bool, related *{{$txt.ForeignTable.NameGo}}) {
	if err := o.Set{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} exec, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and $.AddGlobal $.AddPanic -}}
// Set{{$txt.Function.Name}}GP of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.Name}} to related.
// Adds o to related.R.{{$txt.Function.ForeignName}}.
// Uses the global database handle and panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}GP({{if not $.NoContext}}ctx context.Context, {{end -}} insert bool, related *{{$txt.ForeignTable.NameGo}}) {
	if err := o.Set{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, insert, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Set{{$txt.Function.Name}} of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.Name}} to related.
// Adds o to related.R.{{$txt.Function.ForeignName}}.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error
	if insert {
		if err = related.Insert({{if not $.NoContext}}ctx, {{end -}} exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE {{$schemaTable}} SET %s WHERE %s",
		strmangle.SetParamNames("{{$.LQ}}", "{{$.RQ}}", {{if $.Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.Column}}"{{"}"}}),
		strmangle.WhereClause("{{$.LQ}}", "{{$.RQ}}", {{if $.Dialect.UseIndexPlaceholders}}2{{else}}0{{end}}, {{$varNameSingular}}PrimaryKeyColumns),
	)
	values := []interface{}{related.{{$txt.ForeignTable.ColumnNameGo}}, o.{{$.Table.PKey.Columns | stringMap $.StringFuncs.titleCase | join ", o."}}{{"}"}}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, updateQuery)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	{{if $.NoContext -}}
	if _, err = exec.Exec(updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}
	{{- else -}}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}
	{{- end}}

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

		{{- if .Nullable}}
{{if $.AddGlobal -}}
// Remove{{$txt.Function.Name}}G relationship.
// Sets o.R.{{$txt.Function.Name}} to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle.
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}G({{if not $.NoContext}}ctx context.Context, {{end -}} related *{{$txt.ForeignTable.NameGo}}) error {
	return o.Remove{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, related)
}

{{end -}}

{{if $.AddPanic -}}
// Remove{{$txt.Function.Name}}P relationship.
// Sets o.R.{{$txt.Function.Name}} to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}P({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, related *{{$txt.ForeignTable.NameGo}}) {
	if err := o.Remove{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} exec, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and $.AddGlobal $.AddPanic -}}
// Remove{{$txt.Function.Name}}GP relationship.
// Sets o.R.{{$txt.Function.Name}} to nil.
// Removes o from all passed in related items' relationships struct (Optional).
// Uses the global database handle and panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}GP({{if not $.NoContext}}ctx context.Context, {{end -}} related *{{$txt.ForeignTable.NameGo}}) {
	if err := o.Remove{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, related); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Remove{{$txt.Function.Name}} relationship.
// Sets o.R.{{$txt.Function.Name}} to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	o.{{$txt.LocalTable.ColumnNameGo}}.Valid = false
	{{if $.NoContext -}}
	if {{if not $.NoRowsAffected}}_, {{end -}} err = o.Update(exec, "{{.Column}}"); err != nil {
	{{else -}}
	if {{if not $.NoRowsAffected}}_, {{end -}} err = o.Update(ctx, exec, "{{.Column}}"); err != nil {
	{{end -}}
		o.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	o.R.{{$txt.Function.Name}} = nil
	if related == nil || related.R == nil {
		return nil
	}

	{{if .Unique -}}
	related.R.{{$txt.Function.ForeignName}} = nil
	{{else -}}
	for i, ri := range related.R.{{$txt.Function.ForeignName}} {
		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare(o.{{$txt.Function.LocalAssignment}}, ri.{{$txt.Function.LocalAssignment}}) {
		{{else -}}
		if o.{{$txt.Function.LocalAssignment}} != ri.{{$txt.Function.LocalAssignment}} {
		{{end -}}
			continue
		}

		ln := len(related.R.{{$txt.Function.ForeignName}})
		if ln > 1 && i < ln-1 {
			related.R.{{$txt.Function.ForeignName}}[i] = related.R.{{$txt.Function.ForeignName}}[ln-1]
		}
		related.R.{{$txt.Function.ForeignName}} = related.R.{{$txt.Function.ForeignName}}[:ln-1]
		break
	}
	{{end -}}

	return nil
}
{{end -}}{{/* if foreignkey nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
