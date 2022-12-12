{{- $alias := .Aliases.Table .Table.Name}}
{{- $schemaTable := .Table.Name | .SchemaTable}}
{{- $tableName := .Table.Name }}
{{- $canSoftDelete := .Table.CanSoftDelete $.AutoColumns.Deleted }}
// {{$alias.UpPlural}} retrieves all the records using an executor.
func {{$alias.UpPlural}}(mods ...qm.QueryMod) {{$alias.DownSingular}}Query {
    {{if and .AddSoftDeletes $canSoftDelete -}}
    mods = append(mods, qm.From("{{$schemaTable}}"), qmhelper.WhereIsNull("{{$schemaTable}}.{{or $.AutoColumns.Deleted "deleted_at" | $.Quotes}}"))
    {{else -}}
    mods = append(mods, qm.From("{{$schemaTable}}"))
    {{end -}}

    q := NewQuery(mods...)
    if len(queries.GetSelect(q)) == 0 {
        queries.SetSelect(q, []string{"{{$schemaTable}}.*"})
    }

    return {{$alias.DownSingular}}Query{q}
}

// {{$alias.UpPlural}}WithSchema retrieves all the records using an executor.
func {{$alias.UpPlural}}WithSchema(schema string, mods ...qm.QueryMod) {{$alias.DownSingular}}Query {
    schemaTable := fmt.Sprintf("%s.{{$tableName}}", schema)
    {{if and .AddSoftDeletes $canSoftDelete -}}
    softDelStr := "{{or $.AutoColumns.Deleted "deleted_at" | $.Quotes}}"
    mods = append(mods, qm.From(schemaTable), qmhelper.WhereIsNull(fmt.Sprintf("%s.%s", schemaTable, softDelStr)))
    {{else -}}
    mods = append(mods, qm.From(schemaTable))
    {{end -}}

    q := NewQuery(mods...)
    if len(queries.GetSelect(q)) == 0 {
        allSel := fmt.Sprintf("%s.*", schemaTable)
        queries.SetSelect(q, []string{allSel})
    }

    return {{$alias.DownSingular}}Query{q}
}
