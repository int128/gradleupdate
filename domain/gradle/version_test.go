package gradle

import (
	"fmt"
	"testing"
)

func TestVersion_Compare(t *testing.T) {
	for _, c := range []struct {
		a        Version
		b        Version
		expected int
	}{
		{Version("1"), Version("1"), 0},
		{Version("2"), Version("1"), 1},
		{Version("1.2"), Version("1.2"), 0},
		{Version("1.3"), Version("1.2"), 1},
		{Version("1.2.0"), Version("1.2"), 1},
		{Version("1.2.1"), Version("1.2"), 1},
		{Version("1.10"), Version("1.2"), 1},
		{Version("1.24"), Version("1.23"), 1},
		{Version("1.100"), Version("1.23"), 1},
		{Version("2.0"), Version("1.2"), 1},
		{Version("10.2"), Version("1.2"), 1},

		{Version("1.2-rc1"), Version("1.2-rc1"), 0},
		{Version("1.2-rc3"), Version("1.2-rc1"), 1},
		// TODO:
		// {Version("1.10-rc1"), Version("1.2-rc1"), 1},
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

func TestVersion_GreaterOrEqualThan(t *testing.T) {
	for _, c := range []struct {
		a        Version
		b        Version
		expected bool
	}{
		{Version("1.1"), Version("1.2"), false},
		{Version("1.2"), Version("1.2"), true},
		{Version("1.3"), Version("1.2"), true},
	} {
		t.Run(fmt.Sprintf("%s≧%s", c.a, c.b), func(t *testing.T) {
			actual := c.a.GreaterOrEqualThan(c.b)
			if c.expected != actual {
				t.Errorf("(%s).GreaterOrEqualThan(%s) wants %v but %v", c.a, c.b, c.expected, actual)
			}
		})
	}
}
