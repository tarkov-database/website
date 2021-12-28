package item

import (
	"encoding/json"
	"errors"

	"github.com/tarkov-database/website/core/api"
)

var (
	ErrInvalidKind = errors.New("kind is not valid")
)

type objectID = string

type timestamp = api.Timestamp

type Entity interface {
	GetID() objectID
	GetKind() Kind
	GetName() string
	GetShortName() string
	GetDescription() string
}

type EntityResult interface {
	GetKind() Kind
	GetCount() int64
	GetEntities() []Entity
}

type Kind string

func (k Kind) IsValid() bool {
	if _, err := k.GetEntity(); err != nil {
		return false
	}

	return true
}

func (k Kind) String() string {
	return string(k)
}

func (k *Kind) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func (k *Kind) UnmarshalJSON(b []byte) error {
	var kind string

	err := json.Unmarshal(b, &kind)
	if err != nil {
		return err
	}

	*k = Kind(kind)

	return nil
}

func (k Kind) GetEntity() (Entity, error) {
	var e Entity

	switch k {
	case KindAmmunition:
		e = &Ammunition{}
	case KindArmor:
		e = &Armor{}
	case KindBackpack:
		e = &Backpack{}
	case KindBarter:
		e = &Barter{}
	case KindClothing:
		e = &Clothing{}
	case KindCommon:
		e = &Item{}
	case KindContainer:
		e = &Container{}
	case KindFirearm:
		e = &Firearm{}
	case KindFood:
		e = &Food{}
	case KindGrenade:
		e = &Grenade{}
	case KindHeadphone:
		e = &Headphone{}
	case KindKey:
		e = &Key{}
	case KindMagazine:
		e = &Magazine{}
	case KindMap:
		e = &Map{}
	case KindMedical:
		e = &Medical{}
	case KindMelee:
		e = &Melee{}
	case KindModification:
		e = &Modification{}
	case KindModificationAuxiliary:
		e = &Auxiliary{}
	case KindModificationBarrel:
		e = &Barrel{}
	case KindModificationBipod:
		e = &Bipod{}
	case KindModificationCharge:
		e = &Charge{}
	case KindModificationDevice:
		e = &Device{}
	case KindModificationForegrip:
		e = &Foregrip{}
	case KindModificationGasblock:
		e = &GasBlock{}
	case KindModificationGoggles:
		e = &Goggles{}
	case KindModificationGogglesSpecial:
		e = &GogglesSpecial{}
	case KindModificationHandguard:
		e = &Handguard{}
	case KindModificationLauncher:
		e = &Launcher{}
	case KindModificationMount:
		e = &Mount{}
	case KindModificationMuzzle:
		e = &Muzzle{}
	case KindModificationPistolgrip:
		e = &PistolGrip{}
	case KindModificationReceiver:
		e = &Receiver{}
	case KindModificationSight:
		e = &Sight{}
	case KindModificationSightSpecial:
		e = &SightSpecial{}
	case KindModificationStock:
		e = &Stock{}
	case KindMoney:
		e = &Money{}
	case KindTacticalrig:
		e = &TacticalRig{}
	default:
		return e, ErrInvalidKind
	}

	return e, nil
}

