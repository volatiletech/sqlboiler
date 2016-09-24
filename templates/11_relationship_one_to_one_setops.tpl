{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $localNameSingular := .Table | singular | camelCase}}
// Set{{$txt.Function.Name}} of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.Name}} to related.
// Adds o to related.R.{{$txt.Function.ForeignName}}.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}(exec boil.Executor, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	oldVal := related.{{$txt.Function.ForeignAssignment}}
	related.{{$txt.Function.ForeignAssignment}} = o.{{$txt.Function.LocalAssignment}}
	{{if .ForeignColumnNullable -}}
	related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
	{{- end}}

	if insert {
		if err = related.Insert(exec); err != nil {
			related.{{$txt.Function.ForeignAssignment}} = oldVal
			{{if .ForeignColumnNullable -}}
			related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
			{{- end}}
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	} else {
		if err = related.Update(exec, "{{.ForeignColumn}}"); err != nil {
			related.{{$txt.Function.ForeignAssignment}} = oldVal
			{{if .ForeignColumnNullable -}}
			related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
			{{- end}}
			return errors.Wrap(err, "failed to update foreign table")
		}
	}

	if o.R == nil {
		o.R = &{{$localNameSingular}}R{
			{{$txt.Function.Name}}: related,
		}
	} else {
		o.R.{{$txt.Function.Name}} = related
	}

	if related.R == nil {
		related.R = &{{$varNameSingular}}R{
			{{$txt.Function.ForeignName}}: o,
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = o
	}
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
