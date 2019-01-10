package domain

import (
	"testing"

	"github.com/int128/gradleupdate/domain/testdata"
)

func TestFindGradleWrapperVersion(t *testing.T) {
	version := FindGradleWrapperVersion(testdata.GradleWrapperProperties4102)
	if want := "4.10.2"; version.String() != want {
		t.Errorf("version wants %s but %s", want, version)
	}
}

func TestReplaceGradleWrapperVersion(t *testing.T) {
	replaced := ReplaceGradleWrapperVersion(testdata.GradleWrapperProperties4102, "5.0")
	want := `distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=https\://services.gradle.org/distributions/gradle-5.0-bin.zip
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
`
	if replaced != want {
		t.Errorf("replaced wants %s but %s", want, replaced)
	}
}
