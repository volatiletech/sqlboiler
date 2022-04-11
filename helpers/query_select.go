package helpers

import (
	"context"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v5/boil"
	"github.com/volatiletech/sqlboiler/v5/queries"
)

type SelectQuery[U any, T Table[U], H TableSelectHooks[U], S slice[U]] struct {
	*queries.Query
}

func (SelectQuery[U, T, H, S]) New(q *queries.Query) SelectQuery[U, T, H, S] {
	return SelectQuery[U, T, H, S]{
		Query: q,
	}
}

// One returns a single record from the query.
func (q SelectQuery[U, T, H, S]) One(ctx context.Context, exec boil.ContextExecutor) (U, error) {
	var table T
	var hooks H
	var result = table.New()

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, result)
	if err != nil {
		return result, errors.Wrapf(err, "failed to execute a one query for %s", table.TableInfo().Name)
	}

	if err := DoHooks(ctx, exec, result, hooks.AfterSelectHooks()); err != nil {
		return result, err
	}

	return result, nil
}

// OneP checks if the row exists in the table, and panics on error.
func (q SelectQuery[U, T, H, S]) OneP(ctx context.Context, exec boil.ContextExecutor) U {
	e, err := q.One(ctx, exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// OneG checks if the row exists in the table using the global executor.
func (q SelectQuery[U, T, H, S]) OneG(ctx context.Context) (U, error) {
	return q.One(ctx, boil.GetContextDB())
}

// OneGP checks if the row exists in the table using the global executor, and panics on error.
func (q SelectQuery[U, T, H, S]) OneGP(ctx context.Context) U {
	e, err := q.One(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// All returns all records from the query.
func (q SelectQuery[U, T, H, S]) All(ctx context.Context, exec boil.ContextExecutor) (S, error) {
	var table T
	var hooks H
	var slice S
	var results = make(S, 0)

	err := q.Bind(ctx, exec, &results)
	if err != nil {
		return slice, errors.Wrapf(err, "failed to assign all query results to %s slice", table.TableInfo().Name)
	}

	if len(hooks.AfterSelectHooks()) != 0 {
		for _, result := range results {
			err := DoHooks(ctx, exec, result, hooks.AfterSelectHooks())
			if err != nil {
				return slice, err
			}
		}
	}

	return slice, nil
}

// AllP checks if the row exists in the table, and panics on error.
func (q SelectQuery[U, T, H, S]) AllP(ctx context.Context, exec boil.ContextExecutor) S {
	e, err := q.All(ctx, exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// AllG checks if the row exists in the table using the global executor.
func (q SelectQuery[U, T, H, S]) AllG(ctx context.Context) (S, error) {
	return q.All(ctx, boil.GetContextDB())
}

// AllGP checks if the row exists in the table using the global executor, and panics on error.
func (q SelectQuery[U, T, H, S]) AllGP(ctx context.Context) S {
	e, err := q.All(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
