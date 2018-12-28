package domain

import (
	"fmt"
	"testing"
)

func TestGradleVersion_Compare(t *testing.T) {
	for _, c := range []struct {
		a        GradleVersion
		b        GradleVersion
		expected int
	}{
		{GradleVersion("1"), GradleVersion("1"), 0},
		{GradleVersion("2"), GradleVersion("1"), 1},
		{GradleVersion("1.2"), GradleVersion("1.2"), 0},
		{GradleVersion("1.3"), GradleVersion("1.2"), 1},
		{GradleVersion("1.2.0"), GradleVersion("1.2"), 1},
		{GradleVersion("1.2.1"), GradleVersion("1.2"), 1},
		{GradleVersion("1.10"), GradleVersion("1.2"), 1},
		{GradleVersion("1.24"), GradleVersion("1.23"), 1},
		{GradleVersion("1.100"), GradleVersion("1.23"), 1},
		{GradleVersion("2.0"), GradleVersion("1.2"), 1},
		{GradleVersion("10.2"), GradleVersion("1.2"), 1},

		{GradleVersion("1.2-rc1"), GradleVersion("1.2-rc1"), 0},
		{GradleVersion("1.2-rc3"), GradleVersion("1.2-rc1"), 1},
		// TODO:
		// {GradleVersion("1.10-rc1"), GradleVersion("1.2-rc1"), 1},
	} {
		t.Run(fmt.Sprintf("%s%s%s", c.a, signSymbol(c.expected), c.b), func(t *testing.T) {
			actual := c.a.Compare(c.b)
			if c.expected != actual {
				t.Errorf("(%s).Compare(%s) wants %v but %v", c.a, c.b, c.expected, actual)
			}
		})
		t.Run(fmt.Sprintf("%s%s%s", c.b, signSymbol(c.expected*-1), c.a), func(t *testing.T) {
			actual := c.b.Compare(c.a)
			if c.expected*-1 != actual {
				t.Errorf("(%s).Compare(%s) wants %v but %v", c.b, c.a, c.expected*-1, actual)
			}
		})
	}
}

func signSymbol(s int) string {
	switch s {
	case 1:
		return ">"
	case -1:
		return "<"
	}
	return "="
}

func TestGradleVersion_GreaterOrEqualThan(t *testing.T) {
	for _, c := range []struct {
		a        GradleVersion
		b        GradleVersion
		expected bool
	}{
		{GradleVersion("1.1"), GradleVersion("1.2"), false},
		{GradleVersion("1.2"), GradleVersion("1.2"), true},
		{GradleVersion("1.3"), GradleVersion("1.2"), true},
	} {
		t.Run(fmt.Sprintf("%s≧%s", c.a, c.b), func(t *testing.T) {
			actual := c.a.GreaterOrEqualThan(c.b)
			if c.expected != actual {
				t.Errorf("(%s).GreaterOrEqualThan(%s) wants %v but %v", c.a, c.b, c.expected, actual)
			}
		})
	}
}
