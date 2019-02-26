package usecases

import (
	"fmt"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/pkg/errors"
)

type noSuchRepositoryError struct {
	repository git.RepositoryID
}

func (err *noSuchRepositoryError) Error() string {
	return fmt.Sprintf("no such repository %s", err.repository)
}

type noSuchRepositoryErrorCauser struct{}

func (e *noSuchRepositoryErrorCauser) IsNoSuchRepositoryError(err error) bool {
	_, ok := errors.Cause(err).(*noSuchRepositoryError)
	return ok
}
