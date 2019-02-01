package gradle_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/int128/gradleupdate/domain/git"
	"github.com/int128/gradleupdate/domain/gradle"
	"github.com/int128/gradleupdate/domain/testdata"
)

func TestFindWrapperVersion(t *testing.T) {
	version := gradle.FindWrapperVersion(testdata.GradleWrapperProperties4102)
	if want := "4.10.2"; version.String() != want {
		t.Errorf("version wants %s but %s", want, version)
	}
}

func TestReplaceWrapperVersion(t *testing.T) {
	replaced := gradle.ReplaceWrapperVersion(testdata.GradleWrapperProperties4102, "5.0")
	if replaced.String() != testdata.GradleWrapperProperties50.String() {
		t.Errorf("replaced wants %s but %s", testdata.GradleWrapperProperties50, replaced)
	}
}

func TestCheckUpdatePrecondition(t *testing.T) {
	for _, c := range []struct {
		name string
		in   gradle.UpdatePreconditionIn
		out  gradle.UpdatePreconditionOut
	}{
		{
			"OK",
			gradle.UpdatePreconditionIn{
				Readme:                  git.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &gradle.Release{Version: "5.1"},
			},
			gradle.ReadyToUpdate,
		}, {
			"NoReadmeBadge",
			gradle.UpdatePreconditionIn{
				Readme:                  git.FileContent("INVALID"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &gradle.Release{Version: "5.1"},
			},
			gradle.NoReadmeBadge,
		}, {
			"NoGradleVersion",
			gradle.UpdatePreconditionIn{
				Readme:                  git.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: git.FileContent("INVALID"),
				LatestGradleRelease:     &gradle.Release{Version: "5.1"},
			},
			gradle.NoGradleVersion,
		}, {
			"AlreadyHasLatestGradle",
			gradle.UpdatePreconditionIn{
				Readme:                  git.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &gradle.Release{Version: "5.0"},
			},
			gradle.AlreadyHasLatestGradle,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			out := gradle.CheckUpdatePrecondition(c.in)
			if diff := deep.Equal(c.out, out); diff != nil {
				t.Error(diff)
			}
		})
	}
}
