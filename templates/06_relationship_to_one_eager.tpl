{{- define "relationship_to_one_eager_helper" -}}
  {{- $dot := .Dot -}}{{/* .Dot holds the root templateData struct, passed in through preserveDot */}}
  {{- $varNameSingular := $dot.Table.Name | singular | camelCase -}}
  {{- with .Rel -}}
  {{- $arg := printf "maybe%s" .LocalTable.NameGo -}}
  {{- $slice := printf "%sSlice" .LocalTable.NameGo -}}
// Load{{.Function.Name}} allows an eager lookup of values, cached into the
// loaded structs of the objects.
func ({{$varNameSingular}}L) Load{{.Function.Name}}(e boil.Executor, singular bool, {{$arg}} interface{}) error {
  var slice []*{{.LocalTable.NameGo}}
  var object *{{.LocalTable.NameGo}}

  count := 1
  if singular {
    object = {{$arg}}.(*{{.LocalTable.NameGo}})
  } else {
    slice = *{{$arg}}.(*{{$slice}})
    count = len(slice)
  }

  args := make([]interface{}, count)
  if singular {
    args[0] = object.{{.LocalTable.ColumnNameGo}}
  } else {
    for i, obj := range slice {
      args[i] = obj.{{.LocalTable.ColumnNameGo}}
    }
  }

  query := fmt.Sprintf(
    "select * from {{.ForeignKey.ForeignTable | $dot.SchemaTable}} where {{.ForeignKey.ForeignColumn | $dot.Quotes}} in (%s)",
    strmangle.Placeholders(dialect.IndexPlaceholders, count, 1, 1),
  )

  if boil.DebugMode {
    fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
  }

  results, err := e.Query(query, args...)
  if err != nil {
    return errors.Wrap(err, "failed to eager load {{.ForeignTable.NameGo}}")
  }
  defer results.Close()

  var resultSlice []*{{.ForeignTable.NameGo}}
  if err = boil.Bind(results, &resultSlice); err != nil {
    return errors.Wrap(err, "failed to bind eager loaded slice {{.ForeignTable.NameGo}}")
  }

  {{if not $dot.NoHooks -}}
  if len({{.ForeignTable.Name | singular | camelCase}}AfterSelectHooks) != 0 {
    for _, obj := range resultSlice {
      if err := obj.doAfterSelectHooks(e); err != nil {
        return err
      }
    }
  }
  {{- end}}

  if singular && len(resultSlice) != 0 {
    if object.R == nil {
      object.R = &{{$varNameSingular}}R{}
    }
    object.R.{{.Function.Name}} = resultSlice[0]
    return nil
  }

  for _, foreign := range resultSlice {
    for _, local := range slice {
      if local.{{.Function.LocalAssignment}} == foreign.{{.Function.ForeignAssignment}} {
        if local.R == nil {
          local.R = &{{$varNameSingular}}R{}
        }
        local.R.{{.Function.Name}} = foreign
        break
      }
    }
  }

  return nil
}
  {{- end -}}{{- /* end with */ -}}
{{end -}}{{- /* end define */ -}}

{{- /* Begin execution of template for one-to-one eager load */ -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $txt := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
    {{- template "relationship_to_one_eager_helper" (preserveDot $dot $txt) -}}
  {{- end -}}
{{end}}
