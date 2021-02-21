package boil

import (
	"github.com/volatiletech/strmangle"
)

// Columns kinds
// Note: These are not exported because they should only be used
// internally. To hide the clutter from the public API so boil
// continues to be a package aimed towards users, we add a method
// to query for each kind on the column type itself, see IsInfer
// as example.
const (
	columnsNone int = iota
	columnsInfer
	columnsWhitelist
	columnsGreylist
	columnsBlacklist
)

// Columns is a list of columns and a kind of list.
// Each kind interacts differently with non-zero and default column
// inference to produce a final list of columns for a given query.
// (Typically insert/updates).
type Columns struct {
	Kind int
	Cols []string
}

// None creates an empty column list.
func None() Columns {
	return Columns{
		Kind: columnsNone,
	}
}

// IsNone checks to see if no columns should be inferred.
// This method is here simply to not have to export the columns types.
func (c Columns) IsNone() bool {
	return c.Kind == columnsNone
}

// Infer is a placeholder that means there is no other list, simply
// infer the final list of columns for insert/update etc.
func Infer() Columns {
	return Columns{
		Kind: columnsInfer,
	}
}

// IsInfer checks to see if these columns should be inferred.
// This method is here simply to not have to export the columns types.
func (c Columns) IsInfer() bool {
	return c.Kind == columnsInfer
}

// Whitelist creates a list that completely overrides column inference.
// It becomes the final list for the insert/update etc.
func Whitelist(columns ...string) Columns {
	return Columns{
		Kind: columnsWhitelist,
		Cols: columns,
	}
}

// IsWhitelist checks to see if these columns should be inferred.
// This method is here simply to not have to export the columns types.
func (c Columns) IsWhitelist() bool {
	return c.Kind == columnsWhitelist
}

// Blacklist creates a list that overrides column inference choices
// by excluding a column from the inferred list. In essence, inference
// creates the list of columns, and blacklisted columns are removed from
// that list to produce the final list.
func Blacklist(columns ...string) Columns {
	return Columns{
		Kind: columnsBlacklist,
		Cols: columns,
	}
}

// IsBlacklist checks to see if these columns should be inferred.
// This method is here simply to not have to export the columns types.
func (c Columns) IsBlacklist() bool {
	return c.Kind == columnsBlacklist
}

// Greylist creates a list that adds to the inferred column choices.
// The final list is composed of both inferred columns and greylisted columns.
func Greylist(columns ...string) Columns {
	return Columns{
		Kind: columnsGreylist,
		Cols: columns,
	}
}

// IsGreylist checks to see if these columns should be inferred.
// This method is here simply to not have to export the columns types.
func (c Columns) IsGreylist() bool {
	return c.Kind == columnsGreylist
}

// InsertColumnSet generates the set of columns to insert and return for an
// insert statement. The return columns are used to get values that are
// assigned within the database during the insert to keep the struct in sync
// with what's in the db. The various interactions with the different
// types of Columns list are outlined below.
//
// Note that a default column's zero value is based on the Go type and does
// not take into account the default value in the database.
//
//  None:
//   insert: empty
//   return: empty
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
func (c Columns) InsertColumnSet(cols, defaults, noDefaults, nonZeroDefaults []string) ([]string, []string) {
	switch c.Kind {
	case columnsNone:
		return nil, nil

	case columnsInfer:
		insert := make([]string, len(noDefaults))
		copy(insert, noDefaults)
		insert = append(insert, nonZeroDefaults...)
		insert = strmangle.SortByKeys(cols, insert)
		ret := strmangle.SetComplement(defaults, insert)
		return insert, ret

	case columnsWhitelist:
		return c.Cols, strmangle.SetComplement(defaults, c.Cols)

	case columnsBlacklist:
		insert := make([]string, len(noDefaults))
		copy(insert, noDefaults)
		insert = append(insert, nonZeroDefaults...)
		insert = strmangle.SetComplement(insert, c.Cols)
		insert = strmangle.SortByKeys(cols, insert)
		ret := strmangle.SetComplement(defaults, insert)
		return insert, ret

	case columnsGreylist:
		insert := make([]string, len(noDefaults))
		copy(insert, noDefaults)
		insert = append(insert, nonZeroDefaults...)
		insert = strmangle.SetMerge(insert, c.Cols)
		insert = strmangle.SortByKeys(cols, insert)
		ret := strmangle.SetComplement(defaults, insert)
		return insert, ret
	default:
		panic("not a real column list kind")
	}
}

// UpdateColumnSet generates the set of columns to update for an update
// statement. The various interactions with the different types of Columns
// list are outlined below. In the case of greylist you can only add pkeys,
// which isn't useful in an update since then you can't find the original
// record you were trying to update.
//
//  None:      empty
//  Infer:     all - pkey-columns
//  whitelist: whitelist
//  blacklist: all - pkeys - blacklist
//  greylist:  all - pkeys + greylist
func (c Columns) UpdateColumnSet(allColumns, pkeyCols []string) []string {
	switch c.Kind {
	case columnsNone:
		return nil
	case columnsInfer:
		return strmangle.SetComplement(allColumns, pkeyCols)
	case columnsWhitelist:
		return c.Cols
	case columnsBlacklist:
		return strmangle.SetComplement(strmangle.SetComplement(allColumns, pkeyCols), c.Cols)
	case columnsGreylist:
		// okay to modify return of SetComplement since it's a new slice
		update := append(strmangle.SetComplement(allColumns, pkeyCols), c.Cols...)
		return strmangle.SortByKeys(allColumns, update)
	default:
		panic("not a real column list kind")
	}
}
