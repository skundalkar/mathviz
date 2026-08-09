// Package gradientdescent shows what "learning rate" actually controls: how
// big a step a ball rolling downhill takes at each move. Too small and it
// creeps toward the bottom forever; too large and it overshoots the valley,
// bounces to the other wall, and can fly further away with every step.
package gradientdescent

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "gradient-descent",
		Title: "Gradient descent",
		Blurb: "Drop a ball on the side of a bowl-shaped valley and let physics take over: " +
			"it rolls toward the bottom, the lowest point, guided at every instant by which way " +
			"is downhill. Gradient descent is that same idea turned into an algorithm — at each " +
			"step, move against the slope (the gradient) by an amount set by the 'learning " +
			"rate'. A small learning rate takes tiny, cautious steps: it will get there, but " +
			"slowly. A learning rate that's too large overshoots the bottom entirely, lands " +
			"higher up the opposite wall, and if it's large enough the overshoot gets worse " +
			"each step instead of better — the 'ball' flies further from the minimum with every " +
			"bounce. Drag the learning rate and watch the same starting point converge, crawl, " +
			"or blow up.",
		Params: []concept.ParamSpec{
			{Key: "lr", Label: "Learning rate", Min: 0.01, Max: 1.2, Step: 0.01, Def: 0.3},
			{Key: "steps", Label: "Steps", Min: 1, Max: 30, Step: 1, Def: 12},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	_ = p
	return viz.New(680, 340, -1, 1, -1, 1).String()
}
