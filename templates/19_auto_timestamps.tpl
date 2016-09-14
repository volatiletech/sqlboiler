{{- define "timestamp_insert_helper" -}}
	{{- if not .NoAutoTimestamps -}}
	{{- $colNames := .Table.Columns | columnNames -}}
	{{if containsAny $colNames "created_at" "updated_at"}}
	currTime := time.Now().In(boil.GetLocation())
		{{range $ind, $col := .Table.Columns}}
			{{- if eq $col.Name "created_at" -}}
				{{- if $col.Nullable}}
	if o.CreatedAt.Time.IsZero() {
		o.CreatedAt.Time = currTime
		o.CreatedAt.Valid = true
	}
				{{- else}}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
				{{- end -}}
			{{- end -}}
			{{- if eq $col.Name "updated_at" -}}
				{{- if $col.Nullable}}
	if o.UpdatedAt.Time.IsZero() {
		o.UpdatedAt.Time = currTime
		o.UpdatedAt.Valid = true
	}
				{{- else}}
	if o.UpdatedAt.IsZero() {
		o.UpdatedAt = currTime
	}
				{{- end -}}
			{{- end -}}
		{{end}}
	{{end}}
	{{- end}}
{{- end -}}
{{- define "timestamp_update_helper" -}}
	{{- if not .NoAutoTimestamps -}}
	{{- $colNames := .Table.Columns | columnNames -}}
	{{if containsAny $colNames "updated_at"}}
	currTime := time.Now().In(boil.GetLocation())
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
{{- define "timestamp_upsert_helper" -}}
	{{- if not .NoAutoTimestamps -}}
	{{- $colNames := .Table.Columns | columnNames -}}
	{{if containsAny $colNames "created_at" "updated_at"}}
	currTime := time.Now().In(boil.GetLocation())
		{{range $ind, $col := .Table.Columns}}
			{{- if eq $col.Name "created_at" -}}
				{{- if $col.Nullable}}
	if o.CreatedAt.Time.IsZero() {
		o.CreatedAt.Time = currTime
		o.CreatedAt.Valid = true
	}
				{{- else}}
	if o.CreatedAt.IsZero() {
		o.CreatedAt = currTime
	}
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
{{end -}}
