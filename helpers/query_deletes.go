package helpers

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v5/boil"
	"github.com/volatiletech/sqlboiler/v5/queries"
)

type DeleteQuery[U any, T Table[U]] struct {
	*queries.Query
}

func (DeleteQuery[U, T]) New(q *queries.Query) DeleteQuery[U, T] {
	return DeleteQuery[U, T]{
		Query: q,
	}
}

// DeleteAll deletes all matching rows.
func (q DeleteQuery[U, T]) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var table T

	if q.Query == nil {
		return 0, errors.Errorf("no query provided for %s delete all", table.TableInfo().Name)
	}

	deletionCol := table.TableInfo().DeletionColumnName

	if deletionCol == "" || boil.DoHardDelete(ctx) {
		queries.SetDelete(q.Query)
	} else {
		currTime := time.Now().In(boil.GetLocation())
		queries.SetUpdate(q.Query, M{deletionCol: currTime})
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to delete all from %s", table.TableInfo().Name)
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get rows affected by deleteall for %s", table.TableInfo().Name)
	}

	return rowsAff, nil
}

// DeleteAllP checks if the row exists in the table, and panics on error.
func (q DeleteQuery[U, T]) DeleteAllP(ctx context.Context, exec boil.ContextExecutor) int64 {
	e, err := q.DeleteAll(ctx, exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

// DeleteAllG checks if the row exists in the table using the global executor.
func (q DeleteQuery[U, T]) DeleteAllG(ctx context.Context) (int64, error) {
	return q.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAllGP checks if the row exists in the table using the global executor, and panics on error.
func (q DeleteQuery[U, T]) DeleteAllGP(ctx context.Context) int64 {
	e, err := q.DeleteAll(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}
