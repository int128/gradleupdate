package template

import (
	"io"
	"text/template"
)

var badge = template.Must(template.ParseFiles("template/badge.svg"))

type badgeVars struct {
	Badge

	Width  int
	LeftX  int
	RightX int
}

// Badge renders a badge image.
type Badge struct {
	LeftText   string
	LeftFill   string
	LeftWidth  int // Default to auto
	RightText  string
	RightFill  string
	RightWidth int // Default to auto
}

// Render writes an image in SVG format.
func (b *Badge) Render(w io.Writer) error {
	return badge.Execute(w, b.computeVars())
}

func (b *Badge) computeVars() (v badgeVars) {
	v.Badge = *b

	if v.LeftWidth == 0 {
		v.LeftWidth = computeTextWidth(b.LeftText)
	}
	if v.RightWidth == 0 {
		v.RightWidth = computeTextWidth(b.RightText)
	}

	v.Width = v.LeftWidth + v.RightWidth
	v.LeftX = v.LeftWidth / 2
	v.RightX = v.LeftWidth + v.RightWidth/2
	return
}

func computeTextWidth(s string) int {
	l := len(s)
	return l*l/4 + l*5 + 12
}
