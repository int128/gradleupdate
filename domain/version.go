package domain

import (
	"strconv"
	"strings"
)

// GradleVersion represents version of a Gradle.
type GradleVersion string

// Compare returns an integer comparing two versions by semantic versioning.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func (a GradleVersion) Compare(b GradleVersion) int {
	as := strings.Split(string(a), ".")
	bs := strings.Split(string(b), ".")
	for i := 0; i < minInt(len(as), len(bs)); i++ {
		ai, aerr := strconv.Atoi(as[i])
		bi, berr := strconv.Atoi(bs[i])
		if aerr == nil && berr == nil {
			switch {
			case ai > bi:
				return 1
			case ai < bi:
				return -1
			}
		} else {
			cmp := strings.Compare(as[i], bs[i])
			if cmp != 0 {
				return cmp
			}
		}
	}
	switch {
	case len(as) > len(bs):
		return 1
	case len(as) < len(bs):
		return -1
	}
	return 0
}

func minInt(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func (a GradleVersion) String() string {
	return string(a)
}

func (a GradleVersion) GreaterOrEqualThan(b GradleVersion) bool {
	return a.Compare(b) >= 0
}
