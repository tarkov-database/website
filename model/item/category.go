package item

import (
	"errors"
)

var (
	ErrInvalidCategory = errors.New("category is not valid")
)

const (
	CategoryAmmunition               = "ammunition"
	CategoryArmor                    = "armor"
	CategoryBackpack                 = "backpack"
	CategoryBarter                   = "barter"
	CategoryClothing                 = "clothing"
	CategoryCommon                   = "common"
	CategoryContainer                = "container"
	CategoryFirearm                  = "firearm"
	CategoryFood                     = "food"
	CategoryGrenade                  = "grenade"
	CategoryHeadphone                = "headphone"
	CategoryKey                      = "key"
	CategoryMagazine                 = "magazine"
	CategoryMap                      = "map"
	CategoryMedical                  = "medical"
	CategoryMelee                    = "melee"
	CategoryModification             = "mod-other"
	CategoryModificationAuxiliary    = "auxiliary"
	CategoryModificationBarrel       = "barrel"
	CategoryModificationBipod        = "bipod"
	CategoryModificationCharge       = "charge"
	CategoryModificationDevice       = "device"
	CategoryModificationForegrip     = "foregrip"
	CategoryModificationGasblock     = "gasblock"
	CategoryModificationGoggles      = "goggles"
	CategoryModificationHandguard    = "handguard"
	CategoryModificationLauncher     = "launcher"
	CategoryModificationMount        = "mount"
	CategoryModificationMuzzle       = "muzzle"
	CategoryModificationPistolgrip   = "pistolgrip"
	CategoryModificationReceiver     = "receiver"
	CategoryModificationSight        = "sight"
	CategoryModificationSightSpecial = "sight-special"
	CategoryModificationStock        = "stock"
	CategoryMoney                    = "money"
	CategoryTacticalrig              = "tacticalrig"
)

func CategoryToKind(s string) (Kind, error) {
	var k Kind

	switch s {
	case CategoryAmmunition:
		k = KindAmmunition
	case CategoryArmor:
		k = KindArmor
	case CategoryBackpack:
		k = KindBackpack
	case CategoryBarter:
		k = KindBarter
	case CategoryClothing:
		k = KindClothing
	case CategoryCommon:
		k = KindCommon
	case CategoryContainer:
		k = KindContainer
	case CategoryFirearm:
		k = KindFirearm
	case CategoryFood:
		k = KindFood
	case CategoryGrenade:
		k = KindGrenade
	case CategoryHeadphone:
		k = KindHeadphone
	case CategoryKey:
		k = KindKey
	case CategoryMagazine:
		k = KindMagazine
	case CategoryMap:
		k = KindMap
	case CategoryMedical:
		k = KindMedical
	case CategoryMelee:
		k = KindMelee
	case CategoryModification:
		k = KindModification
	case CategoryModificationAuxiliary:
		k = KindModificationAuxiliary
	case CategoryModificationBarrel:
		k = KindModificationBarrel
	case CategoryModificationBipod:
		k = KindModificationBipod
	case CategoryModificationCharge:
		k = KindModificationCharge
	case CategoryModificationDevice:
		k = KindModificationDevice
	case CategoryModificationForegrip:
		k = KindModificationForegrip
	case CategoryModificationGasblock:
		k = KindModificationGasblock
	case CategoryModificationGoggles:
		k = KindModificationGoggles
	case CategoryModificationHandguard:
		k = KindModificationHandguard
	case CategoryModificationLauncher:
		k = KindModificationLauncher
	case CategoryModificationMount:
		k = KindModificationMount
	case CategoryModificationMuzzle:
		k = KindModificationMuzzle
	case CategoryModificationPistolgrip:
		k = KindModificationPistolgrip
	case CategoryModificationReceiver:
		k = KindModificationReceiver
	case CategoryModificationSight:
		k = KindModificationSight
	case CategoryModificationSightSpecial:
		k = KindModificationSightSpecial
	case CategoryModificationStock:
		k = KindModificationStock
	case CategoryMoney:
		k = KindMoney
	case CategoryTacticalrig:
		k = KindTacticalrig
	default:
		return k, ErrInvalidCategory
	}

	return k, nil
}

