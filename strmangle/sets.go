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

// SetIntersect returns the elements that are common in a and b
func SetIntersect(a []string, b []string) []string {
	c := make([]string, 0, len(a))

	for _, aVal := range a {
		found := false
		for _, bVal := range b {
			if aVal == bVal {
				found = true
				break
			}
		}
		if found {
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
