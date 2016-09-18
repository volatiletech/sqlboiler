{{- define "relationship_to_one_setops_helper" -}}
	{{- $dot := .Dot -}}{{/* .Dot holds the root templateData struct, passed in through preserveDot */}}
	{{- with .Rel -}}
		{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase -}}
		{{- $localNameSingular := .ForeignKey.Table | singular | camelCase}}
// Set{{.Function.Name}} of the {{.ForeignKey.Table | singular}} to the related item.
// Sets {{.Function.Receiver}}.R.{{.Function.Name}} to related.
// Adds {{.Function.Receiver}} to related.R.{{.Function.ForeignName}}.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) Set{{.Function.Name}}(exec boil.Executor, insert bool, related *{{.ForeignTable.NameGo}}) error {
	var err error
	if insert {
		if err = related.Insert(exec); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	{{if .Function.OneToOne -}}
	oldVal := related.{{.Function.ForeignAssignment}}
	related.{{.Function.ForeignAssignment}} = {{.Function.Receiver}}.{{.Function.LocalAssignment}}
		{{if .ForeignKey.ForeignColumnNullable -}}
	related.{{.ForeignTable.ColumnNameGo}}.Valid = true
		{{end -}}
	if err = related.Update(exec, "{{.ForeignKey.ForeignColumn}}"); err != nil {
		related.{{.Function.ForeignAssignment}} = oldVal
		return errors.Wrap(err, "failed to update local table")
	}
	{{else -}}
	oldVal := {{.Function.Receiver}}.{{.Function.LocalAssignment}}
	related.{{.Function.ForeignAssignment}} = {{.Function.Receiver}}.{{.Function.LocalAssignment}}
	{{.Function.Receiver}}.{{.Function.LocalAssignment}} = related.{{.Function.ForeignAssignment}}
	if err = {{.Function.Receiver}}.Update(exec, "{{.ForeignKey.Column}}"); err != nil {
		{{.Function.Receiver}}.{{.Function.LocalAssignment}} = oldVal
		return errors.Wrap(err, "failed to update local table")
	}
	{{end -}}

	if {{.Function.Receiver}}.R == nil {
		{{.Function.Receiver}}.R = &{{$localNameSingular}}R{
			{{.Function.Name}}: related,
		}
	} else {
		{{.Function.Receiver}}.R.{{.Function.Name}} = related
	}

	{{if (or .ForeignKey.Unique .Function.OneToOne) -}}
	if related.R == nil {
		related.R = &{{$varNameSingular}}R{
			{{.Function.ForeignName}}: {{.Function.Receiver}},
		}
	} else {
		related.R.{{.Function.ForeignName}} = {{.Function.Receiver}}
	}
	{{else -}}
	if related.R == nil {
		related.R = &{{$varNameSingular}}R{
			{{.Function.ForeignName}}: {{.LocalTable.NameGo}}Slice{{"{"}}{{.Function.Receiver}}{{"}"}},
		}
	} else {
		related.R.{{.Function.ForeignName}} = append(related.R.{{.Function.ForeignName}}, {{.Function.Receiver}})
	}
	{{end -}}

	{{if .ForeignKey.Nullable}}
	{{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}.Valid = true
	{{end -}}
	return nil
}

		{{- if or (.ForeignKey.Nullable) (and .Function.OneToOne .ForeignKey.ForeignColumnNullable)}}
// Remove{{.Function.Name}} relationship.
// Sets {{.Function.Receiver}}.R.{{.Function.Name}} to nil.
// Removes {{.Function.Receiver}} from all passed in related items' relationships struct (Optional).
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) Remove{{.Function.Name}}(exec boil.Executor, related *{{.ForeignTable.NameGo}}) error {
	var err error

	{{if .Function.OneToOne -}}
	related.{{.ForeignTable.ColumnNameGo}}.Valid = false
	if err = related.Update(exec, "{{.ForeignKey.ForeignColumn}}"); err != nil {
		related.{{.ForeignTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}
	{{else -}}
	{{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}.Valid = false
	if err = {{.Function.Receiver}}.Update(exec, "{{.ForeignKey.Column}}"); err != nil {
		{{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}.Valid = true
		return errors.Wrap(err, "failed to update local table")
	}
	{{end -}}

	{{.Function.Receiver}}.R.{{.Function.Name}} = nil
	if related == nil || related.R == nil {
		return nil
	}

	{{if .ForeignKey.Unique -}}
	related.R.{{.Function.ForeignName}} = nil
	{{else -}}
	for i, ri := range related.R.{{.Function.ForeignName}} {
		{{if .Function.UsesBytes -}}
		if 0 != bytes.Compare({{.Function.Receiver}}.{{.Function.LocalAssignment}}, ri.{{.Function.LocalAssignment}}) {
		{{else -}}
		if {{.Function.Receiver}}.{{.Function.LocalAssignment}} != ri.{{.Function.LocalAssignment}} {
		{{end -}}
			continue
		}

		ln := len(related.R.{{.Function.ForeignName}})
		if ln > 1 && i < ln-1 {
			related.R.{{.Function.ForeignName}}[i] = related.R.{{.Function.ForeignName}}[ln-1]
		}
		related.R.{{.Function.ForeignName}} = related.R.{{.Function.ForeignName}}[:ln-1]
		break
	}
	{{end -}}

	return nil
}
{{end -}}{{/* if foreignkey nullable */}}
{{- end -}}{{/* end with */}}
{{- end -}}{{/* end define */}}

{{- /* Begin execution of template for one-to-one setops */ -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $txt := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
		{{- template "relationship_to_one_setops_helper" (preserveDot $dot $txt) -}}
	{{- end -}}
{{- end -}}
