{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- $table := .Table }}
	{{- range .Table.ToManyRelationships -}}
	{{- $txt := txtsFromToMany $.Tables $table .}}
	{{- $varNameSingular := .Table | singular | camelCase -}}
	{{- $foreignVarNameSingular := .ForeignTable | singular | camelCase -}}
func test{{$txt.LocalTable.NameGo}}ToMany{{$txt.Function.Name}}(t *testing.T) {
	var err error
	{{if not $.NoContext}}ctx := context.Background(){{end}}
	tx := MustTx({{if $.NoContext}}boil.Begin(){{else}}boil.BeginTx(ctx, nil){{end}})
	defer tx.Rollback()

	var a {{$txt.LocalTable.NameGo}}
	var b, c {{$txt.ForeignTable.NameGo}}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$txt.LocalTable.NameGo}} struct: %s", err)
	}

	if err := a.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	randomize.Struct(seed, &b, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}ColumnsWithDefault...)
	randomize.Struct(seed, &c, {{$foreignVarNameSingular}}DBTypes, false, {{$foreignVarNameSingular}}ColumnsWithDefault...)

	{{if not .ToJoinTable -}}
		{{if $txt.Function.UsesPrimitives}}
	b.{{$txt.ForeignTable.ColumnNameGo}} = a.{{$txt.LocalTable.ColumnNameGo}}
	c.{{$txt.ForeignTable.ColumnNameGo}} = a.{{$txt.LocalTable.ColumnNameGo}}
		{{else -}}
	queries.Assign(&b.{{$txt.ForeignTable.ColumnNameGo}}, a.{{$txt.LocalTable.ColumnNameGo}})
	queries.Assign(&c.{{$txt.ForeignTable.ColumnNameGo}}, a.{{$txt.LocalTable.ColumnNameGo}})
		{{- end}}
	{{- end}}
	if err = b.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert({{if not $.NoContext}}ctx, {{end -}} tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	{{if .ToJoinTable -}}
	_, err = tx.Exec("insert into {{.JoinTable | $.SchemaTable}} ({{.JoinLocalColumn | $.Quotes}}, {{.JoinForeignColumn | $.Quotes}}) values {{if $.Dialect.UseIndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}", a.{{$txt.LocalTable.ColumnNameGo}}, b.{{$txt.ForeignTable.ColumnNameGo}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = tx.Exec("insert into {{.JoinTable | $.SchemaTable}} ({{.JoinLocalColumn | $.Quotes}}, {{.JoinForeignColumn | $.Quotes}}) values {{if $.Dialect.UseIndexPlaceholders}}($1, $2){{else}}(?, ?){{end}}", a.{{$txt.LocalTable.ColumnNameGo}}, c.{{$txt.ForeignTable.ColumnNameGo}})
	if err != nil {
		t.Fatal(err)
	}
	{{end}}

	{{$varname := .ForeignTable | singular | camelCase -}}
	{{$varname}}, err := a.{{$txt.Function.Name}}().All({{if not $.NoContext}}ctx, {{end -}} tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range {{$varname}} {
		{{if $txt.Function.UsesPrimitives -}}
		if v.{{$txt.ForeignTable.ColumnNameGo}} == b.{{$txt.ForeignTable.ColumnNameGo}} {
			bFound = true
		}
		if v.{{$txt.ForeignTable.ColumnNameGo}} == c.{{$txt.ForeignTable.ColumnNameGo}} {
			cFound = true
		}
		{{else -}}
		if queries.Equal(v.{{$txt.ForeignTable.ColumnNameGo}}, b.{{$txt.ForeignTable.ColumnNameGo}}) {
			bFound = true
		}
		if queries.Equal(v.{{$txt.ForeignTable.ColumnNameGo}}, c.{{$txt.ForeignTable.ColumnNameGo}}) {
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
	if err = a.L.Load{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} tx, false, (*[]*{{$txt.LocalTable.NameGo}})(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.{{$txt.Function.Name}}); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.{{$txt.Function.Name}} = nil
	if err = a.L.Load{{$txt.Function.Name}}({{if not $.NoContext}}ctx, {{end -}} tx, true, &a, nil); err != nil {
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
