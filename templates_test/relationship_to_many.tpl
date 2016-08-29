{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . }}
  {{- $table := .Table }}
  {{- range .Table.ToManyRelationships -}}
    {{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
{{- template "relationship_to_one_test_helper" (textsFromOneToOneRelationship $dot.PkgName $dot.Tables $table .) -}}
    {{- else -}}
    {{- $rel := textsFromRelationship $dot.Tables $table . -}}
func test{{$rel.LocalTable.NameGo}}ToMany{{$rel.Function.Name}}(t *testing.T) {
  var err error
  tx := MustTx(boil.Begin())
  defer tx.Rollback()

  var a {{$rel.LocalTable.NameGo}}
  var b, c {{$rel.ForeignTable.NameGo}}

  if err := a.Insert(tx); err != nil {
    t.Fatal(err)
  }

  seed := randomize.NewSeed()
  randomize.Struct(seed, &b, {{$rel.ForeignTable.NameSingular | camelCase}}DBTypes, false, "{{.ForeignColumn}}")
  randomize.Struct(seed, &c, {{$rel.ForeignTable.NameSingular | camelCase}}DBTypes, false, "{{.ForeignColumn}}")
  {{if .Nullable -}}
  a.{{.Column | titleCase}}.Valid = true
  {{- end}}
  {{- if .ForeignColumnNullable -}}
  b.{{.ForeignColumn | titleCase}}.Valid = true
  c.{{.ForeignColumn | titleCase}}.Valid = true
  {{- end}}
  {{if not .ToJoinTable -}}
  b.{{$rel.Function.ForeignAssignment}} = a.{{$rel.Function.LocalAssignment}}
  c.{{$rel.Function.ForeignAssignment}} = a.{{$rel.Function.LocalAssignment}}
  {{- end}}
  if err = b.Insert(tx); err != nil {
    t.Fatal(err)
  }
  if err = c.Insert(tx); err != nil {
    t.Fatal(err)
  }

  {{if .ToJoinTable -}}
  _, err = tx.Exec(`insert into "{{.JoinTable}}" ({{.JoinLocalColumn}}, {{.JoinForeignColumn}}) values ($1, $2)`, a.{{$rel.LocalTable.ColumnNameGo}}, b.{{$rel.ForeignTable.ColumnNameGo}})
  if err != nil {
    t.Fatal(err)
  }
  _, err = tx.Exec(`insert into "{{.JoinTable}}" ({{.JoinLocalColumn}}, {{.JoinForeignColumn}}) values ($1, $2)`, a.{{$rel.LocalTable.ColumnNameGo}}, c.{{$rel.ForeignTable.ColumnNameGo}})
  if err != nil {
    t.Fatal(err)
  }
  {{end}}

  {{$varname := .ForeignTable | singular | camelCase -}}
  {{$varname}}, err := a.{{$rel.Function.Name}}(tx).All()
  if err != nil {
    t.Fatal(err)
  }

  bFound, cFound := false, false
  for _, v := range {{$varname}} {
    if v.{{$rel.Function.ForeignAssignment}} == b.{{$rel.Function.ForeignAssignment}} {
      bFound = true
    }
    if v.{{$rel.Function.ForeignAssignment}} == c.{{$rel.Function.ForeignAssignment}} {
      cFound = true
    }
  }

  if !bFound {
    t.Error("expected to find b")
  }
  if !cFound {
    t.Error("expected to find c")
  }

  slice := {{$rel.LocalTable.NameGo}}Slice{&a}
  if err = a.R.Load{{$rel.Function.Name}}(tx, false, &slice); err != nil {
    t.Fatal(err)
  }
  if got := len(a.R.{{$rel.Function.Name}}); got != 2 {
    t.Error("number of eager loaded records wrong, got:", got)
  }

  a.R.{{$rel.Function.Name}} = nil
  if err = a.R.Load{{$rel.Function.Name}}(tx, true, &a); err != nil {
    t.Fatal(err)
  }
  if got := len(a.R.{{$rel.Function.Name}}); got != 2 {
    t.Error("number of eager loaded records wrong, got:", got)
  }

  if t.Failed() {
    t.Logf("%#v", {{$varname}})
  }
}

{{end -}}{{- /* if unique */ -}}
{{- end -}}{{- /* range */ -}}
{{- end -}}{{- /* outer if join table */ -}}
