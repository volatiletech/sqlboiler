{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $txt := txtsFromToMany $.Tables $table . -}}
		{{- $varNameSingular := .Table | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase}}
		{{- $foreignPKeyCols := (getTable $.Tables .ForeignTable).PKey.Columns -}}
		{{- $foreignSchemaTable := .ForeignTable | $.SchemaTable}}
{{if $.AddGlobal -}}
// Add{{$txt.Function.Name}}G adds the given related objects to the existing relationships
// of the {{$table.Name | singular}}, optionally inserting them as new records.
// Appends related to o.R.{{$txt.Function.Name}}.
// Sets related.R.{{$txt.Function.ForeignName}} appropriately.
// Uses the global database handle.
func (o *{{$txt.LocalTable.NameGo}}) Add{{$txt.Function.Name}}G({{if not $.NoContext}}ctx context.Context, {{end -}} insert bool, related ...*{{$txt.ForeignTable.NameGo}}) error {
	return o.Add{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, insert, related...)
}

{{end -}}

{{if $.AddPanic -}}
// Add{{$txt.Function.Name}}P adds the given related objects to the existing relationships
// of the {{$table.Name | singular}}, optionally inserting them as new records.
// Appends related to o.R.{{$txt.Function.Name}}.
// Sets related.R.{{$txt.Function.ForeignName}} appropriately.
// Panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Add{{$txt.Function.Name}}P({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, insert bool, related ...*{{$txt.ForeignTable.NameGo}}) {
	if err := o.Add{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and $.AddGlobal $.AddPanic -}}
// Add{{$txt.Function.Name}}GP adds the given related objects to the existing relationships
// of the {{$table.Name | singular}}, optionally inserting them as new records.
// Appends related to o.R.{{$txt.Function.Name}}.
// Sets related.R.{{$txt.Function.ForeignName}} appropriately.
// Uses the global database handle and panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Add{{$txt.Function.Name}}GP({{if not $.NoContext}}ctx context.Context, {{end -}} insert bool, related ...*{{$txt.ForeignTable.NameGo}}) {
	if err := o.Add{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Add{{$txt.Function.Name}} adds the given related objects to the existing relationships
// of the {{$table.Name | singular}}, optionally inserting them as new records.
// Appends related to o.R.{{$txt.Function.Name}}.
// Sets related.R.{{$txt.Function.ForeignName}} appropriately.
func (o *{{$txt.LocalTable.NameGo}}) Add{{$txt.Function.Name}}({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, insert bool, related ...*{{$txt.ForeignTable.NameGo}}) error {
	var err error
	for _, rel := range related {
		if insert {
			{{if not .ToJoinTable -}}
			rel.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
				{{if .ForeignColumnNullable -}}
			rel.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
				{{end -}}
			{{end -}}

			if err = rel.Insert({{if not $.NoContext}}ctx, {{end -}} exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}{{if not .ToJoinTable}} else {
			updateQuery := fmt.Sprintf(
				"UPDATE {{$foreignSchemaTable}} SET %s WHERE %s",
				strmangle.SetParamNames("{{$.LQ}}", "{{$.RQ}}", {{if $.Dialect.UseIndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.ForeignColumn}}"{{"}"}}),
				strmangle.WhereClause("{{$.LQ}}", "{{$.RQ}}", {{if $.Dialect.UseIndexPlaceholders}}2{{else}}0{{end}}, {{$foreignVarNameSingular}}PrimaryKeyColumns),
			)
			values := []interface{}{o.{{$txt.LocalTable.ColumnNameGo}}, rel.{{$foreignPKeyCols | stringMap $.StringFuncs.titleCase | join ", rel."}}{{"}"}}

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

			rel.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
			{{if .ForeignColumnNullable -}}
			rel.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
			{{end -}}
		}{{end -}}
	}

	{{if .ToJoinTable -}}
	for _, rel := range related {
		query := "insert into {{.JoinTable | $.SchemaTable}} ({{.JoinLocalColumn | $.Quotes}}, {{.JoinForeignColumn | $.Quotes}}) values {{if $.Dialect.UseIndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}"
		values := []interface{}{{"{"}}o.{{$txt.LocalTable.ColumnNameGo}}, rel.{{$txt.ForeignTable.ColumnNameGo}}}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, query)
			fmt.Fprintln(boil.DebugWriter, values)
		}

		{{if $.NoContext -}}
		_, err = exec.Exec(query, values...)
		{{else -}}
		_, err = exec.ExecContext(ctx, query, values...)
		{{end -}}
		if err != nil {
			return errors.Wrap(err, "failed to insert into join table")
		}
	}
	{{end -}}

	if o.R == nil {
		o.R = &{{$varNameSingular}}R{
			{{$txt.Function.Name}}: related,
		}
	} else {
		o.R.{{$txt.Function.Name}} = append(o.R.{{$txt.Function.Name}}, related...)
	}

	{{if .ToJoinTable -}}
	for _, rel := range related {
		if rel.R == nil {
			rel.R = &{{$foreignVarNameSingular}}R{
				{{$txt.Function.ForeignName}}: {{$txt.LocalTable.NameGo}}Slice{{"{"}}o{{"}"}},
			}
		} else {
			rel.R.{{$txt.Function.ForeignName}} = append(rel.R.{{$txt.Function.ForeignName}}, o)
		}
	}
	{{else -}}
	for _, rel := range related {
		if rel.R == nil {
			rel.R = &{{$foreignVarNameSingular}}R{
				{{$txt.Function.ForeignName}}: o,
			}
		} else {
			rel.R.{{$txt.Function.ForeignName}} = o
		}
	}
	{{end -}}

	return nil
}

			{{- if (or .ForeignColumnNullable .ToJoinTable)}}
{{if $.AddGlobal -}}
// Set{{$txt.Function.Name}}G removes all previously related items of the
// {{$table.Name | singular}} replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
// Replaces o.R.{{$txt.Function.Name}} with related.
// Sets related.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
// Uses the global database handle.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}G({{if not $.NoContext}}ctx context.Context, {{end -}} insert bool, related ...*{{$txt.ForeignTable.NameGo}}) error {
	return o.Set{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, insert, related...)
}

