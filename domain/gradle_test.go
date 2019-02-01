package domain_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/int128/gradleupdate/domain"
	"github.com/int128/gradleupdate/domain/testdata"
)

func TestFindGradleWrapperVersion(t *testing.T) {
	version := domain.FindGradleWrapperVersion(testdata.GradleWrapperProperties4102)
	if want := "4.10.2"; version.String() != want {
		t.Errorf("version wants %s but %s", want, version)
	}
}

func TestReplaceGradleWrapperVersion(t *testing.T) {
	replaced := domain.ReplaceGradleWrapperVersion(testdata.GradleWrapperProperties4102, "5.0")
	if replaced.String() != testdata.GradleWrapperProperties50.String() {
		t.Errorf("replaced wants %s but %s", testdata.GradleWrapperProperties50, replaced)
	}
}

func TestCheckGradleUpdatePrecondition(t *testing.T) {
	for _, c := range []struct {
		name string
		in   domain.GradleUpdatePreconditionIn
		out  domain.GradleUpdatePreconditionOut
	}{
		{
			"OK",
			domain.GradleUpdatePreconditionIn{
				Readme:                  domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &domain.GradleRelease{Version: "5.1"},
			},
			domain.ReadyToUpdate,
		}, {
			"NoReadmeBadge",
			domain.GradleUpdatePreconditionIn{
				Readme:                  domain.FileContent("INVALID"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &domain.GradleRelease{Version: "5.1"},
			},
			domain.NoReadmeBadge,
		}, {
			"NoGradleVersion",
			domain.GradleUpdatePreconditionIn{
				Readme:                  domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: domain.FileContent("INVALID"),
				LatestGradleRelease:     &domain.GradleRelease{Version: "5.1"},
			},
			domain.NoGradleVersion,
		}, {
			"AlreadyHasLatestGradle",
			domain.GradleUpdatePreconditionIn{
				Readme:                  domain.FileContent("![Gradle Status](https://example.com/owner/repo/status.svg)"),
				BadgeURL:                "/owner/repo/status.svg",
				GradleWrapperProperties: testdata.GradleWrapperProperties50,
				LatestGradleRelease:     &domain.GradleRelease{Version: "5.0"},
			},
			domain.AlreadyHasLatestGradle,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			out := domain.CheckGradleUpdatePrecondition(c.in)
			if diff := deep.Equal(c.out, out); diff != nil {
				t.Error(diff)
			}
		})
	}
}
