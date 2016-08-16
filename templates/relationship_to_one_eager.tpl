{{- define "relationship_to_one_eager_helper" -}}
  {{- $arg := printf "maybe%s" .LocalTable.NameGo -}}
  {{- $slice := printf "%sSlice" .LocalTable.NameGo}}
// Load{{.Function.Name}} allows an eager lookup of values, cached into the
// relationships structs of the objects.
func (r *{{.LocalTable.NameGo}}Relationships) Load{{.Function.Name}}(e boil.Executor, singular bool, {{$arg}} interface{}) error {
  var slice []*{{.LocalTable.NameGo}}
  var object *{{.LocalTable.NameGo}}

  count := 1
  if singular {
    object = {{$arg}}.(*{{.LocalTable.NameGo}})
  } else {
    slice = {{$arg}}.({{$slice}})
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
    `select * from "{{.ForeignKey.ForeignTable}}" where "{{.ForeignKey.ForeignColumn}}" in (%s)`,
    strmangle.Placeholders(count, 1, 1),
  )

  results, err := e.Query(query, args...)
  if err != nil {
    return errors.Wrap(err, "failed to eager load {{.ForeignTable.NameGo}}")
  }
  defer results.Close()

  var resultSlice []*{{.ForeignTable.NameGo}}
  if err = boil.Bind(results, &resultSlice); err != nil {
    return errors.Wrap(err, "failed to bind eager loaded slice {{.ForeignTable}}")
  }

  if singular && len(resultSlice) != 0 {
    object.Relationships = &{{.LocalTable.NameGo}}Relationships{
      {{.Function.Name}}: resultSlice[0],
    }
    return nil
  }

  for _, foreign := range resultSlice {
    for _, local := range slice {
      if local.{{.Function.LocalAssignment}} == foreign.{{.Function.ForeignAssignment}} {
        if local.Relationships == nil {
          local.Relationships = &{{.LocalTable.NameGo}}Relationships{}
        }
        local.Relationships.{{.Function.Name}} = foreign
        break
      }
    }
  }

  return nil
}
{{- end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
{{- template "relationship_to_one_eager_helper" $rel -}}
{{end -}}
{{- end -}}
