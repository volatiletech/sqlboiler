package boilingcore

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	"github.com/razor-1/sqlboiler/importers"
)

var (
	noEditDisclaimerFmt = `// Code generated by SQLBoiler%s(https://github.com/razor-1/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

`
	noEditDisclaimer = []byte(fmt.Sprintf(noEditDisclaimerFmt, " "))
)

var (
	// templateByteBuffer is re-used by all template construction to avoid
	// allocating more memory than is needed. This will later be a problem for
	// concurrency, address it then.
	templateByteBuffer = &bytes.Buffer{}

	rgxRemoveNumberedPrefix = regexp.MustCompile(`^[0-9]+_`)
	rgxSyntaxError          = regexp.MustCompile(`(\d+):\d+: `)

	testHarnessWriteFile = ioutil.WriteFile
)

type executeTemplateData struct {
	state *State
	data  *templateData

	templates     *templateList
	dirExtensions dirExtMap

	importSet      importers.Set
	importNamedSet importers.Map

	combineImportsOnType bool
	isTest               bool
}

// generateOutput builds the file output and sends it to outHandler for saving
func generateOutput(state *State, dirExts dirExtMap, data *templateData) error {
	return executeTemplates(executeTemplateData{
		state:                state,
		data:                 data,
		templates:            state.Templates,
		importSet:            state.Config.Imports.All,
		combineImportsOnType: true,
		dirExtensions:        dirExts,
	})
}

// generateTestOutput builds the test file output and sends it to outHandler for saving
func generateTestOutput(state *State, dirExts dirExtMap, data *templateData) error {
	return executeTemplates(executeTemplateData{
		state:                state,
		data:                 data,
		templates:            state.TestTemplates,
		importSet:            state.Config.Imports.Test,
		combineImportsOnType: false,
		isTest:               true,
		dirExtensions:        dirExts,
	})
}

// generateSingletonOutput processes the templates that should only be run
// one time.
func generateSingletonOutput(state *State, data *templateData) error {
	return executeSingletonTemplates(executeTemplateData{
		state:          state,
		data:           data,
		templates:      state.Templates,
		importNamedSet: state.Config.Imports.Singleton,
	})
}

// generateSingletonTestOutput processes the templates that should only be run
// one time.
func generateSingletonTestOutput(state *State, data *templateData) error {
	return executeSingletonTemplates(executeTemplateData{
		state:          state,
		data:           data,
		templates:      state.TestTemplates,
		importNamedSet: state.Config.Imports.TestSingleton,
		isTest:         true,
	})
}

func executeTemplates(e executeTemplateData) error {
	if e.data.Table.IsJoinTable {
		return nil
	}

	var imps importers.Set
	imps.Standard = e.importSet.Standard
	imps.ThirdParty = e.importSet.ThirdParty
	if e.combineImportsOnType {
		colTypes := make([]string, len(e.data.Table.Columns))
		nullColTypes := make([]string, len(e.data.Table.Columns))
		for i, ct := range e.data.Table.Columns {
			colTypes[i] = ct.Type
			nullColTypes[i] = ct.NullType
		}
		colTypes = append(colTypes, nullColTypes...)

		imps = importers.AddTypeImports(imps, e.state.Config.Imports.BasedOnType, colTypes)
	}

	for dir, dirExts := range e.dirExtensions {
		for ext, tplNames := range dirExts {
			out := templateByteBuffer
			out.Reset()

			isGo := filepath.Ext(ext) == ".go"
			if isGo {
				pkgName := e.state.Config.PkgName
				if len(dir) != 0 {
					pkgName = filepath.Base(dir)
				}
				writeFileDisclaimer(out)
				writePackageName(out, pkgName)
				writeImports(out, imps)
			}

			for _, tplName := range tplNames {
				if err := executeTemplate(out, e.templates.Template, tplName, e.data); err != nil {
					return err
				}
			}

			fName := e.data.Table.Name
			if e.isTest {
				fName += "_test"
			}
			fName += ext
			if len(dir) != 0 {
				fName = filepath.Join(dir, fName)
			}

			if err := writeFile(e.state.Config.OutFolder, fName, out, isGo); err != nil {
				return err
			}
		}
	}

	return nil
}

