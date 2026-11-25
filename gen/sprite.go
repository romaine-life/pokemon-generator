package gen

import (
	"math"
)

// genome gives each named trait of a family one fixed value, whatever order the traits are
// read in, so every stage of an evolution line reads the same body and they look related.
type genome struct{ key string }

func (ge genome) f(name string, lo, hi float64) float64 { return newRNG(ge.key, name).float(lo, hi) }
func (ge genome) n(name string, count int) int          { return newRNG(ge.key, name).intn(0, count-1) }
func (ge genome) chance(name string, p float64) bool    { return newRNG(ge.key, name).chance(p) }
func (ge genome) weighted(name string, w []float64) int {
	return newRNG(ge.key, name).weighted(w)
}

type spriteSpec struct {
	family    string
	arch      archetype
	primary   elemType
	secondary *elemType
	tier      int // 0 first stage, 1 middle, 2 final
	legendary bool
	shiny     bool
}

// anatomy is where the drawn body's parts landed, for the parts drawn on top of them.
type anatomy struct {
	headX, headY, headRx, headRy float64 // the face goes here
	bodyY, bodyRx, bodyRy        float64
	ground                       float64
	bodyZ                        int
	belly                        func(px, py float64) bool
}

var (
	earStyles  = []string{"none", "pointy", "round", "long", "horns", "antennae", "crest", "unicorn"}
	earWeights = map[string][]float64{
		"biped":     {1, 3, 2, 1.5, 2, 0.5, 1, 0.7},
		"quadruped": {0.5, 4, 2, 1.5, 2, 0, 0.5, 0.8},
		"serpent":   {2, 0.5, 0, 0, 2, 0.5, 3, 1},
		"avian":     {1, 0.5, 0, 0, 0.5, 0, 6, 0},
		"blob":      {3, 0.5, 1, 0.5, 0.5, 2.5, 0.5, 0.5},
		"insect":    {0.5, 0, 0, 0, 1, 7, 0, 1},
		"aquatic":   {3, 1, 0, 0, 0.5, 1, 3, 1},
		"golem":     {3, 0.5, 0, 0, 2, 0.5, 1.5, 1},
		"spirit":    {2.5, 2, 0.5, 0.5, 1, 0.5, 2, 0.5},
	}
	tailChance = map[string]float64{
		"biped": 0.6, "quadruped": 0.9, "serpent": 0.7, "avian": 0.35, "blob": 0.1,
		"insect": 0.35, "aquatic": 0, "golem": 0.15, "spirit": 0,
	}
)

// drawCreature draws the creature at overall scale k into a fresh canvas.
func drawCreature(sp spriteSpec, k float64) *canvas {
	ge := genome{sp.family}
	cv := &canvas{}
	s := k * []float64{0.66, 0.84, 1}[sp.tier]
	// First stages are big-headed.
	headScale := 1 + (1-[]float64{0.66, 0.84, 1}[sp.tier])*0.9

	ground := 37.5
	if sp.arch.floats {
		ground = 34.5
	}

	// The shadow on the ground, behind everything.
	cv.part()
	shadowRx := ge.f("shadow", 7, 11) * s
	if sp.arch.floats {
		cv.ellipse(midX, 38, shadowRx*0.7, 0.9, mShadow)
	} else {
		cv.ellipse(midX, ground+0.3, shadowRx, 1.1, mShadow)
	}

	if sp.legendary {
		cv.part()
		r := ge.f("aura", 13, 16) * k
		cy := ground - r - 2
		for a := 0.0; a < 2*math.Pi; a += 0.03 {
			if int(a*6/math.Pi)%2 == 0 {
				continue
			}
			for t := r - 0.6; t <= r+0.6; t += 0.3 {
				x, y := midX+math.Cos(a)*t, cy+math.Sin(a)*t*1.05
				cv.set(int(x), int(y), mAura)
			}
		}
	}

	var an anatomy
	switch sp.arch.id {
	case "quadruped":
		an = drawQuadruped(cv, ge, sp, s, headScale, ground)
	case "serpent":
		an = drawSerpent(cv, ge, sp, s, headScale, ground)
	case "avian":
		an = drawAvian(cv, ge, sp, s, headScale, ground)
	case "blob":
		an = drawBlob(cv, ge, sp, s, ground)
	case "insect":
		an = drawInsect(cv, ge, sp, s, headScale, ground)
	case "aquatic":
		an = drawAquatic(cv, ge, sp, s, ground)
	case "golem":
		an = drawGolem(cv, ge, sp, s, headScale, ground)
	case "spirit":
		an = drawSpirit(cv, ge, sp, s, ground)
	default:
		an = drawBiped(cv, ge, sp, s, headScale, ground)
	}

	if an.belly != nil && ge.chance("hasBelly", 0.75) {
		z := an.bodyZ
		cv.paint(func(px, py float64) bool {
			return an.belly(px, py) && cv.c[int(py)][int(px)].z == z
		}, mBelly, mBody)
	}
	drawPattern(cv, ge, an)
	drawHeadgear(cv, ge, sp, s, an)
	drawFace(cv, ge, sp, an)
	return cv
}

