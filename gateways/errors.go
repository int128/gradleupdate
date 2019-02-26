package gateways

import (
	"github.com/pkg/errors"
)

type noSuchEntityError struct {
	error
}

type noSuchEntityErrorCauser struct{}

func (e *noSuchEntityErrorCauser) IsNoSuchEntityError(err error) bool {
	_, ok := errors.Cause(err).(*noSuchEntityError)
	return ok
}

type entityAlreadyExistsError struct {
	error
}

type entityAlreadyExistsErrorCauser struct{}

func (e *entityAlreadyExistsErrorCauser) IsEntityAlreadyExistsError(err error) bool {
	_, ok := errors.Cause(err).(*entityAlreadyExistsError)
	return ok
}
