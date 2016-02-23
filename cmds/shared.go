package cmds

import (
	"fmt"
	"os"
	"strings"
)

func errorQuit(err error) {
	fmt.Println(fmt.Sprintf("Error: %s\n---\n", err))
	structCmd.Help()
	os.Exit(-1)
}

func makeGoColName(name string) string {
	s := strings.Split(name, "_")

	for i := 0; i < len(s); i++ {
		if s[i] == "id" {
			s[i] = "ID"
			continue
		}
		s[i] = strings.Title(s[i])
	}

	return strings.Join(s, "")
}

func makeDBColName(tableName, colName string) string {
	return tableName + "_" + colName
}
