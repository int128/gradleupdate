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

// DarkBadge returns a badge of dark color.
func DarkBadge(text string) Badge {
	return Badge{
		LeftText:   "Gradle",
		LeftFill:   "#555",
		LeftWidth:  47,
		RightText:  text,
		RightFill:  "#9f9f9f",
		RightWidth: computeTextWidth(text),
	}
}

// GreenBadge returns a badge of green color.
func GreenBadge(text string) Badge {
	return Badge{
		LeftText:   "Gradle",
		LeftFill:   "#555",
		LeftWidth:  47,
		RightText:  text,
		RightFill:  "#4c1",
		RightWidth: computeTextWidth(text),
	}
}

// RedBadge returns a badge of red color.
func RedBadge(text string) Badge {
	return Badge{
		LeftText:   "Gradle",
		LeftFill:   "#555",
		LeftWidth:  47,
		RightText:  text,
		RightFill:  "#e05d44",
		RightWidth: computeTextWidth(text),
	}
}

func computeTextWidth(s string) int {
	l := len(s)
	return l*l/4 + l*5 + 12
}
