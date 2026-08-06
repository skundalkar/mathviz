// Package viz is a tiny, dependency-free SVG builder. It exists so concepts can
// describe pictures declaratively and stay pure & testable — every helper just
// appends strings. There is intentionally no cleverness here; readability and
// deterministic output (for golden tests) matter more than features.
package viz

import (
	"fmt"
	"math"
	"strings"
)

// Palette — a small, colorblind-friendly set. Kept here so every concept looks
// like one system. Swap these to rebrand the whole gallery at once.
const (
	Ink    = "#1f2933" // near-black text / axes
	Muted  = "#7b8794" // secondary text, gridlines
	Accent = "#2f6fed" // primary series (blue)
	Warm   = "#e8663d" // secondary series / highlight (orange)
	Good   = "#2f9e6f" // positive region (green)
	Bad    = "#d64545" // negative region (red)
	Faint  = "#eef2f7" // panel fill
)

// Canvas is a fixed-size drawing surface with a simple linear mapping from
// data coordinates to pixel coordinates.
type Canvas struct {
	W, H       float64 // pixel size
	PadL, PadR float64 // horizontal padding (room for axis labels)
	PadT, PadB float64 // vertical padding
	XMin, XMax float64 // data-space x range
	YMin, YMax float64 // data-space y range
	b          strings.Builder
}

// New returns a Canvas with sensible padding and the given data ranges.
func New(w, h, xmin, xmax, ymin, ymax float64) *Canvas {
	c := &Canvas{
		W: w, H: h,
		PadL: 48, PadR: 18, PadT: 18, PadB: 34,
		XMin: xmin, XMax: xmax, YMin: ymin, YMax: ymax,
	}
	c.b.WriteString(fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 %g %g" font-family="ui-sans-serif,system-ui,sans-serif">`,
		w, h))
	c.b.WriteString(fmt.Sprintf(`<rect x="0" y="0" width="%g" height="%g" fill="white"/>`, w, h))
	return c
}

// X maps a data-space x to a pixel x.
func (c *Canvas) X(x float64) float64 {
	if c.XMax == c.XMin {
		return c.PadL
	}
	return c.PadL + (x-c.XMin)/(c.XMax-c.XMin)*(c.W-c.PadL-c.PadR)
}

// Y maps a data-space y to a pixel y (flipped so larger y is higher up).
func (c *Canvas) Y(y float64) float64 {
	if c.YMax == c.YMin {
		return c.H - c.PadB
	}
	return c.H - c.PadB - (y-c.YMin)/(c.YMax-c.YMin)*(c.H-c.PadT-c.PadB)
}

// Axes draws a plain L-shaped axis with an optional x baseline at y=0.
func (c *Canvas) Axes() *Canvas {
	x0, y0 := c.PadL, c.H-c.PadB
	c.b.WriteString(fmt.Sprintf(
		`<line x1="%g" y1="%g" x2="%g" y2="%g" stroke="%s" stroke-width="1.5"/>`,
		x0, c.PadT, x0, y0, Ink))
	c.b.WriteString(fmt.Sprintf(
		`<line x1="%g" y1="%g" x2="%g" y2="%g" stroke="%s" stroke-width="1.5"/>`,
		x0, y0, c.W-c.PadR, y0, Ink))
	return c
}

// VLine draws a vertical reference line across the plot at data-x.
func (c *Canvas) VLine(x float64, color string, dash bool) *Canvas {
	px := c.X(x)
	d := ""
	if dash {
		d = ` stroke-dasharray="5 4"`
	}
	c.b.WriteString(fmt.Sprintf(
		`<line x1="%g" y1="%g" x2="%g" y2="%g" stroke="%s" stroke-width="1.5"%s/>`,
		px, c.PadT, px, c.H-c.PadB, color, d))
	return c
}

// Path draws a polyline through data-space points in the given stroke color.
func (c *Canvas) Path(pts [][2]float64, color string, width float64) *Canvas {
	if len(pts) == 0 {
		return c
	}
	var sb strings.Builder
	for i, p := range pts {
		cmd := "L"
		if i == 0 {
			cmd = "M"
		}
		sb.WriteString(fmt.Sprintf("%s%.2f %.2f ", cmd, c.X(p[0]), c.Y(p[1])))
	}
	c.b.WriteString(fmt.Sprintf(
		`<path d="%s" fill="none" stroke="%s" stroke-width="%g" stroke-linejoin="round"/>`,
		strings.TrimSpace(sb.String()), color, width))
	return c
}

// Area fills the region between a curve (data-space points) and the y=0 axis,
// clipped to [xlo, xhi]. Used to shade "the area under the curve" regions.
func (c *Canvas) Area(pts [][2]float64, xlo, xhi float64, fill string, opacity float64) *Canvas {
	var sb strings.Builder
	first := true
	var lastX float64
	for _, p := range pts {
		if p[0] < xlo || p[0] > xhi {
			continue
		}
		cmd := "L"
		if first {
			sb.WriteString(fmt.Sprintf("M%.2f %.2f ", c.X(p[0]), c.Y(0)))
			first = false
		}
		sb.WriteString(fmt.Sprintf("%s%.2f %.2f ", cmd, c.X(p[0]), c.Y(p[1])))
		lastX = p[0]
	}
	if first {
		return c // nothing in range
	}
	sb.WriteString(fmt.Sprintf("L%.2f %.2f Z", c.X(lastX), c.Y(0)))
	c.b.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" opacity="%g"/>`,
		strings.TrimSpace(sb.String()), fill, opacity))
	return c
}

// Rect draws a filled rectangle in pixel space.
func (c *Canvas) Rect(x, y, w, h float64, fill string, opacity float64) *Canvas {
	c.b.WriteString(fmt.Sprintf(`<rect x="%g" y="%g" width="%g" height="%g" fill="%s" opacity="%g"/>`,
		x, y, w, h, fill, opacity))
	return c
}

// Text draws left-anchored text at a pixel position.
func (c *Canvas) Text(x, y float64, s string, size float64, color, anchor string) *Canvas {
	if anchor == "" {
		anchor = "start"
	}
	c.b.WriteString(fmt.Sprintf(
		`<text x="%g" y="%g" font-size="%g" fill="%s" text-anchor="%s">%s</text>`,
		x, y, size, color, anchor, escape(s)))
	return c
}

// Tick draws a short x-axis tick with a label at data-x.
func (c *Canvas) Tick(x float64, label string) *Canvas {
	px, py := c.X(x), c.H-c.PadB
	c.b.WriteString(fmt.Sprintf(`<line x1="%g" y1="%g" x2="%g" y2="%g" stroke="%s" stroke-width="1"/>`,
		px, py, px, py+5, Muted))
	c.Text(px, py+18, label, 12, Muted, "middle")
	return c
}

// String closes and returns the finished SVG document.
func (c *Canvas) String() string {
	return c.b.String() + "</svg>"
}

// Sample evaluates f at n+1 points across [lo, hi], returning data-space points.
// Concepts use this to turn a math function into a drawable curve.
func Sample(lo, hi float64, n int, f func(float64) float64) [][2]float64 {
	if n < 1 {
		n = 1
	}
	pts := make([][2]float64, 0, n+1)
	for i := 0; i <= n; i++ {
		x := lo + (hi-lo)*float64(i)/float64(n)
		y := f(x)
		if math.IsInf(y, 0) || math.IsNaN(y) {
			continue
		}
		pts = append(pts, [2]float64{x, y})
	}
	return pts
}

func escape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")
	return r.Replace(s)
}
