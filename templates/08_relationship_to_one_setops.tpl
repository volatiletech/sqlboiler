{{- /* Begin execution of template for one-to-one setops */ -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $txt := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase -}}
		{{- $localNameSingular := .ForeignKey.Table | singular | camelCase}}
// Set{{$txt.Function.Name}} of the {{.ForeignKey.Table | singular}} to the related item.
// Sets {{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} to related.
// Adds {{$txt.Function.Receiver}} to related.R.{{$txt.Function.ForeignName}}.
func ({{$txt.Function.Receiver}} *{{$txt.LocalTable.NameGo}}) Set{{$txt.Function.Name}}(exec boil.Executor, insert bool, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	{{if$txt.Function.OneToOne -}}
	oldVal := related.{{$txt.Function.ForeignAssignment}}
	related.{{$txt.Function.ForeignAssignment}} = {{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}}
		{{if .ForeignKey.ForeignColumnNullable -}}
	related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		{{end -}}
	if err = related.Update(exec, "{{.ForeignKey.ForeignColumn}}"); err != nil {
		related.{{$txt.Function.ForeignAssignment}} = oldVal
		return errors.Wrap(err, "failed to update local table")
	}
	{{else -}}
	oldVal := {{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}}
	related.{{$txt.Function.ForeignAssignment}} = {{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}}
	{{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}} = related.{{$txt.Function.ForeignAssignment}}
	if err = {{$txt.Function.Receiver}}.Update(exec, "{{.ForeignKey.Column}}"); err != nil {
		{{$txt.Function.Receiver}}.{{$txt.Function.LocalAssignment}} = oldVal
		return errors.Wrap(err, "failed to update local table")
	}
	{{end -}}

	if {{$txt.Function.Receiver}}.R == nil {
		{{$txt.Function.Receiver}}.R = &{{$localNameSingular}}R{
			{{$txt.Function.Name}}: related,
		}
	} else {
		{{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} = related
	}

	{{if (or .ForeignKey.Unique$txt.Function.OneToOne) -}}
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

	{{if .ForeignKey.Nullable}}
	{{$txt.Function.Receiver}}.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
	{{end -}}
	return nil
}

		{{- if or (.ForeignKey.Nullable) (and$txt.Function.OneToOne .ForeignKey.ForeignColumnNullable)}}
// Remove{{$txt.Function.Name}} relationship.
// Sets {{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} to nil.
// Removes {{$txt.Function.Receiver}} from all passed in related items' relationships struct (Optional).
func ({{$txt.Function.Receiver}} *{{$txt.LocalTable.NameGo}}) Remove{{$txt.Function.Name}}(exec boil.Executor, related *{{$txt.ForeignTable.NameGo}}) error {
	var err error

	{{if$txt.Function.OneToOne -}}
	related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = false
	if err = related.Update(exec, "{{.ForeignKey.ForeignColumn}}"); err != nil {
		related.{{$txt.ForeignTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}
	{{else -}}
	{{$txt.Function.Receiver}}.{{$txt.LocalTable.ColumnNameGo}}.Valid = false
	if err = {{$txt.Function.Receiver}}.Update(exec, "{{.ForeignKey.Column}}"); err != nil {
		{{$txt.Function.Receiver}}.{{$txt.LocalTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}
	{{end -}}

	{{$txt.Function.Receiver}}.R.{{$txt.Function.Name}} = nil
	if related == nil || related.R == nil {
		return nil
	}

	{{if .ForeignKey.Unique -}}
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
