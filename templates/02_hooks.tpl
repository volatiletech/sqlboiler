{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
var {{$varNameSingular}}BeforeCreateHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}BeforeUpdateHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterCreateHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterUpdateHooks []{{$tableNameSingular}}Hook

// doBeforeCreateHooks executes all "before create" hooks.
func (o *{{$tableNameSingular}}) doBeforeCreateHooks() (err error) {
  for _, hook := range {{$varNameSingular}}BeforeCreateHooks {
    if err := hook(o); err != nil {
      return err
    }
  }

  return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *{{$tableNameSingular}}) doBeforeUpdateHooks() (err error) {
  for _, hook := range {{$varNameSingular}}BeforeUpdateHooks {
    if err := hook(o); err != nil {
      return err
    }
  }

  return nil
}

// doAfterCreateHooks executes all "after create" hooks.
func (o *{{$tableNameSingular}}) doAfterCreateHooks() (err error) {
  for _, hook := range {{$varNameSingular}}AfterCreateHooks {
    if err := hook(o); err != nil {
      return err
    }
  }

  return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *{{$tableNameSingular}}) doAfterUpdateHooks() (err error) {
  for _, hook := range {{$varNameSingular}}AfterUpdateHooks {
    if err := hook(o); err != nil {
      return err
    }
  }

  return nil
}

func {{$tableNameSingular}}AddHook(hookPoint boil.HookPoint, {{$varNameSingular}}Hook {{$tableNameSingular}}Hook) {
  switch hookPoint {
    case boil.HookBeforeCreate:
      {{$varNameSingular}}BeforeCreateHooks = append({{$varNameSingular}}BeforeCreateHooks, {{$varNameSingular}}Hook)
    case boil.HookBeforeUpdate:
      {{$varNameSingular}}BeforeUpdateHooks = append({{$varNameSingular}}BeforeUpdateHooks, {{$varNameSingular}}Hook)
    case boil.HookAfterCreate:
      {{$varNameSingular}}AfterCreateHooks = append({{$varNameSingular}}AfterCreateHooks, {{$varNameSingular}}Hook)
    case boil.HookAfterUpdate:
      {{$varNameSingular}}AfterUpdateHooks = append({{$varNameSingular}}AfterUpdateHooks, {{$varNameSingular}}Hook)
  }
}
