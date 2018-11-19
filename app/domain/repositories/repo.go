package repositories

import (
	"context"

	"github.com/int128/gradleupdate/app/domain"
)

type Repository interface {
	GetFile(context.Context, domain.RepositoryIdentifier, string) (*domain.File, error)
}
