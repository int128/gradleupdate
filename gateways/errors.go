package gateways

import "github.com/int128/gradleupdate/gateways/interfaces"

type repositoryError struct {
	error
	noSuchEntity  bool
	alreadyExists bool
}

func (err *repositoryError) NoSuchEntity() bool  { return err.noSuchEntity }
func (err *repositoryError) AlreadyExists() bool { return err.alreadyExists }

var _ gateways.RepositoryError = &repositoryError{}
