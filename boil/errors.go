package boil

type boilErr struct {
	error
}

// WrapErr wraps err in a boilErr
func WrapErr(err error) error {
	return boilErr{
		error: err,
	}
}

// Error returns the underlying error string
func (e boilErr) Error() string {
	return e.error.Error()
}

// IsBoilErr checks if err is a boilErr
func IsBoilErr(err error) bool {
	_, ok := err.(boilErr)
	return ok
}
