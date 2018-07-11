// Package strmangle is a collection of string manipulation functions.
// Primarily used by boil and templates for code generation.
// Because it is focused on pipelining inside templates
// you will see some odd parameter ordering.
package strmangle

type Mangler interface {
	SchemaTable(lq, rq string, driver string, schema string, table string) string
	IdentQuote(lq byte, rq byte, s string) string
	IdentQuoteSlice(lq byte, rq byte, s []string) []string
	Identifier(in int) string
	QuoteCharacter(q byte) string
	Plural(name string) string
	Singular(name string) string
	TitleCase(n string) string
	CamelCase(name string) string
	TitleCaseIdentifier(id string) string
	MakeStringMap(types map[string]string) string
	StringMap(modifier func(string) string, strs []string) []string
	PrefixStringSlice(str string, strs []string) []string
	Placeholders(indexPlaceholders bool, count int, start int, group int) string
	SetParamNames(lq, rq string, start int, columns []string) string
	WhereClause(lq, rq string, start int, cols []string) string
	WhereClauseRepeated(lq, rq string, start int, cols []string, count int) string
	JoinSlices(sep string, a, b []string) []string
	StringSliceMatch(a []string, b []string) bool
	ContainsAny(a []string, finds ...string) bool
	GenerateTags(tags []string, columnName string) string
	GenerateIgnoreTags(tags []string) string
	ParseEnumVals(s string) []string
	ParseEnumName(s string) string
	IsEnumNormal(values []string) bool
	ShouldTitleCaseEnum(value string) bool
	ReplaceReservedWords(word string) string

	UpdateColumnSet(allColumns, pkeyCols, whitelist []string) []string
	InsertColumnSet(cols, defaults, noDefaults, nonZeroDefaults, whitelist []string) ([]string, []string)
	SetInclude(str string, slice []string) bool
	SetComplement(a []string, b []string) []string
	SetMerge(a []string, b []string) []string
	SortByKeys(keys []string, strs []string) []string
}

// DefaultMangler is the mangler used by all the
var DefaultMangler = NewDefaultMangler()

// SchemaTable returns a table name with a schema prefixed if
// using a database that supports real schemas, for example,
// for Postgres: "schema_name"."table_name",
// for MS SQL: [schema_name].[table_name], versus
// simply "table_name" for MySQL (because it does not support real schemas)
func SchemaTable(lq, rq string, driver string, schema string, table string) string {
	return DefaultMangler.SchemaTable(lq, rq, driver, schema, table)
}

// IdentQuote attempts to quote simple identifiers in SQL statements
func IdentQuote(lq byte, rq byte, s string) string { return DefaultMangler.IdentQuote(lq, rq, s) }

// IdentQuoteSlice applies IdentQuote to a slice.
func IdentQuoteSlice(lq byte, rq byte, s []string) []string {
	return DefaultMangler.IdentQuoteSlice(lq, rq, s)
}

// Identifier is a base conversion from Base 10 integers to Base 26
// integers that are represented by an alphabet from a-z
// See tests for example outputs.
func Identifier(in int) string { return DefaultMangler.Identifier(in) }

// QuoteCharacter returns a string that allows the quote character
// to be embedded into a Go string that uses double quotes:
func QuoteCharacter(q byte) string { return DefaultMangler.QuoteCharacter(q) }

// Plural converts singular words to plural words (eg: person to people)
func Plural(name string) string { return DefaultMangler.Plural(name) }

// Singular converts plural words to singular words (eg: people to person)
func Singular(name string) string { return DefaultMangler.Singular(name) }

// TitleCase changes a snake-case variable name
// into a go styled object variable name of "ColumnName".
// titleCase also fully uppercases "ID" components of names, for example
// "column_name_id" to "ColumnNameID".
//
// Note: This method is ugly because it has been highly optimized,
// we found that it was a fairly large bottleneck when we were using regexp.
func TitleCase(n string) string { return DefaultMangler.TitleCase(n) }

// CamelCase takes a variable name in the format of "var_name" and converts
// it into a go styled variable name of "varName".
// camelCase also fully uppercases "ID" components of names, for example
// "var_name_id" to "varNameID".
func CamelCase(name string) string { return DefaultMangler.CamelCase(name) }

// TitleCaseIdentifier splits on dots and then titlecases each fragment.
// map titleCase (split c ".")
func TitleCaseIdentifier(id string) string { return DefaultMangler.TitleCaseIdentifier(id) }

// MakeStringMap converts a map[string]string into the format:
// "key": "value", "key": "value"
func MakeStringMap(types map[string]string) string { return DefaultMangler.MakeStringMap(types) }

// StringMap maps a function over a slice of strings.
func StringMap(modifier func(string) string, strs []string) []string {
	return DefaultMangler.StringMap(modifier, strs)
}

// PrefixStringSlice with the given str.
func PrefixStringSlice(str string, strs []string) []string {
	return DefaultMangler.PrefixStringSlice(str, strs)
}

