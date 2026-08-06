// Package all blank-imports every concept package so their init() functions run
// and register them. Both the WASM front-end and any tooling import this single
// package to pull in the whole catalog. Adding a concept = adding one line here.
package all

import (
	_ "mathviz/internal/concepts/logscale"
	_ "mathviz/internal/concepts/precisionrecall"
	_ "mathviz/internal/concepts/stddev"
)
