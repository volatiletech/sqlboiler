package dbdrivers

import (
	"fmt"

	"github.com/nullbio/sqlboiler/strmangle"
)

// ToManyRelationship describes a relationship between two tables where the
// local table has no id, and the foreign table has an id that matches a column
// in the local table.
type ToManyRelationship struct {
	Name          string
	ForeignTable  string
	ForeignColumn string
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

			singularName := strmangle.Singular(table)
			standardColName := fmt.Sprintf("%s_id", singularName)

			relationship := ToManyRelationship{
				ForeignTable:  t.Name,
				ForeignColumn: f.Column,
			}

			if standardColName == f.ForeignColumn {
				relationship.Name = strmangle.TitleCase(strmangle.Plural(name))
			} else {
				relationship.Name = strmangle.TitleCase(strmangle.Plural(name))
			}

			relationships = append(relationships, relationship)
		}
	}

	return relationships
}
