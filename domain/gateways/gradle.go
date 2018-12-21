package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

//go:generate mockgen -destination mock_gateways/gradle.go github.com/int128/gradleupdate/domain/gateways GradleService

type GradleService interface {
	GetCurrentVersion(ctx context.Context) (domain.GradleVersion, error)
}
