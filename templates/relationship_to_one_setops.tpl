{{- define "relationship_to_one_setops_helper" -}}
{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase}}

// Set{{.Function.Name}} of the {{.ForeignKey.Table | singular}} to the related item.
// Sets R.{{.Function.Name}} to related.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) Set{{.Function.Name}}(exec boil.Executor, insert bool, related *{{.ForeignTable.NameGo}}) error {
  var err error
  if insert {
    if err = related.Insert(exec); err != nil {
      return err
    }
  }

  oldVal := {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}
  {{.Function.Receiver}}.{{.Function.LocalAssignment}} = related.{{.Function.ForeignAssignment}}
  if err = {{.Function.Receiver}}.Update(exec, "{{.ForeignKey.Column}}"); err != nil {
    {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}} = oldVal
    return err
  }

  if {{.Function.Receiver}}.R == nil {
    {{.Function.Receiver}}.R = &{{.LocalTable.NameGo}}R{
      {{.Function.Name}}: related,
    }
  } else {
    {{.Function.Receiver}}.R.{{.Function.Name}} = related
  }

  {{if .ForeignKey.Nullable}}
  {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}.Valid = true
  {{end -}}
  return nil
}
{{- if .ForeignKey.Nullable}}

// Remove{{.Function.Name}} relationship.
// Sets R.{{.Function.Name}} to nil.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) Remove{{.Function.Name}}(exec boil.Executor) error {
  var err error

  {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}.Valid = false
  if err = {{.Function.Receiver}}.Update(exec, "{{.ForeignKey.Column}}"); err != nil {
    {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}.Valid = true
    return err
  }

  {{.Function.Receiver}}.R.{{.Function.Name}} = nil
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
