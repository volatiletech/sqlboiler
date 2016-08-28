{{- if eq .NoHooks false -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func {{$varNameSingular}}BeforeCreateHook(o *{{$tableNameSingular}}) error {
  *o = {{$tableNameSingular}}{}
  return nil
}

func {{$varNameSingular}}AfterCreateHook(o *{{$tableNameSingular}}) error {
  *o = {{$tableNameSingular}}{}
  return nil
}

func {{$varNameSingular}}BeforeUpdateHook(o *{{$tableNameSingular}}) error {
  *o = {{$tableNameSingular}}{}
  return nil
}

func {{$varNameSingular}}AfterUpdateHook(o *{{$tableNameSingular}}) error {
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

  {{$tableNameSingular}}AddHook(boil.HookBeforeCreate, {{$varNameSingular}}BeforeCreateHook)
  if err = o.doBeforeCreateHooks(); err != nil {
    t.Errorf("Unable to execute doBeforeCreateHooks: %s", err)
  }
  if !reflect.DeepEqual(o, empty) {
    t.Errorf("Expected BeforeCreateHook function to empty object, but got: %#v", o)
  }

  {{$varNameSingular}}BeforeCreateHooks = []{{$tableNameSingular}}Hook{}
}
{{- end}}
