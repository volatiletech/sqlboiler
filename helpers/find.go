package helpers

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v5/boil"
	"github.com/volatiletech/sqlboiler/v5/queries"
)

func Find[U any, T Table[U], H TableSelectHooks[U], PK queryArgs](ctx context.Context, exec boil.ContextExecutor, pk PK, selectCols ...string) (U, error) {
	var table T
	var hooks H

	result := table.New()
	info := table.TableInfo()
	dialect := info.Dialect

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(dialect.IdentQuoteSlice(selectCols), ",")
	}

	var start int
	if dialect.UseIndexPlaceholders {
		start = 1
	}

	whereClause := dialect.WhereClause(start, info.PrimaryKeyColumns)

	var softDeleteClause string
	if info.DeletionColumnName != "" {
		softDeleteClause = fmt.Sprintf("and %s is null", info.DeletionColumnName)
	}

	query := fmt.Sprintf(
		"select %s from %s where %s %s",
		sel, info.Name, whereClause, softDeleteClause,
	)

	err := queries.Raw(query, pk.Values()...).Bind(ctx, exec, result)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result, sql.ErrNoRows
		}
		return result, errors.Wrapf(err, "unable to select from %s", info.Name)
	}

	if err := DoHooks(ctx, exec, result, hooks.AfterSelectHooks()); err != nil {
		return result, err
	}

	return result, nil
}
