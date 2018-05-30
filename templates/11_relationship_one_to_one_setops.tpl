{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $.Tables $.Table . -}}
		{{- $varNameSingular := .Table | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $foreignPKeyCols := (getTable $.Tables .ForeignTable).PKey.Columns -}}
		{{- $foreignSchemaTable := .ForeignTable | $.SchemaTable}}
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
		related.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
		{{if .ForeignColumnNullable -}}
		related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		{{- end}}

		if err = related.Insert({{if not $.NoContext}}ctx, {{end -}} exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	} else {
		updateQuery := fmt.Sprintf(
			"UPDATE {{$foreignSchemaTable}} SET %s WHERE %s",
			strmangle.SetParamNames("{{$.LQ}}", "{{$.RQ}}", {{if $.Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.ForeignColumn}}"{{"}"}}),
			strmangle.WhereClause("{{$.LQ}}", "{{$.RQ}}", {{if $.Dialect.UseIndexPlaceholders}}2{{else}}0{{end}}, {{$foreignVarNameSingular}}PrimaryKeyColumns),
		)
		values := []interface{}{o.{{$txt.LocalTable.ColumnNameGo}}, related.{{$foreignPKeyCols | stringMap $.StringFuncs.titleCase | join ", related."}}{{"}"}}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, updateQuery)
			fmt.Fprintln(boil.DebugWriter, values)
		}

		{{if $.NoContext -}}
		if _, err = exec.Exec(updateQuery, values...); err != nil {
		{{else -}}
		if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		{{end -}}
			return errors.Wrap(err, "failed to update foreign table")
		}

		related.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
		{{if .ForeignColumnNullable -}}
		related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		{{- end}}
	}


	if o.R == nil {
		o.R = &{{$varNameSingular}}R{
			{{$txt.Function.Name}}: related,
		}
	} else {
		o.R.{{$txt.Function.Name}} = related
	}

	if related.R == nil {
		related.R = &{{$foreignVarNameSingular}}R{
			{{$txt.Function.ForeignName}}: o,
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = o
	}
	return nil
}

		{{- if .ForeignColumnNullable}}
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

	related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
	if {{if not $.NoRowsAffected}}_, {{end -}} err = related.Update({{if not $.NoContext}}ctx, {{end -}} exec, "{{.ForeignColumn}}"); err != nil {
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
