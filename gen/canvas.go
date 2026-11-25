package gen

import "math"

// The sprite grid. Creatures face the viewer and are drawn mirror-symmetric about the
// vertical center line, except for a tail, which curls out on one side.
const (
	gridW = 40
	gridH = 40
	midX  = gridW / 2.0
)

// Materials a cell can be painted with. Each gets its own color, and the colorizer shades
// geometry materials and leaves the face flat.
type material uint8

const (
	mEmpty material = iota
	mBody
	mBelly
	mAccent
	mPattern
	mWing
	mEyeWhite
	mPupil
	mShine
	mMouth
	mCheek
	mTip
	mAura
	mShadow
)

// flat materials aren't shaded and don't form part edges.
func (m material) flat() bool {
	switch m {
	case mEyeWhite, mPupil, mShine, mMouth, mCheek:
		return true
	}
	return false
}

// glow materials sit behind the creature without an outline.
func (m material) glow() bool { return m == mAura || m == mShadow }

type cell struct {
	mat material
	z   int // which drawn part the cell belongs to; later parts are in front
}

type canvas struct {
	c       [gridH][gridW]cell
	z       int
	clipped bool // something was drawn outside the safe area, so the creature must shrink
}

// part starts a new geometric part, in front of everything drawn so far.
func (cv *canvas) part() int {
	cv.z++
	return cv.z
}

// set paints a cell inside the safe area, which leaves a free column or row on every side
// for the outline and a margin between it and the edge.
func (cv *canvas) set(x, y int, m material) {
	if x < 2 || y < 2 || x > gridW-3 || y > gridH-2 {
		if !m.glow() {
			cv.clipped = true
		}
		return
	}
	cv.c[y][x] = cell{m, cv.z}
}

func (cv *canvas) at(x, y int) cell {
	if x < 0 || y < 0 || x >= gridW || y >= gridH {
		return cell{}
	}
	return cv.c[y][x]
}

// fill sets every cell whose center is inside the shape.
func (cv *canvas) fill(inside func(px, py float64) bool, minX, minY, maxX, maxY float64, m material) {
	for y := int(math.Floor(minY)) - 1; y <= int(math.Ceil(maxY))+1; y++ {
		for x := int(math.Floor(minX)) - 1; x <= int(math.Ceil(maxX))+1; x++ {
			if inside(float64(x)+0.5, float64(y)+0.5) {
				cv.set(x, y, m)
			}
		}
	}
}

func (cv *canvas) ellipse(cx, cy, rx, ry float64, m material) {
	rx, ry = math.Max(rx, 0.6), math.Max(ry, 0.6)
	cv.fill(func(px, py float64) bool {
		dx, dy := (px-cx)/rx, (py-cy)/ry
		return dx*dx+dy*dy <= 1
	}, cx-rx, cy-ry, cx+rx, cy+ry, m)
}

// ellipseSym draws an ellipse dx right of center and its mirror image.
func (cv *canvas) ellipseSym(dx, cy, rx, ry float64, m material) {
	cv.ellipse(midX+dx, cy, rx, ry, m)
	if math.Abs(dx) > 0.01 {
		cv.ellipse(midX-dx, cy, rx, ry, m)
	}
}

// capsule is a line segment with round ends and radius r.
func (cv *canvas) capsule(x1, y1, x2, y2, r float64, m material) {
	r = math.Max(r, 0.55)
	vx, vy := x2-x1, y2-y1
	l2 := vx*vx + vy*vy
	cv.fill(func(px, py float64) bool {
		t := 0.0
		if l2 > 0 {
			t = clamp(((px-x1)*vx+(py-y1)*vy)/l2, 0, 1)
		}
		ex, ey := px-(x1+t*vx), py-(y1+t*vy)
		return ex*ex+ey*ey <= r*r
	}, math.Min(x1, x2)-r, math.Min(y1, y2)-r, math.Max(x1, x2)+r, math.Max(y1, y2)+r, m)
}

func (cv *canvas) capsuleSym(dx1, y1, dx2, y2, r float64, m material) {
	cv.capsule(midX+dx1, y1, midX+dx2, y2, r, m)
	cv.capsule(midX-dx1, y1, midX-dx2, y2, r, m)
}

type pt struct{ x, y float64 }

func (cv *canvas) tri(a, b, c pt, m material) {
	sign := func(p, q, r pt) float64 { return (p.x-r.x)*(q.y-r.y) - (q.x-r.x)*(p.y-r.y) }
	cv.fill(func(px, py float64) bool {
		p := pt{px, py}
		d1, d2, d3 := sign(p, a, b), sign(p, b, c), sign(p, c, a)
		neg := d1 < 0 || d2 < 0 || d3 < 0
		pos := d1 > 0 || d2 > 0 || d3 > 0
		return !(neg && pos)
	}, math.Min(a.x, math.Min(b.x, c.x)), math.Min(a.y, math.Min(b.y, c.y)),
		math.Max(a.x, math.Max(b.x, c.x)), math.Max(a.y, math.Max(b.y, c.y)), m)
}

// triSym draws a triangle given in offsets from center, and its mirror image.
func (cv *canvas) triSym(a, b, c pt, m material) {
	cv.tri(pt{midX + a.x, a.y}, pt{midX + b.x, b.y}, pt{midX + c.x, c.y}, m)
	cv.tri(pt{midX - a.x, a.y}, pt{midX - b.x, b.y}, pt{midX - c.x, c.y}, m)
}

// paint recolors cells already drawn with one of the given materials, where inside holds,
// without changing which part they belong to.
func (cv *canvas) paint(inside func(px, py float64) bool, m material, over ...material) {
	for y := range gridH {
		for x := range gridW {
			c := &cv.c[y][x]
			if c.mat == mEmpty || !inside(float64(x)+0.5, float64(y)+0.5) {
				continue
			}
			for _, o := range over {
				if c.mat == o {
					c.mat = m
					break
				}
			}
		}
	}
}

// dot paints one cell and its mirror image, where the creature already is (a face never
// spills off its head).
func (cv *canvas) dot(x, y int, m material) {
	for _, xx := range []int{x, gridW - 1 - x} {
		if xx < 1 || y < 1 || xx > gridW-2 || y > gridH-2 {
			continue
		}
		if c := &cv.c[y][xx]; c.mat != mEmpty && !c.mat.glow() {
			c.mat = m
		}
	}
}
