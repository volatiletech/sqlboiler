{{- define "timestamp_insert_helper" -}}
  {{- if eq .NoAutoTimestamps false -}}
  {{- $colNames := .Table.Columns | columnNames -}}
  {{if containsAny $colNames "created_at" "updated_at"}}
  loc := boil.GetLocation()
  currTime := time.Time{}
  if loc != nil {
    currTime = time.Now().In(boil.GetLocation())
  } else {
    currTime = time.Now()
  }
    {{range $ind, $col := .Table.Columns}}
      {{- if eq $col.Name "created_at" -}}
        {{- if $col.Nullable}}
  o.CreatedAt.Time = currTime
  o.CreatedAt.Valid = true
        {{- else}}
  o.CreatedAt = currTime
        {{- end -}}
      {{- end -}}
      {{- if eq $col.Name "updated_at" -}}
        {{- if $col.Nullable}}
  o.UpdatedAt.Time = currTime
  o.UpdatedAt.Valid = true
        {{- else}}
  o.UpdatedAt = currTime
        {{- end -}}
      {{- end -}}
    {{end}}
  {{end}}
  {{- end}}
{{- end -}}
{{- define "timestamp_update_helper" -}}
  {{- if eq .NoAutoTimestamps false -}}
  {{- $colNames := .Table.Columns | columnNames -}}
  {{if containsAny $colNames "updated_at"}}
  loc := boil.GetLocation()
  currTime := time.Time{}
  if loc != nil {
    currTime = time.Now().In(boil.GetLocation())
  } else {
    currTime = time.Now()
  }
    {{range $ind, $col := .Table.Columns}}
      {{- if eq $col.Name "updated_at" -}}
        {{- if $col.Nullable}}
  o.UpdatedAt.Time = currTime
  o.UpdatedAt.Valid = true
        {{- else}}
  o.UpdatedAt = currTime
        {{- end -}}
      {{- end -}}
    {{end}}
  {{end}}
  {{- end}}
{{end -}}
