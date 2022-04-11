{{- if not .NoHooks -}}
{{- $alias := .Aliases.Table .Table.Name}}

var {{$alias.DownSingular}}AfterSelectHooks helpers.TableHooks[*{{$alias.UpSingular}}]

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert) -}}
var {{$alias.DownSingular}}BeforeInsertHooks helpers.TableHooks[*{{$alias.UpSingular}}]
var {{$alias.DownSingular}}AfterInsertHooks helpers.TableHooks[*{{$alias.UpSingular}}]
{{- end}}

{{if not .Table.IsView -}}
var {{$alias.DownSingular}}BeforeUpdateHooks helpers.TableHooks[*{{$alias.UpSingular}}]
var {{$alias.DownSingular}}AfterUpdateHooks helpers.TableHooks[*{{$alias.UpSingular}}]

var {{$alias.DownSingular}}BeforeDeleteHooks helpers.TableHooks[*{{$alias.UpSingular}}]
var {{$alias.DownSingular}}AfterDeleteHooks helpers.TableHooks[*{{$alias.UpSingular}}]
{{- end}}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanUpsert) -}}
var {{$alias.DownSingular}}BeforeUpsertHooks helpers.TableHooks[*{{$alias.UpSingular}}]
var {{$alias.DownSingular}}AfterUpsertHooks helpers.TableHooks[*{{$alias.UpSingular}}]
{{- end}}

// AfterSelectHooks returns all registered "after Select" hooks for the model
func (t {{$alias.DownSingular}}Hooks) AfterSelectHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}AfterSelectHooks
}

// doAfterSelectHooks executes all "after Select" hooks.
func (o *{{$alias.UpSingular}}) doAfterSelectHooks(ctx context.Context, exec boil.ContextExecutor) error {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}AfterSelectHooks)
}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert) -}}
// BeforeInsertHooks returns all registered "before Insert" hooks for the model
func (t {{$alias.DownSingular}}Hooks) BeforeInsertHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}BeforeInsertHooks
}

// doBeforeInsertHooks executes all "before Insert" hooks.
func (o *{{$alias.UpSingular}}) doBeforeInsertHooks(ctx context.Context, exec boil.ContextExecutor) error {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}BeforeInsertHooks)
}

// AfterInsertHooks returns all registered "after Insert" hooks for the model
func (t {{$alias.DownSingular}}Hooks) AfterInsertHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}AfterInsertHooks
}

// doAfterInsertHooks executes all "after Insert" hooks.
func (o *{{$alias.UpSingular}}) doAfterInsertHooks(ctx context.Context, exec boil.ContextExecutor) error {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}AfterInsertHooks)
}
{{- end}}

{{if not .Table.IsView -}}
// BeforeUpdateHooks returns all registered "before Update" hooks for the model
func (t {{$alias.DownSingular}}Hooks) BeforeUpdateHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}BeforeUpdateHooks
}

// doBeforeUpdateHooks executes all "before Update" hooks.
func (o *{{$alias.UpSingular}}) doBeforeUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}BeforeUpdateHooks)
}

// AfterUpdateHooks returns all registered after Update hooks for the model
func (t {{$alias.DownSingular}}Hooks) AfterUpdateHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}AfterUpdateHooks
}

// doAfterUpdateHooks executes all "after Update" hooks.
func (o *{{$alias.UpSingular}}) doAfterUpdateHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}AfterUpdateHooks)
}

// BeforeDeleteHooks returns all registered "before Delete" hooks for the model
func (t {{$alias.DownSingular}}Hooks) BeforeDeleteHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}BeforeDeleteHooks
}

// doBeforeDeleteHooks executes all "before Delete" hooks.
func (o *{{$alias.UpSingular}}) doBeforeDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}BeforeDeleteHooks)
}

// AfterDeleteHooks returns all registered after Delete hooks for the model
func (t {{$alias.DownSingular}}Hooks) AfterDeleteHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}AfterDeleteHooks
}

// doAfterDeleteHooks executes all "after Delete" hooks.
func (o *{{$alias.UpSingular}}) doAfterDeleteHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}AfterDeleteHooks)
}
{{- end}}

{{if or (not .Table.IsView) (.Table.ViewCapabilities.CanUpsert) -}}
// BeforeUpsertHooks returns all registered "before Upsert" hooks for the model
func (t {{$alias.DownSingular}}Hooks) BeforeUpsertHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}BeforeUpsertHooks
}

// doBeforeUpsertHooks executes all "before Upsert" hooks.
func (o *{{$alias.UpSingular}}) doBeforeUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}BeforeUpsertHooks)
}

// AfterUpsertHooks returns all registered after Upsert hooks for the model
func (t {{$alias.DownSingular}}Hooks) AfterUpsertHooks() helpers.TableHooks[*{{$alias.UpSingular}}] {
	return {{$alias.DownSingular}}AfterUpsertHooks
}

// doAfterUpsertHooks executes all "after Upsert" hooks.
func (o *{{$alias.UpSingular}}) doAfterUpsertHooks(ctx context.Context, exec boil.ContextExecutor) (err error) {
	return helpers.DoHooks(ctx, exec, o, {{$alias.DownSingular}}AfterUpsertHooks)
}
{{- end}}

// Add{{$alias.UpSingular}}Hook registers your hook function for all future operations.
func Add{{$alias.UpSingular}}Hook(hookPoint boil.HookPoint, {{$alias.DownSingular}}Hook {{$alias.UpSingular}}Hook) {
	switch hookPoint {
		case boil.AfterSelectHook:
			{{$alias.DownSingular}}AfterSelectHooks = append({{$alias.DownSingular}}AfterSelectHooks, {{$alias.DownSingular}}Hook)
		{{- if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert)}}
		case boil.BeforeInsertHook:
			{{$alias.DownSingular}}BeforeInsertHooks = append({{$alias.DownSingular}}BeforeInsertHooks, {{$alias.DownSingular}}Hook)
		case boil.AfterInsertHook:
			{{$alias.DownSingular}}AfterInsertHooks = append({{$alias.DownSingular}}AfterInsertHooks, {{$alias.DownSingular}}Hook)
		{{- end}}
		{{- if not .Table.IsView}}
		case boil.BeforeUpdateHook:
			{{$alias.DownSingular}}BeforeUpdateHooks = append({{$alias.DownSingular}}BeforeUpdateHooks, {{$alias.DownSingular}}Hook)
		case boil.AfterUpdateHook:
			{{$alias.DownSingular}}AfterUpdateHooks = append({{$alias.DownSingular}}AfterUpdateHooks, {{$alias.DownSingular}}Hook)
		case boil.BeforeDeleteHook:
			{{$alias.DownSingular}}BeforeDeleteHooks = append({{$alias.DownSingular}}BeforeDeleteHooks, {{$alias.DownSingular}}Hook)
		case boil.AfterDeleteHook:
			{{$alias.DownSingular}}AfterDeleteHooks = append({{$alias.DownSingular}}AfterDeleteHooks, {{$alias.DownSingular}}Hook)
		{{- end}}
		{{- if or (not .Table.IsView) (.Table.ViewCapabilities.CanInsert)}}
		case boil.BeforeUpsertHook:
			{{$alias.DownSingular}}BeforeUpsertHooks = append({{$alias.DownSingular}}BeforeUpsertHooks, {{$alias.DownSingular}}Hook)
		case boil.AfterUpsertHook:
			{{$alias.DownSingular}}AfterUpsertHooks = append({{$alias.DownSingular}}AfterUpsertHooks, {{$alias.DownSingular}}Hook)
		{{- end}}
	}
}
{{- end}}
