package gatewaysTestDoubles

import "github.com/int128/gradleupdate/gateways/interfaces"

type NoSuchEntityError struct{}

func (err *NoSuchEntityError) Error() string       { return "NoSuchEntityError" }
func (err *NoSuchEntityError) NoSuchEntity() bool  { return true }
func (err *NoSuchEntityError) AlreadyExists() bool { return false }

var _ gateways.RepositoryError = &NoSuchEntityError{}
