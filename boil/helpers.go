package boil

import (
	"fmt"
	"reflect"

	"github.com/vattle/sqlboiler/strmangle"
)

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

// NonZeroDefaultSet returns the fields included in the
// defaults slice that are non zero values
func NonZeroDefaultSet(defaults []string, obj interface{}) []string {
	c := make([]string, 0, len(defaults))

	val := reflect.Indirect(reflect.ValueOf(obj))

	for _, d := range defaults {
		fieldName := strmangle.TitleCase(d)
		field := val.FieldByName(fieldName)
		if !field.IsValid() {
			panic(fmt.Sprintf("Could not find field name %s in type %T", fieldName, obj))
		}

		zero := reflect.Zero(field.Type())
		if !reflect.DeepEqual(zero.Interface(), field.Interface()) {
			c = append(c, d)
		}
	}

	return c
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
