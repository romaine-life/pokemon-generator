package gen

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode"
)

func TestSameSeedSameCreature(t *testing.T) {
	a, _ := json.Marshal(Generate(Options{Seed: "Nelson"}))
	b, _ := json.Marshal(Generate(Options{Seed: "  nelson  "}))
	if string(a) != string(b) {
		t.Fatal("seeds differing only in case and spacing gave different creatures")
	}
	c, _ := json.Marshal(Generate(Options{Seed: "nelson2"}))
	if string(a) == string(c) {
		t.Fatal("different seeds gave the same creature")
	}
}

func TestTypeLock(t *testing.T) {
	for _, ty := range Types() {
		for i := range 20 {
			m := Generate(Options{Seed: fmt.Sprint("lock ", i), Type: ty.ID})
			if m.Types[0].ID != ty.ID {
				t.Fatalf("type %s locked, got %s", ty.ID, m.Types[0].ID)
			}
			if len(m.Types) == 2 && m.Types[1].ID == ty.ID {
				t.Fatalf("seed %d: both types are %s", i, ty.ID)
			}
		}
	}
	if ValidType("plasma") || !ValidType("") || !ValidType("fire") {
		t.Fatal("ValidType disagrees with the type list")
	}
}

func TestStagePick(t *testing.T) {
	for i := range 200 {
		seed := fmt.Sprint("stage ", i)
		m := Generate(Options{Seed: seed})
		if m.Stage < 1 || m.Stage > m.Stages || len(m.Evolution) != m.Stages {
			t.Fatalf("%s: stage %d of %d with %d evolutions", seed, m.Stage, m.Stages, len(m.Evolution))
		}
		for st := 1; st <= m.Stages; st++ {
			e := Generate(Options{Seed: seed, Stage: st})
			if e.Stage != st || e.Name != m.Evolution[st-1].Name || e.Number != m.Evolution[st-1].Number {
				t.Fatalf("%s stage %d: got %s #%d, the line says %s #%d", seed, st, e.Name, e.Number, m.Evolution[st-1].Name, m.Evolution[st-1].Number)
			}
		}
	}
}

// Every creature is complete, readable and fits its sprite grid.
func TestCreaturesAreWellFormed(t *testing.T) {
	legendaries, shinies, plans := 0, 0, map[string]bool{}
	for i := range 2000 {
		m := Generate(Options{Seed: fmt.Sprint("well formed ", i)})
		where := fmt.Sprintf("seed %q (%s)", m.Seed, m.Name)
		plans[m.BodyPlan] = true
		if m.Legendary {
			legendaries++
		}
		if m.Shiny {
			shinies++
		}

		if n := len(m.Name); n < 3 || n > 12 || !unicode.IsUpper(rune(m.Name[0])) || strings.IndexFunc(m.Name, func(r rune) bool { return !unicode.IsLetter(r) }) >= 0 {
			t.Fatalf("%s: bad name", where)
		}
		if m.Number < 1 || m.Number > 9999 || m.Category == "" || m.Entry == "" || m.HeightM <= 0 || m.WeightKg <= 0 {
			t.Fatalf("%s: missing facts: %+v", where, m)
		}
		if len(m.Stats) != 6 {
			t.Fatalf("%s: %d stats", where, len(m.Stats))
		}
		sum := 0
		for _, s := range m.Stats {
			if s.Value < 5 || s.Value > 255 {
				t.Fatalf("%s: %s is %d", where, s.Name, s.Value)
			}
			sum += s.Value
		}
		if sum != m.Total || m.Total < 240 || m.Total > 720 {
			t.Fatalf("%s: total %d, stats sum to %d", where, m.Total, sum)
		}
		if len(m.Moves) != 4 {
			t.Fatalf("%s: %d moves", where, len(m.Moves))
		}
		names := map[string]bool{}
		for _, mv := range m.Moves {
			if names[mv.Name] {
				t.Fatalf("%s: move %s twice", where, mv.Name)
			}
			names[mv.Name] = true
			if (mv.Category == "Status") != (mv.Power == 0) || mv.Power > 120 || mv.Accuracy < 70 || mv.Accuracy > 100 || mv.PP < 5 {
				t.Fatalf("%s: bad move %+v", where, mv)
			}
		}
		hidden := 0
		for _, a := range m.Abilities {
			if a.Hidden {
				hidden++
			}
		}
		if hidden != 1 || len(m.Abilities) < 2 {
			t.Fatalf("%s: abilities %+v", where, m.Abilities)
		}
		checkSprite(t, where, m.Sprite)
	}
	if len(plans) != len(archetypes) {
		t.Errorf("only %d of %d body plans came up in 2000 seeds", len(plans), len(archetypes))
	}
	// 3% legendary and 1 in 64 shiny: well inside these bounds over 2000 seeds.
	if legendaries < 20 || legendaries > 110 || shinies < 10 || shinies > 70 {
		t.Errorf("%d legendaries and %d shinies in 2000 seeds", legendaries, shinies)
	}
}

func checkSprite(t *testing.T, where string, s Sprite) {
	t.Helper()
	if s.Width != gridW || s.Height != gridH || len(s.Pixels) != gridW*gridH || s.Palette[0] != "#00000000" {
		t.Fatalf("%s: sprite shape %dx%d, %d pixels", where, s.Width, s.Height, len(s.Pixels))
	}
	solid := 0
	for i, p := range s.Pixels {
		if p < 0 || p >= len(s.Palette) {
			t.Fatalf("%s: pixel %d names color %d of %d", where, i, p, len(s.Palette))
		}
		if p == 0 || s.Palette[p] == shadowColor {
			continue
		}
		solid++
		x, y := i%gridW, i/gridW
		if x == 0 || y == 0 || x == gridW-1 {
			t.Fatalf("%s: drawn at the grid's edge (%d,%d): the creature didn't fit", where, x, y)
		}
	}
	if solid < 80 {
		t.Fatalf("%s: only %d pixels drawn", where, solid)
	}
}

func TestJoinName(t *testing.T) {
	for _, c := range []struct{ prefix, link, ending, want string }{
		{"embr", "a", "pup", "Embrapup"},
		{"aqu", "o", "eel", "Aquel"},
		{"volt", "i", "tot", "Voltot"},
		{"zap", "o", "it", "Zapit"},
		{"thunderstormy", "o", "king", "Thunderstorm"},
	} {
		if got := joinName(c.prefix, c.link, c.ending); got != c.want {
			t.Errorf("joinName(%q, %q, %q) = %q, want %q", c.prefix, c.link, c.ending, got, c.want)
		}
	}
}
