{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $txt := txtsFromFKey $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $localNameSingular := .Table | singular | camelCase}}
// Set{{$txt.Function.Name}} of the {{.Table | singular}} to the related item.
// Sets o.R.{{$txt.Function.Name}} to related.
// Adds o to related.R.{{$txt.Function.ForeignName}}.
func (o *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}(exec boil.Executor, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	oldVal := o.{{$txt.Function.LocalAssignment}}
	o.{{$txt.Function.LocalAssignment}} = related.{{$txt.Function.ForeignAssignment}}
	{{if .Nullable -}}
	o.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
	{{- end}}
	if err = o.Update(exec, "{{.Column}}"); err != nil {
		o.{{$txt.Function.LocalAssignment}} = oldVal
		{{if .Nullable -}}
		o.{{$txt.LocalTable.ColumnNameGo}}.Valid = false
		{{- end}}
		return errors.Wrap(err, "failed to update local table")
	}

	if o.R == nil {
		o.R = &{{$localNameSingular}}R{
			{{$txt.Function.Name}}: related,
		}
	} else {
		o.R.{{$txt.Function.Name}} = related
	}

	{{if .Unique -}}
	if related.R == nil {
		related.R = &{{$varNameSingular}}R{
			{{$txt.Function.ForeignName}}: o,
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = o
	}
	{{else -}}
	if related.R == nil {
		related.R = &{{$varNameSingular}}R{
			{{$txt.Function.ForeignName}}: {{$txt.LocalTable.NameGo}}Slice{{"{"}}o{{"}"}},
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = append(related.R.{{$txt.Function.ForeignName}}, o)
	}
	{{- end}}

	return nil
}

		{{- if .Nullable}}
// Remove{{$txt.Function.Name}} relationship.
// Sets o.R.{{$txt.Function.Name}} to nil.
// Removes o from all passed in related items' relationships struct (Optional).
func (o *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}(exec boil.Executor, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	o.{{$txt.LocalTable.ColumnNameGo}}.Valid = false
	if err = o.Update(exec, "{{.Column}}"); err != nil {
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
