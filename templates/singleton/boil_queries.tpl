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

func mergeModels(tx *sql.Tx, primaryID uint64, secondaryID uint64, relatedFields map[string]string) error {
	if len(relatedFields) < 1 {
    return nil
  }

  for table, column := range relatedFields {
    // TODO: use NewQuery here, not plain sql
    query := "UPDATE " + table + " SET " + column + " = ? WHERE " + column + " = ?"
    _, err := tx.Exec(query, primaryID, secondaryID)
    if err != nil {
      return errors.WithStack(err)
    }
  }
  return checkMerge(tx, relatedFields)
}

func checkMerge(tx *sql.Tx, fields map[string]string) error {
	columns := []interface{}{}
	seenColumns := map[string]bool{}
	placeholders := []string{}
	for _, column := range fields {
		if _, ok := seenColumns[column]; !ok {
			columns = append(columns, column)
			seenColumns[column] = true
			placeholders = append(placeholders, "?")

		}
	}

	placeholder := strings.Join(placeholders, ", ")

	q := `SELECT table_name, column_name FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=DATABASE() AND column_name IN (` + placeholder + `)`
	rows, err := tx.Query(q, columns...)
	defer rows.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	for rows.Next() {
		var tableName string
		var columnName string
		err = rows.Scan(&tableName, &columnName)
		if err != nil {
			return errors.WithStack(err)
		}

		if _, exists := fields[tableName]; !exists {
			return errors.New("Missing merge for " + tableName + "." + columnName)
		}
	}

	return nil
}
