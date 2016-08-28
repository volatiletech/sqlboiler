{{- if eq .NoHooks false -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $tableNamePlural := .Table.Name | plural | titleCase -}}
{{- $varNamePlural := .Table.Name | plural | camelCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
func {{$varNameSingular}}BeforeInsertHook(o *{{$tableNameSingular}}) error {
  *o = {{$tableNameSingular}}{}
  return nil
}

func {{$varNameSingular}}AfterInsertHook(o *{{$tableNameSingular}}) error {
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

  {{$tableNameSingular}}AddHook(boil.HookBeforeInsert, {{$varNameSingular}}BeforeInsertHook)
  if err = o.doBeforeInsertHooks(); err != nil {
    t.Errorf("Unable to execute doBeforeInsertHooks: %s", err)
  }
  if !reflect.DeepEqual(o, empty) {
    t.Errorf("Expected BeforeInsertHook function to empty object, but got: %#v", o)
  }

  {{$varNameSingular}}BeforeInsertHooks = []{{$tableNameSingular}}Hook{}
}
{{- end}}
