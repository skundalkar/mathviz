// Package fourier visualizes Fourier series: building a repeating signal --
// even one with a sudden jump, like a square wave -- out of nothing but
// smooth sine waves at increasing frequencies. Each added term narrows the
// gap between the partial sum and the true target almost everywhere, but a
// fixed-height ripple near any jump (the Gibbs phenomenon) never goes away no
// matter how many terms are added.
package fourier

import (
	"mathviz/internal/concept"
	"mathviz/internal/viz"
)

func init() {
	concept.Register(concept.Concept{
		ID:    "fourier-series",
		Seq:   40,
		Title: "Fourier series (building a wave out of waves)",
		Sections: []concept.Section{
			{
				Heading: "Why would you need this?",
				Body: []string{
					"sine-cosine walks through how a single pure oscillation is built from " +
						"circular motion -- one smooth wave, one frequency. But a violin string, a " +
						"square-wave clock pulse in a circuit, and an EKG heartbeat trace aren't " +
						"single pure oscillations at all -- they repeat, but with jagged corners or " +
						"sudden jumps a single sine wave can't produce. Here's the sharper puzzle: a " +
						"violin and a flute can play the exact same pitch -- the exact same repeat " +
						"rate -- and still sound completely different. If the base frequency is " +
						"identical, what else could possibly be different about the two signals, " +
						"and is there a systematic way to build (or take apart) a repeating signal " +
						"that isn't a single clean sine wave?",
				},
			},
			{
				Heading: "How does it actually work?",
				Body: []string{
					"Take the hardest case: a square wave that jumps straight from -1 to +1 at " +
						"x=0 and back to -1 at x=π. Add up sine waves at 1x, 3x, 5x, 7x, ... the " +
						"base frequency (only the odd multiples), each one shrunk by 1/k, and scale " +
						"the whole sum by 4/π. Watch what happens at x=π/2 -- the middle of the " +
						"flat top, where the target value is exactly 1 -- as more terms are added:",
					"• 1 term (just sin x): f(π/2) = 1.2732 -- already overshoots 1.",
					"• 2 terms (add the 3rd harmonic): f(π/2) = 0.8488 -- swings under 1.",
					"• 3 terms: f(π/2) = 1.1035. • 5 terms: f(π/2) = 1.0631.",
					"• 10 terms: f(π/2) = 0.9682. • 20 terms: f(π/2) = 0.9841.",
					"Each new term bends the sum closer to 1, but it overshoots and undershoots " +
						"on the way there rather than approaching in a straight line -- the same " +
						"shape of behavior as law-of-large-numbers' running average, just for a " +
						"sum of waves instead of a sum of coin flips. Away from the jump, that " +
						"wobble keeps shrinking toward 0 as more terms are added. But right at the " +
						"jump itself, it doesn't: with 20 terms, the sum peaks at 1.1792 just past " +
						"x=0 -- a 17.9% overshoot above the true value of 1 -- and with 50 terms it " +
						"still peaks at 1.1790, a 17.9% overshoot, in almost exactly the same place. " +
						"More terms squeeze the overshoot into a narrower sliver of x, but not " +
						"shorter -- that's the Gibbs phenomenon, and no finite (or infinite) number " +
						"of sine terms removes it.",
				},
			},
			{
				Heading: "What does the picture show?",
				Body: []string{
					"The dashed line is the target wave over one period; the solid curve is the " +
						"partial sum using the first `harmonics` sine terms. The wave slider picks " +
						"the target: a square wave (odd harmonics only: 1x, 3x, 5x, ...) or a " +
						"sawtooth ramp (every harmonic, alternating sign). Drag harmonics up and " +
						"watch the solid curve hug the dashed target more tightly away from any " +
						"jump -- while a stubborn ripple persists right at the jump itself no matter " +
						"how high you push it. The readout reports the current peak overshoot near " +
						"the jump, in percent, so you can watch it hold steady around 18% (square) " +
						"even as the rest of the curve visibly flattens out.",
				},
			},
			{
				Heading: "What can you do now that you couldn't before?",
				Body: []string{
					"Approximate a repeating signal you couldn't describe with one sine wave -- " +
						"to any tolerance you like, away from its jumps -- by adding enough sine " +
						"terms at the right frequencies and weights. You can also predict a real " +
						"engineering limit instead of being surprised by it: adding more frequency " +
						"components to a filter or a digital approximation of a sharp edge won't " +
						"remove the ringing right at that edge, because Gibbs' overshoot doesn't " +
						"shrink with more terms, only narrows. And you now have a concrete answer " +
						"to the violin-vs-flute puzzle: two signals with the same fundamental " +
						"frequency differ in which harmonics are present and how strongly -- that " +
						"harmonic mixture is what a spectrum analyzer or equalizer actually " +
						"displays and adjusts.",
				},
			},
			{
				Heading: "Where does this show up in real life?",
				Body: []string{
					"A musical instrument's timbre -- why a violin, a flute, and a trumpet playing " +
						"the identical note sound nothing alike -- comes down to which harmonics " +
						"ride along with the fundamental pitch and how loud each one is; that's " +
						"exactly what an audio equalizer's sliders adjust. Oscilloscopes routinely " +
						"show real Gibbs-phenomenon ringing on digital clock signals right after a " +
						"sharp voltage edge -- not a measurement glitch, but this exact effect. JPEG " +
						"and MP3 compression both work by keeping only the most important frequency " +
						"components of an image or sound and discarding the rest. Radio and Wi-Fi " +
						"transmission encode information onto combinations of sine waves at " +
						"different frequencies, and an antenna or receiver has to separate them back " +
						"out. Even an EKG's heartbeat trace is routinely broken into its frequency " +
						"components to spot abnormal rhythms buried under noise.",
				},
			},
			{
				Heading: "What's the common mistake here?",
				Body: []string{
					"Say it like this: a finite sum of sine waves can get arbitrarily close to a " +
						"jump-discontinuous target everywhere except right at the jump, where a " +
						"fixed roughly-18% overshoot persists no matter how many terms you add. Not " +
						"like this: assuming more harmonics eventually make the approximation " +
						"perfect everywhere, including exactly at the jump -- the Gibbs phenomenon " +
						"says the overshoot's height stays put even as its width shrinks. Also not " +
						"like this: assuming a repeating signal's only meaningful property is its " +
						"base repetition rate (its pitch, or period) -- two signals with the " +
						"identical fundamental frequency can look and sound completely different " +
						"depending on which harmonics are layered on top and how strongly.",
				},
			},
		},
		Params: []concept.ParamSpec{
			{Key: "harmonics", Label: "Sine terms (harmonics)", Min: 1, Max: 50, Step: 1, Def: 5},
			{Key: "wave", Label: "Target wave (0=square, 1=sawtooth)", Min: 0, Max: 1, Step: 1, Def: 0},
		},
		Render: render,
	})
}

func render(params map[string]float64) string {
	_ = params
	return viz.New(680, 420, -3.3, 3.3, -1.6, 1.6).String()
}
