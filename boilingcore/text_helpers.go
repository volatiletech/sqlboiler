package boilingcore

import (
	"fmt"
	"strings"

	"github.com/ann-kilzer/sqlboiler/bdb"
	"github.com/ann-kilzer/sqlboiler/strmangle"
)

// TxtToOne contains text that will be used by templates for a one-to-many or
// a one-to-one relationship.
type TxtToOne struct {
	ForeignKey bdb.ForeignKey

	LocalTable struct {
		NameGo       string
		ColumnNameGo string
	}

	ForeignTable struct {
		NameGo       string
		NamePluralGo string
		ColumnNameGo string
		ColumnName   string
	}

	Function struct {
		Name        string
		ForeignName string

		UsesBytes bool

		LocalAssignment   string
		ForeignAssignment string
	}
}

func txtsFromFKey(tables []bdb.Table, table bdb.Table, fkey bdb.ForeignKey) TxtToOne {
	r := TxtToOne{}

	r.ForeignKey = fkey

	r.LocalTable.NameGo = strmangle.TitleCase(strmangle.Singular(table.Name))
	r.LocalTable.ColumnNameGo = strmangle.TitleCase(strmangle.Singular(fkey.Column))

	r.ForeignTable.NameGo = strmangle.TitleCase(strmangle.Singular(fkey.ForeignTable))
	r.ForeignTable.NamePluralGo = strmangle.TitleCase(strmangle.Plural(fkey.ForeignTable))
	r.ForeignTable.ColumnName = fkey.ForeignColumn
	r.ForeignTable.ColumnNameGo = strmangle.TitleCase(strmangle.Singular(fkey.ForeignColumn))

	r.Function.Name, r.Function.ForeignName = txtNameToOne(fkey)

	if fkey.Nullable {
		col := table.GetColumn(fkey.Column)
		r.Function.LocalAssignment = fmt.Sprintf("%s.%s", strmangle.TitleCase(fkey.Column), strings.TrimPrefix(col.Type, "null."))
	} else {
		r.Function.LocalAssignment = strmangle.TitleCase(fkey.Column)
	}

	foreignTable := bdb.GetTable(tables, fkey.ForeignTable)
	foreignColumn := foreignTable.GetColumn(fkey.ForeignColumn)

	if fkey.ForeignColumnNullable {
		r.Function.ForeignAssignment = fmt.Sprintf("%s.%s", strmangle.TitleCase(fkey.ForeignColumn), strings.TrimPrefix(foreignColumn.Type, "null."))
	} else {
		r.Function.ForeignAssignment = strmangle.TitleCase(fkey.ForeignColumn)
	}

	r.Function.UsesBytes = foreignColumn.Type == "[]byte"

	return r
}

func txtsFromOneToOne(tables []bdb.Table, table bdb.Table, oneToOne bdb.ToOneRelationship) TxtToOne {
	fkey := bdb.ForeignKey{
		Table:    oneToOne.Table,
		Name:     "none",
		Column:   oneToOne.Column,
		Nullable: oneToOne.Nullable,
		Unique:   oneToOne.Unique,

		ForeignTable:          oneToOne.ForeignTable,
		ForeignColumn:         oneToOne.ForeignColumn,
		ForeignColumnNullable: oneToOne.ForeignColumnNullable,
		ForeignColumnUnique:   oneToOne.ForeignColumnUnique,
	}

	rel := txtsFromFKey(tables, table, fkey)
	col := table.GetColumn(oneToOne.Column)

	// Reverse foreign key
	rel.ForeignKey.Table, rel.ForeignKey.ForeignTable = rel.ForeignKey.ForeignTable, rel.ForeignKey.Table
	rel.ForeignKey.Column, rel.ForeignKey.ForeignColumn = rel.ForeignKey.ForeignColumn, rel.ForeignKey.Column
	rel.ForeignKey.Nullable, rel.ForeignKey.ForeignColumnNullable = rel.ForeignKey.ForeignColumnNullable, rel.ForeignKey.Nullable
	rel.ForeignKey.Unique, rel.ForeignKey.ForeignColumnUnique = rel.ForeignKey.ForeignColumnUnique, rel.ForeignKey.Unique
	rel.Function.UsesBytes = col.Type == "[]byte"
	rel.Function.ForeignName, rel.Function.Name = txtNameToOne(bdb.ForeignKey{
		Table:         oneToOne.ForeignTable,
		Column:        oneToOne.ForeignColumn,
		Unique:        true,
		ForeignTable:  oneToOne.Table,
		ForeignColumn: oneToOne.Column,
	})
	return rel
}

// TxtToMany contains text that will be used by many-to-one relationships.
type TxtToMany struct {
	LocalTable struct {
		NameGo       string
		ColumnNameGo string
	}

	ForeignTable struct {
		NameGo            string
		NamePluralGo      string
		NameHumanReadable string
		ColumnNameGo      string
		Slice             string
	}

	Function struct {
		Name        string
		ForeignName string

		UsesBytes bool

		LocalAssignment   string
		ForeignAssignment string
	}
}

