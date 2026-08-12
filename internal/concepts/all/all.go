// Package all blank-imports every concept package so their init() functions run
// and register them. Both the WASM front-end and any tooling import this single
// package to pull in the whole catalog. Adding a concept = adding one line here.
package all

import (
	_ "mathviz/internal/concepts/bayestheorem"
	_ "mathviz/internal/concepts/calibration"
	_ "mathviz/internal/concepts/clt"
	_ "mathviz/internal/concepts/complexnumbers"
	_ "mathviz/internal/concepts/confint"
	_ "mathviz/internal/concepts/confusionmatrix"
	_ "mathviz/internal/concepts/correlation"
	_ "mathviz/internal/concepts/derivative"
	_ "mathviz/internal/concepts/entropy"
	_ "mathviz/internal/concepts/evalplaybook"
	_ "mathviz/internal/concepts/expgrowth"
	_ "mathviz/internal/concepts/gradientdescent"
	_ "mathviz/internal/concepts/integral"
	_ "mathviz/internal/concepts/logscale"
	_ "mathviz/internal/concepts/meanmedianmode"
	_ "mathviz/internal/concepts/normalskew"
	_ "mathviz/internal/concepts/overfitting"
	_ "mathviz/internal/concepts/prauc"
	_ "mathviz/internal/concepts/precisionrecall"
	_ "mathviz/internal/concepts/pvalue"
	_ "mathviz/internal/concepts/rocauc"
	_ "mathviz/internal/concepts/sigmoidsoftmax"
	_ "mathviz/internal/concepts/sinecosine"
	_ "mathviz/internal/concepts/stddev"
	_ "mathviz/internal/concepts/variancestddev"
	_ "mathviz/internal/concepts/vectors"
)
