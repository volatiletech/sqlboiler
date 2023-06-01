{{- define "where_ilike_override" -}}
    {{$name := printf "whereHelper%s" (goVarname .Type)}}
func (w {{$name}}) ILIKE(x {{.Type}}) qm.QueryMod { return qm.Where(w.field+" ILIKE ?", x) }
func (w {{$name}}) NILIKE(x {{.Type}}) qm.QueryMod { return qm.Where(w.field+" NOT ILIKE ?", x) }
{{- end -}}