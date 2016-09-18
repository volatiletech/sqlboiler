{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $txt := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
		{{- $varNameSingular := .ForeignKey.ForeignTable | singular | camelCase -}}
// {{$txt.Function.Name}}G pointed to by the foreign key.
func ({{$txt.Function.Receiver}} *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}G(mods ...qm.QueryMod) {{$varNameSingular}}Query {
	return {{$txt.Function.Receiver}}.{{$txt.Function.Name}}(boil.GetDB(), mods...)
}

// {{$txt.Function.Name}} pointed to by the foreign key.
func ({{$txt.Function.Receiver}} *{{$txt.LocalTable.NameGo}}) {{$txt.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) ({{$varNameSingular}}Query) {
	queryMods := []qm.QueryMod{
		qm.Where("{{$txt.ForeignTable.ColumnName}}={{if $dot.Dialect.IndexPlaceholders}}$1{{else}}?{{end}}", {{$txt.Function.Receiver}}.{{$txt.LocalTable.ColumnNameGo}}),
	}

	queryMods = append(queryMods, mods...)

	query := {{$txt.ForeignTable.NamePluralGo}}(exec, queryMods...)
	queries.SetFrom(query.Query, "{{$txt.ForeignTable.Name | $dot.SchemaTable}}")

	return query
}
{{- end -}}
{{- end -}}
