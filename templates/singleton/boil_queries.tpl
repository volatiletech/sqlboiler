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
        if len(conflict.columns) == 1 && conflict.columns[0] == conflict.objectIdColumn {
            err = deleteOneToOneConflictsBeforeMerge(tx, conflict, primaryID, secondaryID)
        } else {
            err = deleteOneToManyConflictsBeforeMerge(tx, conflict, primaryID, secondaryID)
        }
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

func deleteOneToOneConflictsBeforeMerge(tx boil.Executor, conflict conflictingUniqueKey, primaryID uint64, secondaryID uint64) error {
	query := fmt.Sprintf(
		"SELECT COUNT(*) FROM %s WHERE %s IN (%s)",
		conflict.table, conflict.objectIdColumn,
		strmangle.Placeholders(dialect.IndexPlaceholders, 2, 1, 1),
	)

	var count int
	err := tx.QueryRow(query, primaryID, secondaryID).Scan(&count)
	if err != nil {
		return errors.Err(err)
	}

	if count > 2 {
		return errors.Err("it should not be possible to have more than two rows here")
	} else if count != 2 {
		return nil // no conflicting rows
	}

	query = fmt.Sprintf(
		"DELETE FROM %s WHERE %s = %s",
		conflict.table, conflict.objectIdColumn, strmangle.Placeholders(dialect.IndexPlaceholders, 1, 1, 1),
	)

	_, err = tx.Exec(query, secondaryID)
	return errors.Err(err)
}

func deleteOneToManyConflictsBeforeMerge(tx boil.Executor, conflict conflictingUniqueKey, primaryID uint64, secondaryID uint64) error {
	conflictingColumns := strmangle.SetComplement(conflict.columns, []string{conflict.objectIdColumn})
    var objectIDIndex int
    for i, column := range conflict.columns {
        if column == conflict.objectIdColumn {
            objectIDIndex = i
        }
    }
	query := fmt.Sprintf(
		"SELECT %s FROM %s WHERE %s IN (%s) GROUP BY %s HAVING count(distinct %s) > 1",
		strings.Join(conflictingColumns, ","), conflict.table, conflict.objectIdColumn,
		strmangle.Placeholders(dialect.IndexPlaceholders, 2, 1, 1),
		strings.Join(conflictingColumns, ","), conflict.objectIdColumn,
	)

	//The selectParams should be the ObjectIDs to search for regarding the conflict.
	rows, err := tx.Query(query, primaryID, secondaryID)
	if err != nil {
		return errors.Err(err)
	}

	//Since we don't don't know if advance how many columns the query returns, we have dynamically assign them to be
	// used in the delete query.
	colNames, err := rows.Columns()
	if err != nil {
		return errors.Err(err)
	}
	//Each row result of the query needs to be removed for being a conflicting row. Store each row's keys in an array.
	var rowsToRemove = [][]interface{}(nil)
	for rows.Next() {
		//Set pointers for dynamic scan
		iColPtrs := make([]interface{}, len(colNames))
		for i := 0; i < len(colNames); i++ {
			s := string("")
			iColPtrs[i] = &s
		}
		//Dynamically scan n columns
		err = rows.Scan(iColPtrs...)
		if err != nil {
			return errors.Err(err)
		}
		//Grab scanned values for query arguments
		iCol := make([]interface{}, len(colNames))
		for i, col := range iColPtrs {
			x := col.(*string)
			iCol[i] = *x
		}
		rowsToRemove = append(rowsToRemove, iCol)
	}
	defer rows.Close()

	//This query will adjust dynamically depending on the number of conflicting keys, adding AND expressions for each
	// key to ensure the right conflicting rows are deleted.
	query = fmt.Sprintf(
		"DELETE FROM %s %s",
		conflict.table,
		"WHERE "+strings.Join(conflict.columns, " = ? AND ")+" = ?",
	)

	//There could be multiple conflicting rows between ObjectIDs. In the SELECT query we grab each row and their column
	// keys to be deleted here in a loop.
	for _, rowToDelete := range rowsToRemove {
		rowToDelete = insert(rowToDelete, objectIDIndex, secondaryID)
		_, err = tx.Exec(query, rowToDelete...)
		if err != nil {
			return errors.Err(err)
		}
	}
	return nil
}

func insert(slice []interface{}, index int, value interface{}) []interface{} {
	return append(slice[:index], append([]interface{}{value}, slice[index:]...)...)
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
