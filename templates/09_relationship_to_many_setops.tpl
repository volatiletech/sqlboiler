{{- /* Begin execution of template for many-to-one or many-to-many setops */ -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- $table := .Table -}}
	{{- range .Table.ToManyRelationships -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase -}}
		{{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
			{{- /* Begin execution of template for many-to-one setops */ -}}
			{{- $txt := textsFromOneToOneRelationship $dot.PkgName $dot.Tables $table . -}}
			{{- template "relationship_to_one_setops_helper" (preserveDot $dot $txt) -}}
		{{- else -}}
			{{- $rel := textsFromRelationship $dot.Tables $table . -}}
			{{- $localNameSingular := .Table | singular | camelCase -}}
			{{- $foreignNameSingular := .ForeignTable | singular | camelCase}}
// Add{{$rel.Function.Name}} adds the given related objects to the existing relationships
// of the {{$table.Name | singular}}, optionally inserting them as new records.
// Appends related to {{$rel.Function.Receiver}}.R.{{$rel.Function.Name}}.
// Sets related.R.{{$rel.Function.ForeignName}} appropriately.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) Add{{$rel.Function.Name}}(exec boil.Executor, insert bool, related ...*{{$rel.ForeignTable.NameGo}}) error {
	var err error
	for _, rel := range related {
		{{if not .ToJoinTable -}}
		rel.{{$rel.Function.ForeignAssignment}} = {{$rel.Function.Receiver}}.{{$rel.Function.LocalAssignment}}
			{{if .ForeignColumnNullable -}}
		rel.{{$rel.ForeignTable.ColumnNameGo}}.Valid = true
			{{end -}}
		{{end -}}
		if insert {
			if err = rel.Insert(exec); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		}{{if not .ToJoinTable}} else {
			if err = rel.Update(exec, "{{.ForeignColumn}}"); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}
		}{{end -}}
	}

	{{if .ToJoinTable -}}
	for _, rel := range related {
		query := "insert into {{.JoinTable | $dot.SchemaTable}} ({{.JoinLocalColumn | $dot.Quotes}}, {{.JoinForeignColumn | $dot.Quotes}}) values {{if $dot.Dialect.IndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}"
		values := []interface{}{{"{"}}{{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}, rel.{{$rel.ForeignTable.ColumnNameGo}}}

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

	if {{$rel.Function.Receiver}}.R == nil {
		{{$rel.Function.Receiver}}.R = &{{$localNameSingular}}R{
			{{$rel.Function.Name}}: related,
		}
	} else {
		{{$rel.Function.Receiver}}.R.{{$rel.Function.Name}} = append({{$rel.Function.Receiver}}.R.{{$rel.Function.Name}}, related...)
	}

	{{if .ToJoinTable -}}
	for _, rel := range related {
		if rel.R == nil {
			rel.R = &{{$foreignNameSingular}}R{
				{{$rel.Function.ForeignName}}: {{$rel.LocalTable.NameGo}}Slice{{"{"}}{{$rel.Function.Receiver}}{{"}"}},
			}
		} else {
			rel.R.{{$rel.Function.ForeignName}} = append(rel.R.{{$rel.Function.ForeignName}}, {{$rel.Function.Receiver}})
		}
	}
	{{else -}}
	for _, rel := range related {
		if rel.R == nil {
			rel.R = &{{$foreignNameSingular}}R{
				{{$rel.Function.ForeignName}}: {{$rel.Function.Receiver}},
			}
		} else {
			rel.R.{{$rel.Function.ForeignName}} = {{$rel.Function.Receiver}}
		}
	}
	{{end -}}

	return nil
}

			{{- if (or .ForeignColumnNullable .ToJoinTable)}}
// Set{{$rel.Function.Name}} removes all previously related items of the
// {{$table.Name | singular}} replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets {{$rel.Function.Receiver}}.R.{{$rel.Function.ForeignName}}'s {{$rel.Function.Name}} accordingly.
// Replaces {{$rel.Function.Receiver}}.R.{{$rel.Function.Name}} with related.
// Sets related.R.{{$rel.Function.ForeignName}}'s {{$rel.Function.Name}} accordingly.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) Set{{$rel.Function.Name}}(exec boil.Executor, insert bool, related ...*{{$rel.ForeignTable.NameGo}}) error {
	{{if .ToJoinTable -}}
	query := "delete from {{.JoinTable | $dot.SchemaTable}} where {{.JoinLocalColumn | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}"
	values := []interface{}{{"{"}}{{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}}
	{{else -}}
	query := "update {{.ForeignTable | $dot.SchemaTable}} set {{.ForeignColumn | $dot.Quotes}} = null where {{.ForeignColumn | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}"
	values := []interface{}{{"{"}}{{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}}
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
	remove{{$rel.LocalTable.NameGo}}From{{$rel.ForeignTable.NameGo}}Slice({{$rel.Function.Receiver}}, related)
	{{$rel.Function.Receiver}}.R.{{$rel.Function.Name}} = nil
	{{else -}}
	if {{$rel.Function.Receiver}}.R != nil {
		for _, rel := range {{$rel.Function.Receiver}}.R.{{$rel.Function.Name}} {
			rel.{{$rel.ForeignTable.ColumnNameGo}}.Valid = false
			if rel.R == nil {
				continue
			}

			rel.R.{{$rel.Function.ForeignName}} = nil
		}

		{{$rel.Function.Receiver}}.R.{{$rel.Function.Name}} = nil
	}
	{{end -}}

	return {{$rel.Function.Receiver}}.Add{{$rel.Function.Name}}(exec, insert, related...)
}

