package drivers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// DriverMain helps dry up the implementation of main.go for drivers
func DriverMain(driver Interface) {
	b, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to read from stdin")
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(b, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse json from stdin: %v\n%s\n", err, b)
		os.Exit(1)
	}

	dinfo, err := driver.Assemble(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	b, err = json.Marshal(dinfo)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to marshal json:", err)
		os.Exit(1)
	}

	os.Stdout.Write(b)
}
