{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $dot := . }}
	{{- $table := .Table }}
	{{- range .Table.ToManyRelationships -}}
	{{- $txt := txtsFromToMany $dot.Tables $table .}}
	{{- $varNameSingular := .Table | singular | camelCase -}}
	{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
func test{{$txt.LocalTable.NameGo}}ToMany{{$txt.Function.Name}}(t *testing.T) {
	var err error
	tx := MustTx(boil.Begin())
	defer tx.Rollback()

	var a {{$txt.LocalTable.NameGo}}
	var b, c {{$txt.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.LocalTable.NameGo}} struct: %s", err)
	}

	if err := a.Insert(tx); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}ColumnsWithDefault...)
	randomize.Struct(seed, &c, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}ColumnsWithDefault...)
	{{if .Nullable -}}
	a.{{.Column | titleCase}}.Valid = true
	{{- end}}
	{{- if .ForeignColumnNullable -}}
	b.{{.ForeignColumn | titleCase}}.Valid = true
	c.{{.ForeignColumn | titleCase}}.Valid = true
	{{- end}}
	{{if not .ToJoinTable -}}
	b.{{$txt.Function.ForeignAssignment}} = a.{{$txt.Function.LocalAssignment}}
	c.{{$txt.Function.ForeignAssignment}} = a.{{$txt.Function.LocalAssignment}}
	{{- end}}
	if err = b.Insert(tx); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(tx); err != nil {
		t.Fatal(err)
	}

	{{if .ToJoinTable -}}
	_, err = tx.Exec("insert into {{.JoinTable | $dot.SchemaTable}} ({{.JoinLocalColumn | $dot.Quotes}}, {{.JoinForeignColumn | $dot.Quotes}}) values {{if $dot.Dialect.IndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}", a.{{$txt.LocalTable.ColumnNameGo}}, b.{{$txt.ForeignTable.ColumnNameGo}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into {{.JoinTable | $dot.SchemaTable}} ({{.JoinLocalColumn | $dot.Quotes}}, {{.JoinForeignColumn | $dot.Quotes}}) values {{if $dot.Dialect.IndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}", a.{{$txt.LocalTable.ColumnNameGo}}, c.{{$txt.ForeignTable.ColumnNameGo}})
	if err != nil {
		t.Fatal(err)
	}
	{{end}}

	{{$varname := .ForeignTable | singular | camelCase -}}
	{{$varname}}, err := a.{{$txt.Function.Name}}(tx).All()
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range {{$varname}} {
		{{if $txt.Function.UsesBytes -}}
		if 0 == bytes.Compare(v.{{$txt.Function.ForeignAssignment}}, b.{{$txt.Function.ForeignAssignment}}) {
			bFound = true
		}
		if 0 == bytes.Compare(v.{{$txt.Function.ForeignAssignment}}, c.{{$txt.Function.ForeignAssignment}}) {
			cFound = true
		}
		{{else -}}
		if v.{{$txt.Function.ForeignAssignment}} == b.{{$txt.Function.ForeignAssignment}} {
			bFound = true
		}
		if v.{{$txt.Function.ForeignAssignment}} == c.{{$txt.Function.ForeignAssignment}} {
			cFound = true
		}
		{{end -}}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := {{$txt.LocalTable.NameGo}}Slice{&a}
	if err = a.L.Load{{$txt.Function.Name}}(tx, false, &slice); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.{{$txt.Function.Name}}); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.{{$txt.Function.Name}} = nil
	if err = a.L.Load{{$txt.Function.Name}}(tx, true, &a); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.{{$txt.Function.Name}}); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", {{$varname}})
	}
}

{{end -}}{{- /* range */ -}}
{{- end -}}{{- /* outer if join table */ -}}
