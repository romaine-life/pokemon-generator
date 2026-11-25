// Package gen invents creatures. Everything about one comes from its seed, so a seed is a
// permanent address for it: the same seed (and type lock) always gives the same creature.
package gen

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

// Version names the generator. Changing what a seed produces is a new version.
const Version = "1"

// Options picks a creature. Seed is any text; Type, when set, locks the primary type; Stage,
// when set, picks a stage of the evolution line instead of the seed's own pick.
type Options struct {
	Seed  string
	Type  string
	Stage int
}

type TypeRef struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Stat struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

type Move struct {
	Name        string  `json:"name"`
	Type        TypeRef `json:"type"`
	Category    string  `json:"category"` // Physical, Special or Status
	Power       int     `json:"power,omitempty"`
	Accuracy    int     `json:"accuracy"`
	PP          int     `json:"pp"`
	Description string  `json:"description"`
}

// Evolution is one stage of the line, with its sprite for the evolution strip.
type Evolution struct {
	Stage  int    `json:"stage"`
	Name   string `json:"name"`
	Number int    `json:"number"`
	How    string `json:"how,omitempty"` // how the previous stage becomes this one
	Sprite Sprite `json:"sprite"`
}

type Monster struct {
	Version   string      `json:"version"`
	Seed      string      `json:"seed"`
	TypeLock  string      `json:"typeLock,omitempty"`
	Number    int         `json:"number"`
	Name      string      `json:"name"`
	Category  string      `json:"category"`
	Types     []TypeRef   `json:"types"`
	BodyPlan  string      `json:"bodyPlan"`
	Stage     int         `json:"stage"`
	Stages    int         `json:"stages"`
	Legendary bool        `json:"legendary"`
	Shiny     bool        `json:"shiny"`
	HeightM   float64     `json:"heightM"`
	WeightKg  float64     `json:"weightKg"`
	Abilities []Ability   `json:"abilities"`
	Stats     []Stat      `json:"stats"`
	Total     int         `json:"total"`
	Moves     []Move      `json:"moves"`
	Entry     string      `json:"entry"`
	Sprite    Sprite      `json:"sprite"`
	Evolution []Evolution `json:"evolution"`
}

// family is what every stage of an evolution line shares.
type family struct {
	key        string
	arch       archetype
	primary    elemType
	secondary  *elemType
	stages     int
	legendary  bool
	shiny      bool
	number     int
	prefix     string
	link       string
	statWeight [6]float64
	abilities  []Ability
	levels     [3]int
}

// NormalizeSeed makes seeds that differ only in case or surrounding space the same seed.
func NormalizeSeed(seed string) string {
	seed = strings.ToLower(strings.TrimSpace(seed))
	seed = strings.Join(strings.Fields(seed), " ")
	if len(seed) > 64 {
		seed = seed[:64]
	}
	return seed
}

// ValidType reports whether id names a type ("" means no lock).
func ValidType(id string) bool {
	_, ok := typeByID(id)
	return id == "" || ok
}

func newFamily(seed, lock string) family {
	key := "v" + Version + "|" + seed + "|" + lock
	g := newRNG("family", key)
	f := family{key: key}
	f.legendary = g.chance(0.03)
	f.shiny = newRNG("shiny", key).chance(1.0 / 64)

	// The body plan, leaning toward ones that suit a locked type.
	weights := make([]float64, len(archetypes))
	for i, a := range archetypes {
		weights[i] = 1
		for _, t := range a.affinity {
			if t == lock {
				weights[i] += 1.5
			}
		}
	}
	f.arch = archetypes[g.weighted(weights)]

	if t, ok := typeByID(lock); ok {
		f.primary = t
	} else {
		tw := make([]float64, len(types))
		for i, t := range types {
			tw[i] = 1
			for _, a := range f.arch.affinity {
				if a == t.id {
					tw[i] += 2.5
				}
			}
		}
		f.primary = types[g.weighted(tw)]
	}
	dual := 0.4
	if f.legendary {
		dual = 0.7
	}
	if g.chance(dual) {
		var second elemType
		if f.arch.id == "avian" && f.primary.id != "flying" && g.chance(0.6) {
			second, _ = typeByID("flying")
		} else {
			for second = pick(g, types); second.id == f.primary.id; second = pick(g, types) {
			}
		}
		f.secondary = &second
	}

	switch {
	case f.legendary:
		f.stages = 1
	default:
		f.stages = 1 + g.weighted([]float64{15, 30, 55})
	}
	f.number = g.intn(1, 9999-f.stages)
	f.prefix = pick(g, f.primary.prefixes)
	f.link = pick(g, []string{"a", "o", "i", "e", "u", "a", "o"})
	f.levels = [3]int{0, g.intn(14, 26), g.intn(32, 44)}

	for i, w := range f.arch.stats {
		f.statWeight[i] = w * g.float(0.65, 1.4)
	}
	f.abilities = []Ability{pick(g, f.primary.abilities)}
	if f.secondary != nil && g.chance(0.6) {
		f.abilities = append(f.abilities, pick(g, f.secondary.abilities))
	} else if g.chance(0.5) {
		for {
			a := pick(g, f.primary.abilities)
			if a.Name != f.abilities[0].Name {
				f.abilities = append(f.abilities, a)
				break
			}
		}
	}
	hidden := pick(g, hiddenAbilities)
	hidden.Hidden = true
	f.abilities = append(f.abilities, hidden)
	return f
}

// tier maps a stage to its size: first, middle or final.
func (f family) tier(stage int) int {
	switch {
	case f.legendary:
		return 2
	case f.stages == 1:
		return 1
	case f.stages == 2:
		return []int{0, 2}[stage-1]
	default:
		return stage - 1
	}
}

