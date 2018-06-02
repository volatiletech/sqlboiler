package boil

import (
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Columns kinds
const (
	ColumnsInfer int = iota
	ColumnsWhitelist
	ColumnsGreylist
	ColumnsBlacklist
)

// Columns is a list of columns and a kind of list.
// Each kind interacts differently with non-zero and default column
// inference to produce a final list of columns for a given query.
// (Typically insert/updates).
type Columns struct {
	Kind int
	Cols []string
}

// Infer is a placeholder that means there is no other list, simply
// infer the final list of columns for insert/update etc.
func Infer() Columns {
	return Columns{}
}

// Whitelist creates a list that completely overrides column inference.
// It becomes the final list for the insert/update etc.
func Whitelist(columns ...string) Columns {
	return Columns{
		Kind: ColumnsWhitelist,
		Cols: columns,
	}
}

// Blacklist creates a list that overrides column inference choices
// by excluding a column from the inferred list. In essence, inference
// creates the list of columns, and blacklisted columns are removed from
// that list to produce the final list.
func Blacklist(columns ...string) Columns {
	return Columns{
		Kind: ColumnsBlacklist,
		Cols: columns,
	}
}

// Greylist creates a list that adds to the inferred column choices.
// The final list is composed of both inferred columns and greylisted columns.
func Greylist(columns ...string) Columns {
	return Columns{
		Kind: ColumnsGreylist,
		Cols: columns,
	}
}

// InsertColumnSet generates the set of columns to insert and return for an insert statement.
// The return columns are used to get values that are assigned within the database during the
// insert to keep the struct in sync with what's in the db. The various interactions with the different
// types of Columns list are outlined below.
//
//  Infer:
//   insert: columns-without-default + non-zero-default-columns
//   return: columns-with-defaults - insert
//
//  Whitelist:
//   insert: whitelist
//   return: columns-with-defaults - whitelist
//
//  Blacklist:
//    insert: columns-without-default + non-zero-default-columns - blacklist
//    return: columns-with-defaults - insert
//
//  Greylist:
//    insert: columns-without-default + non-zero-default-columns + greylist
//    return: columns-with-defaults - insert
func InsertColumnSet(cols, defaults, noDefaults, nonZeroDefaults []string, columns Columns) ([]string, []string) {
	switch columns.Kind {
	case ColumnsInfer:
		insert := make([]string, len(noDefaults))
		copy(insert, noDefaults)
		insert = append(insert, nonZeroDefaults...)
		insert = strmangle.SortByKeys(cols, insert)
		ret := strmangle.SetComplement(defaults, insert)
		return insert, ret

	case ColumnsWhitelist:
		return columns.Cols, strmangle.SetComplement(defaults, columns.Cols)

	case ColumnsBlacklist:
		insert := make([]string, len(noDefaults))
		copy(insert, noDefaults)
		insert = append(insert, nonZeroDefaults...)
		insert = strmangle.SetComplement(insert, columns.Cols)
		insert = strmangle.SortByKeys(cols, insert)
		ret := strmangle.SetComplement(defaults, insert)
		return insert, ret

	case ColumnsGreylist:
		insert := make([]string, len(noDefaults))
		copy(insert, noDefaults)
		insert = append(insert, nonZeroDefaults...)
		insert = strmangle.SetMerge(insert, columns.Cols)
		insert = strmangle.SortByKeys(cols, insert)
		ret := strmangle.SetComplement(defaults, insert)
		return insert, ret
	default:
		panic("not a real column list kind")
	}
}

// UpdateColumnSet generates the set of columns to update for an update statement.
// The various interactions with the different types of Columns list are outlined below.
// In the case of greylist you can only add pkeys, which isn't useful in an update since
// then you can't find the original record you were trying to update.
//
//  Infer:     all - pkey-columns
//  whitelist: whitelist
//  blacklist: all - pkeys - blacklist
//  greylist:  all - pkeys + greylist
func UpdateColumnSet(allColumns, pkeyCols []string, columns Columns) []string {
	switch columns.Kind {
	case ColumnsInfer:
		return strmangle.SetComplement(allColumns, pkeyCols)
	case ColumnsWhitelist:
		return columns.Cols
	case ColumnsBlacklist:
		return strmangle.SetComplement(strmangle.SetComplement(allColumns, pkeyCols), columns.Cols)
	case ColumnsGreylist:
		// okay to modify return of SetComplement since it's a new slice
		update := append(strmangle.SetComplement(allColumns, pkeyCols), columns.Cols...)
		return strmangle.SortByKeys(allColumns, update)
	default:
		panic("not a real column list kind")
	}
}
