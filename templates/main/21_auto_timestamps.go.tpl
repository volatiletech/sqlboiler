{{- define "timestamp_insert_helper" -}}
	{{- if not .NoAutoTimestamps -}}
	{{- $alias := .Aliases.Table .Table.Name -}}
	{{- $colNames := .Table.Columns | columnNames -}}
	{{if containsAny $colNames (or $.AutoColumns.Created "created_at") (or $.AutoColumns.Updated "updated_at")}}
		{{if not .NoContext -}}
	if !boil.TimestampsAreSkipped(ctx) {
		{{end -}}
		currTime := time.Now().In(boil.GetLocation())
		{{range $ind, $col := .Table.Columns}}
		    {{- $colAlias := $alias.Column $col.Name -}}
			{{- if eq $col.Name (or $.AutoColumns.Created "created_at") -}}
				{{- if eq $col.Type "time.Time" }}
		if o.{{$colAlias}}.IsZero() {
			o.{{$colAlias}} = currTime
		}
				{{- else}}
		if queries.MustTime(o.{{$colAlias}}).IsZero() {
			queries.SetScanner(&o.{{$colAlias}}, currTime)
		}
				{{- end -}}
			{{- end -}}
			{{- if eq $col.Name (or $.AutoColumns.Updated "updated_at") -}}
				{{- if eq $col.Type "time.Time"}}
		if o.{{$colAlias}}.IsZero() {
			o.{{$colAlias}} = currTime
		}
				{{- else}}
		if queries.MustTime(o.{{$colAlias}}).IsZero() {
			queries.SetScanner(&o.{{$colAlias}}, currTime)
		}
				{{- end -}}
			{{- end -}}
		{{end}}
		{{if not .NoContext -}}
	}
		{{end -}}
	{{end}}
	{{- end}}
{{- end -}}
{{- define "timestamp_update_helper" -}}
	{{- if not .NoAutoTimestamps -}}
	{{- $alias := .Aliases.Table .Table.Name -}}
	{{- $colNames := .Table.Columns | columnNames -}}
	{{if containsAny $colNames (or $.AutoColumns.Updated "updated_at")}}
		{{if not .NoContext -}}
	if !boil.TimestampsAreSkipped(ctx) {
		{{end -}}
		currTime := time.Now().In(boil.GetLocation())
		{{range $ind, $col := .Table.Columns}}
	        {{- $colAlias := $alias.Column $col.Name -}}
			{{- if eq $col.Name (or $.AutoColumns.Updated "updated_at") -}}
				{{- if eq $col.Type "time.Time"}}
		o.{{$colAlias}} = currTime
				{{- else}}
		queries.SetScanner(&o.{{$colAlias}}, currTime)
				{{- end -}}
			{{- end -}}
		{{end}}
		{{if not .NoContext -}}
	}
		{{end -}}
	{{end}}
	{{- end}}
{{end -}}
{{- define "timestamp_upsert_helper" -}}
	{{- if not .NoAutoTimestamps -}}
	{{- $alias := .Aliases.Table .Table.Name -}}
	{{- $colNames := .Table.Columns | columnNames -}}
	{{if containsAny $colNames (or $.AutoColumns.Created "created_at") (or $.AutoColumns.Updated "updated_at")}}
		{{if not .NoContext -}}
	if !boil.TimestampsAreSkipped(ctx) {
		{{end -}}
	currTime := time.Now().In(boil.GetLocation())
		{{range $ind, $col := .Table.Columns}}
		    {{- $colAlias := $alias.Column $col.Name -}}
			{{- if eq $col.Name (or $.AutoColumns.Created "created_at") -}}
				{{- if eq $col.Type "time.Time"}}
	if o.{{$colAlias}}.IsZero() {
		o.{{$colAlias}} = currTime
	}
				{{- else}}
	if queries.MustTime(o.{{$colAlias}}).IsZero() {
		queries.SetScanner(&o.{{$colAlias}}, currTime)
	}
				{{- end -}}
			{{- end -}}
			{{- if eq $col.Name (or $.AutoColumns.Updated "updated_at") -}}
				{{- if eq $col.Type "time.Time"}}
	o.{{$colAlias}} = currTime
				{{- else}}
	queries.SetScanner(&o.{{$colAlias}}, currTime)
				{{- end -}}
			{{- end -}}
		{{end}}
		{{if not .NoContext -}}
	}
		{{end -}}
	{{end}}
	{{- end}}
{{end -}}
