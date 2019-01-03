{{- $alias := .Aliases.Table .Table.Name -}}

// {{$alias.UpSingular}} is an object representing the database table.
type {{$alias.UpSingular}} struct {
	{{- range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{- if eq $.StructTagCasing "camel"}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name | camelCase}}" yaml:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}"`
	{{- else -}}
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

var {{$alias.UpSingular}}Where = struct {
	{{range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}EQ  func({{$column.Type}}) qm.QueryMod
	{{$colAlias}}NEQ func({{$column.Type}}) qm.QueryMod
	{{if isPrimitive $column.Type -}}
	{{$colAlias}}LT  func({{$column.Type}}) qm.QueryMod
	{{$colAlias}}LTE func({{$column.Type}}) qm.QueryMod
	{{$colAlias}}GT  func({{$column.Type}}) qm.QueryMod
	{{$colAlias}}GTE func({{$column.Type}}) qm.QueryMod
	{{end -}}
	{{end -}}
}{
	{{range $column := .Table.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}EQ: func(x {{$column.Type}}) qm.QueryMod {
		{{- if $column.Nullable -}}
		return qmhelper.WhereNullEQ(`{{$column.Name}}`, false, x)
		{{else -}}
		return qmhelper.Where(`{{$column.Name}}`, qmhelper.EQ, x)
		{{- end -}}
	},
	{{$colAlias}}NEQ: func(x {{$column.Type}}) qm.QueryMod {
		{{- if $column.Nullable -}}
		return qmhelper.WhereNullEQ(`{{$column.Name}}`, true, x)
		{{else -}}
		return qmhelper.Where(`{{$column.Name}}`, qmhelper.NEQ, x)
		{{- end -}}
	},
	{{if isPrimitive $column.Type -}}
	{{$colAlias}}LT: func(x {{$column.Type}}) qm.QueryMod { return qmhelper.Where(`{{$column.Name}}`, qmhelper.LT, x) },
	{{$colAlias}}LTE: func(x {{$column.Type}}) qm.QueryMod { return qmhelper.Where(`{{$column.Name}}`, qmhelper.LTE, x) },
	{{$colAlias}}GT: func(x {{$column.Type}}) qm.QueryMod { return qmhelper.Where(`{{$column.Name}}`, qmhelper.GT, x) },
	{{$colAlias}}GTE: func(x {{$column.Type}}) qm.QueryMod { return qmhelper.Where(`{{$column.Name}}`, qmhelper.GTE, x) },
	{{end -}}
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
