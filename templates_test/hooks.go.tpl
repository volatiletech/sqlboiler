{{- if not .NoHooks -}}
{{- $alias := .Aliases.Table .Table.Name}}
func {{$alias.DownSingular}}BeforeInsertHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterInsertHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterSelectHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}BeforeUpdateHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterUpdateHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}BeforeDeleteHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterDeleteHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}BeforeUpsertHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterUpsertHook({{if .NoContext}}e boil.Executor{{else}}ctx context.Context, e boil.ContextExecutor{{end}}, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func test{{$alias.UpPlural}}Hooks(t *testing.T) {
	t.Parallel()

	var err error

	{{if not .NoContext}}ctx := context.Background(){{end}}
	empty := &{{$alias.UpSingular}}{}
	o := &{{$alias.UpSingular}}{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, false); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} object: %s", err)
	}

	Add{{$alias.UpSingular}}Hook(boil.BeforeInsertHook, {{$alias.DownSingular}}BeforeInsertHook)
	if err = o.doBeforeInsertHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}BeforeInsertHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterInsertHook, {{$alias.DownSingular}}AfterInsertHook)
	if err = o.doAfterInsertHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterInsertHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterSelectHook, {{$alias.DownSingular}}AfterSelectHook)
	if err = o.doAfterSelectHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterSelectHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.BeforeUpdateHook, {{$alias.DownSingular}}BeforeUpdateHook)
	if err = o.doBeforeUpdateHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}BeforeUpdateHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterUpdateHook, {{$alias.DownSingular}}AfterUpdateHook)
	if err = o.doAfterUpdateHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterUpdateHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.BeforeDeleteHook, {{$alias.DownSingular}}BeforeDeleteHook)
	if err = o.doBeforeDeleteHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}BeforeDeleteHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterDeleteHook, {{$alias.DownSingular}}AfterDeleteHook)
	if err = o.doAfterDeleteHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterDeleteHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.BeforeUpsertHook, {{$alias.DownSingular}}BeforeUpsertHook)
	if err = o.doBeforeUpsertHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}BeforeUpsertHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterUpsertHook, {{$alias.DownSingular}}AfterUpsertHook)
	if err = o.doAfterUpsertHooks({{if not .NoContext}}ctx, {{end -}} nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterUpsertHooks = []{{$alias.UpSingular}}Hook{}
}
{{- end}}
