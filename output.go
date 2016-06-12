package main

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
func generateOutput(state *State, data *templateData) error {
	return executeTemplates(executeTemplateData{
		state:                state,
		data:                 data,
		templates:            state.Templates,
		importSet:            defaultTemplateImports,
		combineImportsOnType: true,
		fileSuffix:           ".go",
	})
}

// generateTestOutput builds the test file output and sends it to outHandler for saving
func generateTestOutput(state *State, data *templateData) error {
	return executeTemplates(executeTemplateData{
		state:                state,
		data:                 data,
		templates:            state.TestTemplates,
		importSet:            defaultTestTemplateImports,
		combineImportsOnType: false,
		fileSuffix:           "_test.go",
	})
}

// generateSingletonOutput processes the templates that should only be run
// one time.
func generateSingletonOutput(state *State, data *templateData) error {
	return executeSingletonTemplates(executeTemplateData{
		state:          state,
		data:           data,
		templates:      state.SingletonTemplates,
		importNamedSet: defaultSingletonTemplateImports,
		fileSuffix:     ".go",
	})
}

// generateSingletonTestOutput processes the templates that should only be run
// one time.
func generateSingletonTestOutput(state *State, data *templateData) error {
	return executeSingletonTemplates(executeTemplateData{
		state:          state,
		data:           data,
		templates:      state.SingletonTestTemplates,
		importNamedSet: defaultSingletonTestTemplateImports,
		fileSuffix:     "_test.go",
	})
}

type executeTemplateData struct {
	state *State
	data  *templateData

	templates templateList

	importSet      imports
	importNamedSet map[string]imports

	combineImportsOnType bool

	fileSuffix string
}

func executeTemplates(e executeTemplateData) error {
	var out [][]byte
	var imps imports

	imps.standard = e.importSet.standard
	imps.thirdParty = e.importSet.thirdParty

	for _, template := range e.templates {
		if e.combineImportsOnType {
			imps = combineTypeImports(imps, importsBasedOnType, e.data.Table.Columns)
		}

		resp, err := executeTemplate(template, e.data)
		if err != nil {
			return fmt.Errorf("Error generating template %s: %s", template.Name(), err)
		}
		out = append(out, resp)
	}

	fName := e.data.Table.Name + e.fileSuffix
	err := outHandler(e.state.Config.OutFolder, fName, e.state.Config.PkgName, imps, out)
	if err != nil {
		return err
	}

	return nil
}

func executeSingletonTemplates(e executeTemplateData) error {
	for _, template := range e.templates {
		resp, err := executeTemplate(template, e.data)
		if err != nil {
			return fmt.Errorf("Error generating template %s: %s", template.Name(), err)
		}

		fName := template.Name()
		ext := filepath.Ext(fName)
		fName = fName[0 : len(fName)-len(ext)]

		imps := imports{
			standard:   e.importNamedSet[fName].standard,
			thirdParty: e.importNamedSet[fName].thirdParty,
		}

		err = outHandler(
			e.state.Config.OutFolder,
			fName+e.fileSuffix,
			e.state.Config.PkgName,
			imps,
			[][]byte{resp},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func generateTestMainOutput(state *State, data *templateData) error {
	if state.TestMainTemplate == nil {
		return errors.New("No TestMain template located for generation")
	}

	var out [][]byte
	var imps imports

	imps.standard = defaultTestMainImports[state.Config.DriverName].standard
	imps.thirdParty = defaultTestMainImports[state.Config.DriverName].thirdParty

	resp, err := executeTemplate(state.TestMainTemplate, data)
	if err != nil {
		return err
	}
	out = append(out, resp)

	err = outHandler(state.Config.OutFolder, "main_test.go", state.Config.PkgName, imps, out)
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

// executeTemplate takes a template and returns the output of the template
// execution.
func executeTemplate(t *template.Template, data *templateData) ([]byte, error) {
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
