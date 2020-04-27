package drivers

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/sqlboiler/v4/importers"
)

type binaryDriver string

// Assemble calls out to the binary with JSON
// The contract for error messages is that a plain text error message is delivered
// and the exit status of the process is non-zero
func (b binaryDriver) Assemble(config Config) (*DBInfo, error) {
	var dbInfo DBInfo
	err := execute(string(b), "assemble", config, &dbInfo, os.Stderr)
	if err != nil {
		return nil, err
	}

	return &dbInfo, nil
}

// Templates calls the templates function to get a map of overidden file names
// and their contents in base64
func (b binaryDriver) Templates() (map[string]string, error) {
	var templates map[string]string
	err := execute(string(b), "templates", nil, &templates, os.Stderr)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

// Imports calls the imports function to get imports from the driver
func (b binaryDriver) Imports() (col importers.Collection, err error) {
	err = execute(string(b), "imports", nil, &col, os.Stderr)
	if err != nil {
		return col, err
	}

	return col, nil
}

func execute(executable, method string, input interface{}, output interface{}, errStream io.Writer) error {
	var err error
	var inputBytes []byte
	if input != nil {
		inputBytes, err = json.Marshal(input)
		if err != nil {
			return errors.Wrap(err, "failed to json-ify driver configuration")
		}
	}

	outputBytes := &bytes.Buffer{}
	cmd := exec.Command(executable, method)
	cmd.Stdout = outputBytes
	cmd.Stderr = errStream
	if inputBytes != nil {
		cmd.Stdin = bytes.NewReader(inputBytes)
	}
	err = cmd.Run()

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			if ee.ProcessState.Exited() && !ee.ProcessState.Success() {
				return errors.Wrapf(err, "driver (%s) exited non-zero", executable)
			}
		}

		return errors.Wrapf(err, "something totally unexpected happened when running the binary driver %s", executable)
	}

	if err = json.Unmarshal(outputBytes.Bytes(), &output); err != nil {
		return errors.Wrap(err, "failed to marshal json from binary")
	}

	return nil
}
