package gen

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"testing"
)

// TestPreviewSheet draws a contact sheet of creatures for eyeballing sprite changes. It runs
// only when PREVIEW_OUT names the PNG to write; PREVIEW_LINES=1 draws evolution lines, one
// family's stages side by side, instead of one creature per seed.
func TestPreviewSheet(t *testing.T) {
	out := os.Getenv("PREVIEW_OUT")
	if out == "" {
		t.Skip("set PREVIEW_OUT to write a contact sheet")
	}
	const cols, rows, scale = 9, 6, 4
	cellPx := gridW * scale
	sheet := image.NewNRGBA(image.Rect(0, 0, cols*cellPx, rows*cellPx))
	draw.Draw(sheet, sheet.Bounds(), &image.Uniform{C: parseHex("#f4f1e8")}, image.Point{}, draw.Src)
	lines := os.Getenv("PREVIEW_LINES") != ""
	for i := range cols * rows {
		o := Options{Seed: fmt.Sprint("preview ", i)}
		if lines {
			o = Options{Seed: fmt.Sprint("line ", i/3), Stage: i%3 + 1}
		}
		m := Generate(o)
		x, y := (i%cols)*cellPx, (i/cols)*cellPx
		draw.Draw(sheet, image.Rect(x, y, x+cellPx, y+cellPx), spriteImage(m.Sprite, scale), image.Point{}, draw.Over)
		t.Logf("%-10s %-12s %-9s %d/%d", o.Seed, m.Name, m.BodyPlan, m.Stage, m.Stages)
	}
	f, err := os.Create(out)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, sheet); err != nil {
		t.Fatal(err)
	}
}