// drawTail curls a tail out of the body's right side, behind it, ending in the primary
// type's decoration. Call it before drawing the body.
func drawTail(cv *canvas, ge genome, sp spriteSpec, s float64, rootX, rootY, bodyRx float64) {
	if !ge.chance("tail", tailChance[sp.arch.id]) {
		return
	}
	cv.part()
	length := ge.f("tailLen", 5, 10) * s
	r := ge.f("tailR", 1, 2) * s
	curl := ge.f("tailCurl", 0.4, 1.3)
	p0 := pt{rootX, rootY}
	p1 := pt{midX + bodyRx + length*0.45, rootY - length*0.15}
	p2 := pt{midX + bodyRx + length*0.6, rootY - length*curl}
	cv.capsule(p0.x, p0.y, p1.x, p1.y, r, mBody)
	cv.capsule(p1.x, p1.y, p2.x, p2.y, r*0.8, mBody)
	drawTip(cv, ge, sp.primary.tailTip, p2, s)
}

func drawTip(cv *canvas, ge genome, style string, p pt, s float64) {
	cv.part()
	sz := math.Max(ge.f("tipSize", 1.6, 2.6)*s, 1.2)
	switch style {
	case "flame":
		cv.ellipse(p.x, p.y-sz*0.3, sz, sz*1.1, mTip)
		cv.tri(pt{p.x - sz, p.y - sz*0.4}, pt{p.x + sz, p.y - sz*0.4}, pt{p.x + sz*0.2, p.y - sz*2.8}, mTip)
	case "bolt":
		cv.capsule(p.x, p.y, p.x+sz*1.2, p.y-sz*1.1, 0.7, mTip)
		cv.capsule(p.x+sz*1.2, p.y-sz*1.1, p.x+sz*0.2, p.y-sz*1.6, 0.7, mTip)
		cv.capsule(p.x+sz*0.2, p.y-sz*1.6, p.x+sz*1.6, p.y-sz*3, 0.7, mTip)
	case "leaf", "feather":
		cv.ellipse(p.x+sz*0.3, p.y-sz*1.2, sz*0.9, sz*1.7, mTip)
	case "fin", "spade", "stinger":
		cv.tri(pt{p.x - sz, p.y}, pt{p.x + sz, p.y}, pt{p.x, p.y - sz*2.4}, mTip)
	case "crystal":
		cv.tri(pt{p.x - sz, p.y - sz}, pt{p.x + sz, p.y - sz}, pt{p.x, p.y - sz*3}, mTip)
		cv.tri(pt{p.x - sz, p.y - sz}, pt{p.x + sz, p.y - sz}, pt{p.x, p.y + sz*0.5}, mTip)
	case "star":
		cv.capsule(p.x-sz, p.y-sz, p.x+sz, p.y-sz, 0.6, mTip)
		cv.capsule(p.x, p.y-sz*2, p.x, p.y, 0.6, mTip)
	case "club", "orb", "drip":
		cv.ellipse(p.x, p.y-sz*0.5, sz*1.1, sz*1.1, mTip)
	default: // tuft, wisp
		cv.ellipse(p.x, p.y-sz*0.4, sz, sz*1.2, mTip)
	}
}

