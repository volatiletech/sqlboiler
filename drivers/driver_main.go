package drivers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// DriverMain helps dry up the implementation of main.go for drivers
func DriverMain(driver Interface) {
	method := os.Args[1]
	var config Config

	switch method {
	case "assemble":
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to read from stdin")
			os.Exit(1)
		}

		err = json.Unmarshal(b, &config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse json from stdin: %v\n%s\n", err, b)
			os.Exit(1)
		}
	case "templates":
		// No input for this method
	case "imports":
		// No input for this method
	}

	var output interface{}
	switch method {
	case "assemble":
		dinfo, err := driver.Assemble(config)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		output = dinfo
	case "templates":
		templates, err := driver.Templates()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		output = templates
	case "imports":
		collection, err := driver.Imports()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		output = collection
	}

	b, err := json.Marshal(output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to marshal json:", err)
		os.Exit(1)
	}

	os.Stdout.Write(b)
}
