{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $pkg := .PkgName -}}
  {{- $localTable := .Table.Name -}}
  {{- $ltable := .Table.Name | singular | titleCase -}}
  {{- range $table := .Tables -}}
    {{- if eq $table.Name $localTable -}}
    {{- else -}}
      {{ range $fkey := .FKeys -}}
        {{- if eq $localTable $fkey.ForeignTable -}}
          {{- $ftable := $table.Name | plural | titleCase -}}
          {{- $recv := $localTable | substring 0 1 | toLower -}}
          {{- $fn := $ftable -}}
          {{- $col := $localTable | singular | printf "%s_id" -}}
          {{- if eq $col $fkey.Column -}}
            {{- $col := $localTable -}}
          {{- end -}}
func ({{$recv}} *{{$ltable}}) {{$fn}}
{{ end -}}
      {{- end -}}
    {{- end -}}
  {{- end -}}
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
