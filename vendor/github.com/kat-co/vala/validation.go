/*
Package vala is a simple, extensible, library to make argument
validation in Go palatable.

This package uses the fluent programming style to provide
simultaneously more robust and more terse parameter validation.

	BeginValidation().Validate(
		IsNotNil(a, "a"),
		IsNotNil(b, "b"),
		IsNotNil(c, "c"),
	).CheckAndPanic().Validate( // Panic will occur here if a, b, or c are nil.
		HasLen(a.Items, 50, "a.Items"),
		GreaterThan(b.UserCount, 0, "b.UserCount"),
		Equals(c.Name, "Vala", "c.name"),
		Not(Equals(c.FriendlyName, "Foo", "c.FriendlyName")),
	).Check()

Notice how checks can be tiered.

Vala is also extensible. As long as a function conforms to the Checker
specification, you can pass it into the Validate method:

	func ReportFitsRepository(report *Report, repository *Repository) Checker {
		return func() (passes bool, err error) {

			err = fmt.Errorf("A %s report does not belong in a %s repository.", report.Type, repository.Type)
			passes = (repository.Type == report.Type)
			return passes, err
		}
	}

	func AuthorCanUpload(authorName string, repository *Repository) Checker {
		return func() (passes bool, err error) {
			err = fmt.Errorf("%s does not have access to this repository.", authorName)
			passes = !repository.AuthorCanUpload(authorName)
			return passes, err
		}
	}

	func AuthorIsCollaborator(authorName string, report *Report) Checker {
		return func() (passes bool, err error) {

			err = fmt.Errorf("The given author was not one of the collaborators for this report.")
			for _, collaboratorName := range report.Collaborators() {
				if collaboratorName == authorName {
					passes = true
					break
				}
			}

			return passes, err
		}
	}

	func HandleReport(authorName string, report *Report, repository *Repository) {

		BeginValidation().Validate(
			AuthorIsCollaborator(authorName, report),
			AuthorCanUpload(authorName, repository),
			ReportFitsRepository(report, repository),
		).CheckAndPanic()
	}
*/
package vala

import (
	"fmt"
	"reflect"
	"strings"
)

// Validation contains all the errors from performing Checkers, and is
// the fluent type off which all Validation methods hang.
type Validation struct {
	Errors []string
}

// BeginValidation begins a validation check.
func BeginValidation() *Validation {
	return nil
}

// Check aggregates all checker errors into a single error and returns
// this error.
func (val *Validation) Check() error {
	if val == nil || len(val.Errors) <= 0 {
		return nil
	}

	return val.constructErrorMessage()
}

// CheckAndPanic aggregates all checker errors into a single error and
// panics with this error.
func (val *Validation) CheckAndPanic() *Validation {
	if val == nil || len(val.Errors) <= 0 {
		return val
	}

	panic(val.constructErrorMessage())
}

// CheckSetErrorAndPanic aggregates any Errors produced by the
// Checkers into a single error, and sets the address of retError to
// this, and panics. The canonical use-case of this is to pass in the
// address of an error you would like to return, and then to catch the
// panic and do nothing.
func (val *Validation) CheckSetErrorAndPanic(retError *error) *Validation {
	if val == nil || len(val.Errors) <= 0 {
		return val
	}

	*retError = val.constructErrorMessage()
	panic(*retError)
}

// Validate runs all of the checkers passed in and collects errors
// into an internal collection. To take action on these errors, call
// one of the Check* methods.
func (val *Validation) Validate(checkers ...Checker) *Validation {

	for _, checker := range checkers {
		if pass, msg := checker(); !pass {
			if val == nil {
				val = &Validation{}
			}

			val.Errors = append(val.Errors, msg)
		}
	}

	return val
}

func (val *Validation) constructErrorMessage() error {
	return fmt.Errorf(
		"parameter validation failed:\t%s",
		strings.Join(val.Errors, "\n\t"),
	)
}

//
// Checker functions
//

// Checker defines the type of function which can represent a Vala
// checker.  If the Checker fails, returns false with a corresponding
// error message. If the Checker succeeds, returns true, but _also_
// returns an error message. This helps to support the Not function.
type Checker func() (checkerIsTrue bool, errorMessage string)

// Not returns the inverse of any Checker passed in.
func Not(checker Checker) Checker {

	return func() (passed bool, errorMessage string) {
		if passed, errorMessage = checker(); passed {
			return false, fmt.Sprintf("Not(%s)", errorMessage)
		}

		return true, ""
	}
}

// Equals performs a basic == on the given parameters and fails if
// they are not equal.
func Equals(param, value interface{}, paramName string) Checker {

	return func() (pass bool, errMsg string) {
		return (param == value), fmt.Sprintf("Parameters were not equal: %s(%v) != %v",
			paramName,
			param,
			value)
	}
}

// IsNotNil checks to see if the value passed in is nil. This Checker
// attempts to check the most performant things first, and then
// degrade into the less-performant, but accurate checks for nil.
func IsNotNil(obtained interface{}, paramName string) Checker {
	return func() (isNotNil bool, errMsg string) {

		if obtained == nil {
			isNotNil = false
		} else if str, ok := obtained.(string); ok {
			isNotNil = str != ""
		} else {
			switch v := reflect.ValueOf(obtained); v.Kind() {
			case
				reflect.Chan,
				reflect.Func,
				reflect.Interface,
				reflect.Map,
				reflect.Ptr,
				reflect.Slice:
				isNotNil = !v.IsNil()
			default:
				panic("Vala is unable to check this type for nilability at this time.")
			}
		}

		return isNotNil, "Parameter was nil: " + paramName
	}
}

// HasLen checks to ensure the given argument is the desired length.
func HasLen(param interface{}, desiredLength int, paramName string) Checker {

	return func() (hasLen bool, errMsg string) {
		hasLen = desiredLength == reflect.ValueOf(param).Len()
		return hasLen, "Parameter did not contain the correct number of elements: " + paramName
	}
}

// GreaterThan checks to ensure the given argument is greater than the
// given value.
func GreaterThan(param int, comparativeVal int, paramName string) Checker {

	return func() (isGreaterThan bool, errMsg string) {
		if isGreaterThan = param > comparativeVal; !isGreaterThan {
			errMsg = fmt.Sprintf(
				"Parameter's length was not greater than:  %s(%d) < %d",
				paramName,
				param,
				comparativeVal)
		}

		return isGreaterThan, errMsg
	}
}

// StringNotEmpty checks to ensure the given string is not empty.
func StringNotEmpty(obtained, paramName string) Checker {
	return func() (isNotEmpty bool, errMsg string) {
		isNotEmpty = obtained != ""
		errMsg = fmt.Sprintf("Parameter is an empty string: %s", paramName)
		return
	}
}
