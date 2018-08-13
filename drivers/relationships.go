package drivers

// ToOneRelationship describes a relationship between two tables where the local
// table has no id, and the foreign table has an id that matches a column in the
// local table, that column can also be unique which changes the dynamic into a
// one-to-one style, not a to-many.
type ToOneRelationship struct {
	Name string `json:"name"`

	Table    string `json:"table"`
	Column   string `json:"column"`
	Nullable bool   `json:"nullable"`
	Unique   bool   `json:"unique"`

	ForeignTable          string `json:"foreign_table"`
	ForeignColumn         string `json:"foreign_column"`
	ForeignColumnNullable bool   `json:"foreign_column_nullable"`
	ForeignColumnUnique   bool   `json:"foreign_column_unique"`
}

// ToManyRelationship describes a relationship between two tables where the
// local table has no id, and the foreign table has an id that matches a column
// in the local table.
type ToManyRelationship struct {
	Name string `json:"name"`

	Table    string `json:"table"`
	Column   string `json:"column"`
	Nullable bool   `json:"nullable"`
	Unique   bool   `json:"unique"`

	ForeignTable          string `json:"foreign_table"`
	ForeignColumn         string `json:"foreign_column"`
	ForeignColumnNullable bool   `json:"foreign_column_nullable"`
	ForeignColumnUnique   bool   `json:"foreign_column_unique"`

	ToJoinTable bool   `json:"to_join_table"`
	JoinTable   string `json:"join_table"`

	JoinLocalFKeyName       string `json:"join_local_fkey_name"`
	JoinLocalColumn         string `json:"join_local_column"`
	JoinLocalColumnNullable bool   `json:"join_local_column_nullable"`
	JoinLocalColumnUnique   bool   `json:"join_local_column_unique"`

	JoinForeignFKeyName       string `json:"join_foreign_fkey_name"`
	JoinForeignColumn         string `json:"join_foreign_column"`
	JoinForeignColumnNullable bool   `json:"join_foreign_column_nullable"`
	JoinForeignColumnUnique   bool   `json:"join_foreign_column_unique"`
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
			if f.ForeignTable == table.Name && (t.IsJoinTable || !f.Unique) {
				relationships = append(relationships, buildToManyRelationship(table, f, t, tables))
			}
		}
	}

	return relationships
}

func buildToOneRelationship(localTable Table, foreignKey ForeignKey, foreignTable Table, tables []Table) ToOneRelationship {
	return ToOneRelationship{
		Name:     foreignKey.Name,
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
		return ToManyRelationship{
			Name:                  foreignKey.Name,
			Table:                 localTable.Name,
			Column:                foreignKey.ForeignColumn,
			Nullable:              foreignKey.ForeignColumnNullable,
			Unique:                foreignKey.ForeignColumnUnique,
			ForeignTable:          foreignTable.Name,
			ForeignColumn:         foreignKey.Column,
			ForeignColumnNullable: foreignKey.Nullable,
			ForeignColumnUnique:   foreignKey.Unique,
			ToJoinTable:           false,
		}
	}

	relationship := ToManyRelationship{
		Table:    localTable.Name,
		Column:   foreignKey.ForeignColumn,
		Nullable: foreignKey.ForeignColumnNullable,
		Unique:   foreignKey.ForeignColumnUnique,

		ToJoinTable: true,
		JoinTable:   foreignTable.Name,

		JoinLocalFKeyName:       foreignKey.Name,
		JoinLocalColumn:         foreignKey.Column,
		JoinLocalColumnNullable: foreignKey.Nullable,
		JoinLocalColumnUnique:   foreignKey.Unique,
	}

	for _, fk := range foreignTable.FKeys {
		if fk.Name == foreignKey.Name {
			continue
		}

		relationship.JoinForeignFKeyName = fk.Name
		relationship.JoinForeignColumn = fk.Column
		relationship.JoinForeignColumnNullable = fk.Nullable
		relationship.JoinForeignColumnUnique = fk.Unique

		relationship.ForeignTable = fk.ForeignTable
		relationship.ForeignColumn = fk.ForeignColumn
		relationship.ForeignColumnNullable = fk.ForeignColumnNullable
		relationship.ForeignColumnUnique = fk.ForeignColumnUnique
	}

	return relationship
}
