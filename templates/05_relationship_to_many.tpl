{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . }}
  {{- $table := .Table }}
  {{- $localTableSing := .Table.Name | singular -}}
  {{- $localTable := $localTableSing | titleCase -}}
  {{- $colName := $localTableSing | printf "%s_id" -}}
  {{- $receiver := .Table.Name | toLower | substring 0 1 -}}
  {{- range toManyRelationships .Table.Name .Tables -}}
      {{- $foreignTableSing := .ForeignTable | singular}}
      {{- $foreignTable := $foreignTableSing | titleCase}}
      {{- $foreignSlice := $foreignTableSing | camelCase | printf "%sSlice"}}
      {{- $foreignTableHumanReadable := .ForeignTable | replace "_" " " -}}
      {{- $foreignPluralNoun := .ForeignTable | plural | titleCase -}}
      {{- $isForeignKeySimplyTableName := or (eq $colName .ForeignColumn) .ToJoinTable -}}

      {{- if $isForeignKeySimplyTableName -}}
// {{$foreignPluralNoun}} retrieves all the {{$localTableSing}}'s {{$foreignTableHumanReadable}}.
func ({{$receiver}} *{{$localTable}}) {{$foreignPluralNoun}}(

      {{- else -}}
        {{- $fnName := .ForeignColumn | remove "_id" | titleCase | printf "%[2]s%[1]s" $foreignPluralNoun -}}
// {{$fnName}} retrieves all the {{$localTableSing}}'s {{$foreignTableHumanReadable}} via {{.ForeignColumn}} column.
func ({{$receiver}} *{{$localTable}}) {{$fnName}}(

      {{- end -}}
exec boil.Executor, selectCols ...string) ({{$foreignSlice}}, error) {
  var ret {{$foreignSlice}}

    {{if .ToJoinTable -}}
  query := fmt.Sprintf(`select "%s" from {{.ForeignTable}} "{{id 0}}" inner join {{.JoinTable}} as "{{id 1}}" on "{{id 1}}"."{{.JoinForeignColumn}}" = "{{id 0}}"."{{.ForeignColumn}}" where "{{id 1}}"."{{.JoinLocalColumn}}"=$1`, `"{{id 1}}".` + strings.Join(selectCols, `","{{id 0}}"."`))
    {{else -}}
  query := fmt.Sprintf(`select "%s" from {{.ForeignTable}} where "{{.ForeignColumn}}"=$1`, strings.Join(selectCols, `","`))
    {{end}}
  rows, err := exec.Query(query, {{.Column | titleCase | printf "%s.%s" $receiver }})
  if err != nil {
    return nil, fmt.Errorf(`{{$dot.PkgName}}: unable to select from {{.ForeignTable}}: %v`, err)
  }
  defer rows.Close()

  for rows.Next() {
    next := new({{$foreignTable}})

    err = rows.Scan(boil.GetStructPointers(next, selectCols...)...)
    if err != nil {
      return nil, fmt.Errorf(`{{$dot.PkgName}}: unable to scan into {{$foreignTable}}: %v`, err)
    }

    ret = append(ret, next)
  }

  if err = rows.Err(); err != nil {
    return nil, fmt.Errorf(`{{$dot.PkgName}}: unable to select from {{.ForeignTable}}: %v`, err)
  }

  return ret, nil
}
{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* outer if join table */ -}}
