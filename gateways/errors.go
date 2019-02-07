package gateways

type repositoryError struct {
	error
	noSuchEntity  bool
	alreadyExists bool
}

func (err *repositoryError) NoSuchEntity() bool  { return err.noSuchEntity }
func (err *repositoryError) AlreadyExists() bool { return err.alreadyExists }
