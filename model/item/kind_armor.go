package item

const (
	KindArmor Kind = "armor"
)

type Armor struct {
	Item

	Type           string           `json:"type"`
	Armor          ArmorProps       `json:"armor"`
	Components     []ArmorComponent `json:"components" bson:"components"`
	RicochetChance string           `json:"ricochetChance,omitempty"`
	Penalties      Penalties        `json:"penalties"`
	Blocking       []string         `json:"blocking"`
	Slots          Slots            `json:"slots"`
	Compatibility  ItemList         `json:"compatibility"`
}

func (a Armor) TotalDurability() float64 {
	var total float64
	for _, component := range a.Components {
		total += component.Durability
	}

	return total
}

type ArmorResult struct {
	*Result
	Items []Armor `json:"items"`
}

func (r *ArmorResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type ArmorComponent struct {
	ArmorProps
}

type ArmorProps struct {
	Class           int64         `json:"class"`
	Durability      float64       `json:"durability"`
	Material        ArmorMaterial `json:"material"`
	BluntThroughput float64       `json:"bluntThroughput"`
	Zones           []string      `json:"zones"`
}

type ArmorMaterial struct {
	Name            string  `json:"name"`
	Destructibility float64 `json:"destructibility"`
}

var armorClasses = [...]string{
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
}

var armorMaterial = [...]string{
	"aluminium",
	"aramid",
	"ceramic",
	"combined",
	"glass",
	"steel",
	"titanium",
	"uhmwpe",
}

var armorFilter = Filter{
	"type": {
		"attachment",
		"body",
		"faceCover",
		"helmet",
		"visor",
	},
	"class":    armorClasses[:],
	"material": armorMaterial[:],
}
