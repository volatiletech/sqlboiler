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

  boil.RandomizeStruct(&b, {{$rel.ForeignTable.NameSingular | camelCase}}DBTypes, true, "{{.ForeignColumn}}")
  boil.RandomizeStruct(&c, {{$rel.ForeignTable.NameSingular | camelCase}}DBTypes, true, "{{.ForeignColumn}}")
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
  if err = b.InsertX(tx); err != nil {
    t.Fatal(err)
  }
  if err = c.InsertX(tx); err != nil {
    t.Fatal(err)
  }

  {{if .ToJoinTable -}}
  _, err = tx.Exec(`insert into {{.JoinTable}} ({{.JoinLocalColumn}}, {{.JoinForeignColumn}}) values ($1, $2)`, a.{{.Column | titleCase}}, b.{{.ForeignColumn | titleCase}})
  if err != nil {
    t.Fatal(err)
  }
  _, err = tx.Exec(`insert into {{.JoinTable}} ({{.JoinLocalColumn}}, {{.JoinForeignColumn}}) values ($1, $2)`, a.{{.Column | titleCase}}, c.{{.ForeignColumn | titleCase}})
  if err != nil {
    t.Fatal(err)
  }
  {{end}}

  {{$varname := $rel.ForeignTable.NamePluralGo | toLower -}}
  {{$varname}}, err := a.{{$rel.Function.Name}}X(tx)
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

  if t.Failed() {
    t.Logf("%#v", {{$varname}})
  }
}

{{ end -}}{{- /* range */ -}}
{{- end -}}{{- /* outer if join table */ -}}
