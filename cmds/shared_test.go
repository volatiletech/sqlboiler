package cmds

import (
	"bytes"
	"io"
	"testing"
)

func TestOutHandler(t *testing.T) {
	buf := &bytes.Buffer{}

	saveTestHarnessStdout := testHarnessStdout
	testHarnessStdout = buf
	defer func() {
		testHarnessStdout = saveTestHarnessStdout
	}()

	data := tplData{
		Table: "patrick",
	}

	templateOutputs := [][]byte{[]byte("hello world"), []byte("patrick's dreams")}

	if err := outHandler("", templateOutputs, &data); err != nil {
		t.Error(err)
	}

	if out := buf.String(); out != "hello world\npatrick's dreams\n" {
		t.Errorf("Wrong output: %q", out)
	}
}

type NopWriteCloser struct {
	io.Writer
}

func (NopWriteCloser) Close() error {
	return nil
}

func nopCloser(w io.Writer) io.WriteCloser {
	return NopWriteCloser{w}
}

func TestOutHandlerFiles(t *testing.T) {
	saveTestHarnessFileOpen := testHarnessFileOpen
	defer func() {
		testHarnessFileOpen = saveTestHarnessFileOpen
	}()

	file := &bytes.Buffer{}
	testHarnessFileOpen = func(path string) (io.WriteCloser, error) {
		return nopCloser(file), nil
	}

	data := tplData{
		Table: "patrick",
	}

	templateOutputs := [][]byte{[]byte("hello world"), []byte("patrick's dreams")}

	if err := outHandler("folder", templateOutputs, &data); err != nil {
		t.Error(err)
	}

	if out := file.String(); out != "hello world\npatrick's dreams\n" {
		t.Errorf("Wrong output: %q", out)
	}
}
