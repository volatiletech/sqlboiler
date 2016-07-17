{{- define "relationship_to_one_test_helper"}}
func Test{{.LocalTable.NameGo}}ToOne{{.ForeignTable.NameGo}}_{{.Function.Name}}(t *testing.T) {
  tx := MustTx(boil.Begin())
  defer tx.Rollback()

  var foreign {{.ForeignTable.NameGo}}
  var local {{.LocalTable.NameGo}}
  {{if .ForeignKey.Nullable -}}
  local.{{.ForeignKey.Column | titleCase}}.Valid = true
  {{end}}
  {{- if .ForeignKey.ForeignColumnNullable -}}
  foreign.{{.ForeignKey.ForeignColumn | titleCase}}.Valid = true
  {{end}}

  if err := foreign.InsertX(tx); err != nil {
    t.Fatal(err)
  }

  local.{{.Function.LocalAssignment}} = foreign.{{.Function.ForeignAssignment}}
  if err := local.InsertX(tx); err != nil {
    t.Fatal(err)
  }

  checkForeign, err := local.{{.Function.Name}}X(tx)
  if err != nil {
    t.Fatal(err)
  }

  if checkForeign.{{.Function.ForeignAssignment}} != foreign.{{.Function.ForeignAssignment}} {
    t.Errorf("want: %v, got %v", foreign.{{.Function.ForeignAssignment}}, checkForeign.{{.Function.ForeignAssignment}})
  }
}

{{end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
{{- template "relationship_to_one_test_helper" $rel -}}
{{end -}}
{{- end -}}
