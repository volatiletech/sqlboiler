{{- define "relationship_to_one_test_helper"}}
func test{{.LocalTable.NameGo}}ToOne{{.ForeignTable.NameGo}}_{{.Function.Name}}(t *testing.T) {
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

  {{if not .Function.ReverseInserts -}}
  if err := foreign.Insert(tx); err != nil {
    t.Fatal(err)
  }

  local.{{.Function.LocalAssignment}} = foreign.{{.Function.ForeignAssignment}}
  if err := local.Insert(tx); err != nil {
    t.Fatal(err)
  }
  {{else -}}
  if err := local.Insert(tx); err != nil {
    t.Fatal(err)
  }

  foreign.{{.Function.ForeignAssignment}} = local.{{.Function.LocalAssignment}}
  if err := foreign.Insert(tx); err != nil {
    t.Fatal(err)
  }
  {{end -}}

  check, err := local.{{.Function.Name}}(tx).One()
  if err != nil {
    t.Fatal(err)
  }

  if check.{{.Function.ForeignAssignment}} != foreign.{{.Function.ForeignAssignment}} {
    t.Errorf("want: %v, got %v", foreign.{{.Function.ForeignAssignment}}, check.{{.Function.ForeignAssignment}})
  }

  if err = local.Relationships.Load{{.Function.Name}}(tx, false, {{.LocalTable.NameGo}}Slice{&local}); err != nil {
    t.Fatal(err)
  }
  if local.Relationships.{{.Function.Name}} == nil {
    t.Error("struct should have been eager loaded")
  }

  local.Relationships.{{.Function.Name}} = nil
  if err = local.Relationships.Load{{.Function.Name}}(tx, true, &local); err != nil {
    t.Fatal(err)
  }
  if local.Relationships.{{.Function.Name}} == nil {
    t.Error("struct should have been eager loaded")
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
