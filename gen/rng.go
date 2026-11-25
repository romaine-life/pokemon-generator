package gen

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
	"math/rand/v2"
)

// rng is a deterministic random source: the same key always gives the same sequence, so a
// seed is a permanent address for its creature.
type rng struct{ r *rand.Rand }

func newRNG(parts ...string) *rng {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	sum := h.Sum(nil)
	return &rng{rand.New(rand.NewPCG(binary.LittleEndian.Uint64(sum[:8]), binary.LittleEndian.Uint64(sum[8:16])))}
}

// float returns a value in [lo, hi).
func (g *rng) float(lo, hi float64) float64 { return lo + g.r.Float64()*(hi-lo) }

// intn returns a value in [lo, hi].
func (g *rng) intn(lo, hi int) int { return lo + g.r.IntN(hi-lo+1) }

func (g *rng) chance(p float64) bool { return g.r.Float64() < p }

func pick[T any](g *rng, items []T) T { return items[g.r.IntN(len(items))] }

// weighted picks an index with probability proportional to its weight.
func (g *rng) weighted(weights []float64) int {
	total := 0.0
	for _, w := range weights {
		total += w
	}
	x := g.r.Float64() * total
	for i, w := range weights {
		if x < w {
			return i
		}
		x -= w
	}
	return len(weights) - 1
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
