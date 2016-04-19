package cmds

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

var testHarnessStdout io.Writer = os.Stdout
var testHarnessFileOpen = func(filename string) (io.WriteCloser, error) {
	file, err := os.Create(filename)
	return file, err
}

// generateOutput builds the file output and sends it to outHandler for saving
func generateOutput(cmdData *CmdData, data *tplData) error {
	if len(cmdData.Templates) == 0 {
		return errors.New("No template files located for generation")
	}

	var out [][]byte
	var imps imports

	imps.standard = sqlBoilerImports.standard
	imps.thirdparty = sqlBoilerImports.thirdparty

	for _, template := range cmdData.Templates {
		imps = combineTypeImports(imps, sqlBoilerTypeImports, data.Table.Columns)
		resp, err := generateTemplate(template, data)
		if err != nil {
			return fmt.Errorf("Error generating template %s: %s", template.Name(), err)
		}
		out = append(out, resp)
	}

	fName := data.Table.Name + ".go"
	err := outHandler(cmdData.OutFolder, fName, cmdData.PkgName, imps, out)
	if err != nil {
		return err
	}

	return nil
}

// generateTestOutput builds the test file output and sends it to outHandler for saving
func generateTestOutput(cmdData *CmdData, data *tplData) error {
	if len(cmdData.TestTemplates) == 0 {
		return errors.New("No template files located for generation")
	}

	var out [][]byte
	var imps imports

	imps.standard = sqlBoilerTestImports.standard
	imps.thirdparty = sqlBoilerTestImports.thirdparty

	for _, template := range cmdData.TestTemplates {
		resp, err := generateTemplate(template, data)
		if err != nil {
			return fmt.Errorf("Error generating test template %s: %s", template.Name(), err)
		}
		out = append(out, resp)
	}

	fName := data.Table.Name + "_test.go"
	err := outHandler(cmdData.OutFolder, fName, cmdData.PkgName, imps, out)
	if err != nil {
		return err
	}

	return nil
}

func generateSinglesOutput(cmdData *CmdData) error {
	if cmdData.SingleTemplates == nil {
		return errors.New("No single templates located for generation")
	}

	for _, template := range cmdData.SingleTemplates {
		var imps imports

		resp, err := generateTemplate(template, &tplData{})
		if err != nil {
			return fmt.Errorf("Error generating template %s: %s", template.Name(), err)
		}

		fName := template.Name()
		ext := filepath.Ext(fName)
		fName = fName[0 : len(fName)-len(ext)]

		imps.standard = sqlBoilerSinglesImports[fName].standard
		imps.thirdparty = sqlBoilerSinglesImports[fName].thirdparty

		fName = fName + ".go"

		err = outHandler(cmdData.OutFolder, fName, cmdData.PkgName, imps, [][]byte{resp})
		if err != nil {
			return err
		}
	}

	return nil
}

func generateTestMainOutput(cmdData *CmdData) error {
	if cmdData.TestMainTemplate == nil {
		return errors.New("No TestMain template located for generation")
	}

	var out [][]byte
	var imps imports

	imps.standard = sqlBoilerTestMainImports[cmdData.DriverName].standard
	imps.thirdparty = sqlBoilerTestMainImports[cmdData.DriverName].thirdparty

	resp, err := generateTemplate(cmdData.TestMainTemplate, &tplData{})
	if err != nil {
		return err
	}
	out = append(out, resp)

	err = outHandler(cmdData.OutFolder, "main_test.go", cmdData.PkgName, imps, out)
	if err != nil {
		return err
	}

	return nil
}

func outHandler(outFolder string, fileName string, pkgName string, imps imports, contents [][]byte) error {
	out := testHarnessStdout

	path := filepath.Join(outFolder, fileName)

	outFile, err := testHarnessFileOpen(path)
	if err != nil {
		return fmt.Errorf("Unable to create output file %s: %s", path, err)
	}
	defer outFile.Close()
	out = outFile

	if _, err := fmt.Fprintf(out, "package %s\n\n", pkgName); err != nil {
		return fmt.Errorf("Unable to write package name %s to file: %s", pkgName, path)
	}

	impStr := buildImportString(imps)
	if len(impStr) > 0 {
		if _, err := fmt.Fprintf(out, "%s\n", impStr); err != nil {
			return fmt.Errorf("Unable to write imports to file handle: %v", err)
		}
	}

	for _, templateOutput := range contents {
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