func (f family) spec(stage int) spriteSpec {
	return spriteSpec{family: f.key, arch: f.arch, primary: f.primary, secondary: f.secondary,
		tier: f.tier(stage), legendary: f.legendary, shiny: f.shiny}
}

func (f family) name(stage int) string {
	g := newRNG("name", f.key, fmt.Sprint(stage))
	ending := pick(g, f.arch.endings[f.tier(stage)])
	return joinName(f.prefix, f.link, ending)
}

// joinName joins a type fragment and a body ending into one pronounceable name.
func joinName(prefix, link, ending string) string {
	isVowel := func(r byte) bool { return strings.IndexByte("aeiouy", r) >= 0 }
	last, first := prefix[len(prefix)-1], ending[0]
	switch {
	case isVowel(last) && isVowel(first):
		ending = ending[1:]
	case !isVowel(last) && !isVowel(first) && last != first:
		prefix += link
	case last == first:
		ending = ending[1:]
	}
	name := prefix + ending
	if len(name) > 12 {
		name = name[:12]
	}
	r := []rune(name)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// Generate makes the creature the options pick.
func Generate(o Options) Monster {
	seed := NormalizeSeed(o.Seed)
	f := newFamily(seed, o.Type)
	stage := o.Stage
	if stage < 1 || stage > f.stages {
		stage = 1 + newRNG("stage", f.key).intn(0, f.stages-1)
	}
	g := newRNG("stage", f.key, fmt.Sprint(stage))
	tier := f.tier(stage)
	size := []float64{0.45, 0.75, 1}[tier]

	m := Monster{
		Version:   Version,
		Seed:      seed,
		TypeLock:  o.Type,
		Number:    f.number + stage - 1,
		Name:      f.name(stage),
		BodyPlan:  f.arch.id,
		Stage:     stage,
		Stages:    f.stages,
		Legendary: f.legendary,
		Shiny:     f.shiny,
		Abilities: f.abilities,
	}
	m.Types = []TypeRef{ref(f.primary)}
	if f.secondary != nil {
		m.Types = append(m.Types, ref(*f.secondary))
	}
	m.Category = pick(g, f.primary.category) + " " + pick(g, f.arch.nouns)
	if f.legendary {
		m.Category = pick(g, []string{"Eternal", "Primordial", "Sovereign", "Celestial", "Mythic"}) + " " + pick(g, f.arch.nouns)
	}
	m.HeightM = math.Round(f.arch.height*size*g.float(0.75, 1.3)*10) / 10
	m.WeightKg = math.Round(f.arch.weight*size*size*size*g.float(0.6, 1.4)*10) / 10
	if f.legendary {
		m.HeightM = math.Round(m.HeightM*1.8*10) / 10
		m.WeightKg = math.Round(m.WeightKg*4*10) / 10
	}
	m.HeightM = math.Max(m.HeightM, 0.1)
	m.WeightKg = math.Max(m.WeightKg, 0.1)

	m.Stats, m.Total = stats(f, g, tier)
	m.Moves = moves(f, g, tier, m.Stats)
	m.Entry = entry(f, g, m.Name, stage)
	m.Sprite = renderSprite(f.spec(stage))

	for st := 1; st <= f.stages; st++ {
		e := Evolution{Stage: st, Name: f.name(st), Number: f.number + st - 1}
		if st > 1 {
			e.How = fmt.Sprintf("Lv. %d", f.levels[st-1])
			if newRNG("how", f.key, fmt.Sprint(st)).chance(0.2) {
				e.How = pick(newRNG("item", f.key, fmt.Sprint(st)), evolutionItems(f.primary))
			}
		}
		if st == stage {
			e.Sprite = m.Sprite
		} else {
			e.Sprite = renderSprite(f.spec(st))
		}
		m.Evolution = append(m.Evolution, e)
	}
	return m
}

func ref(t elemType) TypeRef { return TypeRef{ID: t.id, Name: t.name, Color: t.color} }

var statNames = [6][2]string{{"hp", "HP"}, {"atk", "Attack"}, {"def", "Defense"}, {"spa", "Sp. Atk"}, {"spd", "Sp. Def"}, {"spe", "Speed"}}

func stats(f family, g *rng, tier int) ([]Stat, int) {
	var total float64
	switch {
	case f.legendary:
		total = g.float(580, 690)
	case f.stages == 1:
		total = g.float(430, 500)
	default:
		total = []float64{g.float(280, 330), g.float(390, 445), g.float(490, 550)}[tier]
	}
	sum := 0.0
	for _, w := range f.statWeight {
		sum += w
	}
	out := make([]Stat, 6)
	got := 0
	for i, w := range f.statWeight {
		v := int(math.Round(clamp(total*w/sum*g.float(0.92, 1.08), 5, 255)))
		out[i] = Stat{ID: statNames[i][0], Name: statNames[i][1], Value: v}
		got += v
	}
	return out, got
}

func evolutionItems(t elemType) []string {
	return []string{t.moveWords[0] + " Stone", "Moon Shard", "Friendship", "Trade", "Level up at night", t.name + " Charm"}
}

var hiddenAbilities = []Ability{
	{Name: "Late Bloomer", Description: "Its stats rise a little each turn it stays in battle."},
	{Name: "Curious Nose", Description: "Sniffs out the opponent's held item when it enters battle."},
	{Name: "Lucky Charm", Description: "Its moves land critical hits more often."},
	{Name: "Stage Fright", Description: "Speed rises sharply when it's the last one standing."},
	{Name: "Snack Pouch", Description: "Eats its berry earlier than usual."},
	{Name: "Mirror Coat", Description: "Bounces back the first status move used on it."},
	{Name: "Thick Skull", Description: "Takes no recoil and can't be confused."},
	{Name: "Moonwalker", Description: "Raises evasion each time it switches in."},
}