func drawBiped(cv *canvas, ge genome, sp spriteSpec, s, hs, g float64) anatomy {
	legH := ge.f("legH", 2.5, 6) * s
	legR := math.Max(ge.f("legR", 1.3, 2.3)*s, 1)
	bodyRx := ge.f("bodyRx", 5, 8) * s
	bodyRy := ge.f("bodyRy", 4.5, 7) * s
	bodyY := g - legH - bodyRy*0.8
	headR := ge.f("headR", 4.8, 7) * s * hs
	headRx := headR * ge.f("headW", 1, 1.25)
	headY := bodyY - bodyRy*0.7 - headR*0.6

	drawTail(cv, ge, sp, s, midX+bodyRx*0.5, bodyY+bodyRy*0.4, bodyRx)
	cv.part()
	cv.capsuleSym(bodyRx*0.45, bodyY, bodyRx*0.55, g-legR, legR, mBody)
	cv.ellipseSym(bodyRx*0.6, g-legR*0.6, legR*1.4, legR*0.8, mBody)
	bodyZ := cv.part()
	cv.ellipse(midX, bodyY, bodyRx, bodyRy, mBody)
	cv.part()
	shoulderY := bodyY - bodyRy*0.4
	armR := math.Max(ge.f("armR", 1, 1.9)*s, 0.8)
	handDx := bodyRx + ge.f("reach", 1.5, 4)*s
	handY := shoulderY + ge.f("armDrop", 1.5, 6)*s
	cv.capsuleSym(bodyRx*0.7, shoulderY, handDx, handY, armR, mBody)
	if ge.chance("claws", 0.4) {
		cv.dot(int(midX+handDx), int(handY+armR), mAccent)
	}
	cv.part()
	cv.ellipse(midX, headY, headRx, headR, mBody)
	return anatomy{headX: midX, headY: headY, headRx: headRx, headRy: headR, bodyY: bodyY, bodyRx: bodyRx, bodyRy: bodyRy, ground: g, bodyZ: bodyZ,
		belly: ellipseTest(midX, bodyY+bodyRy*0.2, bodyRx*0.62, bodyRy*0.7)}
}

func drawQuadruped(cv *canvas, ge genome, sp spriteSpec, s, hs, g float64) anatomy {
	legH := ge.f("legH", 3, 6.5) * s
	legR := math.Max(ge.f("legR", 1.3, 2.4)*s, 1)
	bodyRx := ge.f("bodyRx", 7.5, 11) * s
	bodyRy := ge.f("bodyRy", 4.5, 7) * s
	bodyY := g - legH - bodyRy*0.7
	headR := ge.f("headR", 4.5, 6.5) * s * hs
	headRx := headR * ge.f("headW", 1.05, 1.35)
	headY := bodyY - bodyRy*0.55 - headR*0.35

	drawTail(cv, ge, sp, s, midX+bodyRx*0.6, bodyY, bodyRx)
	cv.part()
	cv.capsuleSym(bodyRx*0.3, bodyY, bodyRx*0.3, g-legR*0.5, legR*0.9, mBody)
	bodyZ := cv.part()
	cv.ellipse(midX, bodyY, bodyRx, bodyRy, mBody)
	cv.part()
	cv.capsuleSym(bodyRx*0.72, bodyY, bodyRx*0.75, g-legR*0.6, legR, mBody)
	cv.ellipseSym(bodyRx*0.78, g-legR*0.5, legR*1.3, legR*0.7, mBody)
	if ge.chance("mane", 0.35) {
		cv.part()
		cv.ellipse(midX, headY+headR*0.2, headRx*1.35, headR*1.3, mAccent)
	}
	cv.part()
	cv.ellipse(midX, headY, headRx, headR, mBody)
	if ge.chance("snout", 0.6) {
		cv.ellipse(midX, headY+headR*0.55, headRx*0.45, headR*0.38, mBelly)
	}
	return anatomy{headX: midX, headY: headY, headRx: headRx, headRy: headR, bodyY: bodyY, bodyRx: bodyRx, bodyRy: bodyRy, ground: g, bodyZ: bodyZ,
		belly: ellipseTest(midX, bodyY+bodyRy*0.5, bodyRx*0.5, bodyRy*0.6)}
}

