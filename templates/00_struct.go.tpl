{{- $alias := .Aliases.Table .Table.Name -}}

// {{$alias.UpSingular}} is an object representing the database table.
type {{$alias.UpSingular}} struct {
	{{- range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{if eq $.StructTagCasing "title" -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name | titleCase}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name | titleCase}}" yaml:"{{$column.Name | titleCase}}{{if $column.Nullable}},omitempty{{end}}"`
	{{else if eq $.StructTagCasing "camel" -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name | camelCase}}" yaml:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}"`
	{{else -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}"`
	{{end -}}
	{{end -}}
	{{- if .Table.IsJoinTable -}}
	{{- else}}
	R *{{$alias.DownSingular}}R `{{generateIgnoreTags $.Tags}}boil:"-" json:"-" toml:"-" yaml:"-"`
	L {{$alias.DownSingular}}L `{{generateIgnoreTags $.Tags}}boil:"-" json:"-" toml:"-" yaml:"-"`
	{{end -}}
}

var {{$alias.UpSingular}}Columns = struct {
	{{range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}} string
	{{end -}}
}{
	{{range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}: "{{$column.Name}}",
	{{end -}}
}

{{/* Generated where helpers for all types in the database */}}
// Generated where
{{- range .Table.Columns -}}
	{{- if (oncePut $.DBTypes .Type)}}
	{{$name := printf "whereHelper%s" (goVarname .Type)}}
type {{$name}} struct { field string }
func (w {{$name}}) EQ(x {{.Type}}) qm.QueryMod { return qmhelper.Where{{if .Nullable}}NullEQ(w.field, false, x){{else}}(w.field, qmhelper.EQ, x){{end}} }
func (w {{$name}}) NEQ(x {{.Type}}) qm.QueryMod { return qmhelper.Where{{if .Nullable}}NullEQ(w.field, true, x){{else}}(w.field, qmhelper.NEQ, x){{end}} }
{{if .Nullable -}}
func (w {{$name}}) IsNull() qm.QueryMod { return qmhelper.WhereIsNull(w.field) }
func (w {{$name}}) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
{{end -}}
func (w {{$name}}) LT(x {{.Type}}) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w {{$name}}) LTE(x {{.Type}}) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w {{$name}}) GT(x {{.Type}}) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w {{$name}}) GTE(x {{.Type}}) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
{{if or ( or (eq .Type "int") (eq .Type "int64")) (eq .Type "string") -}}
func (w {{$name}}) IN(slice []{{.Type}}) qm.QueryMod {
  values := make([]interface{}, 0, len(slice))
  for _, value := range slice {
    values = append(values, value)
  }
  return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
{{end -}}
	{{- end -}}
{{- end}}

var {{$alias.UpSingular}}Where = struct {
	{{range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}} whereHelper{{goVarname $column.Type}}
	{{end -}}
}{
	{{range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}: whereHelper{{goVarname $column.Type}}{field: "{{$.Table.Name | $.SchemaTable}}.{{$column.Name | $.Quotes}}"},
	{{end -}}
}

{{- if .Table.IsJoinTable -}}
{{- else}}
// {{$alias.UpSingular}}Rels is where relationship names are stored.
var {{$alias.UpSingular}}Rels = struct {
	{{range .Table.FKeys -}}
	{{- $relAlias := $alias.Relationship .Name -}}
	{{$relAlias.Foreign}} string
	{{end -}}

	{{range .Table.ToOneRelationships -}}
	{{- $ftable := $.Aliases.Table .ForeignTable -}}
	{{- $relAlias := $ftable.Relationship .Name -}}
	{{$relAlias.Local}} string
	{{end -}}

	{{range .Table.ToManyRelationships -}}
	{{- $relAlias := $.Aliases.ManyRelationship .ForeignTable .Name .JoinTable .JoinLocalFKeyName -}}
	{{$relAlias.Local}} string
	{{end -}}{{/* range tomany */}}
}{
	{{range .Table.FKeys -}}
	{{- $relAlias := $alias.Relationship .Name -}}
	{{$relAlias.Foreign}}: "{{$relAlias.Foreign}}",
	{{end -}}

	{{range .Table.ToOneRelationships -}}
	{{- $ftable := $.Aliases.Table .ForeignTable -}}
	{{- $relAlias := $ftable.Relationship .Name -}}
	{{$relAlias.Local}}: "{{$relAlias.Local}}",
	{{end -}}

	{{range .Table.ToManyRelationships -}}
	{{- $relAlias := $.Aliases.ManyRelationship .ForeignTable .Name .JoinTable .JoinLocalFKeyName -}}
	{{$relAlias.Local}}: "{{$relAlias.Local}}",
	{{end -}}{{/* range tomany */}}
}

// {{$alias.DownSingular}}R is where relationships are stored.
type {{$alias.DownSingular}}R struct {
	{{range .Table.FKeys -}}
	{{- $ftable := $.Aliases.Table .ForeignTable -}}
	{{- $relAlias := $alias.Relationship .Name -}}
	{{$relAlias.Foreign}} *{{$ftable.UpSingular}}
	{{end -}}

	{{range .Table.ToOneRelationships -}}
	{{- $ftable := $.Aliases.Table .ForeignTable -}}
	{{- $relAlias := $ftable.Relationship .Name -}}
	{{$relAlias.Local}} *{{$ftable.UpSingular}}
	{{end -}}

	{{range .Table.ToManyRelationships -}}
	{{- $ftable := $.Aliases.Table .ForeignTable -}}
	{{- $relAlias := $.Aliases.ManyRelationship .ForeignTable .Name .JoinTable .JoinLocalFKeyName -}}
	{{$relAlias.Local}} {{printf "%sSlice" $ftable.UpSingular}}
	{{end -}}{{/* range tomany */}}
}

// NewStruct creates a new relationship struct
func (*{{$alias.DownSingular}}R) NewStruct() *{{$alias.DownSingular}}R {
	return &{{$alias.DownSingular}}R{}
}

// {{$alias.DownSingular}}L is where Load methods for each relationship are stored.
type {{$alias.DownSingular}}L struct{}
{{end -}}
