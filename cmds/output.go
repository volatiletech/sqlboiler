package cmds

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"os"
	"text/template"
)

var testHarnessStdout io.Writer = os.Stdout
var testHarnessFileOpen = func(filename string) (io.WriteCloser, error) {
	file, err := os.Create(filename)
	return file, err
}

func generateOutput(cmdData *CmdData, data *tplData, testOutput bool) error {
	if (testOutput && len(cmdData.TestTemplates) == 0) || (!testOutput && len(cmdData.Templates) == 0) {
		return errors.New("No template files located for generation")
	}

	var out [][]byte
	var imps imports
	var tpls []*template.Template

	if testOutput {
		imps.standard = sqlBoilerTestImports.standard
		imps.thirdparty = sqlBoilerTestImports.thirdparty
		imps = combineImports(imps, sqlBoilerDriverTestImports[cmdData.DriverName])
		tpls = cmdData.TestTemplates
	} else {
		imps.standard = sqlBoilerImports.standard
		imps.thirdparty = sqlBoilerImports.thirdparty
		tpls = cmdData.Templates
	}

	// Loop through and generate every individual template
	for _, template := range tpls {
		if !testOutput {
			imps = combineTypeImports(imps, sqlBoilerTypeImports, data.Table.Columns)
		}
		resp, err := generateTemplate(template, data)
		if err != nil {
			return err
		}
		out = append(out, resp)
	}

	if err := outHandler(cmdData, out, data, imps, testOutput); err != nil {
		return err
	}

	return nil
}

// outHandler loops over each template in the slice of byte slices and builds an output file.
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

// generateTemplate takes a template and returns the output of the template execution.
func generateTemplate(t *template.Template, data *tplData) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return nil, err
	}

	output, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}

	return output, nil
}