func drawSerpent(cv *canvas, ge genome, sp spriteSpec, s, hs, g float64) anatomy {
	coilRx := ge.f("coilRx", 9, 12.5) * s
	coilRy := ge.f("coilRy", 2.6, 3.6) * s
	coils := 2 + ge.n("coils", 2)
	drawTail(cv, ge, sp, s, midX+coilRx*0.7, g-coilRy, coilRx*0.8)
	var bodyZ int
	y := g - coilRy
	rx := coilRx
	for i := 0; i < coils; i++ {
		z := cv.part()
		if i == 0 {
			bodyZ = z
		}
		cv.ellipse(midX, y, rx, coilRy, mBody)
		y -= coilRy * 1.5
		rx *= 0.78
	}
	neckR := math.Max(ge.f("neckR", 1.8, 3)*s, 1.2)
	headR := ge.f("headR", 4, 5.5) * s * hs
	headRx := headR * ge.f("headW", 1.15, 1.5)
	headY := y - ge.f("neck", 2, 6)*s - headR*0.5
	cv.part()
	cv.capsule(midX, y+coilRy, midX, headY, neckR, mBody)
	if ge.chance("hood", 0.4) {
		cv.part()
		cv.ellipse(midX, headY+headR*0.6, headRx*1.5, headR*1.4, mAccent)
	}
	cv.part()
	cv.ellipse(midX, headY, headRx, headR, mBody)
	return anatomy{headX: midX, headY: headY, headRx: headRx, headRy: headR, bodyY: g - coilRy, bodyRx: coilRx, bodyRy: coilRy, ground: g, bodyZ: bodyZ,
		belly: ellipseTest(midX, g-coilRy*0.4, coilRx*0.7, coilRy*0.6)}
}

func drawAvian(cv *canvas, ge genome, sp spriteSpec, s, hs, g float64) anatomy {
	legH := ge.f("legH", 2.5, 4.5) * s
	bodyRx := ge.f("bodyRx", 5.5, 8) * s
	bodyRy := ge.f("bodyRy", 6, 8.5) * s
	bodyY := g - legH - bodyRy
	headR := ge.f("headR", 4, 5.5) * s * hs
	headY := bodyY - bodyRy*0.8 - headR*0.5
	span := ge.f("span", 5, 10) * s
	lift := ge.f("lift", -2, 6) * s

	drawTail(cv, ge, sp, s, midX+bodyRx*0.4, bodyY+bodyRy*0.6, bodyRx)
	// Wings spread behind the body.
	cv.part()
	sh := pt{bodyRx * 0.4, bodyY - bodyRy*0.35}
	tip := pt{bodyRx + span, bodyY - bodyRy*0.35 - lift}
	low := pt{bodyRx + span*0.55, bodyY + bodyRy*0.55}
	cv.triSym(sh, tip, low, mBody)
	cv.triSym(sh, low, pt{bodyRx * 0.4, bodyY + bodyRy*0.5}, mBody)
	feathers := 2 + ge.n("feathers", 3)
	for i := 0; i < feathers; i++ {
		t := float64(i+1) / float64(feathers+1)
		bx := low.x + (tip.x-low.x)*t
		by := low.y + (tip.y-low.y)*t
		cv.triSym(pt{bx - 1.2, by - 1}, pt{bx + 1.2, by - 1}, pt{bx + 0.4, by + 2*s}, mWing)
	}
	cv.part()
	cv.capsuleSym(bodyRx*0.35, bodyY+bodyRy*0.8, bodyRx*0.35, g-0.8, math.Max(0.7*s, 0.6), mAccent)
	cv.capsuleSym(bodyRx*0.35, g-0.6, bodyRx*0.35+1.6*s, g-0.6, 0.6, mAccent)
	bodyZ := cv.part()
	cv.ellipse(midX, bodyY, bodyRx, bodyRy, mBody)
	cv.part()
	cv.ellipse(midX, headY, headR*1.05, headR, mBody)
	// The beak, below the eyes.
	cv.part()
	bl := math.Max(ge.f("beak", 1.5, 3.2)*s, 1.2)
	by := headY + headR*0.25
	cv.tri(pt{midX - bl*0.7, by}, pt{midX + bl*0.7, by}, pt{midX, by + bl*1.4}, mAccent)
	return anatomy{headX: midX, headY: headY - headR*0.15, headRx: headR * 1.05, headRy: headR, bodyY: bodyY, bodyRx: bodyRx, bodyRy: bodyRy, ground: g, bodyZ: bodyZ,
		belly: ellipseTest(midX, bodyY+bodyRy*0.25, bodyRx*0.6, bodyRy*0.6)}
}

