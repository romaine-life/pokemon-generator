package gen

import (
	"fmt"
	"strings"
)

var (
	physicalWords = []string{"Fang", "Claw", "Rush", "Slam", "Kick", "Lash", "Crash", "Bite", "Strike", "Tackle", "Fist", "Charge"}
	specialWords  = []string{"Beam", "Burst", "Wave", "Pulse", "Storm", "Spiral", "Orb", "Ray", "Cannon", "Breath", "Nova", "Blast"}
	statusWords   = []string{"Dance", "Veil", "Shroud", "Song", "Glare", "Ward", "Focus", "Mist", "Howl", "Lullaby", "Trick", "Bloom"}

	physicalEffects = []string{"May make the target flinch.", "A high critical-hit ratio.", "Hits two times in a row.", "Always goes first.", "The user takes a little recoil.", "Lowers the target's Defense."}
	specialEffects  = []string{"May lower the target's Sp. Def.", "Hits all opponents.", "Charges on the first turn, strikes on the second.", "Never misses.", "Heals the user by half the damage dealt.", "May confuse the target."}
	statusEffects   = []string{"Sharply raises the user's Attack.", "Raises the user's Sp. Atk and Speed.", "Lowers the target's accuracy.", "Puts the target to sleep.", "Protects the user from moves this turn.", "Restores half the user's max HP.", "Raises Defense and Sp. Def."}
)

func moves(f family, g *rng, tier int, st []Stat) []Move {
	// Physical or special, leaning on whichever attack stat is higher.
	physical := 0.5
	if atk, spa := st[1].Value, st[3].Value; atk+spa > 0 {
		physical = clamp(float64(atk)/float64(atk+spa)*1.4-0.2, 0.15, 0.85)
	}
	normal, _ := typeByID("normal")
	moveTypes := []elemType{f.primary, f.primary, normal}
	if f.secondary != nil {
		moveTypes[2] = *f.secondary
	} else if g.chance(0.5) {
		moveTypes[2] = pick(g, types)
	}
	used := map[string]bool{}
	var out []Move
	for i, t := range moveTypes {
		category := "Special"
		if g.chance(physical) {
			category = "Physical"
		}
		out = append(out, attack(g, t, category, tier, i, used))
	}
	out = append(out, status(g, f.primary, used))
	return out
}

func attack(g *rng, t elemType, category string, tier, slot int, used map[string]bool) Move {
	words, effects := specialWords, specialEffects
	if category == "Physical" {
		words, effects = physicalWords, physicalEffects
	}
	name := uniqueName(g, t, words, used)
	power := 5 * g.intn(8+tier*2+slot, min(14+tier*4+slot*2, 24)) // 40 to 120
	m := Move{Name: name, Type: ref(t), Category: category, Power: power, Accuracy: accuracy(g, power), PP: pp(power)}
	if g.chance(0.6) {
		m.Description = pick(g, effects)
	} else {
		m.Description = fmt.Sprintf("A %s attack.", strings.ToLower(describePower(power)))
	}
	return m
}

func status(g *rng, t elemType, used map[string]bool) Move {
	return Move{Name: uniqueName(g, t, statusWords, used), Type: ref(t), Category: "Status", Accuracy: 100, PP: 5 * g.intn(2, 6), Description: pick(g, statusEffects)}
}

func uniqueName(g *rng, t elemType, words []string, used map[string]bool) string {
	for {
		name := pick(g, t.moveWords) + " " + pick(g, words)
		if !used[name] {
			used[name] = true
			return name
		}
	}
}

func accuracy(g *rng, power int) int {
	if power >= 100 {
		return 5 * g.intn(14, 18)
	}
	return 5 * g.intn(18, 20)
}

func pp(power int) int {
	switch {
	case power >= 100:
		return 5
	case power >= 80:
		return 10
	case power >= 60:
		return 15
	default:
		return 25
	}
}

func describePower(power int) string {
	switch {
	case power >= 110:
		return "Devastating"
	case power >= 90:
		return "Powerful"
	case power >= 60:
		return "Solid"
	default:
		return "Quick, light"
	}
}

var (
	openings = []string{"It lives %s.", "It is usually found %s.", "Wild ones gather %s.", "Travelers spot it %s."}
	quirks   = []string{
		"Its %s are softer than they look.",
		"It is very proud of its %s and cleans them constantly.",
		"Researchers still argue about what its %s are for.",
		"Its %s glow faintly when it is happy.",
		"Children in nearby towns trade drawings of its %s.",
		"It shows off its %s to impress others of its kind.",
	}
	stageLines = map[int][]string{
		1: {"It is curious about everything and easily startled.", "It sticks close to its parents until it learns to forage.", "It tumbles over its own feet when it gets excited."},
		2: {"It has grown bold and likes to test its strength.", "It wanders far from home to find stronger rivals.", "Its temper is short, but it is fiercely loyal."},
		3: {"It is calm and wise, and it rarely needs to fight.", "Just its presence is enough to make rivals back down.", "It watches over the young ones in its territory."},
	}
	legendLines = []string{
		"Legends say it shaped the land long before people lived there.",
		"Only a handful of people claim to have seen it, and none agree on what it looked like.",
		"Old songs say it wakes once every thousand years.",
	}
)

func entry(f family, g *rng, name string, stage int) string {
	var b strings.Builder
	fmt.Fprintf(&b, pick(g, openings), pick(g, f.primary.habitats))
	b.WriteString(" ")
	behavior := pick(g, f.primary.behaviors)
	if f.secondary != nil && g.chance(0.5) {
		behavior = pick(g, f.secondary.behaviors)
	}
	fmt.Fprintf(&b, "%s %s.", name, behavior)
	b.WriteString(" ")
	switch {
	case f.legendary:
		b.WriteString(pick(g, legendLines))
	case g.chance(0.5):
		fmt.Fprintf(&b, pick(g, quirks), pick(g, f.arch.parts))
	default:
		line := stage
		if f.stages == 2 && stage == 2 {
			line = 3
		}
		if f.stages == 1 {
			line = 2
		}
		b.WriteString(pick(g, stageLines[line]))
	}
	return b.String()
}
