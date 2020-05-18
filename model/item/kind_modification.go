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

	Ergonomics    float64      `json:"ergonomicsFP"`
	Accuracy      float64      `json:"accuracy"`
	Recoil        float64      `json:"recoil"`
	RaidModdable  int64        `json:"raidModdable"`
	GridModifier  GridModifier `json:"gridModifier"`
	Slots         Slots        `json:"slots"`
	Compatibility ItemList     `json:"compatibility"`
	Conflicts     ItemList     `json:"conflicts"`
}

type ModificationResult struct {
	*Result
	Items []Modification `json:"items"`
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

	Length     float64 `json:"length"`
	Velocity   float64 `json:"velocity"`
	Suppressor bool    `json:"suppressor"`
}

type BarrelResult struct {
	*Result
	Items []Barrel `json:"items"`
}

func (r *BarrelResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var modBarrelFilter = Filter{
	"suppressor": {
		"true",
		"false",
	},
}

type Bipod struct {
	Modification `bson:",inline"`
}

type BipodResult struct {
	*Result
	Items []Bipod `json:"items"`
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
	*Result
	Items []Charge `json:"items"`
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

	Type  string   `json:"type"`
	Modes []string `json:"modes"`
}

type DeviceResult struct {
	*Result
	Items []Device `json:"items"`
}

func (r *DeviceResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var modDeviceFilter = Filter{
	"type": {
		"combo",
		"light",
	},
}

type Foregrip struct {
	Modification `bson:",inline"`
}

type ForegripResult struct {
	*Result
	Items []Foregrip `json:"items"`
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
	*Result
	Items []GasBlock `json:"items"`
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
	*Result
	Items []Handguard `json:"items"`
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

	Caliber string `json:"caliber"`
}

type LauncherResult struct {
	*Result
	Items []Launcher `json:"items"`
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
	*Result
	Items []Mount `json:"items"`
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

	Type     string  `json:"type"`
	Velocity float64 `json:"velocity"`
}

type MuzzleResult struct {
	*Result
	Items []Muzzle `json:"items"`
}

func (r *MuzzleResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var modMuzzleFilter = Filter{
	"type": {
		"brake",
		"combo",
		"compensator",
		"supressor",
	},
}

type PistolGrip struct {
	Modification `bson:",inline"`
}

type PistolGripResult struct {
	*Result
	Items []PistolGrip `json:"items"`
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

	Velocity float64 `json:"velocity"`
}

type ReceiverResult struct {
	*Result
	Items []Receiver `json:"items"`
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

	Type          string   `json:"type"`
	Magnification []string `json:"magnification"`
	VariableZoom  bool     `json:"variableZoom"`
	ZeroDistances []int64  `json:"zeroDistances"`
}

type SightResult struct {
	*Result
	Items []Sight `json:"items"`
}

func (r *SightResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var modSightFilter = Filter{
	"type": {
		"hybrid",
		"iron",
		"reflex",
		"telescopic",
	},
}

type SightSpecial struct {
	Sight        `bson:",inline"`
	OpticSpecial `bson:",inline"`
}

type SightSpecialResult struct {
	*Result
	Items []SightSpecial `json:"items"`
}

func (r *SightSpecialResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var modSightSpecialFilter = Filter{
	"type": {
		"nightVision",
		"thermalVision",
	},
}

type Stock struct {
	Modification `bson:",inline"`

	FoldRectractable bool `json:"foldRectractable"`
}

type StockResult struct {
	*Result
	Items []Stock `json:"items"`
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

	Type string `json:"type"`
}

type GogglesResult struct {
	*Result
	Items []Goggles `json:"items"`
}

func (r *GogglesResult) GetEntities() []Entity {
	e := make([]Entity, len(r.Items))
	for i, item := range r.Items {
		e[i] = item
	}

	return e
}

var modGogglesFilter = Filter{
	"type": {
		"nightVision",
		"thermalVision",
	},
}

type GogglesSpecial struct {
	Goggles      `bson:",inline"`
	OpticSpecial `bson:",inline"`
}

type GogglesSpecialResult struct {
	*Result
	Items []GogglesSpecial `json:"items"`
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
	Modes []string `json:"modes"`
	Color RGBA     `json:"color"`
	Noise string   `json:"noise"`
}
