package helpers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v5/boil"
	"github.com/volatiletech/sqlboiler/v5/queries"
)

type slice[U any] interface {
	~[]U
}

// For views... no methods attached
type ViewSlice[U any] []U

// For tables... with generic methods
type TableSlice[U any, T Table[U], H TableDeleteHooks[U]] []U

func (o TableSlice[U, T, H]) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("update all requires at least one column argument")
	}

	var table T
	info := table.TableInfo()
	dialect := info.Dialect

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	var paramStart, whereStart int
	if dialect.UseIndexPlaceholders {
		paramStart = 1
		whereStart = len(colNames) + 1
	}
	paramClause := dialect.SetParamNames(paramStart, colNames)
	whereClause := dialect.WhereClauseRepeated(whereStart, info.PrimaryKeyColumns, len(o))

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), info.PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s WHERE %s", info.Name, paramClause, whereClause)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to update all in %s slice", table.TableInfo().Name)
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrapf(err, "unable to retrieve rows affected all in update all %s", table.TableInfo().Name)
	}

	return rowsAff, nil
}

func (q TableSlice[U, T, H]) UpdateAllP(ctx context.Context, exec boil.ContextExecutor, cols M) int64 {
	e, err := q.UpdateAll(ctx, exec, cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

func (q TableSlice[U, T, H]) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

func (q TableSlice[U, T, H]) UpdateAllGP(ctx context.Context, cols M) int64 {
	e, err := q.UpdateAll(ctx, boil.GetContextDB(), cols)
	if err != nil {
		panic(boil.WrapErr(err))
	}

	return e
}

func (o *TableSlice[U, T, H]) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	var table T

	info := table.TableInfo()
	dialect := info.Dialect
	slice := make(TableSlice[U, T, H], 0)

	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), info.PrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	var start int
	if dialect.UseIndexPlaceholders {
		start = 1
	}

	whereClause := dialect.WhereClauseRepeated(start, info.PrimaryKeyColumns, len(*o))

	var softDeleteClause string
	if info.DeletionColumnName != "" {
		softDeleteClause = fmt.Sprintf("and %s is null", info.DeletionColumnName)
	}

	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s %s", info.Name, whereClause, softDeleteClause)

	err := queries.Raw(sql, args...).Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrapf(err, "unable to reload all in %s slice", info.Name)
	}

	*o = slice

	return nil
}

func (q TableSlice[U, T, H]) ReloadAllP(ctx context.Context, exec boil.ContextExecutor) {
	err := q.ReloadAll(ctx, exec)
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

func (q TableSlice[U, T, H]) ReloadAllG(ctx context.Context) error {
	return q.ReloadAll(ctx, boil.GetContextDB())
}

func (q TableSlice[U, T, H]) ReloadAllGP(ctx context.Context) {
	err := q.ReloadAll(ctx, boil.GetContextDB())
	if err != nil {
		panic(boil.WrapErr(err))
	}
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TableSlice[U, T, H]) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var table T
	var hooks H

	if len(hooks.BeforeDeleteHooks()) != 0 {
		for _, obj := range o {
			if err := DoHooks(ctx, exec, obj, hooks.BeforeDeleteHooks()); err != nil {
				return 0, err
			}
		}
	}

	info := table.TableInfo()
	dialect := info.Dialect
	deletionCol := info.DeletionColumnName
	hardDelete := deletionCol == "" || boil.DoHardDelete(ctx)

	var (
		sql  string
		args []interface{}
	)
	if hardDelete {
		var start int
		if dialect.UseIndexPlaceholders {
			start = 1
		}

		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), info.PrimaryKeyMapping)
			args = append(args, pkeyArgs...)
		}

		sql = fmt.Sprintf("DELETE FROM %s WHERE %s",
			info.Name,
			dialect.WhereClauseRepeated(start, info.PrimaryKeyColumns, len(o)),
		)
	} else {
		var parmStart, whereStart int
		if dialect.UseIndexPlaceholders {
			parmStart = 1
			whereStart = 2
		}

		currTime := time.Now().In(boil.GetLocation())
		args = []interface{}{currTime}

		for _, obj := range o {
			pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), info.PrimaryKeyMapping)
			args = append(args, pkeyArgs...)
			table.SetAsSoftDeleted(obj, currTime)
		}

		sql = fmt.Sprintf("UPDATE %s SET %s WHERE %s", info.Name,
			dialect.SetParamNames(parmStart, []string{deletionCol}),
			dialect.WhereClauseRepeated(whereStart, info.PrimaryKeyColumns, len(o)),
		)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrapf(err, "unable to delete all from %s slice", table.TableInfo().Name)
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get rows affected by deleteall for %s slice", table.TableInfo().Name)
	}

	if len(hooks.AfterDeleteHooks()) != 0 {
		for _, obj := range o {
			if err := DoHooks(ctx, exec, obj, hooks.AfterDeleteHooks()); err != nil {
				return 0, err
			}
		}
	}

	return rowsAff, nil
}
