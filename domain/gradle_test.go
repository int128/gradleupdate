package domain_test

import (
	"testing"

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