func (k Kind) GetEntityResult() (EntityResult, error) {
	var r EntityResult

	switch k {
	case KindAmmunition:
		r = &AmmunitionResult{Result: &Result{Kind: KindAmmunition}}
	case KindArmor:
		r = &ArmorResult{Result: &Result{Kind: KindArmor}}
	case KindBackpack:
		r = &BackpackResult{Result: &Result{Kind: KindBackpack}}
	case KindBarter:
		r = &BarterResult{Result: &Result{Kind: KindBarter}}
	case KindClothing:
		r = &ClothingResult{Result: &Result{Kind: KindClothing}}
	case KindCommon:
		r = &ItemResult{Result: &Result{Kind: KindCommon}}
	case KindContainer:
		r = &ContainerResult{Result: &Result{Kind: KindContainer}}
	case KindFirearm:
		r = &FirearmResult{Result: &Result{Kind: KindFirearm}}
	case KindFood:
		r = &FoodResult{Result: &Result{Kind: KindFood}}
	case KindGrenade:
		r = &GrenadeResult{Result: &Result{Kind: KindGrenade}}
	case KindHeadphone:
		r = &HeadphoneResult{Result: &Result{Kind: KindHeadphone}}
	case KindKey:
		r = &KeyResult{Result: &Result{Kind: KindKey}}
	case KindMagazine:
		r = &MagazineResult{Result: &Result{Kind: KindMagazine}}
	case KindMap:
		r = &MapResult{Result: &Result{Kind: KindMap}}
	case KindMedical:
		r = &MedicalResult{Result: &Result{Kind: KindMedical}}
	case KindMelee:
		r = &MeleeResult{Result: &Result{Kind: KindMelee}}
	case KindModification:
		r = &ModificationResult{Result: &Result{Kind: KindModification}}
	case KindModificationAuxiliary:
		r = &AuxiliaryResult{Result: &Result{Kind: KindModificationAuxiliary}}
	case KindModificationBarrel:
		r = &BarrelResult{Result: &Result{Kind: KindModificationBarrel}}
	case KindModificationBipod:
		r = &BipodResult{Result: &Result{Kind: KindModificationBipod}}
	case KindModificationCharge:
		r = &ChargeResult{Result: &Result{Kind: KindModificationCharge}}
	case KindModificationDevice:
		r = &DeviceResult{Result: &Result{Kind: KindModificationDevice}}
	case KindModificationForegrip:
		r = &ForegripResult{Result: &Result{Kind: KindModificationForegrip}}
	case KindModificationGasblock:
		r = &GasBlockResult{Result: &Result{Kind: KindModificationGasblock}}
	case KindModificationGoggles:
		r = &GogglesResult{Result: &Result{Kind: KindModificationGoggles}}
	case KindModificationHandguard:
		r = &HandguardResult{Result: &Result{Kind: KindModificationHandguard}}
	case KindModificationLauncher:
		r = &LauncherResult{Result: &Result{Kind: KindModificationLauncher}}
	case KindModificationMount:
		r = &MountResult{Result: &Result{Kind: KindModificationMount}}
	case KindModificationMuzzle:
		r = &MuzzleResult{Result: &Result{Kind: KindModificationMuzzle}}
	case KindModificationPistolgrip:
		r = &PistolGripResult{Result: &Result{Kind: KindModificationPistolgrip}}
	case KindModificationReceiver:
		r = &ReceiverResult{Result: &Result{Kind: KindModificationReceiver}}
	case KindModificationSight:
		r = &SightResult{Result: &Result{Kind: KindModificationSight}}
	case KindModificationSightSpecial:
		r = &SightSpecialResult{Result: &Result{Kind: KindModificationSightSpecial}}
	case KindModificationStock:
		r = &StockResult{Result: &Result{Kind: KindModificationStock}}
	case KindMoney:
		r = &MoneyResult{Result: &Result{Kind: KindMoney}}
	case KindTacticalrig:
		r = &TacticalRigResult{Result: &Result{Kind: KindTacticalrig}}
	default:
		return r, ErrInvalidKind
	}

	return r, nil
}

var KindList = [...]Kind{
	KindAmmunition,
	KindArmor,
	KindBackpack,
	KindBarter,
	KindClothing,
	KindCommon,
	KindContainer,
	KindFirearm,
	KindFood,
	KindGrenade,
	KindHeadphone,
	KindKey,
	KindMagazine,
	KindMap,
	KindMedical,
	KindMelee,
	KindModification,
	KindModificationAuxiliary,
	KindModificationBarrel,
	KindModificationBipod,
	KindModificationCharge,
	KindModificationDevice,
	KindModificationForegrip,
	KindModificationGasblock,
	KindModificationGoggles,
	KindModificationHandguard,
	KindModificationLauncher,
	KindModificationMount,
	KindModificationMuzzle,
	KindModificationPistolgrip,
	KindModificationReceiver,
	KindModificationSight,
	KindModificationSightSpecial,
	KindModificationStock,
	KindMoney,
	KindTacticalrig,
}

type Result struct {
	Kind  Kind  `json:"kind"`
	Count int64 `json:"total"`
}

func (r Result) GetKind() Kind {
	return r.Kind
}

func (r Result) GetCount() int64 {
	return r.Count
}

type Filter map[string][]string

func (f Filter) GetAll() map[string][]string {
	return f
}

func (f Filter) Get(k string) []string {
	return f[k]
}

func (k Kind) GetFilter() Filter {
	switch k {
	case KindAmmunition:
		return ammunitionFilter
	case KindArmor:
		return armorFilter
	case KindClothing:
		return clothingFilter
	case KindFirearm:
		return firearmFilter
	case KindFood:
		return foodFilter
	case KindGrenade:
		return grenadeFilter
	case KindMagazine:
		return magazineFilter
	case KindMedical:
		return medicalFilter
	case KindModificationBarrel:
		return modBarrelFilter
	case KindModificationDevice:
		return modDeviceFilter
	case KindModificationMuzzle:
		return modMuzzleFilter
	case KindModificationSight:
		return modSightFilter
	case KindModificationSightSpecial:
		return modSightSpecialFilter
	case KindModificationGoggles, KindModificationGogglesSpecial:
		return modGogglesFilter
	case KindTacticalrig:
		return tacticalFilter
	default:
		return Filter{}
	}
}