func drawBlob(cv *canvas, ge genome, sp spriteSpec, s, g float64) anatomy {
	rx := ge.f("bodyRx", 8, 12.5) * s
	ry := ge.f("bodyRy", 7, 10.5) * s
	cy := g - ry
	drawTail(cv, ge, sp, s, midX+rx*0.6, cy, rx)
	if ge.chance("nubs", 0.6) {
		cv.part()
		cv.ellipseSym(rx*0.95, cy+ry*0.2, 1.8*s+0.6, 1.4*s+0.6, mBody)
	}
	bodyZ := cv.part()
	cv.ellipse(midX, cy, rx, ry, mBody)
	// A flat bottom where it sits.
	cv.ellipse(midX, g-ry*0.25, rx*0.95, ry*0.3, mBody)
	if ge.chance("drips", 0.45) {
		n := 1 + ge.n("dripN", 2)
		for i := 0; i < n; i++ {
			dx := rx * (0.3 + 0.25*float64(i))
			cv.ellipseSym(dx, g+0.2, 1.1, 1.3, mBody)
		}
	}
	return anatomy{headX: midX, headY: cy - ry*0.15, headRx: rx * 0.8, headRy: ry * 0.6, bodyY: cy, bodyRx: rx, bodyRy: ry, ground: g, bodyZ: bodyZ,
		belly: ellipseTest(midX, cy+ry*0.55, rx*0.6, ry*0.4)}
}

func drawInsect(cv *canvas, ge genome, sp spriteSpec, s, hs, g float64) anatomy {
	abRx := ge.f("abRx", 5.5, 8.5) * s
	abRy := ge.f("abRy", 4.5, 7) * s
	abY := g - abRy - 0.5
	thR := ge.f("thR", 3.2, 4.8) * s
	thY := abY - abRy*0.75 - thR*0.4
	headR := ge.f("headR", 3.5, 5) * s * hs
	headRx := headR * ge.f("headW", 1.1, 1.4)
	headY := thY - thR*0.7 - headR*0.6

	if ge.chance("wings", 0.6) || sp.tier == 2 {
		cv.part()
		wl := ge.f("wingLen", 6, 11) * s
		cv.ellipseSym(thR+wl*0.45, thY-wl*0.35, wl*0.55, wl*0.32, mWing)
		if ge.chance("wings2", 0.5) {
			cv.ellipseSym(thR+wl*0.3, thY+wl*0.25, wl*0.38, wl*0.24, mWing)
		}
	}
	drawTail(cv, ge, sp, s, midX+abRx*0.4, abY+abRy*0.5, abRx*0.6)
	cv.part()
	legR := math.Max(0.65*s, 0.55)
	for i := 0; i < 3; i++ {
		fy := thY - thR*0.2 + float64(i)*thR*0.55
		out := thR + ge.f("legOut", 2.5, 5)*s
		cv.capsuleSym(thR*0.6, fy, out, fy+1.5*s, legR, mBody)
		cv.capsuleSym(out, fy+1.5*s, out+0.5, g-0.8-float64(2-i)*0.4, legR, mBody)
	}
	bodyZ := cv.part()
	cv.ellipse(midX, abY, abRx, abRy, mBody)
	cv.part()
	cv.ellipse(midX, thY, thR*1.1, thR, mBody)
	cv.part()
	cv.ellipse(midX, headY, headRx, headR, mBody)
	if ge.chance("mandibles", 0.4) {
		cv.part()
		cv.capsuleSym(headRx*0.35, headY+headR*0.8, headRx*0.1, headY+headR*1.3, 0.6, mAccent)
	}
	return anatomy{headX: midX, headY: headY, headRx: headRx, headRy: headR, bodyY: abY, bodyRx: abRx, bodyRy: abRy, ground: g, bodyZ: bodyZ,
		belly: func(px, py float64) bool { return int(py)%3 == 0 && ellipseTest(midX, abY, abRx, abRy)(px, py) }}
}

