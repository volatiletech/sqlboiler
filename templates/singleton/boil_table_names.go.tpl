var TableNames = struct {
	{{range $tblAlias := .Aliases.Tables -}}
	{{$tblAlias.UpPlural}} string
	{{end -}}
}{
	{{range $tblName, $tblAlias := .Aliases.Tables -}}
	{{$tblAlias.UpPlural}}: "{{$tblName}}",
	{{end -}}
}
