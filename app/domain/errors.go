package domain

type NotFoundError struct {
	Cause error
}

func (err NotFoundError) Error() string {
	return err.Cause.Error()
}

func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}