// Remove{{$rel.Function.Name}} relationships from objects passed in.
// Removes related items from R.{{$rel.Function.Name}} (uses pointer comparison, removal does not keep order)
// Sets related.R.{{$rel.Function.ForeignName}}.
func ({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}) Remove{{$rel.Function.Name}}(exec boil.Executor, related ...*{{$rel.ForeignTable.NameGo}}) error {
	var err error
	{{if .ToJoinTable -}}
	query := fmt.Sprintf(
		"delete from {{.JoinTable | $dot.SchemaTable}} where {{.JoinLocalColumn | $dot.Quotes}} = {{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}} and {{.JoinForeignColumn | $dot.Quotes}} in (%s)",
		strmangle.Placeholders(dialect.IndexPlaceholders, len(related), 1, 1),
	)
	values := []interface{}{{"{"}}{{$rel.Function.Receiver}}.{{$rel.LocalTable.ColumnNameGo}}}

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
		rel.{{$rel.ForeignTable.ColumnNameGo}}.Valid = false
		{{if not .ToJoinTable -}}
		if rel.R != nil {
			rel.R.{{$rel.Function.ForeignName}} = nil
		}
		{{end -}}
		if err = rel.Update(exec, "{{.ForeignColumn}}"); err != nil {
			return err
		}
	}
	{{end -}}

	{{if .ToJoinTable -}}
	remove{{$rel.LocalTable.NameGo}}From{{$rel.ForeignTable.NameGo}}Slice({{$rel.Function.Receiver}}, related)
	{{end -}}
	if {{$rel.Function.Receiver}}.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range {{$rel.Function.Receiver}}.R.{{$rel.Function.Name}} {
			if rel != ri {
				continue
			}

			ln := len({{$rel.Function.Receiver}}.R.{{$rel.Function.Name}})
			if ln > 1 && i < ln-1 {
				{{$rel.Function.Receiver}}.R.{{$rel.Function.Name}}[i] = {{$rel.Function.Receiver}}.R.{{$rel.Function.Name}}[ln-1]
			}
			{{$rel.Function.Receiver}}.R.{{$rel.Function.Name}} = {{$rel.Function.Receiver}}.R.{{$rel.Function.Name}}[:ln-1]
			break
		}
	}

	return nil
}

				{{if .ToJoinTable -}}
func remove{{$rel.LocalTable.NameGo}}From{{$rel.ForeignTable.NameGo}}Slice({{$rel.Function.Receiver}} *{{$rel.LocalTable.NameGo}}, related []*{{$rel.ForeignTable.NameGo}}) {
	for _, rel := range related {
		if rel.R == nil {
			continue
		}
		for i, ri := range rel.R.{{$rel.Function.ForeignName}} {
			if {{$rel.Function.Receiver}}.{{$rel.Function.LocalAssignment}} != ri.{{$rel.Function.LocalAssignment}} {
				continue
			}

			ln := len(rel.R.{{$rel.Function.ForeignName}})
			if ln > 1 && i < ln-1 {
				rel.R.{{$rel.Function.ForeignName}}[i] = rel.R.{{$rel.Function.ForeignName}}[ln-1]
			}
			rel.R.{{$rel.Function.ForeignName}} = rel.R.{{$rel.Function.ForeignName}}[:ln-1]
			break
		}
	}
}
				{{end -}}{{- /* if ToJoinTable */ -}}
			{{- end -}}{{- /* if nullable foreign key */ -}}
		{{- end -}}{{- /* if unique foreign key */ -}}
	{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if IsJoinTable */ -}}
