package strmangle

import (
	"fmt"
	"strings"

	"github.com/nullbio/sqlboiler/dbdrivers"
)

func RandDBStruct(varName string, table dbdrivers.Table) string {
	return ""
}

func RandDBStructSlice(varName string, num int, table dbdrivers.Table) string {
	var structs []string
	for i := 0; i < num; i++ {
		structs = append(structs, RandDBStruct(varName, table))
	}

	innerStructs := strings.Join(structs, ",")
	return fmt.Sprintf("%s := %s{%s}", varName, TitleCasePlural(table.Name), innerStructs)
}
