// buildUpsertQueryMSSQL builds a SQL statement string using the upsertData provided.
func buildUpsertQueryMSSQL(dia drivers.Dialect, tableName string, primary, update, insert []string, output []string) string {
	insert = strmangle.IdentQuoteSlice(dia.LQ, dia.RQ, insert)

	buf := strmangle.GetBuffer()
	defer strmangle.PutBuffer(buf)

	startIndex := 1

	fmt.Fprintf(buf, "MERGE INTO %s as [t]\n", tableName)
	fmt.Fprintf(buf, "USING (SELECT %s) as [s] ([%s])\n",
		strmangle.Placeholders(dia.UseIndexPlaceholders, len(primary), startIndex, 1),
		strings.Join(primary, string(dia.RQ)+","+string(dia.LQ)))
	fmt.Fprint(buf, "ON (")
	for i, v := range primary {
		if i != 0 {
			fmt.Fprint(buf, " AND ")
		}
		fmt.Fprintf(buf, "[s].[%s] = [t].[%s]", v, v)
	}
	fmt.Fprint(buf, ")\n")

	startIndex += len(primary)

	if len(update) > 0 {
		fmt.Fprint(buf, "WHEN MATCHED THEN ")
		fmt.Fprintf(buf, "UPDATE SET %s\n", strmangle.SetParamNames(string(dia.LQ), string(dia.RQ), startIndex, update))

		startIndex += len(update)
	}

	fmt.Fprint(buf, "WHEN NOT MATCHED THEN ")
	fmt.Fprintf(buf, "INSERT (%s) VALUES (%s)",
		strings.Join(insert, ", "),
		strmangle.Placeholders(dia.UseIndexPlaceholders, len(insert), startIndex, 1))

	if len(output) > 0 {
		fmt.Fprintf(buf, "\nOUTPUT INSERTED.[%s];", strings.Join(output, "],INSERTED.["))
	} else {
		fmt.Fprint(buf, ";")
	}

	return buf.String()
}
