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
			return errors.WithStack(err)
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
			return errors.WithStack(err)
		}
	}
	return checkMerge(tx, foreignKeys)
}

func deleteConflictsBeforeMerge(tx boil.Executor, conflict conflictingUniqueKey, primaryID uint64, secondaryID uint64) error {
	conflictingColumns := strmangle.SetComplement(conflict.columns, []string{conflict.objectIdColumn})

	if len(conflictingColumns) < 1 {
		return nil
	} else if len(conflictingColumns) > 1 {
		return errors.New("this doesnt work for unique keys with more than two columns (yet)")
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
		return errors.WithStack(err)
	}

	args := []interface{}{secondaryID}
	for rows.Next() {
		var value string
		err = rows.Scan(&value)
		if err != nil {
			return errors.WithStack(err)
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
		return errors.WithStack(err)
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
		return errors.WithStack(err)
	}

	for rows.Next() {
		var tableName string
		var columnName string
		err = rows.Scan(&tableName, &columnName)
		if err != nil {
			return errors.WithStack(err)
		}

		if _, exists := handledTablesColumns[tableName+"."+columnName]; !exists {
			return errors.New("Missing merge for " + tableName + "." + columnName)
		}
	}

	return nil
}

// duplicated in queries/query.go
func InterpolateParams(query string, args ...interface{}) (string, error) {
	for i := 0; i < len(args); i++ {
		field := reflect.ValueOf(args[i])

		if value, ok := field.Interface().(time.Time); ok {
			query = strings.Replace(query, "?", `"`+value.Format("2006-01-02 15:04:05")+`"`, 1)
		} else if nullable, ok := field.Interface().(null.Nullable); ok {
			if nullable.IsNull() {
				query = strings.Replace(query, "?", "NULL", 1)
			} else {
				switch field.Type() {
				case reflect.TypeOf(null.Time{}):
					query = strings.Replace(query, "?", `"`+field.Interface().(null.Time).Time.Format("2006-01-02 15:04:05")+`"`, 1)
				case reflect.TypeOf(null.Int{}):
					query = strings.Replace(query, "?", strconv.FormatInt(int64(field.Interface().(null.Int).Int), 10), 1)
				case reflect.TypeOf(null.Int8{}):
					query = strings.Replace(query, "?", strconv.FormatInt(int64(field.Interface().(null.Int8).Int8), 10), 1)
				case reflect.TypeOf(null.Int16{}):
					query = strings.Replace(query, "?", strconv.FormatInt(int64(field.Interface().(null.Int16).Int16), 10), 1)
				case reflect.TypeOf(null.Int32{}):
					query = strings.Replace(query, "?", strconv.FormatInt(int64(field.Interface().(null.Int32).Int32), 10), 1)
				case reflect.TypeOf(null.Int64{}):
					query = strings.Replace(query, "?", strconv.FormatInt(field.Interface().(null.Int64).Int64, 10), 1)
				case reflect.TypeOf(null.Uint{}):
					query = strings.Replace(query, "?", strconv.FormatUint(uint64(field.Interface().(null.Uint).Uint), 10), 1)
				case reflect.TypeOf(null.Uint8{}):
					query = strings.Replace(query, "?", strconv.FormatUint(uint64(field.Interface().(null.Uint8).Uint8), 10), 1)
				case reflect.TypeOf(null.Uint16{}):
					query = strings.Replace(query, "?", strconv.FormatUint(uint64(field.Interface().(null.Uint16).Uint16), 10), 1)
				case reflect.TypeOf(null.Uint32{}):
					query = strings.Replace(query, "?", strconv.FormatUint(uint64(field.Interface().(null.Uint32).Uint32), 10), 1)
				case reflect.TypeOf(null.Uint64{}):
					query = strings.Replace(query, "?", strconv.FormatUint(field.Interface().(null.Uint64).Uint64, 10), 1)
				case reflect.TypeOf(null.String{}):
					query = strings.Replace(query, "?", `"`+field.Interface().(null.String).String+`"`, 1)
				case reflect.TypeOf(null.Bool{}):
					if field.Interface().(null.Bool).Bool {
						query = strings.Replace(query, "?", "1", 1)
					} else {
						query = strings.Replace(query, "?", "0", 1)
					}
				}
			}
		} else {
			switch field.Kind() {
			case reflect.Bool:
				boolString := "0"
				if field.Bool() {
					boolString = "1"
				}
				query = strings.Replace(query, "?", boolString, 1)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				query = strings.Replace(query, "?", strconv.FormatInt(field.Int(), 10), 1)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				query = strings.Replace(query, "?", strconv.FormatUint(field.Uint(), 10), 1)
			case reflect.String:
				query = strings.Replace(query, "?", `"`+field.String()+`"`, 1)
			default:
				return "", errors.New("Dont know how to interpolate type " + field.Type().String())
			}
		}
	}
	return query, nil
}
