package item

const (
	KindModification               Kind = "modification"
	KindModificationBarrel         Kind = "modificationBarrel"
	KindModificationBipod          Kind = "modificationBipod"
	KindModificationCharge         Kind = "modificationCharge"
	KindModificationDevice         Kind = "modificationDevice"
	KindModificationForegrip       Kind = "modificationForegrip"
	KindModificationGasblock       Kind = "modificationGasblock"
	KindModificationHandguard      Kind = "modificationHandguard"
	KindModificationLauncher       Kind = "modificationLauncher"
	KindModificationMount          Kind = "modificationMount"
	KindModificationMuzzle         Kind = "modificationMuzzle"
	KindModificationGoggles        Kind = "modificationGoggles"
	KindModificationGogglesSpecial Kind = "modificationGogglesSpecial"
	KindModificationPistolgrip     Kind = "modificationPistolgrip"
	KindModificationReceiver       Kind = "modificationReceiver"
	KindModificationSight          Kind = "modificationSight"
	KindModificationSightSpecial   Kind = "modificationSightSpecial"
	KindModificationStock          Kind = "modificationStock"
)

type Modification struct {
	Item `bson:",inline"`

	Ergonomics    float64      `json:"ergonomicsFP" bson:"ergonomicsFP"`
	Accuracy      float64      `json:"accuracy" bson:"accuracy"`
	Recoil        float64      `json:"recoil" bson:"recoil"`
	RaidModdable  int64        `json:"raidModdable" bson:"raidModdable"`
	GridModifier  GridModifier `json:"gridModifier" bson:"gridModifier"`
	Slots         Slots        `json:"slots" bson:"slots"`
	Compatibility ItemList     `json:"compatibility" bson:"compatibility"`
	Conflicts     ItemList     `json:"conflicts" bson:"conflicts"`
}

type ModificationResult struct {
	Count int64          `json:"total"`
	Items []Modification `json:"items"`
}

func (r *ModificationResult) GetCount() int64 {
	return r.Count
}