func drawAquatic(cv *canvas, ge genome, sp spriteSpec, s, g float64) anatomy {
	rx := ge.f("bodyRx", 6, 9.5) * s
	ry := ge.f("bodyRy", 7, 10) * s
	cy := g - ry - 1.5
	finL := ge.f("finL", 3.5, 7) * s
	cv.part()
	// The tail fin, below.
	cv.triSym(pt{0, cy + ry*0.7}, pt{finL * 0.9, g + 0.3}, pt{finL * 0.2, g - 0.6}, mWing)
	// The dorsal fin, above.
	cv.tri(pt{midX - finL*0.35, cy - ry*0.7}, pt{midX + finL*0.35, cy - ry*0.7}, pt{midX + finL*0.15*ge.f("dorsalLean", -1, 1), cy - ry - finL*0.9}, mWing)
	bodyZ := cv.part()
	cv.ellipse(midX, cy, rx, ry, mBody)
	cv.part()
	// Side fins.
	cv.triSym(pt{rx * 0.8, cy - ry*0.05}, pt{rx * 0.8, cy + ry*0.4}, pt{rx + finL, cy + ry*0.5 - ge.f("finLift", 0, 4)*s}, mWing)
	if ge.chance("whiskers", 0.35) {
		cv.part()
		cv.capsuleSym(rx*0.35, cy+ry*0.25, rx*0.9, cy+ry*0.5, 0.55, mAccent)
	}
	return anatomy{headX: midX, headY: cy - ry*0.25, headRx: rx * 0.9, headRy: ry * 0.5, bodyY: cy, bodyRx: rx, bodyRy: ry, ground: g, bodyZ: bodyZ,
		belly: ellipseTest(midX, cy+ry*0.55, rx*0.6, ry*0.45)}
}

func drawGolem(cv *canvas, ge genome, sp spriteSpec, s, hs, g float64) anatomy {
	legH := ge.f("legH", 2, 4) * s
	bodyW := ge.f("bodyW", 6.5, 9.5) * s
	bodyH := ge.f("bodyH", 9, 13) * s
	top := g - legH - bodyH
	headR := ge.f("headR", 3.2, 4.5) * s * hs
	headY := top - headR*0.6
	fist := ge.f("fist", 2.5, 4) * s

	drawTail(cv, ge, sp, s, midX+bodyW*0.6, top+bodyH*0.6, bodyW)
	cv.part()
	cv.capsuleSym(bodyW*0.5, top+bodyH, bodyW*0.55, g-1.2*s, math.Max(1.8*s, 1.2), mBody)
	cv.ellipseSym(bodyW*0.6, g-1, 2.4*s+0.4, 1.2*s+0.4, mBody)
	bodyZ := cv.part()
	// A rounded block: a rectangle with ellipse ends.
	cv.fill(func(px, py float64) bool {
		dx := math.Abs(px-midX) / bodyW
		dy := math.Abs(py-(top+bodyH/2)) / (bodyH / 2)
		return math.Pow(dx, 4)+math.Pow(dy, 4) <= 1
	}, midX-bodyW, top, midX+bodyW, top+bodyH, mBody)
	cv.part()
	cv.capsuleSym(bodyW*0.85, top+bodyH*0.2, bodyW+fist*0.6, top+bodyH*0.6, math.Max(1.6*s, 1), mBody)
	cv.ellipseSym(bodyW+fist*0.7, top+bodyH*0.75, fist, fist*0.95, mBody)
	cv.part()
	cv.ellipse(midX, headY, headR*1.15, headR, mBody)
	return anatomy{headX: midX, headY: headY, headRx: headR * 1.15, headRy: headR, bodyY: top + bodyH/2, bodyRx: bodyW, bodyRy: bodyH / 2, ground: g, bodyZ: bodyZ,
		belly: func(px, py float64) bool {
			return math.Abs(px-midX) < bodyW*0.45 && py > top+bodyH*0.3 && py < top+bodyH*0.8
		}}
}

