//go:build js && wasm

// Command wasm is the browser front-end, compiled to WebAssembly. It is a thin
// view layer: it reads the concept registry, draws a sidebar of lessons, and
// for the selected lesson renders a slider per Param. Every time a slider moves
// it calls the concept's pure Render function and drops the returned SVG into
// the page. All the interesting logic lives in the (testable) concept packages;
// this file only touches the DOM.
package main

import (
	"strconv"
	"strings"
	"syscall/js"

	"mathviz/internal/concept"
	_ "mathviz/internal/concepts/all"
)

var (
	doc     = js.Global().Get("document")
	current concept.Concept
	params  = map[string]float64{}
)

func main() {
	buildSidebar()
	if cs := concept.All(); len(cs) > 0 {
		selectConcept(cs[0])
	}
	select {} // keep the Go runtime alive for event callbacks
}

func el(tag string) js.Value { return doc.Call("createElement", tag) }

func byID(id string) js.Value { return doc.Call("getElementById", id) }

func buildSidebar() {
	nav := byID("concept-list")
	nav.Set("innerHTML", "")
	for _, c := range concept.All() {
		c := c
		item := el("button")
		item.Set("className", "concept-link")
		item.Set("textContent", c.Title)
		item.Call("addEventListener", "click", js.FuncOf(func(js.Value, []js.Value) any {
			selectConcept(c)
			return nil
		}))
		nav.Call("appendChild", item)
	}
	byID("concept-count").Set("textContent", strconv.Itoa(concept.Count())+" concepts")
}

func selectConcept(c concept.Concept) {
	current = c
	params = c.Defaults()

	byID("concept-title").Set("textContent", c.Title)
	renderExplanation(c)

	controls := byID("controls")
	controls.Set("innerHTML", "")
	for _, ps := range c.Params {
		ps := ps
		row := el("div")
		row.Set("className", "control-row")

		label := el("label")
		label.Set("className", "control-label")
		valSpan := el("span")
		valSpan.Set("className", "control-value")
		valSpan.Set("textContent", fmtVal(ps.Def, ps.Unit))
		label.Set("textContent", ps.Label+"  ")
		label.Call("appendChild", valSpan)

		slider := el("input")
		slider.Set("type", "range")
		slider.Set("min", ps.Min)
		slider.Set("max", ps.Max)
		slider.Set("step", ps.Step)
		slider.Set("value", ps.Def)
		slider.Call("addEventListener", "input", js.FuncOf(func(this js.Value, _ []js.Value) any {
			v, _ := strconv.ParseFloat(this.Get("value").String(), 64)
			params[ps.Key] = v
			valSpan.Set("textContent", fmtVal(v, ps.Unit))
			draw()
			return nil
		}))

		row.Call("appendChild", label)
		row.Call("appendChild", slider)
		controls.Call("appendChild", row)
	}
	draw()
}

// renderExplanation fills #concept-blurb with the concept's explanation.
// Concepts that set Sections get a sequence of headed blocks — each
// heading a question, so the sections read as a guided path instead of one
// wall of text. Concepts that only set Blurb (not yet migrated to
// Sections) fall back to a single plain paragraph, same as before.
func renderExplanation(c concept.Concept) {
	host := byID("concept-blurb")
	host.Set("innerHTML", "")

	if len(c.Sections) == 0 {
		p := el("p")
		p.Set("className", "blurb-fallback")
		p.Set("textContent", c.Blurb)
		host.Call("appendChild", p)
		return
	}

	for _, sec := range c.Sections {
		block := el("section")
		block.Set("className", "explain-section")

		h := el("h3")
		h.Set("className", "explain-heading")
		h.Set("textContent", sec.Heading)
		block.Call("appendChild", h)

		var list js.Value
		hasList := false
		flushList := func() {
			if hasList {
				block.Call("appendChild", list)
				hasList = false
			}
		}

		for _, para := range sec.Body {
			if strings.HasPrefix(para, "• ") {
				if !hasList {
					list = el("ul")
					list.Set("className", "explain-list")
					hasList = true
				}
				li := el("li")
				li.Set("textContent", strings.TrimPrefix(para, "• "))
				list.Call("appendChild", li)
				continue
			}
			flushList()
			p := el("p")
			p.Set("textContent", para)
			block.Call("appendChild", p)
		}
		flushList()

		host.Call("appendChild", block)
	}
}

func draw() {
	byID("stage").Set("innerHTML", current.Render(params))
}

func fmtVal(v float64, unit string) string {
	return strconv.FormatFloat(v, 'g', 3, 64) + unit
}
