{{- $alias := .Aliases.View .View.Name -}}
{{- $orig_tbl_name := .View.Name -}}

// {{$alias.UpSingular}} is an object representing the database view.
type {{$alias.UpSingular}} struct {
	{{- range $column := .View.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{- $orig_col_name := $column.Name -}}
	{{- range $column.Comment | splitLines -}} // {{ . }}
	{{end -}}
	{{if ignore $orig_tbl_name $orig_col_name $.TagIgnore -}}
	{{$colAlias}} {{$column.Type}} `{{generateIgnoreTags $.Tags}}boil:"{{$column.Name}}" json:"-" toml:"-" yaml:"-"`
	{{else if eq $.StructTagCasing "title" -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name | titleCase}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name | titleCase}}" yaml:"{{$column.Name | titleCase}}{{if $column.Nullable}},omitempty{{end}}"`
	{{else if eq $.StructTagCasing "camel" -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name | camelCase}}" yaml:"{{$column.Name | camelCase}}{{if $column.Nullable}},omitempty{{end}}"`
	{{else if eq $.StructTagCasing "alias" -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $colAlias}}boil:"{{$column.Name}}" json:"{{$colAlias}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$colAlias}}" yaml:"{{$colAlias}}{{if $column.Nullable}},omitempty{{end}}"`
	{{else -}}
	{{$colAlias}} {{$column.Type}} `{{generateTags $.Tags $column.Name}}boil:"{{$column.Name}}" json:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}" toml:"{{$column.Name}}" yaml:"{{$column.Name}}{{if $column.Nullable}},omitempty{{end}}"`
	{{end -}}
	{{end -}}
}

var {{$alias.UpSingular}}Columns = struct {
	{{range $column := .View.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}} string
	{{end -}}
}{
	{{range $column := .View.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}: "{{$column.Name}}",
	{{end -}}
}

var {{$alias.UpSingular}}ViewColumns = struct {
	{{range $column := .View.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}} string
	{{end -}}
}{
	{{range $column := .View.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}: "{{$orig_tbl_name}}.{{$column.Name}}",
	{{end -}}
}

{{/* Generated where helpers for all types in the database */}}
// Generated where
{{- range .View.Columns -}}
	{{- if (oncePut $.DBTypes .Type)}}
		{{$name := printf "whereHelper%s" (goVarname .Type)}}
type {{$name}} struct { field string }
func (w {{$name}}) EQ(x {{.Type}}) qm.QueryMod { return qmhelper.Where{{if .Nullable}}NullEQ(w.field, false, x){{else}}(w.field, qmhelper.EQ, x){{end}} }
func (w {{$name}}) NEQ(x {{.Type}}) qm.QueryMod { return qmhelper.Where{{if .Nullable}}NullEQ(w.field, true, x){{else}}(w.field, qmhelper.NEQ, x){{end}} }
func (w {{$name}}) LT(x {{.Type}}) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LT, x) }
func (w {{$name}}) LTE(x {{.Type}}) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.LTE, x) }
func (w {{$name}}) GT(x {{.Type}}) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GT, x) }
func (w {{$name}}) GTE(x {{.Type}}) qm.QueryMod { return qmhelper.Where(w.field, qmhelper.GTE, x) }
		{{if isPrimitive .Type -}}
func (w {{$name}}) IN(slice []{{.Type}}) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
		values = append(values, value)
	}
	return qm.WhereIn(fmt.Sprintf("%s IN ?", w.field), values...)
}
func (w {{$name}}) NIN(slice []{{.Type}}) qm.QueryMod {
	values := make([]interface{}, 0, len(slice))
	for _, value := range slice {
	  values = append(values, value)
	}
	return qm.WhereNotIn(fmt.Sprintf("%s NOT IN ?", w.field), values...)
}
		{{end -}}
	{{end -}}
	{{if .Nullable -}}
		{{- if (oncePut $.DBTypes (printf "%s.null" .Type))}}
		{{$name := printf "whereHelper%s" (goVarname .Type)}}
func (w {{$name}}) IsNull() qm.QueryMod { return qmhelper.WhereIsNull(w.field) }
func (w {{$name}}) IsNotNull() qm.QueryMod { return qmhelper.WhereIsNotNull(w.field) }
		{{end -}}
	{{end -}}
{{- end}}

var {{$alias.UpSingular}}Where = struct {
	{{range $column := .View.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}} whereHelper{{goVarname $column.Type}}
	{{end -}}
}{
	{{range $column := .View.Columns -}}
	{{- $colAlias := $alias.Column $column.Name -}}
	{{$colAlias}}: whereHelper{{goVarname $column.Type}}{field: "{{$.View.Name | $.SchemaTable}}.{{$column.Name | $.Quotes}}"},
	{{end -}}
}

