package drivers

import "testing"

type testRegistrationDriver struct{}

func (t testRegistrationDriver) Assemble(config map[string]interface{}) (*DBInfo, error) {
	return &DBInfo{
		Tables:  nil,
		Dialect: Dialect{},
	}, nil
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
