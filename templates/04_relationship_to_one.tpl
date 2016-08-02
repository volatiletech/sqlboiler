{{- define "relationship_to_one_helper"}}
// {{.Function.Name}}G pointed to by the foreign key.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}G(mods ...qm.QueryMod) (*{{.ForeignTable.NameGo}}, error) {
  return {{.Function.Receiver}}.{{.Function.Name}}(boil.GetDB(), mods...)
}

// {{.Function.Name}}GP pointed to by the foreign key. Panics on error.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}GP(mods ...qm.QueryMod) *{{.ForeignTable.NameGo}} {
  o, err := {{.Function.Receiver}}.{{.Function.Name}}(boil.GetDB(), mods...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return o
}

// {{.Function.Name}}P pointed to by the foreign key with exeuctor. Panics on error.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}P(exec boil.Executor, mods ...qm.QueryMod) *{{.ForeignTable.NameGo}} {
  o, err := {{.Function.Receiver}}.{{.Function.Name}}(exec, mods...)
  if err != nil {
    panic(boil.WrapErr(err))
  }

  return o
}

// {{.Function.Name}} pointed to by the foreign key.
func ({{.Function.Receiver}} *{{.LocalTable.NameGo}}) {{.Function.Name}}(exec boil.Executor, mods ...qm.QueryMod) (*{{.ForeignTable.NameGo}}, error) {
  queryMods := make([]qm.QueryMod, 2, len(mods)+2)
  queryMods[0] = qm.From("{{.ForeignTable.Name}}")
  queryMods[1] = qm.Where("{{.ForeignTable.ColumnName}}=$1", {{.Function.Receiver}}.{{.LocalTable.ColumnNameGo}})

  queryMods = append(queryMods, mods...)

  return {{.ForeignTable.NamePluralGo}}(exec, queryMods...).One()
}

{{end -}}
{{- if .Table.IsJoinTable -}}
{{- else -}}
  {{- $dot := . -}}
  {{- range .Table.FKeys -}}
    {{- $rel := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
{{- template "relationship_to_one_helper" $rel -}}
{{end -}}
{{- end -}}
