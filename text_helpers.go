package main

import (
	"fmt"
	"strings"

	"github.com/vattle/sqlboiler/bdb"
	"github.com/vattle/sqlboiler/strmangle"
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
		Name         string
		ColumnNameGo string
		ColumnName   string
	}

	Function struct {
		Name        string
		ForeignName string

		Varname   string
		Receiver  string
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

	r.ForeignTable.Name = fkey.ForeignTable
	r.ForeignTable.NameGo = strmangle.TitleCase(strmangle.Singular(fkey.ForeignTable))
	r.ForeignTable.NamePluralGo = strmangle.TitleCase(strmangle.Plural(fkey.ForeignTable))
	r.ForeignTable.ColumnName = fkey.ForeignColumn
	r.ForeignTable.ColumnNameGo = strmangle.TitleCase(strmangle.Singular(fkey.ForeignColumn))

	r.Function.Name = strmangle.TitleCase(strmangle.Singular(trimSuffixes(fkey.Column)))
	plurality := strmangle.Plural
	if fkey.Unique {
		plurality = strmangle.Singular
	}
	r.Function.ForeignName = mkFunctionName(strmangle.Singular(fkey.ForeignTable), strmangle.TitleCase(plurality(fkey.Table)), fkey.Column, false)
	r.Function.Varname = strmangle.CamelCase(strmangle.Singular(fkey.ForeignTable))
	r.Function.Receiver = strings.ToLower(table.Name[:1])

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

	rel.Function.Name = strmangle.TitleCase(strmangle.Singular(oneToOne.ForeignTable))
	rel.Function.ForeignName = mkFunctionName(strmangle.Singular(oneToOne.Table), strmangle.TitleCase(strmangle.Singular(oneToOne.Table)), oneToOne.ForeignColumn, false)
	rel.Function.UsesBytes = col.Type == "[]byte"
	return rel
}

// TxtToMany contains text that will be used by many-to-one relationships.
type TxtToMany struct {
	LocalTable struct {
		NameGo       string
		NameSingular string
		ColumnNameGo string
	}

	ForeignTable struct {
		NameGo            string
		NameSingular      string
		NamePluralGo      string
		NameHumanReadable string
		ColumnNameGo      string
		Slice             string
	}

	Function struct {
		Name        string
		ForeignName string
		Receiver    string

		UsesBytes bool

		LocalAssignment   string
		ForeignAssignment string
	}
}

// txtsFromToMany creates a struct that does a lot of the text
// transformation in advance for a given relationship.
func txtsFromToMany(tables []bdb.Table, table bdb.Table, rel bdb.ToManyRelationship) TxtToMany {
	r := TxtToMany{}
	r.LocalTable.NameSingular = strmangle.Singular(table.Name)
	r.LocalTable.NameGo = strmangle.TitleCase(r.LocalTable.NameSingular)
	r.LocalTable.ColumnNameGo = strmangle.TitleCase(rel.Column)

	r.ForeignTable.NameSingular = strmangle.Singular(rel.ForeignTable)
	r.ForeignTable.NamePluralGo = strmangle.TitleCase(strmangle.Plural(rel.ForeignTable))
	r.ForeignTable.NameGo = strmangle.TitleCase(r.ForeignTable.NameSingular)
	r.ForeignTable.ColumnNameGo = strmangle.TitleCase(rel.ForeignColumn)
	r.ForeignTable.Slice = fmt.Sprintf("%sSlice", strmangle.TitleCase(r.ForeignTable.NameSingular))
	r.ForeignTable.NameHumanReadable = strings.Replace(rel.ForeignTable, "_", " ", -1)

	r.Function.Receiver = strings.ToLower(table.Name[:1])
	r.Function.Name = mkFunctionName(r.LocalTable.NameSingular, r.ForeignTable.NamePluralGo, rel.ForeignColumn, rel.ToJoinTable)
	plurality := strmangle.Singular
	foreignNamingColumn := rel.ForeignColumn
	if rel.ToJoinTable {
		plurality = strmangle.Plural
		foreignNamingColumn = rel.JoinLocalColumn
	}
	r.Function.ForeignName = strmangle.TitleCase(plurality(trimSuffixes(foreignNamingColumn)))

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
