package drivers

import (
	"testing"

	"github.com/aarondl/sqlboiler/v4/importers"
)

type testRegistrationDriver struct{}

func (t testRegistrationDriver) Assemble(config Config) (*DBInfo, error) {
	return &DBInfo{
		Tables:  nil,
		Dialect: Dialect{},
	}, nil
}

func (t testRegistrationDriver) Templates() (map[string]string, error) {
	return nil, nil
}

func (t testRegistrationDriver) Imports() (importers.Collection, error) {
	return importers.Collection{}, nil
}

func TestRegistration(t *testing.T) {
	mock := testRegistrationDriver{}
	RegisterFromInit("mock1", mock)

	if d, ok := registeredDrivers["mock1"]; !ok {
		t.Error("driver was not found")
	} else if d != mock {
		t.Error("got the wrong driver back")
	}
}

func TestBinaryRegistration(t *testing.T) {
	RegisterBinary("mock2", "/bin/true")

	if d, ok := registeredDrivers["mock2"]; !ok {
		t.Error("driver was not found")
	} else if string(d.(binaryDriver)) != "/bin/true" {
		t.Error("got the wrong driver back")
	}
}

func TestBinaryFromArgRegistration(t *testing.T) {
	RegisterBinaryFromCmdArg("/bin/true/mock5")

	if d, ok := registeredDrivers["mock5"]; !ok {
		t.Error("driver was not found")
	} else if string(d.(binaryDriver)) != "/bin/true/mock5" {
		t.Error("got the wrong driver back")
	}
}

func TestGetDriver(t *testing.T) {
	didYouPanic := false

	RegisterBinary("mock4", "/bin/true")

	func() {
		defer func() {
			if r := recover(); r != nil {
				didYouPanic = true
			}
		}()

		_ = GetDriver("mock4")
	}()

	if didYouPanic {
		t.Error("expected not to panic when fetching a driver that's known")
	}

	func() {
		defer func() {
			if r := recover(); r != nil {
				didYouPanic = true
			}
		}()

		_ = GetDriver("notpresentdriver")
	}()

	if !didYouPanic {
		t.Error("expected to recover from a panic")
	}
}

func TestReregister(t *testing.T) {
	didYouPanic := false

	func() {
		defer func() {
			if r := recover(); r != nil {
				didYouPanic = true
			}
		}()

		RegisterBinary("mock3", "/bin/true")
		RegisterBinary("mock3", "/bin/true")
	}()

	if !didYouPanic {
		t.Error("expected to recover from a panic")
	}
}
