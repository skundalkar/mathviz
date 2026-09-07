// Package convolution visualizes 1D discrete convolution: sliding a small
// kernel across a signal and, at every position, multiplying the
// overlapping values and summing them -- the same "chop into pieces and
// add up" idea integral used to turn a whole curve into one number, except
// here the sliding window produces a brand-new signal instead of a single
// total. The same sliding-sum operation, with different kernels, either
// smooths a noisy signal or picks out exactly where it jumps.
package convolution

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "convolution",
		Seq:   94,
		Title: "Convolution (sliding, flipping, and summing overlap)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"integral chopped a curve into tiny slabs and summed them down to a single " +
						"number: the total area. Convolution reuses that exact 'chop into pieces " +
						"and sum the overlap' idea, but slides a small second function across a " +
						"first one instead of collapsing everything into one total. Say you have " +
						"a blocky sensor reading that jumps sharply from 0 up to 1 and back down " +
						"— a pulse — and you want two very different things from it at once: a " +
						"smoothed-out version that isn't so jumpy, and a precise marker of exactly " +
						"where the sharp jumps happen. Averaging the whole signal into one number, " +
						"the way integral would, throws away *where* anything happens. What you " +
						"need is an operation that looks at a small local neighborhood around each " +
						"point and produces a whole new signal — and, it turns out, the exact same " +
						"operation can either smooth the signal or find its jumps, just by " +
						"swapping out one small ingredient.",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the pulse signal x = [0, 0, 1, 1, 1, 1, 0, 0]. Slide a 3-wide window " +
						"across it; at each position, multiply the window's 3 values against a " +
						"small kernel and add the products — a tiny 3-term Riemann sum, computed " +
						"once per position instead of once total. First kernel: a moving average, " +
						"k = [1/3, 1/3, 1/3]:",
					"• Position 0: average(x0,x1,x2) = average(0,0,1) = 0.3333.",
					"• Position 1: average(0,1,1) = 0.6667. Position 2: average(1,1,1) = 1.0000.",
					"• Position 3: average(1,1,1) = 1.0000. Position 4: average(1,1,0) = 0.6667. " +
						"Position 5: average(1,0,0) = 0.3333.",
					"Full smoothed output: [0.33, 0.67, 1.00, 1.00, 0.67, 0.33] — the sharp jump " +
						"from 0 to 1 is now a gentle two-step ramp.",
					"Now the formal definition matters: true convolution reads the kernel " +
						"*backwards* as it slides — (x*k)[i] = Σⱼ x[i+j]·k[len(k)-1-j] — not " +
						"forwards. It doesn't matter for the symmetric averaging kernel above " +
						"(reversed, [1/3,1/3,1/3] is identical to itself), but it matters a lot " +
						"for a lopsided kernel. Swap in an edge-detecting kernel, k = [-1, 0, 1], " +
						"reversed to [1, 0, -1] before sliding:",
					"• Position 0: x0·1 + x1·0 + x2·(-1) = 0 + 0 − 1 = −1. Position 1: 0+0−1 = −1.",
					"• Position 2: x2·1+x3·0+x4·(-1) = 1+0−1 = 0. Position 3: 1+0−1 = 0.",
					"• Position 4: x4·1+x5·0+x6·(-1) = 1+0−0 = 1. Position 5: 1+0−0 = 1.",
					"Full edge output: [−1, −1, 0, 0, 1, 1] — a negative spike exactly where the " +
						"signal jumps UP, a positive spike exactly where it jumps DOWN.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The top row of bars is the fixed input signal x. The 3-wide window " +
						"highlighted in orange is the current slide position, set by the t " +
						"slider; the small labels above it show the kernel weight actually " +
						"multiplying each highlighted bar (which is the kernel read backwards " +
						"when 'flip' is on — true convolution — or forwards when it's off — " +
						"cross-correlation). The readout spells out that position's exact " +
						"arithmetic. The bottom row of bars is the output signal, revealed one " +
						"bar at a time as t increases, each bar the sum computed at that slide " +
						"position. Switch the kernel slider between smoothing and edge-detecting " +
						"and watch the bottom row change shape completely from the same top row " +
						"and the same sliding mechanism.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Build an entire family of signal-processing operations — smoothing, sharp" +
						"ening, edge-finding, echo, and more — out of one single mechanism (slide, " +
						"multiply, sum) just by changing which small set of numbers goes in the " +
						"kernel. That's also exactly what a convolutional neural network's " +
						"'filters' are: small kernels, except instead of being hand-picked like " +
						"the averaging and edge-detecting ones here, their numbers are learned " +
						"from data so the network discovers whatever local patterns (edges, " +
						"textures, and further in, shapes) turn out to be useful for the task.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"Photo-editing 'blur' and 'sharpen' tools are convolutions with a smoothing " +
						"or sharpening kernel. Computer vision's classic edge-detection filters " +
						"(Sobel, Prewitt) are 2D versions of the exact [-1,0,1]-style kernel used " +
						"above. Audio echo and equalizer effects convolve a sound wave with a " +
						"kernel that repeats or reshapes it over time. Every convolutional layer " +
						"in an image-recognition neural network is this same sliding-window sum, " +
						"run with many learned kernels at once across a 2D image instead of a 1D " +
						"signal. Even a stock market's 'moving average' line is the smoothing " +
						"kernel above, applied to daily prices instead of sensor readings.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: the *sign convention* of convolution depends on whether " +
						"the kernel gets flipped before sliding — true mathematical convolution " +
						"flips it, giving [−1,−1,0,0,1,1] for the edge kernel above (negative at " +
						"a rising edge, positive at a falling edge). Not like this: assuming " +
						"'convolution' always means that flipped version — most machine-learning " +
						"frameworks, and the phrase 'CNN filter,' actually implement " +
						"cross-correlation (no flip), which gives the exact sign-reversed " +
						"[1,1,0,0,−1,−1] for the same kernel and signal here (positive at a rising " +
						"edge instead). Neither convention is 'wrong' — they only disagree on " +
						"which is the kernel and which is a mirror image of it — but mixing the " +
						"two up mid-calculation, or assuming a paper's formula and a library's " +
						"implementation are automatically using the same one, silently flips the " +
						"sign of every asymmetric result.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "t", Label: "Slide position (t)", Min: 0, Max: 5, Step: 1, Def: 0},
			{Key: "kernel", Label: "Kernel (0=smoothing, 1=edge detector)", Min: 0, Max: 1, Step: 1, Def: 0},
			{Key: "flip", Label: "Flip kernel (0=cross-correlation, 1=true convolution)", Min: 0, Max: 1, Step: 1, Def: 1},
		},
		Render: render,
	})
}

func render(p map[string]float64) string {
	c := viz.New(680, 460, 0, 1, 0, 1)
	return c.String()
}
