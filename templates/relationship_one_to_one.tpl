{{- define "relationship_to_one_helper" -}}
	{{- $dot := .Dot -}}{{/* .Dot holds the root templateData struct, passed in through preserveDot */}}
	{{- with .Rel -}}{{/* Rel holds the text helper data, passed in through preserveDot */}}
		{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase -}}
// {{.Function.Name}}G pointed to by the foreign key.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}G(mods ...qm.QueryMod) {{$varNameSingular}}Query {
	return {{.Function.Receiver}}.{{.Function.Name}}(boil.GetDB(), mods...)
}

// {{.Function.Name}} pointed to by the foreign key.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) ({{$varNameSingular}}Query) {
	queryMods := []qm.QueryMod{
		qm.Where("{{.ForeignTable.ColumnName}}={{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}", {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}}),
	}

	queryMods = append(queryMods, mods...)

	query := {{.ForeignTable.NamePluralGo}}(exec, queryMods...)
	queries.SetFrom(query.Query, "{{.ForeignTable.Name | $dot.SchemaTable}}")

	return query
}
	{{- end -}}{{/* end with */}}
{{end -}}{{/* end define */}}

{{- /* Begin execution of template for one-to-one relationship */ -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $txt := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
{{- template "relationship_to_one_helper" (preserveDot $dot $txt) -}}
{{- end -}}
{{- end -}}
