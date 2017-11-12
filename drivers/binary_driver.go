package drivers

import (
	"bytes"
	"encoding/json"
	"os/exec"

	"github.com/pkg/errors"
)

type binaryDriver string

// Assemble calls out to the binary with JSON
// The contract for error messages is that a plain text error message is delivered
// and the exit status of the process is non-zero
func (b binaryDriver) Assemble(config map[string]interface{}) (*DBInfo, error) {
	configuration, err := json.Marshal(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to json-ify driver configuration")
	}

	executable := string(b)

	output := &bytes.Buffer{}
	cmd := exec.Command(executable)
	cmd.Stdout = output
	cmd.Stderr = output
	cmd.Stdin = bytes.NewReader(configuration)
	err = cmd.Run()

	if err != nil {
		if cmd.ProcessState.Exited() && !cmd.ProcessState.Success() {
			return nil, errors.Wrapf(err, "driver (%s) returned an error:\n%s", executable, output.Bytes())
		}

		return nil, errors.Wrapf(err, "something totally unexpected happened when running the binary driver %s", executable)
	}

	var info DBInfo
	if err = json.Unmarshal(output.Bytes(), &info); err != nil {
		return nil, errors.Wrap(err, "failed to marshal json from binary")
	}

	return &info, nil
}
