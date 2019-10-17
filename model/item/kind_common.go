package item

type ItemList map[Kind][]objectID

const KindCommon Kind = "common"

type Item struct {
	ID          objectID       `json:"_id" bson:"_id"`
	Name        string         `json:"name" bson:"name"`
	ShortName   string         `json:"shortName" bson:"shortName"`
	Description string         `json:"description" bson:"description"`
	Price       int64          `json:"price" bson:"price"`
	Weight      float64        `json:"weight" bson:"weight"`
	MaxStack    int64          `json:"maxStack" bson:"maxStack"`
	Rarity      string         `json:"rarity" bson:"rarity"`
	Grid        GridProperties `json:"grid" bson:"grid"`
	Modified    timestamp      `json:"_modified" bson:"_modified"`
	Kind        Kind           `json:"_kind" bson:"_kind"`
}

func (i Item) GetID() objectID {
	return i.ID
}

func (i Item) GetKind() Kind {
	return i.Kind
}

func (i Item) GetName() string {
	return i.Name
}

func (i Item) GetShortName() string {
	return i.ShortName
}

func (i Item) GetDescription() string {
	return i.Description
}

type ItemResult struct {
	Count int64  `json:"total"`
	Items []Item `json:"items"`
}

func (r *ItemResult) GetCount() int64 {
	return r.Count
}

func (r *ItemResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type GridProperties struct {
	Color  RGBA  `json:"color" bson:"color"`
	Height int64 `json:"height" bson:"height"`
	Width  int64 `json:"width" bson:"width"`
}

type GridModifier struct {
	Height int64 `json:"height" bson:"height"`
	Width  int64 `json:"width" bson:"width"`
}

type RGBA struct {
	R uint `json:"r" bson:"r"`
	G uint `json:"g" bson:"g"`
	B uint `json:"b" bson:"b"`
	A uint `json:"a" bson:"a"`
}
