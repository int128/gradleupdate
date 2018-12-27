package gateways

import (
	"context"

	"github.com/int128/gradleupdate/domain"
)

type BadgeLastAccessRepository interface {
	Get(context.Context, domain.RepositoryIdentifier) (*domain.BadgeLastAccess, error)
	Put(context.Context, domain.BadgeLastAccess) error
}
