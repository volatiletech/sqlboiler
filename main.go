package main

import (
	"fmt"
	"os"

	"github.com/pobri19/sqlboiler/cmds"
)

func main() {
	if err := cmds.SQLBoiler.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
