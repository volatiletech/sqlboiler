{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . -}}
	{{- range .Table.FKeys -}}
		{{- $txt := textsFromForeignKey $dot.PkgName $dot.Tables $dot.Table . -}}
func test{{$txt.LocalTable.NameGo}}ToOne{{$txt.ForeignTable.NameGo}}_{{$txt.Function.Name}}(t *testing.T) {
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var foreign {{$txt.ForeignTable.NameGo}}
	var local {{$txt.LocalTable.NameGo}}
	{{if .ForeignKey.Nullable -}}
	local.{{.ForeignKey.Column | titleCase}}.Valid = true
	{{end}}
	{{- if .ForeignKey.ForeignColumnNullable -}}
	foreign.{{.ForeignKey.ForeignColumn | titleCase}}.Valid = true
	{{end}}

	{{if not $txt.Function.OneToOne -}}
	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}

	local.{{$txt.Function.LocalAssignment}} = foreign.{{$txt.Function.ForeignAssignment}}
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}
	{{else -}}
	if err := local.Insert(tx); err != nil {
		t.Fatal(err)
	}

	foreign.{{$txt.Function.ForeignAssignment}} = local.{{$txt.Function.LocalAssignment}}
	if err := foreign.Insert(tx); err != nil {
		t.Fatal(err)
	}
	{{end -}}

	check, err := local.{{$txt.Function.Name}}(tx).One()
	if err != nil {
		t.Fatal(err)
	}

	{{if $txt.Function.UsesBytes -}}
	if 0 != bytes.Compare(check.{{$txt.Function.ForeignAssignment}}, foreign.{{$txt.Function.ForeignAssignment}}) {
	{{else -}}
	if check.{{$txt.Function.ForeignAssignment}} != foreign.{{$txt.Function.ForeignAssignment}} {
	{{end -}}
		t.Errorf("want: %v, got %v", foreign.{{$txt.Function.ForeignAssignment}}, check.{{$txt.Function.ForeignAssignment}})
	}

	slice := {{$txt.LocalTable.NameGo}}Slice{&local}
	if err = local.L.Load{{$txt.Function.Name}}(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$txt.Function.Name}} == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.{{$txt.Function.Name}} = nil
	if err = local.L.Load{{$txt.Function.Name}}(tx, true, &local); err != nil {
		t.Fatal(err)
	}
	if local.R.{{$txt.Function.Name}} == nil {
		t.Error("struct should have been eager loaded")
	}
}
{{end -}}{{/* range */}}
{{- end -}}{{/* join table */}}
