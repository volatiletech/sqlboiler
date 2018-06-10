package boilingcore

import (
	"fmt"
	"strings"

	"github.com/volatiletech/sqlboiler/drivers"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// txtNameToOne creates the local and foreign function names for
// one-to-many and one-to-one relationships, where local is the side with
// the foreign key.
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
func txtNameToOne(fk drivers.ForeignKey) (localFn, foreignFn string) {
	foreignFn = strmangle.Singular(trimSuffixes(fk.Column))
	fkeyNotTableName := foreignFn != strmangle.Singular(fk.ForeignTable)
	foreignFn = strmangle.TitleCase(foreignFn)

	if fkeyNotTableName {
		localFn = foreignFn
	}

	plurality := strmangle.Plural
	if fk.Unique {
		plurality = strmangle.Singular
	}
	localFn += strmangle.TitleCase(plurality(fk.Table))

	return localFn, foreignFn
}

// txtNameToMany creates the local and foreign function names for
// many-to-many relationships, where local refers to the table
// who would have the foreign key if it weren't for the join table.
//
// That's to say: If we had tags and videos, ordinarily if it were a
// one to many, tags would have the video_id, and so for the video_id
// fkey local means "tags", and foreign means "videos".
//
//   | tags |  | tags_videos      |  | videos |
//   | id   |  | tag_id, video_id |  | id     |
//
// In this setup, if we were able to not have a join table, it would look
// like this:
//
//   | tags     |  | videos |
//   | id       |  | id     |
//   | video_id |  | tag_id |
//
// Hence when looking from the perspective of the "video_id" foreign key
// local = tags, foreign = videos.
//
// cases:
// = many-to-many
// sponsors - contests
// sponsor_id contest_id
// fk == table = sponsor.Contests | contest.Sponsors
//
// = many-to-many
// sponsors - contests
// wiggle_id jiggle_id
// fk != table = sponsor.JiggleSponsors | contest.WiggleContests
//
// = many-to-many
// industries - industries
// industry_id  mapped_industry_id
//
// fk == table = industry.Industries
// fk != table = industry.MappedIndustryIndustry
func txtNameToMany(toMany drivers.ToManyRelationship) (localFn, foreignFn string) {
	if !toMany.ToJoinTable {
		panic(fmt.Sprintf("this method is only for join tables: %s <-> %s, %s", toMany.Table, toMany.ForeignTable, toMany.Name))
	}

	localFkey := strmangle.Singular(trimSuffixes(toMany.JoinLocalColumn))
	foreignFkey := strmangle.Singular(trimSuffixes(toMany.JoinForeignColumn))

	if localFkey != strmangle.Singular(toMany.Table) {
		localFn = strmangle.TitleCase(localFkey)
	}
	localFn += strmangle.TitleCase(strmangle.Plural(toMany.Table))

	if foreignFkey != strmangle.Singular(toMany.ForeignTable) {
		foreignFn = strmangle.TitleCase(foreignFkey)
	}
	foreignFn += strmangle.TitleCase(strmangle.Plural(toMany.ForeignTable))

	return localFn, foreignFn
}

func usesPrimitives(tables []drivers.Table, table, column, foreignTable, foreignColumn string) bool {
	local := drivers.GetTable(tables, table)
	foreign := drivers.GetTable(tables, foreignTable)

	col := local.GetColumn(column)
	foreignCol := foreign.GetColumn(foreignColumn)

	return isPrimitive(col.Type) && isPrimitive(foreignCol.Type)
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

func isPrimitive(typ string) bool {
	switch typ {
	// Numeric
	case "int", "int8", "int16", "int32", "int64":
		return true
	case "uint", "uint8", "uint16", "uint32", "uint64":
		return true
	case "float32", "float64":
		return true
	case "byte", "rune", "string":
		return true
	}

	return false
}
