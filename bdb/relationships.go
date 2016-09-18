package bdb

// ToOneRelationship describes a relationship between two tables where the local
// table has no id, and the foregin table has an id that matches a column in the
// local table, that column is also unique which changes the dynamic into a
// one-to-one style, not a to-many.
type ToOneRelationship struct {
	Table    string
	Column   string
	Nullable bool
	Unique   bool

	ForeignTable          string
	ForeignColumn         string
	ForeignColumnNullable bool
	ForeignColumnUnique   bool
}

// ToManyRelationship describes a relationship between two tables where the
// local table has no id, and the foreign table has an id that matches a column
// in the local table.
type ToManyRelationship struct {
	Table    string
	Column   string
	Nullable bool
	Unique   bool

	ForeignTable          string
	ForeignColumn         string
	ForeignColumnNullable bool
	ForeignColumnUnique   bool

	ToJoinTable bool
	JoinTable   string

	JoinLocalColumn         string
	JoinLocalColumnNullable bool
	JoinLocalColumnUnique   bool

	JoinForeignColumn         string
	JoinForeignColumnNullable bool
	JoinForeignColumnUnique   bool
}

// ToOneRelationships relationship lookups
// Input should be the sql name of a table like: videos
func ToOneRelationships(table string, tables []Table) []ToOneRelationship {
	localTable := GetTable(tables, table)
	return toOneRelationships(localTable, tables)
}

// ToManyRelationships relationship lookups
// Input should be the sql name of a table like: videos
func ToManyRelationships(table string, tables []Table) []ToManyRelationship {
	localTable := GetTable(tables, table)
	return toManyRelationships(localTable, tables)
}

func toOneRelationships(table Table, tables []Table) []ToOneRelationship {
	var relationships []ToOneRelationship

	for _, t := range tables {
		for _, f := range t.FKeys {
			if f.ForeignTable == table.Name && !t.IsJoinTable && f.Unique {
				relationships = append(relationships, buildToOneRelationship(table, f, t, tables))
			}

		}
	}

	return relationships
}

func toManyRelationships(table Table, tables []Table) []ToManyRelationship {
	var relationships []ToManyRelationship

	for _, t := range tables {
		for _, f := range t.FKeys {
			if f.ForeignTable == table.Name && !f.Unique {
				relationships = append(relationships, buildToManyRelationship(table, f, t, tables))
			}
		}
	}

	return relationships
}

func buildToOneRelationship(localTable Table, foreignKey ForeignKey, foreignTable Table, tables []Table) ToOneRelationship {
	return ToOneRelationship{
		Table:    localTable.Name,
		Column:   foreignKey.ForeignColumn,
		Nullable: foreignKey.ForeignColumnNullable,
		Unique:   foreignKey.ForeignColumnUnique,

		ForeignTable:          foreignTable.Name,
		ForeignColumn:         foreignKey.Column,
		ForeignColumnNullable: foreignKey.Nullable,
		ForeignColumnUnique:   foreignKey.Unique,
	}
}

func buildToManyRelationship(localTable Table, foreignKey ForeignKey, foreignTable Table, tables []Table) ToManyRelationship {
	if !foreignTable.IsJoinTable {
		col := localTable.GetColumn(foreignKey.ForeignColumn)
		return ToManyRelationship{
			Table:                 localTable.Name,
			Column:                foreignKey.ForeignColumn,
			Nullable:              col.Nullable,
			Unique:                col.Unique,
			ForeignTable:          foreignTable.Name,
			ForeignColumn:         foreignKey.Column,
			ForeignColumnNullable: foreignKey.Nullable,
			ForeignColumnUnique:   foreignKey.Unique,
			ToJoinTable:           false,
		}
	}

	col := foreignTable.GetColumn(foreignKey.Column)
	relationship := ToManyRelationship{
		Table:       localTable.Name,
		Column:      foreignKey.ForeignColumn,
		Nullable:    col.Nullable,
		Unique:      col.Unique,
		ToJoinTable: true,
		JoinTable:   foreignTable.Name,
	}

	for _, fk := range foreignTable.FKeys {
		if fk.ForeignTable != localTable.Name {
			relationship.JoinForeignColumn = fk.Column
			relationship.JoinForeignColumnNullable = fk.Nullable
			relationship.JoinForeignColumnUnique = fk.Unique

			foreignTable := GetTable(tables, fk.ForeignTable)
			foreignCol := foreignTable.GetColumn(fk.ForeignColumn)
			relationship.ForeignTable = fk.ForeignTable
			relationship.ForeignColumn = fk.ForeignColumn
			relationship.ForeignColumnNullable = foreignCol.Nullable
			relationship.ForeignColumnUnique = foreignCol.Unique
		} else {
			relationship.JoinLocalColumn = fk.Column
			relationship.JoinLocalColumnNullable = fk.Nullable
			relationship.JoinLocalColumnUnique = fk.Unique
		}
	}

	return relationship
}
