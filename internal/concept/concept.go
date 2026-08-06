// Package concept defines the plugin contract every visualization implements
// and a global registry the WASM front-end iterates over to build the gallery.
//
// A concept is deliberately split so that ALL of the interesting logic (the
// math and the SVG it produces) lives in plain, testable Go. The WASM layer
// only reads Params to draw sliders and calls Render on every slider change.
// Nothing here imports syscall/js, so `go test ./...` exercises every concept.
package concept

import "sort"

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

// Concept is one self-contained lesson: a title, an explanation, the knobs the
// learner can turn, and a pure function that turns knob values into an SVG.
type Concept struct {
	ID     string      // url-safe id, also the folder name
	Title  string      // gallery title
	Blurb  string      // one-paragraph plain-language explanation
	Params []ParamSpec // interactive controls
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
// package so every init runs. Registering a duplicate ID panics loudly — that
// is a programming error the build should never ship.
func Register(c Concept) {
	if c.ID == "" {
		panic("concept: Register called with empty ID")
	}
	if _, dup := registry[c.ID]; dup {
		panic("concept: duplicate ID " + c.ID)
	}
	registry[c.ID] = c
}

// All returns every registered concept, sorted by ID for a stable gallery order.
func All() []Concept {
	out := make([]Concept, 0, len(registry))
	for _, c := range registry {
		out = append(out, c)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
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
