package item

type ItemList map[Kind][]objectID

const KindCommon Kind = "common"

type Item struct {
	ID          objectID       `json:"_id"`
	Name        string         `json:"name"`
	ShortName   string         `json:"shortName"`
	Description string         `json:"description"`
	Weight      float64        `json:"weight"`
	MaxStack    int64          `json:"maxStack"`
	Grid        GridProperties `json:"grid"`
	Modified    timestamp      `json:"_modified"`
	Kind        Kind           `json:"_kind"`
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
	*Result
	Items []Item `json:"items"`
}

func (r *ItemResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type GridProperties struct {
	Color  RGBA  `json:"color"`
	Height int64 `json:"height"`
	Width  int64 `json:"width"`
}

type GridModifier struct {
	Height int64 `json:"height"`
	Width  int64 `json:"width"`
}

type RGBA struct {
	R uint `json:"r"`
	G uint `json:"g"`
	B uint `json:"b"`
	A uint `json:"a"`
}
