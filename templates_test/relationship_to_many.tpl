{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . }}
  {{- $table := .Table }}
  {{- range toManyRelationships .Table.Name .Tables -}}
    {{- $rel := textsFromRelationship $dot.Tables $table . -}}
func Test{{$rel.LocalTable.NameGo}}ToMany{{$rel.Function.Name}}(t *testing.T) {
  var err error
  tx := MustTx(boil.Begin())
  defer tx.Rollback()

  var a {{$rel.LocalTable.NameGo}}
  var b, c {{$rel.ForeignTable.NameGo}}

  if err := a.InsertX(tx); err != nil {
    t.Fatal(err)
  }

  boil.RandomizeStruct(&b, {{$rel.ForeignTable.NameSingular | camelCase}}DBTypes, true)
  boil.RandomizeStruct(&c, {{$rel.ForeignTable.NameSingular | camelCase}}DBTypes, true)
  b.{{$rel.Function.ForeignAssignment}} = a.{{$rel.Function.LocalAssignment}}
  c.{{$rel.Function.ForeignAssignment}} = a.{{$rel.Function.LocalAssignment}}
  if err := b.InsertX(tx); err != nil {
    t.Fatal(err)
  }
  if err := c.InsertX(tx); err != nil {
    t.Fatal(err)
  }

  {{$varname := $rel.ForeignTable.NamePluralGo | toLower -}}
  {{$varname}}, err := a.{{$rel.Function.Name}}X(tx)

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

  if t.Failed() {
    t.Logf("%#v", {{$varname}})
  }
}

{{ end -}}{{- /* range */ -}}
{{- end -}}{{- /* outer if join table */ -}}
