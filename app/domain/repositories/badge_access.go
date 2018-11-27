package repositories

import (
	"context"
	"github.com/int128/gradleupdate/app/domain"
)

type BadgeLastAccess interface {
	Get(context.Context, domain.RepositoryIdentifier) (domain.BadgeLastAccess, error)
	Put(context.Context, domain.BadgeLastAccess) error
}
