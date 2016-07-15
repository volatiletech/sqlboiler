{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.Tables $dot.Table . -}}
func Test{{$rel.LocalTable.NameGo}}ToOne{{$rel.ForeignTable.NameGo}}_{{$rel.LocalTable.ColumnNameGo}}(t *testing.T) {
  tx, err := boil.Begin()
  if err != nil {
    t.Fatal(err)
  }
  defer tx.Rollback()

  var foreign {{$rel.ForeignTable.NameGo}}
  var local {{$rel.LocalTable.NameGo}}

  if err := foreign.InsertX(tx); err != nil {
    t.Fatal(err)
  }

  local.{{$rel.Function.LocalAssignment}} = foreign.{{$rel.Function.ForeignAssignment}} 
  if err := local.InsertX(tx); err != nil {
    t.Fatal(err)
  }

  checkForeign, err := local.{{$rel.LocalTable.ColumnNameGo}}X(tx)
  if err != nil {
    t.Fatal(err)
  }

  if checkForeign.{{$rel.Function.ForeignAssignment}} != foreign.{{$rel.Function.ForeignAssignment}} {
    t.Errorf("want: %v, got %v", foreign.{{$rel.Function.ForeignAssignment}}, checkForeign.{{$rel.Function.ForeignAssignment}})
  }
}

{{end -}}
{{- end -}}
