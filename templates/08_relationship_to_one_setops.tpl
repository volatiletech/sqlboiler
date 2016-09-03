{{- define "relationship_to_one_setops_helper" -}}
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

  oldVal := {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}
  {{.Function.Receiver}}.{{.Function.LocalAssignment}} = related.{{.Function.ForeignAssignment}}
  if err = {{.Function.Receiver}}.Update(exec, "{{.ForeignKey.Column}}"); err != nil {
    {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}} = oldVal
    return errors.Wrap(err, "failed to update local table")
  }

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
{{- if .ForeignKey.Nullable}}

// Remove{{.Function.Name}} relationship.
// Sets {{.Function.Receiver}}.R.{{.Function.Name}} to nil.
// Removes {{.Function.Receiver}} from all passed in related items' relationships struct (Optional).
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) Remove{{.Function.Name}}(exec boil.Executor, related *{{.ForeignTable.NameGo}}) error {
  var err error

  {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}.Valid = false
  if err = {{.Function.Receiver}}.Update(exec, "{{.ForeignKey.Column}}"); err != nil {
    {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}.Valid = true
    return errors.Wrap(err, "failed to update local table")
  }

  {{.Function.Receiver}}.R.{{.Function.Name}} = nil
  if related == nil || related.R == nil {
    return nil
  }

  {{if .ForeignKey.Unique -}}
  related.R.{{.Function.ForeignName}} = nil
  {{else -}}
  for i, ri := range related.R.{{.Function.ForeignName}} {
    if {{.Function.Receiver}}.{{.Function.LocalAssignment}} != ri.{{.Function.LocalAssignment}} {
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
{{end -}}
{{- end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
{{- template "relationship_to_one_setops_helper" $rel -}}
{{- end -}}
{{- end -}}
