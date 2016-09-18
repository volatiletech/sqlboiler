{{- /* Begin execution of template for one-to-one setops */ -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $txt := txtsFromFKey $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .ForeignTable | singular | camelCase -}}
		{{- $localNameSingular := .Table | singular | camelCase}}
// Set{{$txt.Function.Name}} of the {{.Table | singular}} to the related item.
// Sets {{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} to related.
// Adds {{$txt.Function.Receiver}} to related.R.{{$txt.Function.ForeignName}}.
func ({{$txt.Function.Receiver}} *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}(exec boil.Executor, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	oldVal := {{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}}
	related.{{$txt.Function.ForeignAssignment}} = {{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}}
	{{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}} = related.{{$txt.Function.ForeignAssignment}}
	if err = {{$txt.Function.Receiver}}.Update(exec, "{{.Column}}"); err != nil {
		{{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}} = oldVal
		return errors.Wrap(err, "failed to update local table")
	}

	if {{$txt.Function.Receiver}}.R == nil {
		{{$txt.Function.Receiver}}.R = &{{$localNameSingular}}R{
			{{$txt.Function.Name}}: related,
		}
	} else {
		{{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} = related
	}

	{{if .Unique -}}
	if related.R == nil {
		related.R = &{{$varNameSingular}}R{
			{{$txt.Function.ForeignName}}: {{$txt.Function.Receiver}},
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = {{$txt.Function.Receiver}}
	}
	{{else -}}
	if related.R == nil {
		related.R = &{{$varNameSingular}}R{
			{{$txt.Function.ForeignName}}: {{$txt.LocalTable.NameGo}}Slice{{"{"}}{{$txt.Function.Receiver}}{{"}"}},
		}
	} else {
		related.R.{{$txt.Function.ForeignName}} = append(related.R.{{$txt.Function.ForeignName}}, {{$txt.Function.Receiver}})
	}
	{{end -}}

	{{if .Nullable}}
	{{$txt.Function.Receiver}}.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
	{{end -}}
	return nil
}

		{{- if .Nullable}}
// Remove{{$txt.Function.Name}} relationship.
// Sets {{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} to nil.
// Removes {{$txt.Function.Receiver}} from all passed in related items' relationships struct (Optional).
func ({{$txt.Function.Receiver}} *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}(exec boil.Executor, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	{{$txt.Function.Receiver}}.{{$txt.LocalTable.ColumnNameGo}}.Valid = false
	if err = {{$txt.Function.Receiver}}.Update(exec, "{{.Column}}"); err != nil {
		{{$txt.Function.Receiver}}.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}

	{{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} = nil
	if related == nil || related.R == nil {
		return nil
	}

	{{if .Unique -}}
	related.R.{{$txt.Function.ForeignName}} = nil
	{{else -}}
	for i, ri := range related.R.{{$txt.Function.ForeignName}} {
		{{if $txt.Function.UsesBytes -}}
		if 0 != bytes.Compare({{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}}, ri.{{$txt.Function.LocalAssignment}}) {
		{{else -}}
		if {{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}} != ri.{{$txt.Function.LocalAssignment}} {
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