// txtsFromToMany creates a struct that does a lot of the text
// transformation in advance for a given relationship.
func txtsFromToMany(tables []bdb.Table, table bdb.Table, rel bdb.ToManyRelationship) TxtToMany {
	r := TxtToMany{}
	r.LocalTable.NameGo = strmangle.TitleCase(strmangle.Singular(table.Name))
	r.LocalTable.ColumnNameGo = strmangle.TitleCase(rel.Column)

	foreignNameSingular := strmangle.Singular(rel.ForeignTable)
	r.ForeignTable.NamePluralGo = strmangle.TitleCase(strmangle.Plural(rel.ForeignTable))
	r.ForeignTable.NameGo = strmangle.TitleCase(foreignNameSingular)
	r.ForeignTable.ColumnNameGo = strmangle.TitleCase(rel.ForeignColumn)
	r.ForeignTable.Slice = fmt.Sprintf("%sSlice", strmangle.TitleCase(foreignNameSingular))
	r.ForeignTable.NameHumanReadable = strings.Replace(rel.ForeignTable, "_", " ", -1)

	r.Function.Name, r.Function.ForeignName = txtNameToMany(rel)

	col := table.GetColumn(rel.Column)
	if rel.Nullable {
		r.Function.LocalAssignment = fmt.Sprintf("%s.%s", strmangle.TitleCase(rel.Column), strings.TrimPrefix(col.Type, "null."))
	} else {
		r.Function.LocalAssignment = strmangle.TitleCase(rel.Column)
	}

	if rel.ForeignColumnNullable {
		foreignTable := bdb.GetTable(tables, rel.ForeignTable)
		foreignColumn := foreignTable.GetColumn(rel.ForeignColumn)
		r.Function.ForeignAssignment = fmt.Sprintf("%s.%s", strmangle.TitleCase(rel.ForeignColumn), strings.TrimPrefix(foreignColumn.Type, "null."))
	} else {
		r.Function.ForeignAssignment = strmangle.TitleCase(rel.ForeignColumn)
	}

	r.Function.UsesBytes = col.Type == "[]byte"

	return r
}

// txtNameToOne creates the local and foreign function names for
// one-to-many and one-to-one relationships, where local == lhs (one).
//
// = many-to-one
// users - videos : user_id
// users - videos : producer_id
//
// fk == table = user.Videos         | video.User
// fk != table = user.ProducerVideos | video.Producer
//
// = many-to-one
// industries - industries : parent_id
//
// fk == table = industry.Industries | industry.Industry
// fk != table = industry.ParentIndustries | industry.Parent
//
// = one-to-one
// users - videos : user_id
// users - videos : producer_id
//
// fk == table = user.Video         | video.User
// fk != table = user.ProducerVideo | video.Producer
//
// = one-to-one
// industries - industries : parent_id
//
// fk == table = industry.Industry | industry.Industry
// fk != table = industry.ParentIndustry | industry.Industry
func txtNameToOne(fk bdb.ForeignKey) (localFn, foreignFn string) {
	localFn = strmangle.Singular(trimSuffixes(fk.Column))
	fkeyIsTableName := localFn != strmangle.Singular(fk.ForeignTable)
	localFn = strmangle.TitleCase(localFn)

	if fkeyIsTableName {
		foreignFn = localFn
	}

	plurality := strmangle.Plural
	if fk.Unique {
		plurality = strmangle.Singular
	}
	foreignFn += strmangle.TitleCase(plurality(fk.Table))

	return localFn, foreignFn
}

// txtNameToMany creates the local and foreign function names for
// many-to-one and many-to-many relationship, where local == lhs (many)
//
// cases:
// = many-to-many
// sponsors - constests
// sponsor_id contest_id
// fk == table = sponsor.Contests | contest.Sponsors
//
// = many-to-many
// sponsors - constests
// wiggle_id jiggle_id
// fk != table = sponsor.JiggleSponsors | contest.WiggleContests
//
// = many-to-many
// industries - industries
// industry_id  mapped_industry_id
//
// fk == table = industry.Industries
// fk != table = industry.MappedIndustryIndustry
func txtNameToMany(toMany bdb.ToManyRelationship) (localFn, foreignFn string) {
	if toMany.ToJoinTable {
		localFkey := strmangle.Singular(trimSuffixes(toMany.JoinLocalColumn))
		foreignFkey := strmangle.Singular(trimSuffixes(toMany.JoinForeignColumn))

		if localFkey != strmangle.Singular(toMany.Table) {
			foreignFn = strmangle.TitleCase(localFkey)
		}
		foreignFn += strmangle.TitleCase(strmangle.Plural(toMany.Table))

		if foreignFkey != strmangle.Singular(toMany.ForeignTable) {
			localFn = strmangle.TitleCase(foreignFkey)
		}
		localFn += strmangle.TitleCase(strmangle.Plural(toMany.ForeignTable))

		return localFn, foreignFn
	}

	fkeyName := strmangle.Singular(trimSuffixes(toMany.ForeignColumn))
	if fkeyName != strmangle.Singular(toMany.Table) {
		localFn = strmangle.TitleCase(fkeyName)
	}
	localFn += strmangle.TitleCase(strmangle.Plural(toMany.ForeignTable))
	foreignFn = strmangle.TitleCase(strmangle.Singular(fkeyName))
	return localFn, foreignFn
}

// mkFunctionName checks to see if the foreign key name is the same as the local table name (minus _id suffix)
// Simple case: yes - we can name the function the same as the plural table name
// Not simple case: We have to name the function based off the foreign key and the foreign table name
func mkFunctionName(fkeyTableSingular, foreignTablePluralGo, fkeyColumn string, toJoinTable bool) string {
	colName := trimSuffixes(fkeyColumn)
	if toJoinTable || fkeyTableSingular == colName {
		return foreignTablePluralGo
	}

	return strmangle.TitleCase(colName) + foreignTablePluralGo
}

var identifierSuffixes = []string{"_id", "_uuid", "_guid", "_oid"}

// trimSuffixes from the identifier
func trimSuffixes(str string) string {
	ln := len(str)
	for _, s := range identifierSuffixes {
		str = strings.TrimSuffix(str, s)
		if len(str) != ln {
			break
		}
	}

	return str
}
