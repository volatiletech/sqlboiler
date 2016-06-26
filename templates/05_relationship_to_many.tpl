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
    {{- $isNormal := eq $colName .ForeignColumn -}}

    {{- if $isNormal -}}
// {{$foreignPluralNoun}} retrieves all the {{$localTableSing}}'s {{$foreignTableHumanReadable}}.
func ({{$receiver}} *{{$localTable}}) {{$foreignPluralNoun}}(

    {{- else -}}
      {{- $fnName := .ForeignColumn | remove "_id" | titleCase | printf "%[2]s%[1]s" $foreignPluralNoun -}}
// {{$fnName}} retrieves all the {{$localTableSing}}'s {{$foreignTableHumanReadable}} via {{.ForeignColumn}} column.
func ({{$receiver}} *{{$localTable}}) {{$fnName}}(
    {{- end -}}

exec boil.Executor, selectCols ...string) ({{$foreignSlice}}, error) {
  var ret {{$foreignSlice}}

  query := fmt.Sprintf(`select "%s" from {{.ForeignTable}} where "{{.ForeignColumn}}"=$1`, strings.Join(selectCols, `","`))
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
{{end -}}
{{- end -}}
{{/*
// Challengee fetches the Video pointed to by the foreign key.
func (b *Battle) Challengee(exec boil.Executor, selectCols ...string) (*Video, error) {
	video := &Video{}

	query := fmt.Sprintf(`select %s from videos where id = $1`, strings.Join(selectCols, `,`))
	err := exec.QueryRow(query, b.ChallengeeID).Scan(boil.GetStructPointers(video, selectCols...)...)
	if err != nil {
		return nil, fmt.Errorf(`models: unable to select from videos: %v`, err)
	}

	return video, nil
}

  {{- range .Table.FKeys -}}
    {{- $localColumn := .Column | remove "_id" | singular | titleCase -}}
    {{- $foreignColumn := .Column | remove "_id" | singular | titleCase -}}
    {{- $foreignTable := .ForeignTable | singular | titleCase -}}
    {{- $varname := .ForeignTable | singular | camelCase -}}
    {{- $receiver := $localTable | toLower | substring 0 1 -}}
// {{$foreignColumn}} fetches the {{$foreignTable}} pointed to by the foreign key.
func ({{$receiver}} *{{$localTable}}) {{$foreignColumn}}(exec boil.Executor, selectCols ...string) (*{{$foreignTable}}, error) {
  {{$varname}} := &{{$foreignTable}}{}

  query := fmt.Sprintf(`select %s from {{.ForeignTable}} where id = $1`, strings.Join(selectCols, `,`))
  err := exec.QueryRow(query, {{$receiver}}.{{titleCase .Column}}).Scan(boil.GetStructPointers({{$varname}}, selectCols...)...)
  if err != nil {
    return nil, fmt.Errorf(`{{$pkg}}: unable to select from {{.ForeignTable}}: %v`, err)
  }

  return {{$varname}}, nil
}
*/}}