{{end -}}

{{if $.AddPanic -}}
// Set{{$txt.Function.Name}}P removes all previously related items of the
// {{$table.Name | singular}} replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
// Replaces o.R.{{$txt.Function.Name}} with related.
// Sets related.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
// Panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}P({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, insert bool, related ...*{{$txt.ForeignTable.NameGo}}) {
	if err := o.Set{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} exec, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and $.AddGlobal $.AddPanic -}}
// Set{{$txt.Function.Name}}GP removes all previously related items of the
// {{$table.Name | singular}} replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
// Replaces o.R.{{$txt.Function.Name}} with related.
// Sets related.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
// Uses the global database handle and panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}GP({{if not $.NoContext}}ctx context.Context, {{end -}} insert bool, related ...*{{$txt.ForeignTable.NameGo}}) {
	if err := o.Set{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, insert, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Set{{$txt.Function.Name}} removes all previously related items of the
// {{$table.Name | singular}} replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
// Replaces o.R.{{$txt.Function.Name}} with related.
// Sets related.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, insert bool, related ...*{{$txt.ForeignTable.NameGo}}) error {
	{{if .ToJoinTable -}}
	query := "delete from {{.JoinTable | $.SchemaTable}} where {{.JoinLocalColumn | $.Quotes}} = {{if $.Dialect.UseIndexPlaceholders}}$1{{else}}?{{end}}"
	values := []interface{}{{"{"}}o.{{$txt.LocalTable.ColumnNameGo}}}
	{{else -}}
	query := "update {{.ForeignTable | $.SchemaTable}} set {{.ForeignColumn | $.Quotes}} = null where {{.ForeignColumn | $.Quotes}} = {{if $.Dialect.UseIndexPlaceholders}}$1{{else}}?{{end}}"
	values := []interface{}{{"{"}}o.{{$txt.LocalTable.ColumnNameGo}}}
	{{end -}}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	{{if $.NoContext -}}
	_, err := exec.Exec(query, values...)
	{{else -}}
	_, err := exec.ExecContext(ctx, query, values...)
	{{end -}}
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	{{if .ToJoinTable -}}
	remove{{$txt.Function.Name}}From{{$txt.Function.ForeignName}}Slice(o, related)
	if o.R != nil {
		o.R.{{$txt.Function.Name}} = nil
	}
	{{else -}}
	if o.R != nil {
		for _, rel := range o.R.{{$txt.Function.Name}} {
			rel.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.{{$txt.Function.ForeignName}} = nil
		}

		o.R.{{$txt.Function.Name}} = nil
	}
	{{end -}}

	return o.Add{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} exec, insert, related...)
}

{{if $.AddGlobal -}}
// Remove{{$txt.Function.Name}}G relationships from objects passed in.
// Removes related items from R.{{$txt.Function.Name}} (uses pointer comparison, removal does not keep order)
// Sets related.R.{{$txt.Function.ForeignName}}.
// Uses the global database handle.
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}G({{if not $.NoContext}}ctx context.Context, {{end -}} related ...*{{$txt.ForeignTable.NameGo}}) error {
	return o.Remove{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, related...)
}

{{end -}}

