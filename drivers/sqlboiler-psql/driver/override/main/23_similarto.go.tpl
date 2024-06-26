{{- define "where_similarto_override" -}}
    {{$name := printf "whereHelper%s" (goVarname .Type)}}
func (w {{$name}}) SIMILAR(x {{.Type}}) qm.QueryMod { return qm.Where(w.field+" SIMILAR TO ?", x) }
func (w {{$name}}) NSIMILAR(x {{.Type}}) qm.QueryMod { return qm.Where(w.field+" NOT SIMILAR TO ?", x) }
{{- end -}}