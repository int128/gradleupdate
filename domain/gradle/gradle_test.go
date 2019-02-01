package gradle_test

import (
	"testing"

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
