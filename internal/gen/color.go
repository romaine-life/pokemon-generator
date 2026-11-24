package gen

import (
	"fmt"
	"math"
)

// Sprite is a creature's pixel art: Pixels holds Width*Height palette indexes, row by row,
// and index 0 is transparent.
type Sprite struct {
	Width   int      `json:"width"`
	Height  int      `json:"height"`
	Palette []string `json:"palette"` // #rrggbb or #rrggbbaa
	Pixels  []int    `json:"pixels"`
}

// shadowColor is the translucent shadow under a creature.
const shadowColor = "#0000002e"

type hsl struct{ h, s, l float64 }

func (c hsl) lighter(d float64) hsl { return hsl{c.h, c.s, clamp(c.l+d, 0, 0.97)} }

func (c hsl) hex() string {
	h := math.Mod(c.h, 360)
	if h < 0 {
		h += 360
	}
	s, l := clamp(c.s, 0, 1), clamp(c.l, 0, 1)
	ch := (1 - math.Abs(2*l-1)) * s
	x := ch * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - ch/2
	var r, g, b float64
	switch {
	case h < 60:
		r, g, b = ch, x, 0
	case h < 120:
		r, g, b = x, ch, 0
	case h < 180:
		r, g, b = 0, ch, x
	case h < 240:
		r, g, b = 0, x, ch
	case h < 300:
		r, g, b = x, 0, ch
	default:
		r, g, b = ch, 0, x
	}
	to := func(v float64) int { return int(math.Round((v + m) * 255)) }
	return fmt.Sprintf("#%02x%02x%02x", to(r), to(g), to(b))
}

// palette is the base color of each material.
type palette map[material]hsl

func typeColor(ge genome, name string, t elemType) hsl {
	return hsl{ge.f(name+"H", t.hue[0], t.hue[1]), ge.f(name+"S", t.sat[0], t.sat[1]), ge.f(name+"L", t.light[0], t.light[1])}
}

func makePalette(sp spriteSpec) palette {
	ge := genome{sp.family}
	body := typeColor(ge, "body", sp.primary)
	var accent hsl
	switch {
	case sp.secondary != nil:
		accent = typeColor(ge, "accent", *sp.secondary)
	case ge.chance("ivory", 0.4):
		accent = hsl{ge.f("ivoryH", 35, 55), 0.45, 0.82}
	default:
		accent = hsl{body.h + ge.f("accentShift", 150, 210), clamp(body.s+0.1, 0.4, 0.9), ge.f("accentL", 0.45, 0.65)}
	}
	if sp.shiny {
		shift := ge.f("shinyShift", 100, 260)
		body.h += shift
		accent.h += shift
		body.s = clamp(body.s+0.15, 0.3, 1)
	}
	belly := body.lighter(0.22)
	if ge.chance("cream", 0.6) {
		belly = hsl{ge.f("creamH", 30, 55), 0.45, 0.86}
	}
	pattern := body.lighter(-0.18)
	if ge.chance("accentPattern", 0.5) {
		pattern = accent
	}
	tip := accent
	switch sp.primary.tailTip {
	case "flame":
		tip = hsl{ge.f("flameH", 28, 45), 1, 0.56}
	case "bolt":
		tip = hsl{52, 1, 0.55}
	case "leaf":
		tip = hsl{ge.f("leafH", 95, 130), 0.6, 0.42}
	case "crystal":
		tip = hsl{190, 0.75, 0.8}
	}
	return palette{
		mBody:     body,
		mBelly:    belly,
		mAccent:   accent,
		mPattern:  pattern,
		mWing:     hsl{body.h, body.s * 0.7, clamp(body.l+0.25, 0, 0.9)},
		mEyeWhite: {0, 0, 0.97},
		mPupil:    {body.h, 0.35, 0.1},
		mShine:    {0, 0, 1},
		mMouth:    {350, 0.45, 0.28},
		mCheek:    {ge.f("cheekH", 340, 365), 0.8, 0.68},
		mTip:      tip,
		mAura:     accent.lighter(0.15),
	}
}

// colorize turns the drawn canvas into a sprite: each part is lit along its top edge and
// shaded along its bottom, and the creature and its parts are outlined.
func colorize(cv *canvas, pal palette) Sprite {
	body := pal[mBody]
	outline := hsl{body.h, body.s * 0.6, 0.13}.hex()
	inner := hsl{body.h, body.s * 0.5, 0.24}.hex()
	index := map[string]int{"": 0}
	sp := Sprite{Width: gridW, Height: gridH, Palette: []string{"#00000000"}, Pixels: make([]int, gridW*gridH)}
	color := func(hex string) int {
		if i, ok := index[hex]; ok {
			return i
		}
		index[hex] = len(sp.Palette)
		sp.Palette = append(sp.Palette, hex)
		return index[hex]
	}
	solid := func(c cell) bool { return c.mat != mEmpty && !c.mat.glow() }
	for y := range gridH {
		for x := range gridW {
			c := cv.at(x, y)
			var hex string
			switch {
			case c.mat == mShadow:
				hex = shadowColor
			case c.mat == mAura && !touches(cv, x, y, solid):
				hex = pal[mAura].hex() + "cc"
			case c.mat == mEmpty || c.mat == mAura:
				if touches(cv, x, y, solid) {
					hex = outline
				}
			case c.mat.flat():
				hex = pal[c.mat].hex()
			default:
				// An edge against a part in front of this one is drawn as a line.
				if touches(cv, x, y, func(n cell) bool { return solid(n) && n.z > c.z && !n.mat.flat() }) {
					hex = inner
					break
				}
				base := pal[c.mat]
				above, below := cv.at(x, y-1), cv.at(x, y+1)
				switch {
				case !solid(above) || above.z != c.z:
					hex = base.lighter(0.1).hex()
				case !solid(below) || below.z != c.z || (x+y)%2 == 0 && solid(cv.at(x, y+2)) && cv.at(x, y+2).z != c.z:
					hex = base.lighter(-0.13).hex()
				default:
					hex = base.hex()
				}
			}
			sp.Pixels[y*gridW+x] = color(hex)
		}
	}
	return sp
}

func touches(cv *canvas, x, y int, match func(cell) bool) bool {
	return match(cv.at(x-1, y)) || match(cv.at(x+1, y)) || match(cv.at(x, y-1)) || match(cv.at(x, y+1))
}

// renderSprite draws the creature, shrinking it until it fits the grid.
func renderSprite(sp spriteSpec) Sprite {
	k := 1.0
	var cv *canvas
	for range 8 {
		cv = drawCreature(sp, k)
		if !cv.clipped {
			break
		}
		k *= 0.92
	}
	return colorize(cv, makePalette(sp))
}
