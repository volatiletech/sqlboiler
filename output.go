package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"text/template"

	"github.com/pkg/errors"
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

	templates *templateList

	importSet      imports
	importNamedSet map[string]imports

	combineImportsOnType bool

	fileSuffix string
}

func executeTemplates(e executeTemplateData) error {
	if e.data.Table.IsJoinTable {
		return nil
	}

	var out [][]byte
	var imps imports

	imps.standard = e.importSet.standard
	imps.thirdParty = e.importSet.thirdParty

	for _, tplName := range e.templates.Templates() {
		if e.combineImportsOnType {
			imps = combineTypeImports(imps, importsBasedOnType, e.data.Table.Columns)
		}

		resp, err := executeTemplate(e.templates.Template, tplName, e.data)
		if err != nil {
			return fmt.Errorf("Error generating template %s: %s", tplName, err)
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
	if e.data.Table.IsJoinTable {
		return nil
	}

	rgxRemove := regexp.MustCompile(`[0-9]+_`)

	for _, tplName := range e.templates.Templates() {
		resp, err := executeTemplate(e.templates.Template, tplName, e.data)
		if err != nil {
			return fmt.Errorf("Error generating template %s: %s", tplName, err)
		}

		fName := tplName
		ext := filepath.Ext(fName)
		fName = rgxRemove.ReplaceAllString(fName[:len(fName)-len(ext)], "")

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

	resp, err := executeTemplate(state.TestMainTemplate, state.TestMainTemplate.Name(), data)
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

var rgxSyntaxError = regexp.MustCompile(`(\d+):\d+: `)

// executeTemplate takes a template and returns the output of the template
// execution.
func executeTemplate(t *template.Template, name string, data *templateData) ([]byte, error) {
	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, name, data); err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	output, err := format.Source(buf.Bytes())
	if err != nil {
		matches := rgxSyntaxError.FindStringSubmatch(err.Error())
		if matches == nil {
			return nil, errors.Wrap(err, "failed to format template")
		}

		lineNum, _ := strconv.Atoi(matches[1])
		scanner := bufio.NewScanner(&buf)
		errBuf := &bytes.Buffer{}
		line := 0
		for ; scanner.Scan(); line++ {
			if delta := line - lineNum; delta < -5 || delta > 5 {
				continue
			}

			if line == lineNum {
				errBuf.WriteString(">>> ")
			} else {
				fmt.Fprintf(errBuf, "% 3d ", line)
			}
			errBuf.Write(scanner.Bytes())
			errBuf.WriteByte('\n')
		}

		return nil, errors.Wrapf(err, "failed to format template\n\n%s\n", errBuf.Bytes())
	}

	return output, nil
}
