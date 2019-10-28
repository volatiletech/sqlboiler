{{- $alias := .Aliases.Table .Table.Name}}

func (o *{{$alias.UpSingular}}) Editor() *{{$alias.UpSingular}}E {
    return &{{$alias.UpSingular}}E{S:o}
}

func {{$alias.UpSingular}}Editor() *{{$alias.UpSingular}}E {
    return &{{$alias.UpSingular}}E{S:&{{$alias.UpSingular}}{}}
}

type {{$alias.UpSingular}}E struct {
    S *{{$alias.UpSingular}}
    columns []string
}

func (e *{{$alias.UpSingular}}E) Insert({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) error {
    return e.S.Insert(ctx, exec, boil.Whitelist(e.columns...))
}

func (e *{{$alias.UpSingular}}E) Update({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) {{if .NoRowsAffected}}error{{else}}(int64, error){{end -}} {
    return e.S.Update(ctx, exec, boil.Whitelist(e.columns...))
}

{{range $column := .Table.Columns -}}
{{- $colAlias := $alias.Column $column.Name -}}
func (e *{{$alias.UpSingular}}E) {{$colAlias}}({{$column.Name | camelCase}} {{$column.Type}}) *{{$alias.UpSingular}}E {
    e.S.{{$colAlias}} = {{$column.Name | camelCase}}
    e.columns = append(e.columns, {{$alias.UpSingular}}Columns.{{$colAlias}})
    return e
}

{{end -}}
