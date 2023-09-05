{{- if not .NoHooks -}}
{{- $alias := .Aliases.Table .Table.Name}}

var {{$alias.DownSingular}}AfterSelectMu sync.Mutex
var {{$alias.DownSingular}}AfterSelectHooks []{{$alias.UpSingular}}Hook

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert) -}}
var {{$alias.DownSingular}}BeforeInsertMu sync.Mutex
var {{$alias.DownSingular}}BeforeInsertHooks []{{$alias.UpSingular}}Hook
var {{$alias.DownSingular}}AfterInsertMu sync.Mutex
var {{$alias.DownSingular}}AfterInsertHooks []{{$alias.UpSingular}}Hook
{{- end}}

{{if not .Table.IsView -}}
var {{$alias.DownSingular}}BeforeUpdateMu sync.Mutex
var {{$alias.DownSingular}}BeforeUpdateHooks []{{$alias.UpSingular}}Hook
var {{$alias.DownSingular}}AfterUpdateMu sync.Mutex
var {{$alias.DownSingular}}AfterUpdateHooks []{{$alias.UpSingular}}Hook

var {{$alias.DownSingular}}BeforeDeleteMu sync.Mutex
var {{$alias.DownSingular}}BeforeDeleteHooks []{{$alias.UpSingular}}Hook
var {{$alias.DownSingular}}AfterDeleteMu sync.Mutex
var {{$alias.DownSingular}}AfterDeleteHooks []{{$alias.UpSingular}}Hook
{{- end}}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanUpsert) -}}
var {{$alias.DownSingular}}BeforeUpsertMu sync.Mutex
var {{$alias.DownSingular}}BeforeUpsertHooks []{{$alias.UpSingular}}Hook
var {{$alias.DownSingular}}AfterUpsertMu sync.Mutex
var {{$alias.DownSingular}}AfterUpsertHooks []{{$alias.UpSingular}}Hook
{{- end}}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *{{$alias.UpSingular}}) doAfterSelectHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterSelectHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert) -}}
// doBeforeInsertHooks executes all "before insert" hooks.
func (o *{{$alias.UpSingular}}) doBeforeInsertHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}BeforeInsertHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *{{$alias.UpSingular}}) doAfterInsertHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterInsertHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}
{{- end}}

{{if not .Table.IsView -}}
// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *{{$alias.UpSingular}}) doBeforeUpdateHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}BeforeUpdateHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *{{$alias.UpSingular}}) doAfterUpdateHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterUpdateHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *{{$alias.UpSingular}}) doBeforeDeleteHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}BeforeDeleteHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *{{$alias.UpSingular}}) doAfterDeleteHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterDeleteHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}
{{- end}}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanUpsert) -}}
// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *{{$alias.UpSingular}}) doBeforeUpsertHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}BeforeUpsertHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *{{$alias.UpSingular}}) doAfterUpsertHooks({{if .NoContext}}exec boil.Executor{{else}}ctx context.Context, exec boil.ContextExecutor{{end}}) (err error) {
	{{if not .NoContext -}}
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	{{end -}}
	for _, hook := range {{$alias.DownSingular}}AfterUpsertHooks {
		if err := hook({{if not .NoContext}}ctx, {{end -}} exec, o); err != nil {
			return err
		}
	}

	return nil
}
{{- end}}

// Add{{$alias.UpSingular}}Hook registers your hook function for all future operations.
func Add{{$alias.UpSingular}}Hook(hookPoint boil.HookPoint, {{$alias.DownSingular}}Hook {{$alias.UpSingular}}Hook) {
	switch hookPoint {
		case boil.AfterSelectHook:
			{{$alias.DownSingular}}AfterSelectMu.Lock()
			{{$alias.DownSingular}}AfterSelectHooks = append({{$alias.DownSingular}}AfterSelectHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}AfterSelectMu.Unlock()
		{{- if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert)}}
		case boil.BeforeInsertHook:
			{{$alias.DownSingular}}BeforeInsertMu.Lock()
			{{$alias.DownSingular}}BeforeInsertHooks = append({{$alias.DownSingular}}BeforeInsertHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}BeforeInsertMu.Unlock()
		case boil.AfterInsertHook:
			{{$alias.DownSingular}}AfterInsertMu.Lock()
			{{$alias.DownSingular}}AfterInsertHooks = append({{$alias.DownSingular}}AfterInsertHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}AfterInsertMu.Unlock()
		{{- end}}
		{{- if not .Table.IsView}}
		case boil.BeforeUpdateHook:
			{{$alias.DownSingular}}BeforeUpdateMu.Lock()
			{{$alias.DownSingular}}BeforeUpdateHooks = append({{$alias.DownSingular}}BeforeUpdateHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}BeforeUpdateMu.Unlock()
		case boil.AfterUpdateHook:
			{{$alias.DownSingular}}AfterUpdateMu.Lock()
			{{$alias.DownSingular}}AfterUpdateHooks = append({{$alias.DownSingular}}AfterUpdateHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}AfterUpdateMu.Unlock()
		case boil.BeforeDeleteHook:
			{{$alias.DownSingular}}BeforeDeleteMu.Lock()
			{{$alias.DownSingular}}BeforeDeleteHooks = append({{$alias.DownSingular}}BeforeDeleteHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}BeforeDeleteMu.Unlock()
		case boil.AfterDeleteHook:
			{{$alias.DownSingular}}AfterDeleteMu.Lock()
			{{$alias.DownSingular}}AfterDeleteHooks = append({{$alias.DownSingular}}AfterDeleteHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}AfterDeleteMu.Unlock()
		{{- end}}
		{{- if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert)}}
		case boil.BeforeUpsertHook:
			{{$alias.DownSingular}}BeforeUpsertMu.Lock()
			{{$alias.DownSingular}}BeforeUpsertHooks = append({{$alias.DownSingular}}BeforeUpsertHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}BeforeUpsertMu.Unlock()
		case boil.AfterUpsertHook:
			{{$alias.DownSingular}}AfterUpsertMu.Lock()
			{{$alias.DownSingular}}AfterUpsertHooks = append({{$alias.DownSingular}}AfterUpsertHooks, {{$alias.DownSingular}}Hook)
			{{$alias.DownSingular}}AfterUpsertMu.Unlock()
		{{- end}}
	}
}
{{- end}}
