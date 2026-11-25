package gen

// Ability is a passive trait a creature can have.
type Ability struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Hidden      bool   `json:"hidden,omitempty"`
}

// elemType is one of the eighteen elemental types and everything generated from it.
type elemType struct {
	id    string
	name  string
	color string // badge color

	// The creature's body colors: a hue range in degrees, with saturation and lightness ranges.
	hue        [2]float64
	sat, light [2]float64

	prefixes  []string // name fragments
	moveWords []string // first word of a move's name
	category  []string // first word of a "<word> Pokémon" category
	habitats  []string
	behaviors []string // "It <behavior>."
	abilities []Ability
	tailTip   string // decoration at the end of a tail
}

var types = []elemType{
	{id: "normal", name: "Normal", color: "#a8a77a", hue: [2]float64{20, 45}, sat: [2]float64{0.15, 0.45}, light: [2]float64{0.55, 0.78},
		prefixes:  []string{"fluff", "scamp", "nib", "bun", "tuft", "pud", "chub", "mop", "dapp", "snug"},
		moveWords: []string{"Hyper", "Quick", "Body", "Mega", "Wild", "Double", "Lucky", "Rapid"},
		category:  []string{"Fluff", "Scamper", "Snuggle", "Nibble", "Plain", "Tuft"},
		habitats:  []string{"in tall meadow grass", "near busy village roads", "under porches and barns", "along quiet country paths"},
		behaviors: []string{"follows travelers for hours hoping to be fed", "naps in sunbeams and moves only to follow them", "collects shiny pebbles and hides them in its fur", "greets every creature it meets with a loud chirp"},
		abilities: []Ability{{Name: "Cozy Coat", Description: "Halves damage from the first hit it takes in battle."}, {Name: "Copycat Heart", Description: "Becomes the type of the last move that hit it."}, {Name: "Stubborn", Description: "Can't have its stats lowered by opponents."}},
		tailTip:   "tuft"},
	{id: "fire", name: "Fire", color: "#ee8130", hue: [2]float64{0, 32}, sat: [2]float64{0.7, 0.95}, light: [2]float64{0.48, 0.6},
		prefixes:  []string{"pyr", "embr", "cind", "blaz", "scor", "igni", "flar", "sear", "smold", "kindl"},
		moveWords: []string{"Ember", "Cinder", "Blaze", "Flare", "Scorch", "Magma", "Pyre", "Smolder"},
		category:  []string{"Ember", "Cinder", "Blaze", "Kindling", "Magma", "Spark"},
		habitats:  []string{"near volcanic vents", "in sun-baked canyons", "around abandoned forges", "on the slopes of smoking mountains"},
		behaviors: []string{"warms its paws on hot stones before dawn", "sneezes sparks when it gets excited", "keeps its den warm all winter with its own body heat", "dances around campfires and steals the toasted bits"},
		abilities: []Ability{{Name: "Kindling Heart", Description: "Fire moves grow stronger when its HP is low."}, {Name: "Ash Cloak", Description: "Opponents that touch it may be burned."}, {Name: "Hearth Keeper", Description: "Heals a little each turn in harsh sunlight."}},
		tailTip:   "flame"},
	{id: "water", name: "Water", color: "#6390f0", hue: [2]float64{185, 225}, sat: [2]float64{0.55, 0.85}, light: [2]float64{0.48, 0.65},
		prefixes:  []string{"aqu", "hydr", "tid", "splas", "rippl", "mar", "brin", "puddl", "wav", "drizz"},
		moveWords: []string{"Tide", "Torrent", "Bubble", "Ripple", "Surge", "Brine", "Spray", "Undertow"},
		category:  []string{"Tide", "Bubble", "Ripple", "Puddle", "Current", "Brine"},
		habitats:  []string{"in clear mountain streams", "along rocky coastlines", "in deep, still lakes", "in rain puddles after storms"},
		behaviors: []string{"blows bubbles to talk with others of its kind", "rides river currents for fun and walks all the way back", "can hold its breath for an entire afternoon", "sprays water at anything that looks too dry"},
		abilities: []Ability{{Name: "Riptide", Description: "Raises Speed sharply in the rain."}, {Name: "Slick Scales", Description: "Contact moves sometimes slide right off it."}, {Name: "Tidal Memory", Description: "Water moves gain power each turn they're used in a row."}},
		tailTip:   "fin"},
	{id: "grass", name: "Grass", color: "#7ac74c", hue: [2]float64{80, 145}, sat: [2]float64{0.45, 0.75}, light: [2]float64{0.4, 0.58},
		prefixes:  []string{"leaf", "fern", "sprou", "bloss", "thorn", "moss", "verd", "bram", "petal", "clov"},
		moveWords: []string{"Leaf", "Vine", "Petal", "Thorn", "Seed", "Bloom", "Bramble", "Pollen"},
		category:  []string{"Sprout", "Bloom", "Leaf", "Moss", "Bramble", "Seedling"},
		habitats:  []string{"in sunny forest clearings", "among overgrown ruins", "in flower fields", "deep in damp, mossy woods"},
		behaviors: []string{"stands perfectly still in the sun to soak up energy", "grows a new flower every time it is happy", "plants seeds wherever it naps", "smells faintly of fresh rain and cut grass"},
		abilities: []Ability{{Name: "Photosynth", Description: "Restores HP each turn in sunlight."}, {Name: "Bramble Hide", Description: "Damages opponents that make contact."}, {Name: "Deep Roots", Description: "Can't be forced out of battle."}},
		tailTip:   "leaf"},
	{id: "electric", name: "Electric", color: "#f7d02c", hue: [2]float64{44, 58}, sat: [2]float64{0.8, 1}, light: [2]float64{0.5, 0.62},
		prefixes:  []string{"volt", "zap", "spark", "amp", "jolt", "stat", "ohm", "flick", "buzz", "wattl"},
		moveWords: []string{"Volt", "Thunder", "Spark", "Static", "Jolt", "Plasma", "Arc", "Zap"},
		category:  []string{"Spark", "Static", "Volt", "Battery", "Thunder", "Jolt"},
		habitats:  []string{"near power lines", "on stormy hilltops", "inside abandoned power plants", "in fields struck often by lightning"},
		behaviors: []string{"makes everyone's hair stand on end when it's nearby", "charges up by rubbing its cheeks against wool", "lights up its home like a lamp at night", "gets jittery and hard to catch before thunderstorms"},
		abilities: []Ability{{Name: "Static Coat", Description: "Opponents that touch it may be paralyzed."}, {Name: "Overcharge", Description: "Electric moves never miss, but cost a little HP."}, {Name: "Lightning Rod", Description: "Draws in Electric moves and raises Sp. Atk."}},
		tailTip:   "bolt"},
	{id: "ice", name: "Ice", color: "#96d9d6", hue: [2]float64{170, 200}, sat: [2]float64{0.35, 0.65}, light: [2]float64{0.68, 0.84},
		prefixes:  []string{"frost", "glaci", "cryo", "sleet", "rim", "chill", "hail", "snow", "brrr", "fjor"},
		moveWords: []string{"Frost", "Glacier", "Hail", "Blizzard", "Rime", "Icicle", "Sleet", "Aurora"},
		category:  []string{"Frost", "Snowdrift", "Icicle", "Glacier", "Chill", "Sleet"},
		habitats:  []string{"on frozen mountain peaks", "in snowfields that never melt", "inside glacier caves", "on drifting sea ice"},
		behaviors: []string{"breathes out tiny snowflakes when it sleeps", "builds dens out of packed snow", "freezes puddles just to slide on them", "turns sluggish and grumpy in warm weather"},
		abilities: []Ability{{Name: "Snow Cloak", Description: "Raises evasion during hail."}, {Name: "Permafrost", Description: "Takes less damage from Fire moves."}, {Name: "Cold Snap", Description: "Its Ice moves may freeze the target solid."}},
		tailTip:   "crystal"},
	{id: "fighting", name: "Fighting", color: "#c22e28", hue: [2]float64{0, 25}, sat: [2]float64{0.4, 0.7}, light: [2]float64{0.38, 0.55},
		prefixes:  []string{"brawl", "pumm", "grapp", "kik", "doj", "kara", "jab", "rumbl", "scrap", "tusl"},
		moveWords: []string{"Power", "Iron Fist", "Combo", "Counter", "Rising", "Focus", "Dynamic", "Knockout"},
		category:  []string{"Brawler", "Sparring", "Dojo", "Grapple", "Punch", "Training"},
		habitats:  []string{"in mountain dojos", "near waterfalls where it trains", "in rocky arenas", "along old pilgrimage trails"},
		behaviors: []string{"trains from sunrise to sunset without a break", "bows politely before every battle", "punches boulders to toughen its fists", "challenges anything that looks strong"},
		abilities: []Ability{{Name: "Guts", Description: "Attack rises when it has a status problem."}, {Name: "Iron Will", Description: "Can't be made to flinch."}, {Name: "Second Wind", Description: "Recovers HP once when it drops below half."}},
		tailTip:   "tuft"},
	{id: "poison", name: "Poison", color: "#a33ea1", hue: [2]float64{270, 310}, sat: [2]float64{0.4, 0.7}, light: [2]float64{0.42, 0.6},
		prefixes:  []string{"tox", "venn", "blig", "smog", "ooz", "mir", "fum", "sludg", "nox", "bil"},
		moveWords: []string{"Toxic", "Venom", "Sludge", "Acid", "Blight", "Miasma", "Ooze", "Noxious"},
		category:  []string{"Toxic", "Sludge", "Venom", "Fume", "Bog", "Spore"},
		habitats:  []string{"in murky swamps", "around polluted rivers", "in damp sewer tunnels", "near bubbling tar pits"},
		behaviors: []string{"leaves a trail of fizzing ooze wherever it walks", "puffs up and changes color when threatened", "is avoided by every other creature in its swamp", "smells terrible but is surprisingly affectionate"},
		abilities: []Ability{{Name: "Toxic Skin", Description: "Opponents that touch it may be poisoned."}, {Name: "Corrode", Description: "Can poison even Steel-type creatures."}, {Name: "Stench", Description: "Its attacks may make the target flinch."}},
		tailTip:   "drip"},
	{id: "ground", name: "Ground", color: "#e2bf65", hue: [2]float64{25, 45}, sat: [2]float64{0.35, 0.6}, light: [2]float64{0.42, 0.6},
		prefixes:  []string{"terr", "dust", "mud", "dun", "quak", "clod", "sand", "burr", "loam", "gully"},
		moveWords: []string{"Earth", "Mud", "Quake", "Sand", "Dust", "Bedrock", "Tremor", "Dune"},
		category:  []string{"Burrow", "Dune", "Mudslide", "Tremor", "Dust", "Clay"},
		habitats:  []string{"in dry desert dunes", "in deep underground burrows", "along muddy riverbanks", "in dusty badlands"},
		behaviors: []string{"digs tunnels so long they cross whole valleys", "stomps the ground to talk with its family", "rolls in mud to stay cool on hot days", "can sense footsteps from a mile away"},
		abilities: []Ability{{Name: "Sand Veil", Description: "Raises evasion in a sandstorm."}, {Name: "Burrower", Description: "Moves first when it uses a Ground move."}, {Name: "Bedrock", Description: "Survives any hit from full HP with 1 HP left."}},
		tailTip:   "club"},
	{id: "flying", name: "Flying", color: "#a98ff3", hue: [2]float64{195, 250}, sat: [2]float64{0.3, 0.6}, light: [2]float64{0.6, 0.75},
		prefixes:  []string{"gust", "aer", "zeph", "swoop", "plum", "glid", "sky", "breez", "soar", "wisp"},
		moveWords: []string{"Gale", "Sky", "Wing", "Aerial", "Cyclone", "Feather", "Updraft", "Tailwind"},
		category:  []string{"Gust", "Feather", "Skyward", "Breeze", "Glider", "Tailwind"},
		habitats:  []string{"on windy cliffs", "above the clouds", "in the tallest trees of old forests", "along sea cliffs where the wind is strong"},
		behaviors: []string{"rides warm air currents for days without landing", "sings at dawn to mark its territory", "drops feathers that bring good luck", "can spot a berry from a mile up"},
		abilities: []Ability{{Name: "Tailwind Rider", Description: "Raises Speed when the wind blows."}, {Name: "Keen Eye", Description: "Its accuracy can't be lowered."}, {Name: "Featherweight", Description: "Takes no damage from Ground moves or spikes."}},
		tailTip:   "feather"},
	{id: "psychic", name: "Psychic", color: "#f95587", hue: [2]float64{300, 345}, sat: [2]float64{0.45, 0.8}, light: [2]float64{0.55, 0.72},
		prefixes:  []string{"psi", "ment", "mind", "aur", "cosm", "dre", "omen", "lumi", "sooth", "trance"},
		moveWords: []string{"Psy", "Mind", "Dream", "Astral", "Prism", "Telekinetic", "Omen", "Lucid"},
		category:  []string{"Dream", "Mind", "Aura", "Omen", "Trance", "Lucid"},
		habitats:  []string{"in ancient temples", "in quiet libraries", "where the northern lights are brightest", "in places where dreams feel real"},
		behaviors: []string{"floats a few inches off the ground when it's deep in thought", "can tell what you're about to say before you say it", "shows up in people's dreams before it meets them", "bends spoons when it's bored"},
		abilities: []Ability{{Name: "Foresight", Description: "Senses the opponent's strongest move."}, {Name: "Mind Ward", Description: "Immune to confusion and taunts."}, {Name: "Dream Eater", Description: "Heals when an opponent falls asleep."}},
		tailTip:   "orb"},
	{id: "bug", name: "Bug", color: "#a6b91a", hue: [2]float64{55, 100}, sat: [2]float64{0.5, 0.8}, light: [2]float64{0.4, 0.55},
		prefixes:  []string{"skitt", "chit", "buzz", "larv", "mand", "weev", "gnat", "cric", "pupa", "beetl"},
		moveWords: []string{"Swarm", "Sting", "Silk", "Chitin", "Buzz", "Larva", "Mandible", "Hive"},
		category:  []string{"Beetle", "Larva", "Swarm", "Silk", "Cricket", "Hive"},
		habitats:  []string{"under fallen logs", "in tall summer grass", "in giant hives in the forest canopy", "around glowing night flowers"},
		behaviors: []string{"chirps a different song for every season", "spins silk hammocks to sleep in", "marches in long lines with others of its kind", "glows faintly when it finds something sweet"},
		abilities: []Ability{{Name: "Swarm", Description: "Bug moves get stronger when its HP is low."}, {Name: "Shed Skin", Description: "May cure its own status problems each turn."}, {Name: "Compound Eyes", Description: "Raises the accuracy of its moves."}},
		tailTip:   "stinger"},
	{id: "rock", name: "Rock", color: "#b6a136", hue: [2]float64{25, 50}, sat: [2]float64{0.1, 0.35}, light: [2]float64{0.4, 0.58},
		prefixes:  []string{"ston", "boul", "peb", "crag", "shal", "gran", "flint", "quar", "slat", "rubbl"},
		moveWords: []string{"Rock", "Boulder", "Stone", "Avalanche", "Crag", "Pebble", "Landslide", "Granite"},
		category:  []string{"Pebble", "Boulder", "Crag", "Quarry", "Fossil", "Pillar"},
		habitats:  []string{"in rocky mountain passes", "deep in quarries", "inside caves full of crystals", "among ancient stone ruins"},
		behaviors: []string{"is often mistaken for an ordinary boulder", "eats gravel and seems to enjoy it", "rolls downhill when it wants to travel fast", "sleeps for weeks at a time without moving"},
		abilities: []Ability{{Name: "Sturdy", Description: "Can't be knocked out in one hit from full HP."}, {Name: "Rock Head", Description: "Takes no recoil damage."}, {Name: "Sediment", Description: "Defense rises each time it's hit."}},
		tailTip:   "club"},
	{id: "ghost", name: "Ghost", color: "#735797", hue: [2]float64{245, 285}, sat: [2]float64{0.25, 0.55}, light: [2]float64{0.35, 0.55},
		prefixes:  []string{"spect", "phant", "wisp", "gloo", "shad", "haun", "gast", "murk", "bool", "eer"},
		moveWords: []string{"Shadow", "Phantom", "Spirit", "Haunt", "Hex", "Wraith", "Gloom", "Specter"},
		category:  []string{"Haunt", "Phantom", "Wisp", "Lantern", "Specter", "Shroud"},
		habitats:  []string{"in abandoned mansions", "around foggy graveyards", "in forgotten attics", "where lanterns flicker for no reason"},
		behaviors: []string{"hides in shadows and giggles at passersby", "steals socks and never gives them back", "can pass through walls but prefers doors", "makes lights flicker when it wants attention"},
		abilities: []Ability{{Name: "Cursed Body", Description: "May disable a move that hits it."}, {Name: "Phase Out", Description: "Sometimes dodges physical attacks entirely."}, {Name: "Lingering Chill", Description: "Lowers the Speed of opponents that touch it."}},
		tailTip:   "wisp"},
	{id: "dragon", name: "Dragon", color: "#6f35fc", hue: [2]float64{220, 275}, sat: [2]float64{0.5, 0.8}, light: [2]float64{0.42, 0.58},
		prefixes:  []string{"drak", "wyrm", "scal", "dra", "wyv", "serp", "dyn", "ryu", "draco", "tyr"},
		moveWords: []string{"Dragon", "Draco", "Wyrm", "Scale", "Outrage", "Legend", "Primal", "Dragonfire"},
		category:  []string{"Wyrm", "Scale", "Dragon", "Ancient", "Hoard", "Crown"},
		habitats:  []string{"in hidden mountain lairs", "on islands no map shows", "in caves full of treasure", "above stormy seas"},
		behaviors: []string{"guards a hoard of shiny things it has collected for centuries", "is said to appear only once a generation", "roars so loudly that mountains echo for hours", "sheds scales that are harder than steel"},
		abilities: []Ability{{Name: "Dragon Pride", Description: "Attack rises when an opponent faints."}, {Name: "Ancient Scales", Description: "Takes less damage from super-effective moves."}, {Name: "Hoarder", Description: "Always finds an item after battle."}},
		tailTip:   "spade"},
	{id: "dark", name: "Dark", color: "#705746", hue: [2]float64{0, 40}, sat: [2]float64{0.1, 0.35}, light: [2]float64{0.22, 0.38},
		prefixes:  []string{"nox", "umbr", "grim", "sly", "dusk", "nigh", "murk", "sabl", "crow", "vex"},
		moveWords: []string{"Night", "Shadow", "Dusk", "Sly", "Feint", "Ambush", "Midnight", "Sucker"},
		category:  []string{"Dusk", "Trickster", "Shadow", "Prowler", "Moonless", "Bandit"},
		habitats:  []string{"in moonless forests", "in back alleys of big cities", "in caves where light never reaches", "on rooftops at midnight"},
		behaviors: []string{"steals snacks from campers and leaves a feather behind", "only comes out on moonless nights", "plays tricks on anyone who walks alone", "watches from the shadows before deciding if you're a friend"},
		abilities: []Ability{{Name: "Prankster", Description: "Status moves always go first."}, {Name: "Night Stalker", Description: "Moves first at night."}, {Name: "Pickpocket", Description: "Steals the item of an opponent that touches it."}},
		tailTip:   "spade"},
	{id: "steel", name: "Steel", color: "#b7b7ce", hue: [2]float64{200, 230}, sat: [2]float64{0.05, 0.2}, light: [2]float64{0.55, 0.72},
		prefixes:  []string{"ferr", "chrom", "bolt", "rivet", "alloy", "mech", "iron", "gear", "cobal", "titan"},
		moveWords: []string{"Iron", "Steel", "Chrome", "Magnet", "Gear", "Alloy", "Gyro", "Bullet"},
		category:  []string{"Iron", "Gear", "Rivet", "Alloy", "Magnet", "Armor"},
		habitats:  []string{"in old factories", "near iron mines", "in junkyards", "on the hulls of sunken ships"},
		behaviors: []string{"polishes itself every morning until it shines", "makes a faint clanking sound when it walks", "eats scrap metal and hums while it does", "rusts if it's left in the rain too long"},
		abilities: []Ability{{Name: "Heavy Metal", Description: "Its weight doubles, powering up heavy moves."}, {Name: "Magnet Pull", Description: "Steel-type opponents can't escape."}, {Name: "Clear Body", Description: "Its stats can't be lowered."}},
		tailTip:   "spade"},
	{id: "fairy", name: "Fairy", color: "#d685ad", hue: [2]float64{320, 355}, sat: [2]float64{0.5, 0.8}, light: [2]float64{0.7, 0.84},
		prefixes:  []string{"pix", "fae", "glim", "twinkl", "sprit", "lull", "pip", "dain", "flutt", "shimm"},
		moveWords: []string{"Fairy", "Moon", "Dazzle", "Sparkle", "Charm", "Starlight", "Wish", "Glimmer"},
		category:  []string{"Twinkle", "Wish", "Moonbeam", "Pixie", "Glimmer", "Lullaby"},
		habitats:  []string{"in moonlit flower meadows", "inside rings of mushrooms", "where rainbows touch the ground", "in quiet glades at dusk"},
		behaviors: []string{"leaves a trail of sparkles that fades by morning", "hums lullabies that make everyone nearby sleepy", "can only be seen by people who believe in it", "grants small wishes to those who are kind to it"},
		abilities: []Ability{{Name: "Pixie Dust", Description: "Its moves may lower the target's Attack."}, {Name: "Sweet Veil", Description: "It and its allies can't fall asleep."}, {Name: "Moonlit Grace", Description: "Restores HP at night."}},
		tailTip:   "star"},
}

func typeByID(id string) (elemType, bool) {
	for _, t := range types {
		if t.id == id {
			return t, true
		}
	}
	return elemType{}, false
}

// Types lists the eighteen types in order.
func Types() []TypeRef {
	out := make([]TypeRef, len(types))
	for i, t := range types {
		out[i] = ref(t)
	}
	return out
}
