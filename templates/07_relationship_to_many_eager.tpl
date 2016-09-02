{{- if .Table.IsJoinTable -}}
{{- else -}}
{{- $dot := . -}}
{{- range .Table.ToManyRelationships -}}
{{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
  {{- $txt := textsFromOneToOneRelationship $dot.PkgName $dot.Tables $dot.Table . -}}
  {{- template "relationship_to_one_eager_helper" (preserveDot $dot $txt) -}}
{{- else -}}
  {{- $varNameSingular := $dot.Table.Name | singular | camelCase -}}
  {{- $txt := textsFromRelationship $dot.Tables $dot.Table . -}}
  {{- $arg := printf "maybe%s" $txt.LocalTable.NameGo -}}
  {{- $slice := printf "%sSlice" $txt.LocalTable.NameGo -}}
// Load{{$txt.Function.Name}} allows an eager lookup of values, cached into the
// loaded structs of the objects.
func (t *{{$varNameSingular}}L) Load{{$txt.Function.Name}}(e boil.Executor, singular bool, {{$arg}} interface{}) error {
  var slice []*{{$txt.LocalTable.NameGo}}
  var object *{{$txt.LocalTable.NameGo}}

  count := 1
  if singular {
    object = {{$arg}}.(*{{$txt.LocalTable.NameGo}})
  } else {
    slice = *{{$arg}}.(*{{$slice}})
    count = len(slice)
  }

  args := make([]interface{}, count)
  if singular {
    args[0] = object.{{.Column | titleCase}}
  } else {
    for i, obj := range slice {
      args[i] = obj.{{.Column | titleCase}}
    }
  }

    {{if .ToJoinTable -}}
  query := fmt.Sprintf(
    `select "{{id 0}}".*, "{{id 1}}"."{{.JoinLocalColumn}}" from "{{.ForeignTable}}" as "{{id 0}}" inner join "{{.JoinTable}}" as "{{id 1}}" on "{{id 0}}"."{{.ForeignColumn}}" = "{{id 1}}"."{{.JoinForeignColumn}}" where "{{id 1}}"."{{.JoinLocalColumn}}" in (%s)`,
    strmangle.Placeholders(count, 1, 1),
  )
    {{else -}}
  query := fmt.Sprintf(
    `select * from "{{.ForeignTable}}" where "{{.ForeignColumn}}" in (%s)`,
    strmangle.Placeholders(count, 1, 1),
  )
    {{end -}}

  if boil.DebugMode {
    fmt.Fprintf(boil.DebugWriter, "%s\n%v\n", query, args)
  }

  results, err := e.Query(query, args...)
  if err != nil {
    return errors.Wrap(err, "failed to eager load {{.ForeignTable}}")
  }
  defer results.Close()

  var resultSlice []*{{$txt.ForeignTable.NameGo}}
  {{if .ToJoinTable -}}
  {{- $foreignTable := getTable $dot.Tables .ForeignTable -}}
  {{- $joinTable := getTable $dot.Tables .JoinTable -}}
  {{- $localCol := $joinTable.GetColumn .JoinLocalColumn}}
  var localJoinCols []{{$localCol.Type}}
  for results.Next() {
    one := new({{$txt.ForeignTable.NameGo}})
    var localJoinCol {{$localCol.Type}}

    err = results.Scan({{$foreignTable.Columns | columnNames | stringMap $dot.StringFuncs.titleCase | prefixStringSlice "&one." | join ", "}}, &localJoinCol)
    if err = results.Err(); err != nil {
      return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignTable}}")
    }

    resultSlice = append(resultSlice, one)
    localJoinCols = append(localJoinCols, localJoinCol)
  }

  if err = results.Err(); err != nil {
    return errors.Wrap(err, "failed to plebian-bind eager loaded slice {{.ForeignTable}}")
  }
  {{else -}}
  if err = boil.BindFast(results, &resultSlice, {{.ForeignTable | singular | camelCase}}TitleCases); err != nil {
    return errors.Wrap(err, "failed to bind eager loaded slice {{.ForeignTable}}")
  }
  {{end}}

  {{if not $dot.NoHooks -}}
  if len({{.ForeignTable | singular | camelCase}}AfterSelectHooks) != 0 {
    for _, obj := range resultSlice {
      if err := obj.doAfterSelectHooks(e); err != nil {
        return err
      }
    }
  }

  {{- end}}
  if singular {
    if object.R == nil {
      object.R = &{{$varNameSingular}}R{}
    }
    object.R.{{$txt.Function.Name}} = resultSlice
    return nil
  }

  {{if .ToJoinTable -}}
  for i, foreign := range resultSlice {
    localJoinCol := localJoinCols[i]
    for _, local := range slice {
      if local.{{$txt.Function.LocalAssignment}} == localJoinCol {
        if local.R == nil {
          local.R = &{{$varNameSingular}}R{}
        }
        local.R.{{$txt.Function.Name}} = append(local.R.{{$txt.Function.Name}}, foreign)
        break
      }
    }
  }
  {{else -}}
  for _, foreign := range resultSlice {
    for _, local := range slice {
      if local.{{$txt.Function.LocalAssignment}} == foreign.{{$txt.Function.ForeignAssignment}} {
        if local.R == nil {
          local.R = &{{$varNameSingular}}R{}
        }
        local.R.{{$txt.Function.Name}} = append(local.R.{{$txt.Function.Name}}, foreign)
        break
      }
    }
  }
  {{end}}

  return nil
}

{{end -}}{{/* if ForeignColumnUnique */}}
{{- end -}}{{/* range tomany */}}
{{- end -}}{{/* if isjointable */}}
