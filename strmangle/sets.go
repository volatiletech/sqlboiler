package strmangle

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
	merged := make([]string, 0, len(a)+len(b))

	for _, aVal := range a {
		found := false
		for _, mVal := range merged {
			if aVal == mVal {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, aVal)
		}
	}

	for _, bVal := range b {
		found := false
		for _, mVal := range merged {
			if bVal == mVal {
				found = true
				break
			}
		}

		if !found {
			merged = append(merged, bVal)
		}
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
