package boilingcore

import (
	"fmt"

	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Aliases defines aliases for the generation run
type Aliases struct {
	Tables        map[string]TableAlias        `toml:"tables,omitempty" json:"tables,omitempty"`
	Relationships map[string]RelationshipAlias `toml:"relationships,omitempty" json:"relationships,omitempty"`
}

// TableAlias defines the spellings for a table name in Go
type TableAlias struct {
	UpPlural     string `toml:"up_plural,omitempty" json:"up_plural,omitempty"`
	UpSingular   string `toml:"up_singular,omitempty" json:"up_singular,omitempty"`
	DownPlural   string `toml:"down_plural,omitempty" json:"down_plural,omitempty"`
	DownSingular string `toml:"down_singular,omitempty" json:"down_singular,omitempty"`

	Columns map[string]string `toml:"columns,omitempty" json:"columns,omitempty"`
}

// RelationshipAlias defines the naming for both sides of
// a foreign key.
type RelationshipAlias struct {
	Local   string `toml:"local,omitempty" json:"local,omitempty"`
	Foreign string `toml:"foreign,omitempty" json:"foreign,omitempty"`
}

// FillAliases takes the table information from the driver
// and fills in aliases where the user has provided none.
//
// This leaves us with a complete list of Go names for all tables,
// columns, and relationships.
func FillAliases(a *Aliases, tables []drivers.Table) {
	for _, t := range tables {
		if t.IsJoinTable {
			continue
		}

		if a.Tables == nil {
			a.Tables = make(map[string]TableAlias)
		}

		table := a.Tables[t.Name]

		if len(table.UpPlural) == 0 {
			table.UpPlural = strmangle.TitleCase(strmangle.Plural(t.Name))
		}
		if len(table.UpSingular) == 0 {
			table.UpSingular = strmangle.TitleCase(strmangle.Singular(t.Name))
		}
		if len(table.DownPlural) == 0 {
			table.DownPlural = strmangle.CamelCase(strmangle.Plural(t.Name))
		}
		if len(table.DownSingular) == 0 {
			table.DownSingular = strmangle.CamelCase(strmangle.Singular(t.Name))
		}

		if table.Columns == nil {
			table.Columns = make(map[string]string)
		}

		for _, c := range t.Columns {
			if _, ok := table.Columns[c.Name]; !ok {
				table.Columns[c.Name] = strmangle.TitleCase(c.Name)
			}
		}

		a.Tables[t.Name] = table

		if a.Relationships == nil {
			a.Relationships = make(map[string]RelationshipAlias)
		}

		for _, k := range t.FKeys {
			r := a.Relationships[k.Name]
			if len(r.Local) != 0 && len(r.Foreign) != 0 {
				continue
			}

			local, foreign := txtNameToOne(k)
			if len(r.Local) == 0 {
				r.Local = local
			}
			if len(r.Foreign) == 0 {
				r.Foreign = foreign
			}

			a.Relationships[k.Name] = r
		}

		for _, rel := range t.ToManyRelationships {
			localFacingAlias, okLocal := a.Relationships[rel.JoinLocalFKeyName]
			foreignFacingAlias, okForeign := a.Relationships[rel.JoinForeignFKeyName]

			if okLocal && okForeign {
				continue
			}

			local, foreign := txtNameToMany(rel)

			switch {
			case !okLocal && !okForeign:
				localFacingAlias.Local = local
				localFacingAlias.Foreign = foreign
				foreignFacingAlias.Local = foreign
				foreignFacingAlias.Foreign = local
			case okLocal:
				if len(localFacingAlias.Local) == 0 {
					localFacingAlias.Local = local
				}
				if len(localFacingAlias.Foreign) == 0 {
					localFacingAlias.Foreign = foreign
				}

				foreignFacingAlias.Local = localFacingAlias.Foreign
				foreignFacingAlias.Foreign = localFacingAlias.Local
			case okForeign:
				if len(foreignFacingAlias.Local) == 0 {
					foreignFacingAlias.Local = foreign
				}
				if len(foreignFacingAlias.Foreign) == 0 {
					foreignFacingAlias.Foreign = local
				}

				localFacingAlias.Foreign = foreignFacingAlias.Local
				localFacingAlias.Local = foreignFacingAlias.Foreign
			}

			a.Relationships[rel.JoinLocalFKeyName] = localFacingAlias
			a.Relationships[rel.JoinForeignFKeyName] = foreignFacingAlias
		}
	}
}

// Table gets a table alias, panics if not found.
func (a Aliases) Table(table string) TableAlias {
	t, ok := a.Tables[table]
	if !ok {
		panic("could not find table aliases for: " + table)
	}

	return t
}

// Column gets a column's aliased name, panics if not found.
func (a Aliases) Column(table, column string) string {
	t := a.Table(table)

	c, ok := t.Columns[column]
	if !ok {
		panic(fmt.Sprintf("could not find column alias for: %s.%s", table, column))
	}

	return c
}
