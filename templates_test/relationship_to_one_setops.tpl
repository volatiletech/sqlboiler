{{- define "relationship_to_one_setops_test_helper" -}}
{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase -}}
func test{{.LocalTable.NameGo}}ToOneSetOp{{.ForeignTable.NameGo}}_{{.Function.Name}}(t *testing.T) {
  var err error

  tx := MustTx(boil.Begin())
  defer tx.Rollback()

  var a {{.LocalTable.NameGo}}
  var b, c {{.ForeignTable.NameGo}}

  seed := randomize.NewSeed()
  if err = randomize.Struct(seed, &a, {{.ForeignKey.Table | singular | camelCase}}DBTypes, false); err != nil {
    t.Fatal(err)
  }
  if err = randomize.Struct(seed, &b, {{$varNameSingular}}DBTypes, false); err != nil {
    t.Fatal(err)
  }
  if err = randomize.Struct(seed, &c, {{$varNameSingular}}DBTypes, false); err != nil {
    t.Fatal(err)
  }

  if err := a.Insert(tx); err != nil {
    t.Fatal(err)
  }
  if err = b.Insert(tx); err != nil {
    t.Fatal(err)
  }

  for i, x := range []*{{.ForeignTable.NameGo}}{&b, &c} {
    err = a.Set{{.Function.Name}}(tx, i != 0, x)
    if err != nil {
      t.Fatal(err)
    }

    if a.{{.Function.LocalAssignment}} != x.{{.Function.ForeignAssignment}} {
      t.Error("foreign key was wrong value", a.{{.Function.LocalAssignment}})
    }
    if a.R.{{.Function.Name}} != x {
      t.Error("relationship struct not set to correct value")
    }

    zero := reflect.Zero(reflect.TypeOf(a.{{.Function.LocalAssignment}}))
    reflect.Indirect(reflect.ValueOf(&a.{{.Function.LocalAssignment}})).Set(zero)

    if err = a.Reload(tx); err != nil {
      t.Fatal("failed to reload", err)
    }

    if a.{{.Function.LocalAssignment}} != x.{{.Function.ForeignAssignment}} {
      t.Error("foreign key was wrong value", a.{{.Function.LocalAssignment}}, x.{{.Function.ForeignAssignment}})
    }
  }
}
{{- if .ForeignKey.Nullable}}

func test{{.LocalTable.NameGo}}ToOneRemoveOp{{.ForeignTable.NameGo}}_{{.Function.Name}}(t *testing.T) {
}
{{end -}}
{{- end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table .}}

{{template "relationship_to_one_setops_test_helper" $rel -}}
{{- end -}}
{{- end -}}