func executeSingletonTemplates(e executeTemplateData) error {
	if e.data.Table.IsJoinTable {
		return nil
	}

	out := templateByteBuffer
	for _, tplName := range e.templates.Templates() {
		normalized, isSingleton, isGo, usePkg := outputFilenameParts(tplName)
		if !isSingleton {
			continue
		}

		dir, fName := filepath.Split(normalized)
		fName = fName[:strings.IndexByte(fName, '.')]

		out.Reset()

		if isGo {
			imps := importers.Set{
				Standard:   e.importNamedSet[denormalizeSlashes(fName)].Standard,
				ThirdParty: e.importNamedSet[denormalizeSlashes(fName)].ThirdParty,
			}

			pkgName := e.state.Config.PkgName
			if !usePkg {
				pkgName = filepath.Base(dir)
			}
			writeFileDisclaimer(out)
			writePackageName(out, pkgName)
			writeImports(out, imps)
		}

		if err := executeTemplate(out, e.templates.Template, tplName, e.data); err != nil {
			return err
		}

		if err := writeFile(e.state.Config.OutFolder, normalized, out, isGo); err != nil {
			return err
		}
	}

	return nil
}

// writeFileDisclaimer writes the disclaimer at the top with a trailing
// newline so the package name doesn't get attached to it.
func writeFileDisclaimer(out *bytes.Buffer) {
	_, _ = out.Write(noEditDisclaimer)
}

// writePackageName writes the package name correctly, ignores errors
// since it's to the concrete buffer type which produces none
func writePackageName(out *bytes.Buffer, pkgName string) {
	_, _ = fmt.Fprintf(out, "package %s\n\n", pkgName)
}

// writeImports writes the package imports correctly, ignores errors
// since it's to the concrete buffer type which produces none
func writeImports(out *bytes.Buffer, imps importers.Set) {
	if impStr := imps.Format(); len(impStr) > 0 {
		_, _ = fmt.Fprintf(out, "%s\n", impStr)
	}
}

// writeFile writes to the given folder and filename, formatting the buffer
// given.
func writeFile(outFolder string, fileName string, input *bytes.Buffer, format bool) error {
	var byt []byte
	var err error
	if format {
		byt, err = formatBuffer(input)
		if err != nil {
			return err
		}
	} else {
		byt = input.Bytes()
	}

	path := filepath.Join(outFolder, fileName)
	if err := testHarnessWriteFile(path, byt, 0664); err != nil {
		return errors.Wrapf(err, "failed to write output file %s", path)
	}

	return nil
}

// executeTemplate takes a template and returns the output of the template
// execution.
func executeTemplate(buf *bytes.Buffer, t *template.Template, name string, data *templateData) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("failed to execute template: %s\npanic: %+v\n", name, r)
		}
	}()

	if err := t.ExecuteTemplate(buf, name, data); err != nil {
		return errors.Wrapf(err, "failed to execute template: %s", name)
	}
	return nil
}

func formatBuffer(buf *bytes.Buffer) ([]byte, error) {
	output, err := format.Source(buf.Bytes())
	if err == nil {
		return output, nil
	}

	matches := rgxSyntaxError.FindStringSubmatch(err.Error())
	if matches == nil {
		return nil, errors.Wrap(err, "failed to format template")
	}

	lineNum, _ := strconv.Atoi(matches[1])
	scanner := bufio.NewScanner(buf)
	errBuf := &bytes.Buffer{}
	line := 1
	for ; scanner.Scan(); line++ {
		if delta := line - lineNum; delta < -5 || delta > 5 {
			continue
		}

		if line == lineNum {
			errBuf.WriteString(">>>> ")
		} else {
			fmt.Fprintf(errBuf, "% 4d ", line)
		}
		errBuf.Write(scanner.Bytes())
		errBuf.WriteByte('\n')
	}

	return nil, errors.Wrapf(err, "failed to format template\n\n%s\n", errBuf.Bytes())
}

func getLongExt(filename string) string {
	index := strings.IndexByte(filename, '.')
	return filename[index:]
}

// fileFragments will take something of the form:
// templates/singleton/hello.go.tpl
// templates_test/js/hello.js.tpl
func outputFilenameParts(filename string) (normalized string, isSingleton, isGo, usePkg bool) {
	fragments := strings.Split(filename, string(os.PathSeparator))
	isSingleton = fragments[len(fragments)-2] == "singleton"

	var remainingFragments []string
	for _, f := range fragments[1:] {
		if f != "singleton" {
			remainingFragments = append(remainingFragments, f)
		}
	}

	newFilename := remainingFragments[len(remainingFragments)-1]
	newFilename = strings.TrimSuffix(newFilename, ".tpl")
	newFilename = rgxRemoveNumberedPrefix.ReplaceAllString(newFilename, "")
	remainingFragments[len(remainingFragments)-1] = newFilename

	ext := filepath.Ext(newFilename)
	isGo = ext == ".go"

	usePkg = len(remainingFragments) == 1
	normalized = strings.Join(remainingFragments, string(os.PathSeparator))

	return normalized, isSingleton, isGo, usePkg
}
