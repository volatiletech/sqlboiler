{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $txt := txtsFromToMany $dot.Tables $table . -}}
		{{- $varNameSingular := .Table | singular | camelCase -}}
		{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase}}
		{{- $foreignPKeyCols := (getTable $dot.Tables .ForeignTable).PKey.Columns -}}
		{{- $foreignSchemaTable := .ForeignTable | $dot.SchemaTable}}
// Add{{$txt.Function.Name}} adds the given related objects to the existing relationships
// of the {{$table.Name | singular}}, optionally inserting them as new records.
// Appends related to o.R.{{$txt.Function.Name}}.
// Sets related.R.{{$txt.Function.ForeignName}} appropriately.
func (o *{{$txt.LocalTable.NameGo}}) Add{{$txt.Function.Name}}(exec boil.Executor, insert bool, related ...*{{$txt.ForeignTable.NameGo}}) error {
	var err error
	for _, rel := range related {
		if insert {
			{{if not .ToJoinTable -}}
			rel.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
				{{if .ForeignColumnNullable -}}
			rel.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
				{{end -}}
			{{end -}}

			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}{{if not .ToJoinTable}} else {
			updateQuery := fmt.Sprintf(
				"UPDATE {{$foreignSchemaTable}} SET %s WHERE %s",
				strmangle.SetParamNames("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}1{{else}}0{{end}}, []string{{"{"}}"{{.ForeignColumn}}"{{"}"}}),
				strmangle.WhereClause("{{$dot.LQ}}", "{{$dot.RQ}}", {{if $dot.Dialect.IndexPlaceholders}}2{{else}}0{{end}}, {{$foreignVarNameSingular}}PrimaryKeyColumns),
			)
			values := []interface{}{o.{{$txt.LocalTable.ColumnNameGo}}, rel.{{$foreignPKeyCols | stringMap $dot.StringFuncs.titleCase | join ", rel."}}{{"}"}}

			if boil.DebugMode {
				fmt.Fprintln(boil.DebugWriter, updateQuery)
				fmt.Fprintln(boil.DebugWriter, values)
			}

			if _, err = exec.Exec(updateQuery, values...); err != nil {
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
		query := "insert into {{.JoinTable | $dot.SchemaTable}} ({{.JoinLocalColumn | $dot.Quotes}}, {{.JoinForeignColumn | $dot.Quotes}}) values {{if $dot.Dialect.IndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}"
		values := []interface{}{{"{"}}o.{{$txt.LocalTable.ColumnNameGo}}, rel.{{$txt.ForeignTable.ColumnNameGo}}}

		if boil.DebugMode {
			fmt.Fprintln(boil.DebugWriter, query)
			fmt.Fprintln(boil.DebugWriter, values)
		}

		_, err = exec.Exec(query, values...)
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
// Set{{$txt.Function.Name}} removes all previously related items of the
// {{$table.Name | singular}} replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
// Replaces o.R.{{$txt.Function.Name}} with related.
// Sets related.R.{{$txt.Function.ForeignName}}'s {{$txt.Function.Name}} accordingly.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}(exec boil.Executor, insert bool, related ...*{{$txt.ForeignTable.NameGo}}) error {
	{{if .ToJoinTable -}}
	query := "delete from {{.JoinTable | $dot.SchemaTable}} where {{.JoinLocalColumn | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}"
	values := []interface{}{{"{"}}o.{{$txt.LocalTable.ColumnNameGo}}}
	{{else -}}
	query := "update {{.ForeignTable | $dot.SchemaTable}} set {{.ForeignColumn | $dot.Quotes}} = null where {{.ForeignColumn | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}"
	values := []interface{}{{"{"}}o.{{$txt.LocalTable.ColumnNameGo}}}
	{{end -}}
	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err := exec.Exec(query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	{{if .ToJoinTable -}}
	remove{{$txt.Function.Name}}From{{$txt.Function.ForeignName}}Slice(o, related)
	o.R.{{$txt.Function.Name}} = nil
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

	return o.Add{{$txt.Function.Name}}(exec, insert, related...)
}

// Remove{{$txt.Function.Name}} relationships from objects passed in.
// Removes related items from R.{{$txt.Function.Name}} (uses pointer comparison, removal does not keep order)
// Sets related.R.{{$txt.Function.ForeignName}}.
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}(exec boil.Executor, related ...*{{$txt.ForeignTable.NameGo}}) error {
	var err error
	{{if .ToJoinTable -}}
	query := fmt.Sprintf(
		"delete from {{.JoinTable | $dot.SchemaTable}} where {{.JoinLocalColumn | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}} and {{.JoinForeignColumn | $dot.Quotes}} in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(related), 1, 1),
	)
	values := []interface{}{{"{"}}o.{{$txt.LocalTable.ColumnNameGo}}}

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	_, err = exec.Exec(query, values...)
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
		if err = rel.Update(exec, "{{.ForeignColumn}}"); err != nil {
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
