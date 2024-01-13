package item

import "fmt"

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
	Confilcts      ItemList         `json:"conflicts"`
}

func (a Armor) TotalDurability() float64 {
	if a.Armor.Durability > 0 {
		return a.Armor.Durability
	}

	var total float64
	for _, component := range a.Components {
		total += component.Durability
	}

	return total
}

func (a Armor) ClassRange() ClassRange {
	if a.Armor.Class > 0 {
		return ClassRange{
			Min: a.Armor.Class,
			Max: a.Armor.Class,
		}
	}

	if len(a.Components) == 0 {
		return ClassRange{}
	}

	min, max := a.Components[0].Class, a.Components[0].Class
	for _, component := range a.Components {
		if component.Class < min {
			min = component.Class
		}
		if component.Class > max {
			max = component.Class
		}
	}

	return ClassRange{Min: min, Max: max}
}

type ClassRange struct {
	Min int64
	Max int64
}

func (c ClassRange) String() string {
	if c.Min == c.Max {
		return fmt.Sprintf("%d", c.Min)
	}

	return fmt.Sprintf("%d-%d", c.Min, c.Max)
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
		"plate",
	},
	"class":    armorClasses[:],
	"material": armorMaterial[:],
}
