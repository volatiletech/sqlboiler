{{- if not .NoHooks -}}
{{- $tableNameSingular := .Table.Name | singular | titleCase -}}
{{- $varNameSingular := .Table.Name | singular | camelCase -}}
var {{$varNameSingular}}BeforeInsertHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}BeforeUpdateHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}BeforeDeleteHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}BeforeUpsertHooks []{{$tableNameSingular}}Hook

var {{$varNameSingular}}AfterInsertHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterSelectHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterUpdateHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterDeleteHooks []{{$tableNameSingular}}Hook
var {{$varNameSingular}}AfterUpsertHooks []{{$tableNameSingular}}Hook

// doBeforeInsertHooks executes all "before insert" hooks.
func (o *{{$tableNameSingular}}) doBeforeInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}BeforeInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *{{$tableNameSingular}}) doBeforeUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}BeforeUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *{{$tableNameSingular}}) doBeforeDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}BeforeDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *{{$tableNameSingular}}) doBeforeUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}BeforeUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *{{$tableNameSingular}}) doAfterInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}AfterInsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *{{$tableNameSingular}}) doAfterSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}AfterSelectHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *{{$tableNameSingular}}) doAfterUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}AfterUpdateHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *{{$tableNameSingular}}) doAfterDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}AfterDeleteHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *{{$tableNameSingular}}) doAfterUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}AfterUpsertHooks {
		if err := hook(exec, o); err != nil {
			return err
		}
	}

	return nil
}

// Add{{$tableNameSingular}}Hook registers your hook function for all future operations.
func Add{{$tableNameSingular}}Hook(hookPoint boil.HookPoint, {{$varNameSingular}}Hook {{$tableNameSingular}}Hook) {
	switch hookPoint {
		case boil.BeforeInsertHook:
			{{$varNameSingular}}BeforeInsertHooks = append({{$varNameSingular}}BeforeInsertHooks, {{$varNameSingular}}Hook)
		case boil.BeforeUpdateHook:
			{{$varNameSingular}}BeforeUpdateHooks = append({{$varNameSingular}}BeforeUpdateHooks, {{$varNameSingular}}Hook)
		case boil.BeforeDeleteHook:
			{{$varNameSingular}}BeforeDeleteHooks = append({{$varNameSingular}}BeforeDeleteHooks, {{$varNameSingular}}Hook)
		case boil.BeforeUpsertHook:
			{{$varNameSingular}}BeforeUpsertHooks = append({{$varNameSingular}}BeforeUpsertHooks, {{$varNameSingular}}Hook)
		case boil.AfterInsertHook:
			{{$varNameSingular}}AfterInsertHooks = append({{$varNameSingular}}AfterInsertHooks, {{$varNameSingular}}Hook)
		case boil.AfterSelectHook:
			{{$varNameSingular}}AfterSelectHooks = append({{$varNameSingular}}AfterSelectHooks, {{$varNameSingular}}Hook)
		case boil.AfterUpdateHook:
			{{$varNameSingular}}AfterUpdateHooks = append({{$varNameSingular}}AfterUpdateHooks, {{$varNameSingular}}Hook)
		case boil.AfterDeleteHook:
			{{$varNameSingular}}AfterDeleteHooks = append({{$varNameSingular}}AfterDeleteHooks, {{$varNameSingular}}Hook)
		case boil.AfterUpsertHook:
			{{$varNameSingular}}AfterUpsertHooks = append({{$varNameSingular}}AfterUpsertHooks, {{$varNameSingular}}Hook)
	}
}

var {{$varNameSingular}}QueryInsertHooks []{{$tableNameSingular}}QueryHook
var {{$varNameSingular}}QuerySelectHooks []{{$tableNameSingular}}QueryHook
var {{$varNameSingular}}QueryUpdateHooks []{{$tableNameSingular}}QueryHook
var {{$varNameSingular}}QueryDeleteHooks []{{$tableNameSingular}}QueryHook
var {{$varNameSingular}}QueryUpsertHooks []{{$tableNameSingular}}QueryHook

// doInsertHooks executes all "insert" hooks.
func (q {{$varNameSingular}}Query) doInsertHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}QueryInsertHooks {
		if err := hook(exec, q.Query); err != nil {
			return err
		}
	}

	return nil
}

// doSelectHooks executes all "select" hooks.
func (q {{$varNameSingular}}Query) doSelectHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}QuerySelectHooks {
		if err := hook(exec, q.Query); err != nil {
			return err
		}
	}

	return nil
}

// doUpdateHooks executes all "update" hooks.
func (q {{$varNameSingular}}Query) doUpdateHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}QueryUpdateHooks {
		if err := hook(exec, q.Query); err != nil {
			return err
		}
	}

	return nil
}

// doDeleteHooks executes all "delete" hooks.
func (q {{$varNameSingular}}Query) doDeleteHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}QueryDeleteHooks {
		if err := hook(exec, q.Query); err != nil {
			return err
		}
	}

	return nil
}

// doUpsertHooks executes all "upsert" hooks.
func (q {{$varNameSingular}}Query) doUpsertHooks(exec boil.Executor) (err error) {
	for _, hook := range {{$varNameSingular}}QueryUpsertHooks {
		if err := hook(exec, q.Query); err != nil {
			return err
		}
	}

	return nil
}

// Add{{$tableNameSingular}}QueryHook registers your query hook function for all future operations.
func Add{{$tableNameSingular}}QueryHook(hookPoint boil.QueryHookPoint, {{$varNameSingular}}QueryHook {{$tableNameSingular}}QueryHook) {
	switch hookPoint {
		case boil.InsertHook:
			{{$varNameSingular}}QueryInsertHooks = append({{$varNameSingular}}QueryInsertHooks, {{$varNameSingular}}QueryHook)
		case boil.SelectHook:
			{{$varNameSingular}}QuerySelectHooks = append({{$varNameSingular}}QuerySelectHooks, {{$varNameSingular}}QueryHook)
		case boil.UpdateHook:
			{{$varNameSingular}}QueryUpdateHooks = append({{$varNameSingular}}QueryUpdateHooks, {{$varNameSingular}}QueryHook)
		case boil.DeleteHook:
			{{$varNameSingular}}QueryDeleteHooks = append({{$varNameSingular}}QueryDeleteHooks, {{$varNameSingular}}QueryHook)
		case boil.UpsertHook:
			{{$varNameSingular}}QueryUpsertHooks = append({{$varNameSingular}}QueryUpsertHooks, {{$varNameSingular}}QueryHook)
	}
}
{{- end}}
