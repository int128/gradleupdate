package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

type RepositoryRepository interface {
	Get(context.Context, domain.RepositoryIdentifier) (*domain.Repository, error)
	GetFileContent(context.Context, domain.RepositoryIdentifier, string) (domain.FileContent, error)
	Fork(context.Context, domain.RepositoryIdentifier) (*domain.Repository, error)
}
