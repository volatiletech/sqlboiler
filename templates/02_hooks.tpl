{{- if eq .NoHooks false -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
var {{$varNameSingular}}BeforeInsertHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}BeforeUpdateHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}BeforeUpsertHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterInsertHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterUpdateHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterUpsertHooks []{{$tableNameSingular}}Hook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *{{$tableNameSingular}}) doBeforeInsertHooks() (err error) {
  for _, hook := range {{$varNameSingular}}BeforeInsertHooks {
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

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *{{$tableNameSingular}}) doBeforeUpsertHooks() (err error) {
  for _, hook := range {{$varNameSingular}}BeforeUpsertHooks {
    if err := hook(o); err != nil {
      return err
    }
  }

  return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *{{$tableNameSingular}}) doAfterInsertHooks() (err error) {
  for _, hook := range {{$varNameSingular}}AfterInsertHooks {
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

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *{{$tableNameSingular}}) doAfterUpsertHooks() (err error) {
  for _, hook := range {{$varNameSingular}}AfterUpsertHooks {
    if err := hook(o); err != nil {
      return err
    }
  }

  return nil
}

func {{$tableNameSingular}}AddHook(hookPoint boil.HookPoint, {{$varNameSingular}}Hook {{$tableNameSingular}}Hook) {
  switch hookPoint {
    case boil.HookBeforeInsert:
      {{$varNameSingular}}BeforeInsertHooks = append({{$varNameSingular}}BeforeInsertHooks, {{$varNameSingular}}Hook)
    case boil.HookBeforeUpdate:
      {{$varNameSingular}}BeforeUpdateHooks = append({{$varNameSingular}}BeforeUpdateHooks, {{$varNameSingular}}Hook)
    case boil.HookBeforeUpsert:
      {{$varNameSingular}}BeforeUpsertHooks = append({{$varNameSingular}}BeforeUpsertHooks, {{$varNameSingular}}Hook)
    case boil.HookAfterInsert:
      {{$varNameSingular}}AfterInsertHooks = append({{$varNameSingular}}AfterInsertHooks, {{$varNameSingular}}Hook)
    case boil.HookAfterUpdate:
      {{$varNameSingular}}AfterUpdateHooks = append({{$varNameSingular}}AfterUpdateHooks, {{$varNameSingular}}Hook)
    case boil.HookAfterUpsert:
      {{$varNameSingular}}AfterUpsertHooks = append({{$varNameSingular}}AfterUpsertHooks, {{$varNameSingular}}Hook)
  }
}
{{- end}}
