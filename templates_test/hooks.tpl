{{- if not .NoHooks -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func {{$varNameSingular}}BeforeInsertHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterInsertHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterSelectHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func {{$varNameSingular}}BeforeUpdateHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterUpdateHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func {{$varNameSingular}}BeforeDeleteHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterDeleteHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func {{$varNameSingular}}BeforeUpsertHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func {{$varNameSingular}}AfterUpsertHook(e boil.Executor, o *{{$tableNameSingular}}) error {
	*o = {{$tableNameSingular}}{}
	return nil
}

func test{{$tableNamePlural}}Hooks(t *testing.T) {
	t.Parallel()

	var err error

	empty := &{{$tableNameSingular}}{}
	o := &{{$tableNameSingular}}{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, {{$varNameSingular}}DBTypes, false); err != nil {
		t.Errorf("Unable to randomize {{$tableNameSingular}} object: %s", err)
	}

	Add{{$tableNameSingular}}Hook(boil.BeforeInsertHook, {{$varNameSingular}}BeforeInsertHook)
	if err = o.doBeforeInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}BeforeInsertHooks = []{{$tableNameSingular}}Hook{}

	Add{{$tableNameSingular}}Hook(boil.AfterInsertHook, {{$varNameSingular}}AfterInsertHook)
	if err = o.doAfterInsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterInsertHooks = []{{$tableNameSingular}}Hook{}

	Add{{$tableNameSingular}}Hook(boil.AfterSelectHook, {{$varNameSingular}}AfterSelectHook)
	if err = o.doAfterSelectHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterSelectHooks = []{{$tableNameSingular}}Hook{}

	Add{{$tableNameSingular}}Hook(boil.BeforeUpdateHook, {{$varNameSingular}}BeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}BeforeUpdateHooks = []{{$tableNameSingular}}Hook{}

	Add{{$tableNameSingular}}Hook(boil.AfterUpdateHook, {{$varNameSingular}}AfterUpdateHook)
	if err = o.doAfterUpdateHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterUpdateHooks = []{{$tableNameSingular}}Hook{}

	Add{{$tableNameSingular}}Hook(boil.BeforeDeleteHook, {{$varNameSingular}}BeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}BeforeDeleteHooks = []{{$tableNameSingular}}Hook{}

	Add{{$tableNameSingular}}Hook(boil.AfterDeleteHook, {{$varNameSingular}}AfterDeleteHook)
	if err = o.doAfterDeleteHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterDeleteHooks = []{{$tableNameSingular}}Hook{}

	Add{{$tableNameSingular}}Hook(boil.BeforeUpsertHook, {{$varNameSingular}}BeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}BeforeUpsertHooks = []{{$tableNameSingular}}Hook{}

	Add{{$tableNameSingular}}Hook(boil.AfterUpsertHook, {{$varNameSingular}}AfterUpsertHook)
	if err = o.doAfterUpsertHooks(nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	{{$varNameSingular}}AfterUpsertHooks = []{{$tableNameSingular}}Hook{}
}
{{- end}}
