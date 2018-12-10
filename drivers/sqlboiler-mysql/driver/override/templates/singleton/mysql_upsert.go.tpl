// buildUpsertQueryMySQL builds a SQL statement string using the upsertData provided.
func buildUpsertQueryMySQL(dia drivers.Dialect, tableName string, update, whitelist []string) string {
	whitelist = strmangle.IdentQuoteSlice(dia.LQ, dia.RQ, whitelist)
	tableName = strmangle.IdentQuote(dia.LQ, dia.RQ, tableName)

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	var columns string
	if len(whitelist) != 0 {
		columns = strings.Join(whitelist, ",")
	}

	if len(update) == 0 {
		fmt.Fprintf(
			buf,
			"INSERT IGNORE INTO %s (%s) VALUES (%s)",
			tableName,
			columns,
			strmangle.Placeholders(dia.UseIndexPlaceholders, len(whitelist), 1, 1),
		)
		return buf.String()
	}

	fmt.Fprintf(
		buf,
		"INSERT INTO %s (%s) VALUES (%s) ON DUPLICATE KEY UPDATE ",
		tableName,
		columns,
		strmangle.Placeholders(dia.UseIndexPlaceholders, len(whitelist), 1, 1),
	)

	for i, v := range update {
		if i != 0 {
			buf.WriteByte(',')
		}
		quoted := strmangle.IdentQuote(dia.LQ, dia.RQ, v)
		buf.WriteString(quoted)
		buf.WriteString(" = VALUES(")
		buf.WriteString(quoted)
		buf.WriteByte(')')
	}

	return buf.String()
}
