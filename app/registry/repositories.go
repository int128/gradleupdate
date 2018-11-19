package registry

import (
	"context"

	"github.com/int128/gradleupdate/app/domain/repositories"
	"github.com/int128/gradleupdate/app/infrastructure"
	impl "github.com/int128/gradleupdate/app/infrastructure/repositories"
)

type Repositories interface {
	Repository(ctx context.Context) repositories.Repository
}

type defaultRepositories struct{}

func (*defaultRepositories) Repository(ctx context.Context) repositories.Repository {
	return &impl.Repository{GitHubClient: infrastructure.GitHubClient(ctx)}
}

func NewRepositoriesRegistry() Repositories {
	return &defaultRepositories{}
}