func drawSpirit(cv *canvas, ge genome, sp spriteSpec, s, g float64) anatomy {
	r := ge.f("bodyR", 6, 9) * s
	cy := g - r - ge.f("wisp", 6, 10)*s
	wispX := ge.f("wispCurl", -4, 4) * s
	bodyZ := cv.part()
	// A teardrop: the round head and a wisp tapering down from it.
	cv.ellipse(midX, cy, r*ge.f("bodyW", 0.95, 1.15), r, mBody)
	cv.tri(pt{midX - r*0.85, cy + r*0.3}, pt{midX + r*0.85, cy + r*0.3}, pt{midX + wispX, g}, mBody)
	cv.part()
	armR := math.Max(1.1*s, 0.8)
	cv.capsuleSym(r*0.8, cy+r*0.4, r+ge.f("reach", 1.5, 4)*s, cy+r*0.6+ge.f("armDrop", -3, 2)*s, armR, mBody)
	return anatomy{headX: midX, headY: cy, headRx: r, headRy: r * 0.9, bodyY: cy, bodyRx: r, bodyRy: r, ground: g, bodyZ: bodyZ,
		belly: ellipseTest(midX, cy+r*0.6, r*0.55, r*0.45)}
}

// drawHeadgear puts ears, horns, antennae or a crest on top of the head. Later stages grow
// them longer.
func drawHeadgear(cv *canvas, ge genome, sp spriteSpec, s float64, an anatomy) {
	w := earWeights[sp.arch.id]
	style := earStyles[ge.weighted("ears", w)]
	if sp.legendary && style == "none" {
		style = "crest"
	}
	top := an.headY - an.headRy
	rx := an.headRx
	grow := 1 + 0.25*float64(sp.tier)
	length := ge.f("earLen", 3, 6.5) * s * grow
	lean := ge.f("earLean", -1, 3) * s
	cv.part()
	switch style {
	case "pointy":
		cv.triSym(pt{rx * 0.25, top + an.headRy*0.3}, pt{rx * 0.85, top + an.headRy*0.55}, pt{rx*0.6 + lean, top - length}, mBody)
		if ge.chance("innerEar", 0.6) {
			cv.triSym(pt{rx * 0.42, top + an.headRy*0.2}, pt{rx * 0.72, top + an.headRy*0.32}, pt{rx*0.6 + lean*0.85, top - length*0.6}, mBelly)
		}
	case "round":
		r := math.Max(ge.f("earR", 1.8, 3.2)*s*grow, 1.2)
		cv.ellipseSym(rx*0.7, top+r*0.4, r, r, mBody)
	case "long":
		cv.capsuleSym(rx*0.35, top+1.5, rx*0.45+lean, top-length*1.4, math.Max(1.4*s, 1), mBody)
	case "horns":
		cv.triSym(pt{rx * 0.25, top + 1.8}, pt{rx * 0.6, top + 1.8}, pt{rx*0.55 + lean*1.3, top - length}, mAccent)
	case "antennae":
		tx := rx*0.5 + lean*1.2
		cv.capsuleSym(rx*0.2, top+1, tx, top-length, 0.55, mBody)
		cv.ellipseSym(tx, top-length, 1.3, 1.3, mAccent)
	case "crest":
		n := 1 + ge.n("crestN", 3)
		for i := 0; i < n; i++ {
			dx := float64(i) * 1.6 * s
			l := length * (1 - 0.2*float64(i))
			cv.triSym(pt{dx - 1.3, top + 1.5}, pt{dx + 1.3, top + 1.5}, pt{dx + lean*0.3*float64(i), top - l}, mAccent)
		}
	case "unicorn":
		cv.tri(pt{midX - 1.3, top + 1.2}, pt{midX + 1.3, top + 1.2}, pt{midX, top - length*1.2}, mAccent)
	}
	// Final stages and legendaries grow spikes along the shoulders.
	if (sp.tier == 2 || sp.legendary) && ge.chance("spikes", 0.5) && !sp.arch.floats {
		cv.part()
		for i := 0; i < 2; i++ {
			dx := an.bodyRx * (0.45 + 0.25*float64(i))
			y := an.bodyY - an.bodyRy*(0.85-0.25*float64(i))
			cv.triSym(pt{dx - 1.2, y + 1}, pt{dx + 1.2, y + 1}, pt{dx + 0.8, y - 2.5*s}, mAccent)
		}
	}
}