func (r *ModificationResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

// Weapon modifications //

type Barrel struct {
	Modification `bson:",inline"`

	Length     float64 `json:"length" bson:"length"`
	Velocity   float64 `json:"velocity" bson:"velocity"`
	Suppressor bool    `json:"suppressor" bson:"suppressor"`
}

type BarrelResult struct {
	Count int64    `json:"total"`
	Items []Barrel `json:"items"`
}

func (r *BarrelResult) GetCount() int64 {
	return r.Count
}

func (r *BarrelResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Bipod struct {
	Modification `bson:",inline"`
}

type BipodResult struct {
	Count int64   `json:"total"`
	Items []Bipod `json:"items"`
}

func (r *BipodResult) GetCount() int64 {
	return r.Count
}

func (r *BipodResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Charge struct {
	Modification `bson:",inline"`
}

type ChargeResult struct {
	Count int64    `json:"total"`
	Items []Charge `json:"items"`
}

func (r *ChargeResult) GetCount() int64 {
	return r.Count
}

func (r *ChargeResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Device struct {
	Modification `bson:",inline"`

	Type  string   `json:"type" bson:"type"`
	Modes []string `json:"modes" bson:"modes"`
}

type DeviceResult struct {
	Count int64    `json:"total"`
	Items []Device `json:"items"`
}

func (r *DeviceResult) GetCount() int64 {
	return r.Count
}

func (r *DeviceResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Foregrip struct {
	Modification `bson:",inline"`
}

type ForegripResult struct {
	Count int64      `json:"total"`
	Items []Foregrip `json:"items"`
}

func (r *ForegripResult) GetCount() int64 {
	return r.Count
}

func (r *ForegripResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type GasBlock struct {
	Modification `bson:",inline"`
}

type GasBlockResult struct {
	Count int64      `json:"total"`
	Items []GasBlock `json:"items"`
}

func (r *GasBlockResult) GetCount() int64 {
	return r.Count
}

func (r *GasBlockResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Handguard struct {
	Modification `bson:",inline"`
}

type HandguardResult struct {
	Count int64       `json:"total"`
	Items []Handguard `json:"items"`
}

func (r *HandguardResult) GetCount() int64 {
	return r.Count
}

func (r *HandguardResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Launcher struct {
	Modification `bson:",inline"`

	Caliber string `json:"caliber" bson:"caliber"`
}

type LauncherResult struct {
	Count int64      `json:"total"`
	Items []Launcher `json:"items"`
}

func (r *LauncherResult) GetCount() int64 {
	return r.Count
}

func (r *LauncherResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Mount struct {
	Modification `bson:",inline"`
}

type MountResult struct {
	Count int64   `json:"total"`
	Items []Mount `json:"items"`
}

func (r *MountResult) GetCount() int64 {
	return r.Count
}

func (r *MountResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Muzzle struct {
	Modification `bson:",inline"`

	Type     string  `json:"type" bson:"type"`
	Velocity float64 `json:"velocity" bson:"velocity"`
}

type MuzzleResult struct {
	Count int64    `json:"total"`
	Items []Muzzle `json:"items"`
}

func (r *MuzzleResult) GetCount() int64 {
	return r.Count
}

func (r *MuzzleResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type PistolGrip struct {
	Modification `bson:",inline"`
}

type PistolGripResult struct {
	Count int64        `json:"total"`
	Items []PistolGrip `json:"items"`
}

func (r *PistolGripResult) GetCount() int64 {
	return r.Count
}

func (r *PistolGripResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Receiver struct {
	Modification `bson:",inline"`

	Velocity float64 `json:"velocity" bson:"velocity"`
}

type ReceiverResult struct {
	Count int64      `json:"total"`
	Items []Receiver `json:"items"`
}

func (r *ReceiverResult) GetCount() int64 {
	return r.Count
}

func (r *ReceiverResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Sight struct {
	Modification `bson:",inline"`

	Type          string   `json:"type" bson:"type"`
	Magnification []string `json:"magnification" bson:"magnification"`
	VariableZoom  bool     `json:"variableZoom" bson:"variableZoom"`
	ZeroDistances []int64  `json:"zeroDistances" bson:"zeroDistances"`
}

type SightResult struct {
	Count int64   `json:"total"`
	Items []Sight `json:"items"`
}

func (r *SightResult) GetCount() int64 {
	return r.Count
}

func (r *SightResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type SightSpecial struct {
	Sight        `bson:",inline"`
	OpticSpecial `bson:",inline"`
}

type SightSpecialResult struct {
	Count int64          `json:"total"`
	Items []SightSpecial `json:"items"`
}

func (r *SightSpecialResult) GetCount() int64 {
	return r.Count
}

func (r *SightSpecialResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type Stock struct {
	Modification `bson:",inline"`

	FoldRectractable bool `json:"foldRectractable" bson:"foldRectractable"`
}

type StockResult struct {
	Count int64   `json:"total"`
	Items []Stock `json:"items"`
}

func (r *StockResult) GetCount() int64 {
	return r.Count
}

func (r *StockResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

// Gear modifications //

type Goggles struct {
	Modification `bson:",inline"`

	Type string `json:"type" bson:"type"`
}

type GogglesResult struct {
	Count int64     `json:"total"`
	Items []Goggles `json:"items"`
}

func (r *GogglesResult) GetCount() int64 {
	return r.Count
}

func (r *GogglesResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

type GogglesSpecial struct {
	Goggles      `bson:",inline"`
	OpticSpecial `bson:",inline"`
}

type GogglesSpecialResult struct {
	Count int64            `json:"total"`
	Items []GogglesSpecial `json:"items"`
}

func (r *GogglesSpecialResult) GetCount() int64 {
	return r.Count
}

func (r *GogglesSpecialResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

// Properties //

type OpticSpecial struct {
	Modes []string `json:"modes" bson:"modes"`
	Color RGBA     `json:"color" bson:"color"`
	Noise string   `json:"noise" bson:"noise"`
}
