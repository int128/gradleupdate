package di

import (
	impl "github.com/int128/gradleupdate/usecases"
	"github.com/int128/gradleupdate/usecases/interfaces"
)

var usecasesDependencies = []interface{}{
	func(i impl.GetRepository) usecases.GetRepository { return &i },
	func(i impl.GetBadge) usecases.GetBadge { return &i },
	func(i impl.SendUpdate) usecases.SendUpdate { return &i },
	func(i impl.BatchSendUpdates) usecases.BatchSendUpdates { return &i },
	func(i impl.SendPullRequest) usecases.SendPullRequest { return &i },
}
