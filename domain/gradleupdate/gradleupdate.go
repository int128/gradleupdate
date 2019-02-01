package gradleupdate

import (
	"bytes"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
)

type Precondition struct {
	Readme                  git.FileContent
	BadgeURL                string
	GradleWrapperProperties git.FileContent
	LatestGradleRelease     *gradle.Release
}

type PreconditionViolation int

const (
	ReadyToUpdate             = PreconditionViolation(0)
	AlreadyHasLatestGradle    = PreconditionViolation(1)
	NoGradleWrapperProperties = PreconditionViolation(51)
	NoGradleVersion           = PreconditionViolation(52)
	NoReadme                  = PreconditionViolation(53)
	NoReadmeBadge             = PreconditionViolation(54)
)

func CheckPrecondition(precondition Precondition) PreconditionViolation {
	if precondition.GradleWrapperProperties == nil {
		return NoGradleWrapperProperties
	}
	currentGradleVersion := gradle.FindWrapperVersion(precondition.GradleWrapperProperties)
	if currentGradleVersion == "" {
		return NoGradleVersion
	}
	if precondition.Readme == nil {
		return NoReadme
	}
	if !bytes.Contains(precondition.Readme, []byte(precondition.BadgeURL)) {
		return NoReadmeBadge
	}
	if currentGradleVersion.GreaterOrEqualThan(precondition.LatestGradleRelease.Version) {
		return AlreadyHasLatestGradle
	}
	return ReadyToUpdate
}
