package templates

// Badge renders a badge image.
type Badge struct {
	LeftText   string
	LeftFill   string
	LeftWidth  int
	RightText  string
	RightFill  string
	RightWidth int
}

// BadgeTextWidth returns width for the text.
func BadgeTextWidth(s string) int {
	l := len(s)
	return l*l/4 + l*5 + 12
}
