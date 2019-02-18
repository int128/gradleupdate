package gradleupdate

import (
	"fmt"

	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
)

func BranchFor(owner string, version gradle.Version) git.BranchName {
	return git.BranchName(fmt.Sprintf("gradle-%s-%s", version, owner))
}
