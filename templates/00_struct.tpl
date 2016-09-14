{{- define "relationship_to_one_struct_helper" -}}
	{{.Function.Name}} *{{.ForeignTable.NameGo}}
{{- end -}}

{{- $dot := . -}}
{{- $tableNameSingular := .Table.Name | singular -}}
{{- $modelName := $tableNameSingular | titleCase -}}
{{- $modelNameCamel := $tableNameSingular | camelCase -}}
// {{$modelName}} is an object representing the database table.
type {{$modelName}} struct {
	{{range $column := .Table.Columns -}}
	{{titleCase $column.Name}} {{$column.Type}} `{{generateTags $dot.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}"`
	{{end -}}
	{{- if .Table.IsJoinTable -}}
	{{- else}}
	R *{{$modelNameCamel}}R `{{generateIgnoreTags $dot.Tags}}boil:"-" json:"-" toml:"-" yaml:"-"`
	L {{$modelNameCamel}}L `{{generateIgnoreTags $dot.Tags}}boil:"-" json:"-" toml:"-" yaml:"-"`
	{{end -}}
}

{{- if .Table.IsJoinTable -}}
{{- else}}
// {{$modelNameCamel}}R is where relationships are stored.
type {{$modelNameCamel}}R struct {
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

// {{$modelNameCamel}}L is where Load methods for each relationship are stored.
type {{$modelNameCamel}}L struct{}
{{end -}}
