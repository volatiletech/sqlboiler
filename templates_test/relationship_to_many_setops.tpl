{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- $table := .Table -}}
  {{- range .Table.ToManyRelationships -}}
    {{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
{{- template "relationship_to_one_setops_test_helper" (textsFromOneToOneRelationship $dot.PkgName $dot.Tables $table .) -}}
    {{- else -}}
    {{- $varNameSingular := .Table | singular | camelCase -}}
    {{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
    {{- $rel := textsFromRelationship $dot.Tables $table .}}

func test{{$rel.LocalTable.NameGo}}ToManyAddOp{{$rel.Function.Name}}(t *testing.T) {
  var err error

  tx := MustTx(boil.Begin())
  defer tx.Rollback()

  var a {{$rel.LocalTable.NameGo}}
  var b, c, d, e {{$rel.ForeignTable.NameGo}}

  seed := randomize.NewSeed()
  if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}PrimaryKeyColumns...); err != nil {
    t.Fatal(err)
  }
  foreigners := []*{{$rel.ForeignTable.NameGo}}{&b, &c, &d, &e}
  for _, x := range foreigners {
    if err = randomize.Struct(seed, x, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}PrimaryKeyColumns...); err != nil {
      t.Fatal(err)
    }
  }

  if err := a.Insert(tx); err != nil {
    t.Fatal(err)
  }
  if err = b.Insert(tx); err != nil {
    t.Fatal(err)
  }
  if err = c.Insert(tx); err != nil {
    t.Fatal(err)
  }

  foreignersSplitByInsertion := [][]*{{$rel.ForeignTable.NameGo}}{
    {&b, &c},
    {&d, &e},
  }

  for i, x := range foreignersSplitByInsertion {
    err = a.Add{{$rel.Function.Name}}(tx, i != 0, x...)
    if err != nil {
      t.Fatal(err)
    }

    first := x[0]
    second := x[1]
    {{- if .ToJoinTable}}

    if first.R.{{$rel.Function.ForeignName}}[0] != &a {
      t.Error("relationship was not added properly to the slice")
    }
    if second.R.{{$rel.Function.ForeignName}}[0] != &a {
      t.Error("relationship was not added properly to the slice")
    }
    {{- else}}

    if a.{{$rel.Function.LocalAssignment}} != first.{{$rel.Function.ForeignAssignment}} {
      t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, first.{{$rel.Function.ForeignAssignment}})
    }
    if a.{{$rel.Function.LocalAssignment}} != second.{{$rel.Function.ForeignAssignment}} {
      t.Error("foreign key was wrong value", a.{{$rel.Function.LocalAssignment}}, second.{{$rel.Function.ForeignAssignment}})
    }

    if first.R.{{$rel.Function.ForeignName}} != &a {
      t.Error("relationship was not added properly to the foreign slice")
    }
    if second.R.{{$rel.Function.ForeignName}} != &a {
      t.Error("relationship was not added properly to the foreign slice")
    }
    {{- end}}

    if a.R.{{$rel.Function.Name}}[i*2] != first {
      t.Error("relationship struct slice not set to correct value")
    }
    if a.R.{{$rel.Function.Name}}[i*2+1] != second {
      t.Error("relationship struct slice not set to correct value")
    }

    count, err := a.{{$rel.Function.Name}}(tx).Count()
    if err != nil {
      t.Fatal(err)
    }
    if want := int64((i+1)*2); count != want {
      t.Error("want", want, "got", count)
    }
  }
}
{{if .ForeignColumnNullable}}

func test{{$rel.LocalTable.NameGo}}ToManySetOp{{$rel.Function.Name}}(t *testing.T) {
}

func test{{$rel.LocalTable.NameGo}}ToManyRemoveOp{{$rel.Function.Name}}(t *testing.T) {
}
{{end -}}
{{- end -}}{{- /* if unique foreign key */ -}}
{{- end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join table */ -}}
