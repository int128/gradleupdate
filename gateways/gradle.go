package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/infrastructure"
	"github.com/pkg/errors"
)

type GradleService struct {
	GradleClient *infrastructure.GradleClient
}

func (s *GradleService) GetCurrentVersion(ctx context.Context) (domain.GradleVersion, error) {
	cvr, err := s.GradleClient.GetCurrentVersion(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "error while getting current version")
	}
	return domain.GradleVersion(cvr.Version), nil
}
