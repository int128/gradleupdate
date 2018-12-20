package repositories

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

type Repository interface {
	Get(context.Context, domain.RepositoryIdentifier) (domain.Repository, error)
	GetFile(context.Context, domain.RepositoryIdentifier, string) (domain.File, error)
	Fork(context.Context, domain.RepositoryIdentifier) (domain.Repository, error)
}
