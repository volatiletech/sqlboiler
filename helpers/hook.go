package helpers

import (
	"context"

	"github.com/volatiletech/sqlboiler/v5/boil"
)

// TableHook is the signature for custom table hook methods
type TableHook[T any] func(context.Context, boil.ContextExecutor, T) error
type TableHooks[T any] []TableHook[T]

func DoHooks[T any](ctx context.Context, exec boil.ContextExecutor, o T, hooks TableHooks[T]) error {
	if boil.HooksAreSkipped(ctx) {
		return nil
	}

	for _, hook := range hooks {
		if err := hook(ctx, exec, o); err != nil {
			return err
		}
	}

	return nil
}

type TableSelectHooks[T any] interface {
	AfterSelectHooks() TableHooks[T]
}

type TableInsertHooks[T any] interface {
	BeforeInsertHooks() TableHooks[T]
	AfterInsertHooks() TableHooks[T]
}

type TableUpdateHooks[T any] interface {
	BeforeUpdateHooks() TableHooks[T]
	AfterUpdateHooks() TableHooks[T]
}

type TableDeleteHooks[T any] interface {
	BeforeDeleteHooks() TableHooks[T]
	AfterDeleteHooks() TableHooks[T]
}

type TableUpsertHooks[T any] interface {
	BeforeUpsertHooks() TableHooks[T]
	AfterUpsertHooks() TableHooks[T]
}

// NoOpHooks return nil for every hook method.
// this is used when hooks are disabled
type NoOpHooks[T any] struct{}

func (NoOpHooks[T]) AfterSelectHooks() TableHooks[T] {
	return nil
}

func (NoOpHooks[T]) BeforeInsertHooks() TableHooks[T] {
	return nil
}

func (NoOpHooks[T]) AfterInsertHooks() TableHooks[T] {
	return nil
}

func (NoOpHooks[T]) BeforeUpdateHooks() TableHooks[T] {
	return nil
}

func (NoOpHooks[T]) AfterUpdateHooks() TableHooks[T] {
	return nil
}

func (NoOpHooks[T]) BeforeDeleteHooks() TableHooks[T] {
	return nil
}

func (NoOpHooks[T]) AfterDeleteHooks() TableHooks[T] {
	return nil
}

func (NoOpHooks[T]) BeforeUpsertHooks() TableHooks[T] {
	return nil
}

func (NoOpHooks[T]) AfterUpsertHooks() TableHooks[T] {
	return nil
}
