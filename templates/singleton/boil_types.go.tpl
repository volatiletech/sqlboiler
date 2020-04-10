import (
	"fmt"
	"strings"
	"github.com/razor-1/sqlboiler/v3/queries/qm"
)

// M type is for providing columns and column values to UpdateAll.
type M map[string]interface{}

// ErrSyncFail occurs during insert when the record could not be retrieved in
// order to populate default value information. This usually happens when LastInsertId
// fails or there was a primary key configuration that was not resolvable.
var ErrSyncFail = errors.New("{{.PkgName}}: failed to synchronize data after insert")

type insertCache struct {
	query        string
	retQuery     string
	valueMapping []uint64
	retMapping   []uint64
}

type updateCache struct {
	query        string
	valueMapping []uint64
}

type fkRelationship struct {
	Table  string
	Column string
}

type relationMap map[string]fkRelationship

func makeCacheKey(cols boil.Columns, nzDefaults []string) string {
	buf := strmangle.GetBuffer()

	buf.WriteString(strconv.Itoa(cols.Kind))
	for _, w := range cols.Cols {
		buf.WriteString(w)
	}

	if len(nzDefaults) != 0 {
		buf.WriteByte('.')
	}
	for _, nz := range nzDefaults {
		buf.WriteString(nz)
	}

	str := buf.String()
	strmangle.PutBuffer(buf)
	return str
}

// bind needs the fully qualified column names to handle possible duplicate column names in different tables
// this is a a helper that takes a slice of column names and a table name and creates a query that makes
// them all unique, in the form of `table`.`column` AS "table.column"
func FullyQualifiedColumns(tableName string) string {
	columns, ok := TableNameToTableColumns[tableName]
	if !ok {
		return ""
	}
	components := make([]string, len(columns))
	for i, column := range columns {
		components[i] = fmt.Sprintf("`%s`.`%s` AS \"%s.%s\"", tableName, column, tableName, column)
	}

	return strings.Join(components, ",")
}

func GetJoinClause(targetTable, targetColumn, sourceTable, sourceColumn string) string {
    return fmt.Sprintf("`%s` ON `%s`.`%s`=`%s`.`%s`", targetTable, sourceTable, sourceColumn, targetTable, targetColumn)
}

// construct a simple join clause
func JoinClause(sourceTable string, sourceColumn string, relationMap relationMap) (clause string, err error) {
	relation, ok := relationMap[sourceColumn]
	if !ok {
		err = errors.New("Cannot find source column in relationMap")
		return
	}

	clause = GetJoinClause(relation.Table, relation.Column, sourceTable, sourceColumn)
	return
}

// given the name of a type, find which column in our table has a foreign key relationship
// this is done by finding the table name and then looking for it in the relationmap
func getSourceColumn(typeName string, rColumns relationMap) (sourceColumn string, err error) {
	table, ok := TypeNameToTableName[typeName]
	if !ok {
		err = errors.New("No table name for that type name")
		return
	}

	for sc, relation := range rColumns {
		if relation.Table == table {
			sourceColumn = sc
			return
		}
	}
	err = errors.New("No source column found for that type name")
	return
}


// given a type that the relationship points to, build a query mod
// it's an inner join if the relationship isn't nullable, otherwise it's a left outer join
func loadJoinQueryMod(sourceTable, typeName string, rColumns relationMap, nullable bool) qm.QueryMod {
	sourceColumn, err := getSourceColumn(typeName, rColumns)
	if err != nil {
		return nil
	}

	clause, err := JoinClause(sourceTable, sourceColumn, rColumns)
	if err != nil {
		return nil
	}

	if nullable {
		return qm.LeftOuterJoin(clause)
	} else {
		return qm.InnerJoin(clause)
	}
}

{{/*
The following is a little bit of black magic and deserves some explanation

Because postgres and mysql define enums completely differently (one at the
database level as a custom datatype, and one at the table column level as
a unique thing per table)... There's a chance the enum is named (postgres)
and not (mysql). So we can't do this per table so this code is here.

We loop through each table and column looking for enums. If it's named, we
then use some disgusting magic to write state during the template compile to
the "once" map. This lets named enums only be defined once if they're referenced
multiple times in many (or even the same) tables.

Then we check if all it's values are normal, if they are we create the enum
output, if not we output a friendly error message as a comment to aid in
debugging.

Postgres output looks like: EnumNameEnumValue = "enumvalue"
MySQL output looks like:    TableNameColNameEnumValue = "enumvalue"

It only titlecases the EnumValue portion if it's snake-cased.
*/}}
{{$once := onceNew}}
{{- range $table := .Tables -}}
	{{- range $col := $table.Columns | filterColumnsByEnum -}}
		{{- $name := parseEnumName $col.DBType -}}
		{{- $vals := parseEnumVals $col.DBType -}}
		{{- $isNamed := ne (len $name) 0}}
		{{- if and $isNamed (onceHas $once $name) -}}
		{{- else -}}
			{{- if $isNamed -}}
				{{$_ := oncePut $once $name}}
			{{- end -}}
{{- if and (gt (len $vals) 0) (isEnumNormal $vals)}}
// Enum values for {{if $isNamed}}{{$name}}{{else}}{{$table.Name}}.{{$col.Name}}{{end}}
const (
	{{- range $val := $vals -}}
		{{- $valStripped := stripWhitespace $val -}}
	{{- if $isNamed}}{{titleCase $name}}{{else}}{{titleCase $table.Name}}{{titleCase $col.Name}}{{end -}}
	{{if shouldTitleCaseEnum $valStripped}}{{titleCase $valStripped}}{{else}}{{$valStripped}}{{end}} = "{{$val}}"
	{{end -}}
)
{{- else}}
// Enum values for {{if $isNamed}}{{$name}}{{else}}{{$table.Name}}.{{$col.Name}}{{end}} are not proper Go identifiers, cannot emit constants
{{- end -}}
		{{- end -}}
	{{- end -}}
{{- end -}}
