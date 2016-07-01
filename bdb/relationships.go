package bdb

// ToManyRelationship describes a relationship between two tables where the
// local table has no id, and the foreign table has an id that matches a column
// in the local table.
type ToManyRelationship struct {
	Column        string
	ForeignTable  string
	ForeignColumn string

	ToJoinTable       bool
	JoinTable         string
	JoinLocalColumn   string
	JoinForeignColumn string
}

// ToManyRelationships relationship lookups
// Input should be the sql name of a table like: videos
func ToManyRelationships(table string, tables []Table) []ToManyRelationship {
	var relationships []ToManyRelationship

	for _, t := range tables {
		if t.Name == table {
			continue
		}

		for _, f := range t.FKeys {
			if f.ForeignTable != table {
				continue
			}

			relationships = append(relationships, buildRelationship(table, f, t))
		}
	}

	return relationships
}

func buildRelationship(localTable string, foreignKey ForeignKey, foreignTable Table) ToManyRelationship {
	if !foreignTable.IsJoinTable {
		return ToManyRelationship{
			Column:        foreignKey.ForeignColumn,
			ForeignTable:  foreignTable.Name,
			ForeignColumn: foreignKey.Column,
			ToJoinTable:   foreignTable.IsJoinTable,
		}
	}

	relationship := ToManyRelationship{
		Column:      foreignKey.ForeignColumn,
		ToJoinTable: true,
		JoinTable:   foreignTable.Name,
	}

	for _, fk := range foreignTable.FKeys {
		if fk.ForeignTable != localTable {
			relationship.JoinForeignColumn = fk.Column
			relationship.ForeignTable = fk.ForeignTable
			relationship.ForeignColumn = fk.ForeignColumn
		} else {
			relationship.JoinLocalColumn = fk.Column
		}
	}

	return relationship
}
