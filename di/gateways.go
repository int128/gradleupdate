package di

import (
	impl "github.com/int128/gradleupdate/gateways"
	"github.com/int128/gradleupdate/gateways/interfaces"
)

var gatewaysDependencies = []interface{}{
	func(i impl.RepositoryRepository) gateways.RepositoryRepository { return &i },
	func(i impl.PullRequestRepository) gateways.PullRequestRepository { return &i },
	func(i impl.GitService) gateways.GitService { return &i },
	func(i impl.BadgeLastAccessRepository) gateways.BadgeLastAccessRepository { return &i },
	func(i impl.RepositoryLastUpdateRepository) gateways.RepositoryLastUpdateRepository { return &i },
	func(i impl.GradleReleaseRepository) gateways.GradleReleaseRepository { return &i },
	impl.NewToggles,
	impl.NewCredentials,
	func(i impl.Time) gateways.Time { return &i },
	func(i impl.HTTPCacheRepository) gateways.HTTPCacheRepository { return &i },
	func(i impl.AELogger) gateways.Logger { return &i },
}
