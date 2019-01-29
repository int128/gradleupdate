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
	func(i impl.RepositoryLastScanRepository) gateways.RepositoryLastScanRepository { return &i },
	func(i impl.GradleService) gateways.GradleService { return &i },
	func(i impl.HTTPCacheRepository) gateways.HTTPCacheRepository { return &i },
	func(i impl.AELogger) gateways.Logger { return &i },
	impl.NewConfigRepository,
}
