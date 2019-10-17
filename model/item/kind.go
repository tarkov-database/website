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
		r = &AmmunitionResult{}
	case KindArmor:
		r = &ArmorResult{}
	case KindBackpack:
		r = &BackpackResult{}
	case KindBarter:
		r = &BarterResult{}
	case KindClothing:
		r = &ClothingResult{}
	case KindCommon:
		r = &ItemResult{}
	case KindContainer:
		r = &ContainerResult{}
	case KindFirearm:
		r = &FirearmResult{}
	case KindFood:
		r = &FoodResult{}
	case KindGrenade:
		r = &GrenadeResult{}
	case KindHeadphone:
		r = &HeadphoneResult{}
	case KindKey:
		r = &KeyResult{}
	case KindMagazine:
		r = &MagazineResult{}
	case KindMap:
		r = &MapResult{}
	case KindMedical:
		r = &MedicalResult{}
	case KindMelee:
		r = &MeleeResult{}
	case KindModification:
		r = &ModificationResult{}
	case KindModificationBarrel:
		r = &BarrelResult{}
	case KindModificationBipod:
		r = &BipodResult{}
	case KindModificationCharge:
		r = &ChargeResult{}
	case KindModificationDevice:
		r = &DeviceResult{}
	case KindModificationForegrip:
		r = &ForegripResult{}
	case KindModificationGasblock:
		r = &GasBlockResult{}
	case KindModificationGoggles:
		r = &GogglesResult{}
	case KindModificationHandguard:
		r = &HandguardResult{}
	case KindModificationLauncher:
		r = &LauncherResult{}
	case KindModificationMount:
		r = &MountResult{}
	case KindModificationMuzzle:
		r = &MuzzleResult{}
	case KindModificationPistolgrip:
		r = &PistolGripResult{}
	case KindModificationReceiver:
		r = &ReceiverResult{}
	case KindModificationSight:
		r = &SightResult{}
	case KindModificationSightSpecial:
		r = &SightSpecialResult{}
	case KindModificationStock:
		r = &StockResult{}
	case KindMoney:
		r = &MoneyResult{}
	case KindTacticalrig:
		r = &TacticalRigResult{}
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
