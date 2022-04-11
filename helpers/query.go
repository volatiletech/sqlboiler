package helpers

import (
	"context"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v5/boil"
	"github.com/volatiletech/sqlboiler/v5/queries"
)

type BaseQuery[U any, T Table[U]] struct {
	*queries.Query
}

// Count returns the count of all records in the query.
func (q BaseQuery[U, T]) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var table T
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to count %s rows", table.TableInfo().Name)
	}

	return count, nil
}

// CountP checks if the row exists in the table, and panics on error.
func (q BaseQuery[U, T]) CountP(ctx context.Context, exec boil.ContextExecutor) int64 {
	e, err := q.Count(ctx, exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// CountG checks if the row exists in the table using the global executor.
func (q BaseQuery[U, T]) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// CountGP checks if the row exists in the table using the global executor, and panics on error.
func (q BaseQuery[U, T]) CountGP(ctx context.Context) int64 {
	e, err := q.Count(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// Exists checks if the row exists in the table.
func (q BaseQuery[U, T]) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var table T
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrapf(err, "failed to check if %s exists", table.TableInfo().Name)
	}

	return count > 0, nil
}

// ExistsP checks if the row exists in the table, and panics on error.
func (q BaseQuery[U, T]) ExistsP(ctx context.Context, exec boil.ContextExecutor) bool {
	e, err := q.Exists(ctx, exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// ExistsG checks if the row exists in the table using the global executor.
func (q BaseQuery[U, T]) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// ExistsGP checks if the row exists in the table using the global executor, and panics on error.
func (q BaseQuery[U, T]) ExistsGP(ctx context.Context) bool {
	e, err := q.Exists(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

type M map[string]interface{}

// UpdateAll updates all rows with the specified column values.
func (q BaseQuery[U, T]) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	var table T
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)

	if err != nil {
		return 0, errors.Wrapf(err, "unable to update all for %s", table.TableInfo().Name)
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrapf(err, "unable to retrieve rows affected for %s", table.TableInfo().Name)
	}

	return rowsAff, nil
}

// UpdateAllP checks if the row exists in the table, and panics on error.
func (q BaseQuery[U, T]) UpdateAllP(ctx context.Context, exec boil.ContextExecutor, cols M) int64 {
	e, err := q.UpdateAll(ctx, exec, cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// UpdateAllG checks if the row exists in the table using the global executor.
func (q BaseQuery[U, T]) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAllGP checks if the row exists in the table using the global executor, and panics on error.
func (q BaseQuery[U, T]) UpdateAllGP(ctx context.Context, cols M) int64 {
	e, err := q.UpdateAll(ctx, boil.GetContextDB(), cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
