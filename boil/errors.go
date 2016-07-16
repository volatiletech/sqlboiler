package boil

type boilErr struct {
	error
}

func WrapErr(err error) error {
	return boilErr{
		error: err,
	}
}

func (e boilErr) Error() string {
	return e.error.Error()
}

func IsBoilErr(err error) bool {
	_, ok := err.(boilErr)
	return ok
}
