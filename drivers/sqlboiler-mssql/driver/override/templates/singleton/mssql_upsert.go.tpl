// buildUpsertQueryMSSQL builds a SQL statement string using the upsertData provided.
func buildUpsertQueryMSSQL(dia drivers.Dialect, tableName string, primary, update, insert []string, output []string) string {
	insert = strmangle.IdentQuoteSlice(dia.LQ, dia.RQ, insert)

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	startIndex := 1

	_, _ = fmt.Fprintf(buf, "MERGE INTO %s as [t]\n", tableName)
	_, _ = fmt.Fprintf(buf, "USING (SELECT %s) as [s] ([%s])\n",
		strmangle.Placeholders(dia.UseIndexPlaceholders, len(primary), startIndex, 1),
		strings.Join(primary, string(dia.RQ)+","+string(dia.LQ)))
	_, _ = fmt.Fprint(buf, "ON (")
	for i, v := range primary {
		if i != 0 {
			_, _ = fmt.Fprint(buf, " AND ")
		}
		_, _ = fmt.Fprintf(buf, "[s].[%s] = [t].[%s]", v, v)
	}
	_, _ = fmt.Fprint(buf, ")\n")

	startIndex += len(primary)

	_, _ = fmt.Fprint(buf, "WHEN MATCHED THEN ")
	_, _ = fmt.Fprintf(buf, "UPDATE SET %s\n", strmangle.SetParamNames(string(dia.LQ), string(dia.RQ), startIndex, update))

	startIndex += len(update)

	_, _ = fmt.Fprint(buf, "WHEN NOT MATCHED THEN ")
	_, _ = fmt.Fprintf(buf, "INSERT (%s) VALUES (%s)",
		strings.Join(insert, ", "),
		strmangle.Placeholders(dia.UseIndexPlaceholders, len(insert), startIndex, 1))

	if len(output) > 0 {
		_, _ = fmt.Fprintf(buf, "\nOUTPUT INSERTED.[%s];", strings.Join(output, "],INSERTED.["))
	} else {
		_, _ = fmt.Fprint(buf, ";")
	}

	return buf.String()
}
