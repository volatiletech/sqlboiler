{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.ToOneRelationships -}}
		{{- $txt := txtsFromOneToOne $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $localNameSingular := .Table | singular | camelCase}}
// Set{{$txt.Function.Name}} of the {{.Table | singular}} to the related item.
// Sets {{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} to related.
// Adds {{$txt.Function.Receiver}} to related.R.{{$txt.Function.ForeignName}}.
func ({{$txt.Function.Receiver}} *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}(exec boil.Executor, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	oldVal := related.{{$txt.Function.ForeignAssignment}}
	related.{{$txt.Function.ForeignAssignment}} = {{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}}
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

	if {{$txt.Function.Receiver}}.R == nil {
		{{$txt.Function.Receiver}}.R = &{{$localNameSingular}}R{
			{{$txt.Function.Name}}: related,
		}
	} else {
		{{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} = related
	}

	if related.R == nil {
		related.R = &{{$varNameSingular}}R{
			{{$txt.Function.ForeignName}}: {{$txt.Function.Receiver}},
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = {{$txt.Function.Receiver}}
	}
	return nil
}

		{{- if .ForeignColumnNullable}}
// Remove{{$txt.Function.Name}} relationship.
// Sets {{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} to nil.
// Removes {{$txt.Function.Receiver}} from all passed in related items' relationships struct (Optional).
func ({{$txt.Function.Receiver}} *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}(exec boil.Executor, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
	if err = related.Update(exec, "{{.ForeignColumn}}"); err != nil {
		related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	{{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} = nil
	if related == nil || related.R == nil {
		return nil
	}

	related.R.{{$txt.Function.ForeignName}} = nil
	return nil
}
{{end -}}{{/* if foreignkey nullable */}}
{{- end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