// Placeholders generates the SQL statement placeholders for in queries.
// For example, ($1,$2,$3),($4,$5,$6) etc.
// It will start counting placeholders at "start".
// If indexPlaceholders is false, it will convert to ? instead of $1 etc.
func Placeholders(indexPlaceholders bool, count int, start int, group int) string {
	return DefaultMangler.Placeholders(indexPlaceholders, count, start, group)
}

// SetParamNames takes a slice of columns and returns a comma separated
// list of parameter names for a template statement SET clause.
// eg: "col1"=$1, "col2"=$2, "col3"=$3
func SetParamNames(lq, rq string, start int, columns []string) string {
	return DefaultMangler.SetParamNames(lq, rq, start, columns)
}

// WhereClause returns the where clause using start as the $ flag index
// For example, if start was 2 output would be: "colthing=$2 AND colstuff=$3"
func WhereClause(lq, rq string, start int, cols []string) string {
	return DefaultMangler.WhereClause(lq, rq, start, cols)
}

// WhereClauseRepeated returns the where clause repeated with OR clause using start as the $ flag index
// For example, if start was 2 output would be: "(colthing=$2 AND colstuff=$3) OR (colthing=$4 AND colstuff=$5)"
func WhereClauseRepeated(lq, rq string, start int, cols []string, count int) string {
	return DefaultMangler.WhereClauseRepeated(lq, rq, start, cols, count)
}

// JoinSlices merges two string slices of equal length
func JoinSlices(sep string, a, b []string) []string { return DefaultMangler.JoinSlices(sep, a, b) }

// StringSliceMatch returns true if the length of both
// slices is the same, and the elements of both slices are the same.
// The elements can be in any order.
func StringSliceMatch(a []string, b []string) bool { return DefaultMangler.StringSliceMatch(a, b) }

// ContainsAny returns true if any of the passed in strings are
// found in the passed in string slice
func ContainsAny(a []string, finds ...string) bool { return DefaultMangler.ContainsAny(a, finds...) }

// GenerateTags converts a slice of tag strings into tags that
// can be passed onto the end of a struct, for example:
// tags: ["xml", "db"] convert to: xml:"column_name" db:"column_name"
func GenerateTags(tags []string, columnName string) string {
	return DefaultMangler.GenerateTags(tags, columnName)
}

// GenerateIgnoreTags converts a slice of tag strings into
// ignore tags that can be passed onto the end of a struct, for example:
// tags: ["xml", "db"] convert to: xml:"-" db:"-"
func GenerateIgnoreTags(tags []string) string { return DefaultMangler.GenerateIgnoreTags(tags) }

// ParseEnumVals returns the values from an enum string
//
// Postgres and MySQL drivers return different values
// psql:  enum.enum_name('values'...)
// mysql: enum('values'...)
func ParseEnumVals(s string) []string { return DefaultMangler.ParseEnumVals(s) }

// ParseEnumName returns the name portion of an enum if it exists
//
// Postgres and MySQL drivers return different values
// psql:  enum.enum_name('values'...)
// mysql: enum('values'...)
// In the case of mysql, the name will never return anything
func ParseEnumName(s string) string { return DefaultMangler.ParseEnumName(s) }

// IsEnumNormal checks a set of eval values to see if they're "normal"
func IsEnumNormal(values []string) bool { return DefaultMangler.IsEnumNormal(values) }

// ShouldTitleCaseEnum checks a value to see if it's title-case-able
func ShouldTitleCaseEnum(value string) bool { return DefaultMangler.ShouldTitleCaseEnum(value) }

// ReplaceReservedWords takes a word and replaces it with word_ if it's found
// in the list of reserved words.
func ReplaceReservedWords(word string) string { return DefaultMangler.ReplaceReservedWords(word) }

// UpdateColumnSet generates the set of columns to update for an update statement.
// if a whitelist is supplied, it's returned
// if a whitelist is missing then we begin with all columns
// then we remove the primary key columns
func UpdateColumnSet(allColumns, pkeyCols, whitelist []string) []string {
	return DefaultMangler.UpdateColumnSet(allColumns, pkeyCols, whitelist)
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
	return DefaultMangler.InsertColumnSet(cols, defaults, noDefaults, nonZeroDefaults, whitelist)
}

// SetInclude checks to see if the string is found in the string slice
func SetInclude(str string, slice []string) bool { return DefaultMangler.SetInclude(str, slice) }

// SetComplement subtracts the elements in b from a
func SetComplement(a []string, b []string) []string { return DefaultMangler.SetComplement(a, b) }

// SetMerge will return a merged slice without duplicates
func SetMerge(a []string, b []string) []string { return DefaultMangler.SetMerge(a, b) }

// SortByKeys returns a new ordered slice based on the keys ordering
func SortByKeys(keys []string, strs []string) []string { return DefaultMangler.SortByKeys(keys, strs) }
