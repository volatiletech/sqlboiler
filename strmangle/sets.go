package strmangle

// UpdateColumnSet generates the set of columns to update for an update statement.
// if a whitelist is supplied, it's returned
// if a whitelist is missing then we begin with all columns
// then we remove the primary key columns
func UpdateColumnSet(allColumns, pkeyCols, whitelist []string) []string {
	if len(whitelist) != 0 {
		return whitelist
	}

	return SetComplement(allColumns, pkeyCols)
}

// InsertColumnSet generates the set of columns to insert and return for an insert statement
// the return columns are used to get values that are assigned within the database during the
// insert to keep the struct in sync with what's in the db.
// with a whitelist:
// - the whitelist is used for the insert columns
// - the return columns are the result of (columns with default values - the whitelist)
// without a whitelist:
// - start with columns without a default as these always need to be inserted
// - add all columns that have a default in the database but that are non-zero in the struct
// - the return columns are the result of (columns with default values - the previous set)
func InsertColumnSet(cols, defaults, noDefaults, nonZeroDefaults, whitelist []string) ([]string, []string) {
	if len(whitelist) > 0 {
		return whitelist, SetComplement(defaults, whitelist)
	}

	var wl []string
	wl = append(wl, noDefaults...)
	wl = SetMerge(nonZeroDefaults, wl)
	wl = SortByKeys(cols, wl)

	// Only return the columns with default values that are not in the insert whitelist
	rc := SetComplement(defaults, wl)

	return wl, rc
}

// SetInclude checks to see if the string is found in the string slice
func SetInclude(str string, slice []string) bool {
	for _, s := range slice {
		if str == s {
			return true
		}
	}

	return false
}

// SetComplement subtracts the elements in b from a
func SetComplement(a []string, b []string) []string {
	c := make([]string, 0, len(a))

	for _, aVal := range a {
		found := false
		for _, bVal := range b {
			if aVal == bVal {
				found = true
				break
			}
		}
		if !found {
			c = append(c, aVal)
		}
	}

	return c
}

// SetMerge will return a merged slice without duplicates
func SetMerge(a []string, b []string) []string {
	var x, merged []string

	x = append(x, a...)
	x = append(x, b...)

	check := map[string]bool{}
	for _, v := range x {
		if check[v] == true {
			continue
		}

		merged = append(merged, v)
		check[v] = true
	}

	return merged
}

// SortByKeys returns a new ordered slice based on the keys ordering
func SortByKeys(keys []string, strs []string) []string {
	c := make([]string, len(strs))

	index := 0
Outer:
	for _, v := range keys {
		for _, k := range strs {
			if v == k {
				c[index] = v
				index++

				if index > len(strs)-1 {
					break Outer
				}
				break
			}
		}
	}

	return c
}
