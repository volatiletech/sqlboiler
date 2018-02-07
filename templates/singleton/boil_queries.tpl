var dialect = queries.Dialect{
	LQ: 0x{{printf "%x" .Dialect.LQ}},
	RQ: 0x{{printf "%x" .Dialect.RQ}},
	IndexPlaceholders: {{.Dialect.IndexPlaceholders}},
	UseTopClause: {{.Dialect.UseTopClause}},
}

// NewQueryG initializes a new Query using the passed in QueryMods
func NewQueryG(mods ...qm.QueryMod) *queries.Query {
	return NewQuery(boil.GetDB(), mods...)
}

// NewQuery initializes a new Query using the passed in QueryMods
func NewQuery(exec boil.Executor, mods ...qm.QueryMod) *queries.Query {
	q := &queries.Query{}
	queries.SetExecutor(q, exec)
	queries.SetDialect(q, &dialect)
	qm.Apply(q, mods...)

	return q
}

func mergeModels(tx boil.Executor, primaryID uint64, secondaryID uint64, foreignKeys []foreignKey, conflictingKeys []conflictingUniqueKey) error {
	if len(foreignKeys) < 1 {
		return nil
	}
	var err error

	for _, conflict := range conflictingKeys {
		err = deleteConflictsBeforeMerge(tx, conflict, primaryID, secondaryID)
		if err != nil {
			return err
		}
	}

	for _, fk := range foreignKeys {
		// TODO: use NewQuery here, not plain sql
		query := fmt.Sprintf(
			"UPDATE %s SET %s = %s WHERE %s = %s",
			fk.foreignTable, fk.foreignColumn, strmangle.Placeholders(dialect.IndexPlaceholders, 1, 1, 1),
			fk.foreignColumn, strmangle.Placeholders(dialect.IndexPlaceholders, 1, 2, 1),
		)
		_, err = tx.Exec(query, primaryID, secondaryID)
		if err != nil {
			return errors.Err(err)
		}
	}
	return checkMerge(tx, foreignKeys)
}

func deleteConflictsBeforeMerge(tx boil.Executor, conflict conflictingUniqueKey, primaryID uint64, secondaryID uint64) error {
	conflictingColumns := strmangle.SetComplement(conflict.columns, []string{conflict.objectIdColumn})

	if len(conflictingColumns) < 1 {
		return nil
	} else if len(conflictingColumns) > 1 {
		return errors.Err("this doesnt work for unique keys with more than two columns (yet)")
	}

	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s IN (%s) GROUP BY %s HAVING count(distinct %s) > 1",
		conflictingColumns[0], conflict.table, conflict.objectIdColumn,
		strmangle.Placeholders(dialect.IndexPlaceholders, 2, 1, 1),
		conflictingColumns[0], conflict.objectIdColumn,
	)

	rows, err := tx.Query(query, primaryID, secondaryID)
	defer rows.Close()
	if err != nil {
		return errors.Err(err)
	}

	args := []interface{}{secondaryID}
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			return errors.Err(err)
		}
		args = append(args, value)
	}

	// if no rows found, no need to delete anything
	if len(args) < 2 {
		return nil
	}

	query = fmt.Sprintf(
		"DELETE FROM %s WHERE %s = %s AND %s IN (%s)",
		conflict.table, conflict.objectIdColumn, strmangle.Placeholders(dialect.IndexPlaceholders, 1, 1, 1),
		conflictingColumns[0], strmangle.Placeholders(dialect.IndexPlaceholders, len(args)-1, 2, 1),
	)

	_, err = tx.Exec(query, args...)
	if err != nil {
		return errors.Err(err)
	}
	return nil
}

func checkMerge(tx boil.Executor, foreignKeys []foreignKey) error {
	uniqueColumns := []interface{}{}
	uniqueColumnNames := map[string]bool{}
	handledTablesColumns := map[string]bool{}

	for _, fk := range foreignKeys {
		handledTablesColumns[fk.foreignTable+"."+fk.foreignColumn] = true
		if _, ok := uniqueColumnNames[fk.foreignColumn]; !ok {
			uniqueColumns = append(uniqueColumns, fk.foreignColumn)
			uniqueColumnNames[fk.foreignColumn] = true
		}
	}

	q := fmt.Sprintf(
		`SELECT table_name, column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND column_name IN (%s)`,
		strmangle.Placeholders(dialect.IndexPlaceholders, len(uniqueColumns), 1, 1),
	)
	rows, err := tx.Query(q, uniqueColumns...)
	defer rows.Close()
	if err != nil {
		return errors.Err(err)
	}

	for rows.Next() {
		var tableName string
		var columnName string
		err = rows.Scan(&tableName, &columnName)
		if err != nil {
			return errors.Err(err)
		}

		if _, exists := handledTablesColumns[tableName+"."+columnName]; !exists {
			return errors.Err("missing merge for " + tableName + "." + columnName)
		}
	}

	return nil
}
