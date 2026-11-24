package gen

// archetype is a body plan: how the sprite is built, what its name sounds like, and how its
// stats lean.
type archetype struct {
	id string
	// Name endings by size: index 0 for a first stage, 2 for a final one.
	endings [3][]string
	nouns   []string // second word of the category
	// Stat weights: HP, Attack, Defense, Sp. Atk, Sp. Def, Speed.
	stats [6]float64
	// Types this body plan suits, picked more often with it.
	affinity []string
	height   float64 // meters at full size
	weight   float64 // kilograms at full size
	floats   bool
	parts    []string // body parts the dex entry can mention, all plural
}

var archetypes = []archetype{
	{id: "biped",
		endings:  [3][]string{{"kin", "tot", "pip", "imp", "let", "bit"}, {"chap", "ling", "kid", "scout", "nik", "o"}, {"lord", "king", "champ", "titan", "mancer", "ius"}},
		nouns:    []string{"Imp", "Scout", "Sprite", "Knight", "Rascal", "Monarch"},
		stats:    [6]float64{1, 1.1, 1, 1.1, 1, 1},
		affinity: []string{"fighting", "psychic", "fire", "fairy", "dark", "normal"},
		height:   1.5, weight: 55, parts: []string{"hands", "arms", "eyes"}},
	{id: "quadruped",
		endings:  [3][]string{{"pup", "kit", "cub", "fawn", "paw", "nip"}, {"hound", "lynx", "stag", "boar", "pard", "ox"}, {"fang", "mane", "horn", "beast", "thane", "ram"}},
		nouns:    []string{"Pup", "Hound", "Fawn", "Beast", "Mane", "Prowler"},
		stats:    [6]float64{1, 1.2, 1, 0.9, 0.9, 1.2},
		affinity: []string{"fire", "ground", "normal", "electric", "dark", "ice", "grass"},
		height:   1.4, weight: 90, parts: []string{"paws", "whiskers", "legs"}},
	{id: "serpent",
		endings:  [3][]string{{"sss", "ling", "worm", "eel", "coil", "nake"}, {"asp", "viper", "coil", "adder", "slith", "ra"}, {"boa", "wyrm", "hydra", "naga", "coatl", "basilisk"}},
		nouns:    []string{"Coil", "Serpent", "Viper", "Eel", "Wyrm", "Slither"},
		stats:    [6]float64{1, 1.1, 0.9, 1.1, 0.9, 1.1},
		affinity: []string{"poison", "dragon", "water", "ground", "dark"},
		height:   3.5, weight: 70, parts: []string{"coils", "scales", "fangs"}},
	{id: "avian",
		endings:  [3][]string{{"chick", "peep", "tweet", "fledge", "wren", "it"}, {"wing", "hawk", "crow", "jay", "swift", "lark"}, {"talon", "condor", "raptor", "aerie", "phoenix", "gale"}},
		nouns:    []string{"Chick", "Songbird", "Raptor", "Wing", "Talon", "Plume"},
		stats:    [6]float64{0.9, 1.05, 0.85, 1, 0.9, 1.4},
		affinity: []string{"flying", "flying", "flying", "fire", "ice", "electric", "psychic"},
		height:   1.2, weight: 30, parts: []string{"wings", "feathers", "talons"}},
	{id: "blob",
		endings:  [3][]string{{"blob", "drop", "glop", "bean", "mochi", "puff"}, {"goo", "jelly", "mound", "squish", "lump", "plop"}, {"gloop", "slurm", "colossus", "mass", "ooze", "dome"}},
		nouns:    []string{"Blob", "Jelly", "Goo", "Puff", "Squish", "Dumpling"},
		stats:    [6]float64{1.5, 0.8, 1.1, 1, 1.2, 0.6},
		affinity: []string{"poison", "water", "fairy", "normal", "psychic", "ice"},
		height:   0.9, weight: 40, parts: []string{"bubbles", "squishy sides", "drips"}},
	{id: "insect",
		endings:  [3][]string{{"grub", "mite", "pupa", "nit", "larv", "bug"}, {"beetle", "mantis", "moth", "tick", "cada", "chafer"}, {"scarab", "hornet", "moth", "carapax", "queen", "pede"}},
		nouns:    []string{"Grub", "Beetle", "Moth", "Mantis", "Hornet", "Scarab"},
		stats:    [6]float64{0.8, 1.2, 1.1, 0.9, 0.9, 1.2},
		affinity: []string{"bug", "bug", "bug", "poison", "steel", "flying", "grass"},
		height:   1.0, weight: 25, parts: []string{"antennae", "wings", "legs"}},
	{id: "aquatic",
		endings:  [3][]string{{"fry", "guppy", "minnow", "fin", "squid", "kelp"}, {"fin", "carp", "ray", "gill", "pike", "tuna"}, {"leviathan", "shark", "marlin", "kraken", "orca", "maw"}},
		nouns:    []string{"Minnow", "Fin", "Reef", "Gill", "Deep Sea", "Lagoon"},
		stats:    [6]float64{1.1, 1, 1, 1.05, 1, 1.1},
		affinity: []string{"water", "water", "water", "ice", "dark", "electric", "poison"},
		height:   1.6, weight: 60, parts: []string{"fins", "scales", "gills"}},
	{id: "golem",
		endings:  [3][]string{{"bit", "chip", "brick", "nugget", "pebb", "lump"}, {"block", "golem", "tote", "lith", "cairn", "ward"}, {"colossus", "monolith", "bastion", "fortress", "titan", "ogre"}},
		nouns:    []string{"Brick", "Golem", "Monolith", "Bastion", "Guardian", "Fortress"},
		stats:    [6]float64{1.2, 1.2, 1.5, 0.7, 1, 0.5},
		affinity: []string{"rock", "rock", "steel", "steel", "ground", "ice", "fighting"},
		height:   2.0, weight: 300, parts: []string{"arms", "armor plates", "fists"}},
	{id: "spirit",
		endings:  [3][]string{{"wisp", "boo", "glim", "mote", "puff", "ette"}, {"shade", "veil", "geist", "lantern", "shroud", "wraith"}, {"phantasm", "revenant", "banshee", "eidolon", "specter", "requiem"}},
		nouns:    []string{"Wisp", "Lantern", "Veil", "Spirit", "Shroud", "Phantasm"},
		stats:    [6]float64{0.9, 0.8, 0.9, 1.35, 1.2, 1.1},
		affinity: []string{"ghost", "ghost", "ghost", "psychic", "fairy", "dark", "fire", "ice"},
		height:   1.1, weight: 4, floats: true, parts: []string{"wisps", "eyes", "glowing edges"}},
}

func archetypeByID(id string) (archetype, bool) {
	for _, a := range archetypes {
		if a.id == id {
			return a, true
		}
	}
	return archetype{}, false
}
