package gateways

type repositoryError struct {
	error
	noSuchEntity bool
}

func (err *repositoryError) NoSuchEntity() bool { return err.noSuchEntity }
