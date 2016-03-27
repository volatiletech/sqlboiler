package cmds

import (
	"fmt"
	"io"
	"os"
)

var testHarnessStdout io.Writer = os.Stdout
var testHarnessFileOpen = func(filename string) (io.WriteCloser, error) {
	file, err := os.Create(filename)
	return file, err
}

// outHandler loops over the slice of byte slices, outputting them to either
// the OutFile if it is specified with a flag, or to Stdout if no flag is specified.
func outHandler(cmdData *CmdData, output [][]byte, data *tplData, imps imports, testTemplate bool) error {
	out := testHarnessStdout

	var path string
	if len(cmdData.OutFolder) != 0 {
		if testTemplate {
			path = cmdData.OutFolder + "/" + data.Table.Name + "_test.go"
		} else {
			path = cmdData.OutFolder + "/" + data.Table.Name + ".go"
		}

		outFile, err := testHarnessFileOpen(path)
		if err != nil {
			return fmt.Errorf("Unable to create output file %s: %s", path, err)
		}
		defer outFile.Close()
		out = outFile
	}

	if _, err := fmt.Fprintf(out, "package %s\n\n", cmdData.PkgName); err != nil {
		return fmt.Errorf("Unable to write package name %s to file: %s", cmdData.PkgName, path)
	}

	impStr := buildImportString(imps)
	if len(impStr) > 0 {
		if _, err := fmt.Fprintf(out, "%s\n", impStr); err != nil {
			return fmt.Errorf("Unable to write imports to file handle: %v", err)
		}
	}

	for _, templateOutput := range output {
		if _, err := fmt.Fprintf(out, "%s\n", templateOutput); err != nil {
			return fmt.Errorf("Unable to write template output to file handle: %v", err)
		}
	}

	return nil
}

func combineStringSlices(a, b []string) []string {
	c := make([]string, len(a)+len(b))
	if len(a) > 0 {
		copy(c, a)
	}
	if len(b) > 0 {
		copy(c[len(a):], b)
	}

	return c
}

func removeDuplicates(dedup []string) []string {
	if len(dedup) <= 1 {
		return dedup
	}

	for i := 0; i < len(dedup)-1; i++ {
		for j := i + 1; j < len(dedup); j++ {
			if dedup[i] != dedup[j] {
				continue
			}

			if j != len(dedup)-1 {
				dedup[j] = dedup[len(dedup)-1]
				j--
			}
			dedup = dedup[:len(dedup)-1]
		}
	}

	return dedup
}