func KindToCategory(k Kind) (string, error) {
	var c string

	switch k {
	case KindAmmunition:
		c = CategoryAmmunition
	case KindArmor:
		c = CategoryArmor
	case KindBackpack:
		c = CategoryBackpack
	case KindBarter:
		c = CategoryBarter
	case KindClothing:
		c = CategoryClothing
	case KindCommon:
		c = CategoryCommon
	case KindContainer:
		c = CategoryContainer
	case KindFirearm:
		c = CategoryFirearm
	case KindFood:
		c = CategoryFood
	case KindGrenade:
		c = CategoryGrenade
	case KindHeadphone:
		c = CategoryHeadphone
	case KindKey:
		c = CategoryKey
	case KindMagazine:
		c = CategoryMagazine
	case KindMap:
		c = CategoryMap
	case KindMedical:
		c = CategoryMedical
	case KindMelee:
		c = CategoryMelee
	case KindModification:
		c = CategoryModification
	case KindModificationAuxiliary:
		c = CategoryModificationAuxiliary
	case KindModificationBarrel:
		c = CategoryModificationBarrel
	case KindModificationBipod:
		c = CategoryModificationBipod
	case KindModificationCharge:
		c = CategoryModificationCharge
	case KindModificationDevice:
		c = CategoryModificationDevice
	case KindModificationForegrip:
		c = CategoryModificationForegrip
	case KindModificationGasblock:
		c = CategoryModificationGasblock
	case KindModificationGoggles:
		c = CategoryModificationGoggles
	case KindModificationHandguard:
		c = CategoryModificationHandguard
	case KindModificationLauncher:
		c = CategoryModificationLauncher
	case KindModificationMount:
		c = CategoryModificationMount
	case KindModificationMuzzle:
		c = CategoryModificationMuzzle
	case KindModificationPistolgrip:
		c = CategoryModificationPistolgrip
	case KindModificationReceiver:
		c = CategoryModificationReceiver
	case KindModificationSight:
		c = CategoryModificationSight
	case KindModificationSightSpecial:
		c = CategoryModificationSightSpecial
	case KindModificationStock:
		c = CategoryModificationStock
	case KindMoney:
		c = CategoryMoney
	case KindTacticalrig:
		c = CategoryTacticalrig
	default:
		return c, ErrInvalidKind
	}

	return c, nil
}

func CategoryToDisplayName(s string) (string, error) {
	var c string

	switch s {
	case CategoryAmmunition:
		c = "ammunition"
	case CategoryArmor:
		c = "armor"
	case CategoryBackpack:
		c = "backpack"
	case CategoryBarter:
		c = "barter"
	case CategoryClothing:
		c = "clothing"
	case CategoryCommon:
		c = "common"
	case CategoryContainer:
		c = "container"
	case CategoryFirearm:
		c = "firearm"
	case CategoryFood:
		c = "food"
	case CategoryGrenade:
		c = "grenade"
	case CategoryHeadphone:
		c = "headphone"
	case CategoryKey:
		c = "key"
	case CategoryMagazine:
		c = "magazine"
	case CategoryMap:
		c = "map"
	case CategoryMedical:
		c = "medical"
	case CategoryMelee:
		c = "melee"
	case CategoryModification:
		c = "modification other"
	case CategoryModificationBarrel:
		c = "barrel"
	case CategoryModificationAuxiliary:
		c = "auxiliary"
	case CategoryModificationBipod:
		c = "bipod"
	case CategoryModificationCharge:
		c = "charging handle"
	case CategoryModificationDevice:
		c = "device"
	case CategoryModificationForegrip:
		c = "foregrip"
	case CategoryModificationGasblock:
		c = "gas block"
	case CategoryModificationGoggles:
		c = "goggles"
	case CategoryModificationHandguard:
		c = "handguard"
	case CategoryModificationLauncher:
		c = "launcher"
	case CategoryModificationMount:
		c = "mount"
	case CategoryModificationMuzzle:
		c = "muzzle"
	case CategoryModificationPistolgrip:
		c = "pistol grip"
	case CategoryModificationReceiver:
		c = "receiver"
	case CategoryModificationSight:
		c = "sight"
	case CategoryModificationSightSpecial:
		c = "special sight"
	case CategoryModificationStock:
		c = "stock"
	case CategoryMoney:
		c = "money"
	case CategoryTacticalrig:
		c = "tactical rig"
	default:
		return c, ErrInvalidCategory
	}

	return c, nil
}
