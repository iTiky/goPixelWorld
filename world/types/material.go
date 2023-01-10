package types

import (
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
)

type MaterialType int

const (
	MaterialTypeBorder MaterialType = iota
	MaterialTypeWater
	MaterialTypeSand
	MaterialTypeWood
	MaterialTypeSmoke
	MaterialTypeFire
	MaterialTypeSteam
	MaterialTypeGrass
)

type MaterialFlag int

const (
	MaterialFlagIsSand MaterialFlag = iota
	MaterialFlagIsLiquid
	MaterialFlagIsGas
	MaterialFlagIsFire
	MaterialFlagIsFlammable
	MaterialFlagIsUnremovable
)

type Material interface {
	MaterialI

	Type() MaterialType
	ColorAdjusted(health float64) color.Color
	IsFlagged(flag MaterialFlag) bool
	Mass() float64
	InitialHealth() float64
	ProcessInternal(env TileEnvironment)
	ProcessCollision(env CollisionEnvironment)
}

type TileEnvironment interface {
	Actions() []Action

	Health() float64
	Position() Position

	ReduceHealth(step float64) bool
	ReduceEnvHealthByFlag(step float64, flagFilters ...MaterialFlag) int
	ReduceEnvHealthByType(step float64, typeFilters ...MaterialType) int
	RemoveHealthSelfReductions() bool

	AddGravity() bool
	AddReverseGravity() bool

	MoveGas() bool

	ReplaceTile(newMaterial Material, flagFilters ...MaterialFlag) bool
	ReplaceSelf(newMaterial Material) bool
	AddTile(newMaterial Material, dirFilters ...pkg.Direction) bool
	AddTileGrassStyle(newMaterial Material) bool
}

type CollisionEnvironment interface {
	Actions() []Action

	IsFlagged(flag MaterialFlag) bool

	DampSourceForce(k float64) bool

	ReflectSourceForce() bool
	ReflectSourceTargetForces() bool

	MoveSandSource() bool
	MoveLiquidSource() bool

	SwapSourceTarget() bool
}
