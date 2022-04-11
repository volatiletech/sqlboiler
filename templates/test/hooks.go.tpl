{{- $alias := .Aliases.Table .Table.Name}}
func {{$alias.DownSingular}}BeforeInsertHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterInsertHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterSelectHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}BeforeUpdateHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterUpdateHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}BeforeDeleteHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterDeleteHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}BeforeUpsertHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func {{$alias.DownSingular}}AfterUpsertHook(ctx context.Context, e boil.ContextExecutor, o *{{$alias.UpSingular}}) error {
	*o = {{$alias.UpSingular}}{}
	return nil
}

func test{{$alias.UpPlural}}Hooks(t *testing.T) {
	t.Parallel()

	var err error

	ctx := context.Background()
	empty := &{{$alias.UpSingular}}{}
	o := &{{$alias.UpSingular}}{}

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, o, {{$alias.DownSingular}}DBTypes, false); err != nil {
		t.Errorf("Unable to randomize {{$alias.UpSingular}} object: %s", err)
	}

	Add{{$alias.UpSingular}}Hook(boil.BeforeInsertHook, {{$alias.DownSingular}}BeforeInsertHook)
	if err = o.doBeforeInsertHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}BeforeInsertHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterInsertHook, {{$alias.DownSingular}}AfterInsertHook)
	if err = o.doAfterInsertHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doAfterInsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterInsertHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterInsertHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterSelectHook, {{$alias.DownSingular}}AfterSelectHook)
	if err = o.doAfterSelectHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doAfterSelectHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterSelectHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterSelectHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.BeforeUpdateHook, {{$alias.DownSingular}}BeforeUpdateHook)
	if err = o.doBeforeUpdateHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpdateHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}BeforeUpdateHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterUpdateHook, {{$alias.DownSingular}}AfterUpdateHook)
	if err = o.doAfterUpdateHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doAfterUpdateHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpdateHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterUpdateHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.BeforeDeleteHook, {{$alias.DownSingular}}BeforeDeleteHook)
	if err = o.doBeforeDeleteHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doBeforeDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeDeleteHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}BeforeDeleteHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterDeleteHook, {{$alias.DownSingular}}AfterDeleteHook)
	if err = o.doAfterDeleteHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doAfterDeleteHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterDeleteHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterDeleteHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.BeforeUpsertHook, {{$alias.DownSingular}}BeforeUpsertHook)
	if err = o.doBeforeUpsertHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doBeforeUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected BeforeUpsertHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}BeforeUpsertHooks = []{{$alias.UpSingular}}Hook{}

	Add{{$alias.UpSingular}}Hook(boil.AfterUpsertHook, {{$alias.DownSingular}}AfterUpsertHook)
	if err = o.doAfterUpsertHooks(ctx,  nil); err != nil {
		t.Errorf("Unable to execute doAfterUpsertHooks: %s", err)
	}
	if !reflect.DeepEqual(o, empty) {
		t.Errorf("Expected AfterUpsertHook function to empty object, but got: %#v", o)
	}
	{{$alias.DownSingular}}AfterUpsertHooks = []{{$alias.UpSingular}}Hook{}
}
