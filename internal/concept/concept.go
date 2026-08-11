// Package concept defines the plugin contract every visualization implements
// and a global registry the WASM front-end iterates over to build the gallery.
//
// A concept is deliberately split so that ALL of the interesting logic (the
// math and the SVG it produces) lives in plain, testable Go. The WASM layer
// only reads Params to draw sliders and calls Render on every slider change.
// Nothing here imports syscall/js, so `go test ./...` exercises every concept.
package concept

import (
	"sort"
	"strconv"
)

// ParamSpec describes one interactive slider for a concept.
type ParamSpec struct {
	Key   string  // stable identifier used as the map key passed to Render
	Label string  // human label shown next to the slider
	Min   float64 // slider minimum
	Max   float64 // slider maximum
	Step  float64 // slider granularity
	Def   float64 // default / initial value
	Unit  string  // optional suffix shown after the value (e.g. "σ", "")
}

// Section is one labeled block of a concept's explanation. Heading is phrased
// as a question ("Why would you need this?", "How does it actually work?")
// so the sequence of Sections reads as a guided path that gradually builds
// understanding, rather than one undifferentiated paragraph. Body holds one
// or more plain-language paragraphs; an entry prefixed with "• " renders as
// a bullet-list item instead of a paragraph, so a short worked sequence
// (e.g. sweeping a threshold through a few settings) can read as steps. An
// entry shaped like "|cell|cell|cell|" renders as one row of a table —
// consecutive table rows form one table, with the first row as the header
// — for reference material that was worked out as a table in conversation
// (pattern/diagnosis/action, reading/diagnosis/action) and reads better as
// one than as prose.
type Section struct {
	Heading string
	Body    []string
}

// Concept is one self-contained lesson: a title, an explanation, the knobs the
// learner can turn, and a pure function that turns knob values into an SVG.
type Concept struct {
	ID    string // url-safe id, also the folder name
	Title string // gallery title
	// Blurb is a single-paragraph fallback explanation, used when Sections
	// is empty. Prefer Sections for new concepts — see Section's doc
	// comment — but Blurb keeps older concepts compiling and rendering
	// (as one plain paragraph, same as before) without forcing every
	// existing entry to be restructured at once.
	Blurb string
	// Sections is the preferred, structured explanation: an ordered set of
	// question-headed blocks the WASM front-end renders as separate,
	// visually distinct sections instead of one wall of text. When set,
	// the front-end renders these and ignores Blurb.
	Sections []Section
	Params   []ParamSpec // interactive controls
	// Seq is the build sequence number: 1 for the first concept ever built,
	// incrementing by one for each concept after it, never reused or
	// renumbered. It exists purely to order the gallery sidebar newest-first
	// (see All()) so a freshly built concept is easy to find instead of
	// buried alphabetically. Set it explicitly in the Register call to
	// concept.Count()+1 at the time the concept is scaffolded.
	Seq int
	// Render maps the current parameter values to a complete <svg> string.
	// It must be pure: same input -> same output, no globals, no time, no rand.
	Render func(p map[string]float64) string
}

// Defaults returns a fresh map of every param at its default value. Handy for
// tests and for the WASM layer's initial render.
func (c Concept) Defaults() map[string]float64 {
	m := make(map[string]float64, len(c.Params))
	for _, ps := range c.Params {
		m[ps.Key] = ps.Def
	}
	return m
}

var registry = map[string]Concept{}

// Register adds a concept to the global registry. Concepts call this from an
// init() func; the cmd/wasm and cmd/build binaries blank-import the concepts
// package so every init runs. Registering a duplicate ID, or one missing a
// positive Seq, panics loudly — those are programming errors the build
// should never ship.
func Register(c Concept) {
	if c.ID == "" {
		panic("concept: Register called with empty ID")
	}
	if c.Seq <= 0 {
		panic("concept: Register called with no Seq for ID " + c.ID)
	}
	if _, dup := registry[c.ID]; dup {
		panic("concept: duplicate ID " + c.ID)
	}
	for _, existing := range registry {
		if existing.Seq == c.Seq {
			panic("concept: duplicate Seq " + strconv.Itoa(c.Seq) + " (" + existing.ID + " and " + c.ID + ")")
		}
	}
	registry[c.ID] = c
}

// All returns every registered concept, newest first (highest Seq first) so
// the gallery surfaces whatever was just built instead of burying it
// alphabetically. Ties (which Register prevents) fall back to ID for a
// deterministic order.
func All() []Concept {
	out := make([]Concept, 0, len(registry))
	for _, c := range registry {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Seq != out[j].Seq {
			return out[i].Seq > out[j].Seq
		}
		return out[i].ID < out[j].ID
	})
	return out
}

// Get returns a concept by ID.
func Get(id string) (Concept, bool) {
	c, ok := registry[id]
	return c, ok
}

// Count reports how many concepts are registered. The digest uses this to show
// progress ("14 concepts and counting").
func Count() int { return len(registry) }
