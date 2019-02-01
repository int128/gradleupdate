package gradleupdate_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/gradleupdate"
	"github.com/int128/gradleupdate/domain/testdata"
)

func TestCheckUpdatePrecondition(t *testing.T) {
	for _, c := range []struct {
		name                  string
		precondition          gradleupdate.Precondition
		preconditionViolation gradleupdate.PreconditionViolation
	}{
		{
			"OK",
			gradleupdate.Precondition{
				Readme:                  git.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &gradle.Release{Version: "5.1"},
			},
			gradleupdate.ReadyToUpdate,
		}, {
			"NoReadmeBadge",
			gradleupdate.Precondition{
				Readme:                  git.FileContent("INVALID"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &gradle.Release{Version: "5.1"},
			},
			gradleupdate.NoReadmeBadge,
		}, {
			"NoGradleVersion",
			gradleupdate.Precondition{
				Readme:                  git.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: git.FileContent("INVALID"),
				LatestGradleRelease:     &gradle.Release{Version: "5.1"},
			},
			gradleupdate.NoGradleVersion,
		}, {
			"AlreadyHasLatestGradle",
			gradleupdate.Precondition{
				Readme:                  git.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &gradle.Release{Version: "5.0"},
			},
			gradleupdate.AlreadyHasLatestGradle,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			out := gradleupdate.CheckPrecondition(c.precondition)
			if diff := deep.Equal(c.preconditionViolation, out); diff != nil {
				t.Error(diff)
			}
		})
	}
}