func drawPattern(cv *canvas, ge genome, an anatomy) {
	switch []string{"none", "none", "spots", "stripes", "mask", "socks"}[ge.n("pattern", 6)] {
	case "spots":
		n := 1 + ge.n("spotN", 3)
		for i := range n {
			dx := ge.f("spotX"+string(rune('a'+i)), 1, an.bodyRx*0.8)
			y := an.bodyY + ge.f("spotY"+string(rune('a'+i)), -an.bodyRy*0.6, an.bodyRy*0.6)
			r := ge.f("spotR"+string(rune('a'+i)), 0.8, 1.6)
			cv.paint(ellipseTest(midX+dx, y, r, r), mPattern, mBody)
			cv.paint(ellipseTest(midX-dx, y, r, r), mPattern, mBody)
		}
	case "stripes":
		period := 3 + ge.n("stripeP", 2)
		cv.paint(func(px, py float64) bool {
			return int(py)%period == 0 && math.Abs(py-an.bodyY) < an.bodyRy*0.8
		}, mPattern, mBody)
	case "mask":
		cv.paint(func(px, py float64) bool {
			return math.Abs(py-an.headY) < 1.6 && math.Abs(px-an.headX) < an.headRx
		}, mPattern, mBody)
	case "socks":
		cv.paint(func(px, py float64) bool { return py > an.ground-2.2 }, mPattern, mBody)
	}
}

func drawFace(cv *canvas, ge genome, sp spriteSpec, an anatomy) {
	size := 2
	if sp.tier == 0 || an.headRx > 6.5 {
		size = 3
	}
	// ex is the left eye's outer column; its inner column stays left of 19 so the two eyes
	// never meet in the middle.
	ex := min(int(math.Round(midX-an.headRx*ge.f("eyeSpread", 0.32, 0.5)-float64(size)/2)), 18-size)
	ey := int(math.Round(an.headY + an.headRy*ge.f("eyeY", -0.15, 0.12)))
	style := []string{"round", "round", "dot", "fierce", "sleepy"}[ge.n("eyes", 5)]
	if sp.tier == 0 && style == "fierce" {
		style = "round"
	}
	// x runs inward from the outer edge of the left eye; dot mirrors it to the right eye.
	switch style {
	case "dot":
		for dy := 0; dy < size; dy++ {
			for dx := 0; dx < size-1; dx++ {
				cv.dot(ex+dx, ey+dy, mPupil)
			}
		}
		cv.dot(ex, ey, mShine)
	case "sleepy":
		for dx := -1; dx < size; dx++ {
			cv.dot(ex+dx, ey+1, mPupil)
		}
	default:
		for dy := 0; dy <= size; dy++ {
			for dx := 0; dx < size; dx++ {
				cv.dot(ex+dx, ey+dy, mEyeWhite)
			}
		}
		for dy := 1; dy <= size; dy++ {
			cv.dot(ex+size-1, ey+dy, mPupil)
			if size == 3 {
				cv.dot(ex+size-2, ey+dy, mPupil)
			}
		}
		cv.dot(ex+size-1, ey+1, mShine)
		if style == "fierce" {
			cv.dot(ex+size-1, ey, mPupil)
			cv.dot(ex+size-2, ey, mPupil)
		}
	}
	my := ey + size + 1
	if sp.arch.id == "avian" {
		my = -1 // the beak is the mouth
	}
	if my > 0 {
		switch []string{"smile", "line", "fang", "open", "none"}[ge.n("mouth", 5)] {
		case "smile":
			cv.dot(19, my, mMouth)
			cv.dot(18, my-1, mMouth)
		case "line":
			cv.dot(19, my, mMouth)
		case "fang":
			cv.dot(19, my, mMouth)
			cv.dot(18, my, mMouth)
			if c := cv.at(18, my+1); c.mat != mEmpty {
				cv.c[my+1][18].mat = mEyeWhite
			}
		case "open":
			cv.dot(19, my, mMouth)
			cv.dot(19, my+1, mCheek)
		}
	}
	if ge.chance("cheeks", 0.3) || sp.primary.id == "electric" || sp.primary.id == "fairy" {
		cv.dot(ex-1, ey+size+1, mCheek)
		cv.dot(ex, ey+size+1, mCheek)
	}
}

func ellipseTest(cx, cy, rx, ry float64) func(px, py float64) bool {
	return func(px, py float64) bool {
		dx, dy := (px-cx)/rx, (py-cy)/ry
		return dx*dx+dy*dy <= 1
	}
}
