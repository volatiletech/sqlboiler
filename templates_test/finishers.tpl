{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func test{{$tableNamePlural}}Bind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}} := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert(tx); err != nil {
		t.Error(err)
	}

	if err = {{$tableNamePlural}}(tx).Bind({{$varNameSingular}}); err != nil {
		t.Error(err)
	}
}

func test{{$tableNamePlural}}One(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}} := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}, {{$varNameSingular}}DBTypes, true, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}.Insert(tx); err != nil {
		t.Error(err)
	}

	if x, err := {{$tableNamePlural}}(tx).One(); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func test{{$tableNamePlural}}All(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	{{$varNameSingular}}One := &{{$tableNameSingular}}{}
	{{$varNameSingular}}Two := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}One, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}
	if err = randomize.Struct(seed, {{$varNameSingular}}Two, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}One.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = {{$varNameSingular}}Two.Insert(tx); err != nil {
		t.Error(err)
	}

	slice, err := {{$tableNamePlural}}(tx).All()
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func test{{$tableNamePlural}}Count(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	{{$varNameSingular}}One := &{{$tableNameSingular}}{}
	{{$varNameSingular}}Two := &{{$tableNameSingular}}{}
	if err = randomize.Struct(seed, {{$varNameSingular}}One, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}
	if err = randomize.Struct(seed, {{$varNameSingular}}Two, {{$varNameSingular}}DBTypes, false, {{$varNameSingular}}ColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} struct: %s", err)
	}

	tx := MustTx(boil.Begin())
	defer tx.Rollback()
	if err = {{$varNameSingular}}One.Insert(tx); err != nil {
		t.Error(err)
	}
	if err = {{$varNameSingular}}Two.Insert(tx); err != nil {
		t.Error(err)
	}

	count, err := {{$tableNamePlural}}(tx).Count()
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}
