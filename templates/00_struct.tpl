{{- define "relationship_to_one_struct_helper" -}}
  {{.Function.Name}} *{{.ForeignTable.NameGo}}
{{- end -}}

{{- $tableNameSingular := .Table.Name | singular -}}
{{- $modelName := $tableNameSingular | titleCase -}}
// {{$modelName}} is an object representing the database table.
type {{$modelName}} struct {
  {{range $column := .Table.Columns -}}
  {{titleCase $column.Name}} {{$column.Type}} `boil:"{{$column.Name}}" json:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}"`
  {{end -}}
  {{- if .Table.IsJoinTable -}}
  {{- else}}
  //Relationships *{{$modelName}}Relationships `boil:"-" json:"-" toml:"-" yaml:"-"`
  {{end -}}
}

{{- $dot := . -}}
{{- if .Table.IsJoinTable -}}
{{- else}}
// {{$modelName}}Relationships are where relationships are both cached
// and eagerly loaded.
type {{$modelName}}Relationships struct {
  {{range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
    {{- template "relationship_to_one_struct_helper" $rel}}
  {{end -}}
  {{- range .Table.ToManyRelationships -}}
    {{- if (and .ForeignColumnUnique (not .ToJoinTable)) -}}
      {{- template "relationship_to_one_struct_helper" (textsFromOneToOneRelationship $dot.PkgName $dot.Tables $dot.Table .)}}
    {{else -}}
    {{- $rel := textsFromRelationship $dot.Tables $dot.Table . -}}
  {{$rel.Function.Name}} {{$rel.ForeignTable.Slice}}
{{end -}}{{/* if ForeignColumnUnique */}}
{{- end -}}{{/* range tomany */}}
}

{{end -}}
