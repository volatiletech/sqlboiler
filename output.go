package sqlboiler

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
	if len(state.Templates) == 0 {
		return errors.New("No template files located for generation")
	}
	var out [][]byte
	var imps imports

	imps.standard = defaultTemplateImports.standard
	imps.thirdParty = defaultTemplateImports.thirdParty

	for _, template := range state.Templates {
		imps = combineTypeImports(imps, importsBasedOnType, data.Table.Columns)
		resp, err := generateTemplate(template, data)
		if err != nil {
			return fmt.Errorf("Error generating template %s: %s", template.Name(), err)
		}
		out = append(out, resp)
	}

	fName := data.Table.Name + ".go"
	err := outHandler(state.Config.OutFolder, fName, state.Config.PkgName, imps, out)
	if err != nil {
		return err
	}

	return nil
}

// generateTestOutput builds the test file output and sends it to outHandler for saving
func generateTestOutput(state *State, data *templateData) error {
	if len(state.TestTemplates) == 0 {
		return errors.New("No template files located for generation")
	}

	var out [][]byte
	var imps imports

	imps.standard = defaultTestTemplateImports.standard
	imps.thirdParty = defaultTestTemplateImports.thirdParty

	for _, template := range state.TestTemplates {
		resp, err := generateTemplate(template, data)
		if err != nil {
			return fmt.Errorf("Error generating test template %s: %s", template.Name(), err)
		}
		out = append(out, resp)
	}

	fName := data.Table.Name + "_test.go"
	err := outHandler(state.Config.OutFolder, fName, state.Config.PkgName, imps, out)
	if err != nil {
		return err
	}

	return nil
}

// generateSingletonOutput processes the templates that should only be run
// one time.
func generateSingletonOutput(state *State) error {
	if state.SingletonTemplates == nil {
		return errors.New("No singleton templates located for generation")
	}

	templateData := &templateData{
		PkgName:    state.Config.PkgName,
		DriverName: state.Config.DriverName,
	}

	for _, template := range state.SingletonTemplates {
		var imps imports

		resp, err := generateTemplate(template, templateData)
		if err != nil {
			return fmt.Errorf("Error generating template %s: %s", template.Name(), err)
		}

		fName := template.Name()
		ext := filepath.Ext(fName)
		fName = fName[0 : len(fName)-len(ext)]

		imps.standard = defaultSingletonTemplateImports[fName].standard
		imps.thirdParty = defaultSingletonTemplateImports[fName].thirdParty

		fName = fName + ".go"

		err = outHandler(state.Config.OutFolder, fName, state.Config.PkgName, imps, [][]byte{resp})
		if err != nil {
			return err
		}
	}

	return nil
}

// generateSingletonTestOutput processes the templates that should only be run
// one time.
func generateSingletonTestOutput(state *State) error {
	if state.SingletonTestTemplates == nil {
		return errors.New("No singleton test templates located for generation")
	}

	templateData := &templateData{
		PkgName:    state.Config.PkgName,
		DriverName: state.Config.DriverName,
	}

	for _, template := range state.SingletonTestTemplates {
		var imps imports

		resp, err := generateTemplate(template, templateData)
		if err != nil {
			return fmt.Errorf("Error generating test template %s: %s", template.Name(), err)
		}

		fName := template.Name()
		ext := filepath.Ext(fName)
		fName = fName[0 : len(fName)-len(ext)]

		imps.standard = defaultSingletonTestTemplateImports[fName].standard
		imps.thirdParty = defaultSingletonTestTemplateImports[fName].thirdParty

		fName = fName + "_test.go"

		err = outHandler(state.Config.OutFolder, fName, state.Config.PkgName, imps, [][]byte{resp})
		if err != nil {
			return err
		}
	}

	return nil
}

func generateTestMainOutput(state *State) error {
	if state.TestMainTemplate == nil {
		return errors.New("No TestMain template located for generation")
	}

	var out [][]byte
	var imps imports

	imps.standard = defaultTestMainImports[state.Config.DriverName].standard
	imps.thirdParty = defaultTestMainImports[state.Config.DriverName].thirdParty

	templateData := &templateData{
		Tables:     state.Tables,
		PkgName:    state.Config.PkgName,
		DriverName: state.Config.DriverName,
	}

	resp, err := generateTemplate(state.TestMainTemplate, templateData)
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

// generateTemplate takes a template and returns the output of the template execution.
func generateTemplate(t *template.Template, data *templateData) ([]byte, error) {
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