{{if $.AddPanic -}}
// Remove{{$txt.Function.Name}}P relationships from objects passed in.
// Removes related items from R.{{$txt.Function.Name}} (uses pointer comparison, removal does not keep order)
// Sets related.R.{{$txt.Function.ForeignName}}.
// Panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}P({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, related ...*{{$txt.ForeignTable.NameGo}}) {
	if err := o.Remove{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} exec, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

{{if and $.AddGlobal $.AddPanic -}}
// Remove{{$txt.Function.Name}}GP relationships from objects passed in.
// Removes related items from R.{{$txt.Function.Name}} (uses pointer comparison, removal does not keep order)
// Sets related.R.{{$txt.Function.ForeignName}}.
// Uses the global database handle and panics on error.
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}GP({{if not $.NoContext}}ctx context.Context, {{end -}} related ...*{{$txt.ForeignTable.NameGo}}) {
	if err := o.Remove{{$txt.Function.Name}}({{if $.NoContext}}boil.GetDB(){{else}}ctx, boil.GetContextDB(){{end}}, related...); err != nil {
		panic(boil.WrapErr(err))
	}
}

{{end -}}

// Remove{{$txt.Function.Name}} relationships from objects passed in.
// Removes related items from R.{{$txt.Function.Name}} (uses pointer comparison, removal does not keep order)
// Sets related.R.{{$txt.Function.ForeignName}}.
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}({{if $.NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}, related ...*{{$txt.ForeignTable.NameGo}}) error {
	var err error
	{{if .ToJoinTable -}}
	query := fmt.Sprintf(
		"delete from {{.JoinTable | $.SchemaTable}} where {{.JoinLocalColumn | $.Quotes}} = {{if $.Dialect.UseIndexPlaceholders}}$1{{else}}?{{end}} and {{.JoinForeignColumn | $.Quotes}} in (%s)",
		strmangle.Placeholders(dialect.UseIndexPlaceholders, len(related), 2, 1),
	)
	values := []interface{}{{"{"}}o.{{$txt.LocalTable.ColumnNameGo}}}
	for _, rel := range related {
		values = append(values, rel.{{$txt.ForeignTable.ColumnNameGo}})
	}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	{{if $.NoContext -}}
	_, err = exec.Exec(query, values...)
	{{else -}}
	_, err = exec.ExecContext(ctx, query, values...)
	{{end -}}
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}
	{{else -}}
	for _, rel := range related {
		rel.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
		{{if not .ToJoinTable -}}
		if rel.R != nil {
			rel.R.{{$txt.Function.ForeignName}} = nil
		}
		{{end -}}
		if {{if not $.NoRowsAffected}}_, {{end -}} err = rel.Update({{if not $.NoContext}}ctx, {{end -}} exec, "{{.ForeignColumn}}"); err != nil {
			return err
		}
	}
	{{end -}}

	{{if .ToJoinTable -}}
	remove{{$txt.Function.Name}}From{{$txt.Function.ForeignName}}Slice(o, related)
	{{end -}}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.{{$txt.Function.Name}} {
			if rel != ri {
				continue
			}

			ln := len(o.R.{{$txt.Function.Name}})
			if ln > 1 && i < ln-1 {
				o.R.{{$txt.Function.Name}}[i] = o.R.{{$txt.Function.Name}}[ln-1]
			}
			o.R.{{$txt.Function.Name}} = o.R.{{$txt.Function.Name}}[:ln-1]
			break
		}
	}

	return nil
}

				{{if .ToJoinTable -}}
func remove{{$txt.Function.Name}}From{{$txt.Function.ForeignName}}Slice(o *{{$txt.LocalTable.NameGo}}, related []*{{$txt.ForeignTable.NameGo}}) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.{{$txt.Function.ForeignName}} {
			{{if $txt.Function.UsesBytes -}}
			if 0 != bytes.Compare(o.{{$txt.Function.LocalAssignment}}, ri.{{$txt.Function.LocalAssignment}}) {
			{{else -}}
			if o.{{$txt.Function.LocalAssignment}} != ri.{{$txt.Function.LocalAssignment}} {
			{{end -}}
				continue
			}

			ln := len(rel.R.{{$txt.Function.ForeignName}})
			if ln > 1 && i < ln-1 {
				rel.R.{{$txt.Function.ForeignName}}[i] = rel.R.{{$txt.Function.ForeignName}}[ln-1]
			}
			rel.R.{{$txt.Function.ForeignName}} = rel.R.{{$txt.Function.ForeignName}}[:ln-1]
			break
		}
	}
}
				{{end -}}{{- /* if ToJoinTable */ -}}
			{{- end -}}{{- /* if nullable foreign key */ -}}
	{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if IsJoinTable */ -}}
