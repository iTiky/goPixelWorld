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
	MaterialTypeMetal
	MaterialTypeRock
	MaterialTypeGraviton
)

type MaterialFlag int

const (
	MaterialFlagIsSand MaterialFlag = iota
	MaterialFlagIsLiquid
	MaterialFlagIsGas
	MaterialFlagIsFire
	MaterialFlagIsFlammable
	MaterialFlagIsUnremovable
	MaterialFlagIsUnmovable
)

type MaterialCloseRangeType int

const (
	MaterialCloseRangeTypeNone MaterialCloseRangeType = iota
	MaterialCloseRangeTypeSelfOnly
	MaterialCloseRangeTypeSurrounding
	MaterialCloseRangeTypeInCircleRange
)

type Material interface {
	MaterialI

	Type() MaterialType
	ColorAdjusted(health float64) color.Color
	IsFlagged(flags ...MaterialFlag) bool
	Mass() float64
	InitialHealth() float64

	CloseRangeType() MaterialCloseRangeType
	CloseRangeCircleRadius() int

	ProcessInternal(env TileEnvironment)
	ProcessCollision(env CollisionEnvironment)
}

type TileEnvironment interface {
	Actions() []Action

	Health() float64
	Position() Position

	StateParam(key string) int
	UpdateStateParam(paramKey string, paramValue int) bool

	DampSelfHealth(step float64) bool
	DampEnvHealthByFlag(step float64, flagFilters ...MaterialFlag) int
	DampEnvHealthByType(step float64, typeFilters ...MaterialType) int
	RemoveSelfHealthDamps() bool

	AddGravity() bool
	AddReverseGravity() bool
	AddForceInRange(mag float64, notFlagFilters ...MaterialFlag) bool

	MoveGas() bool

	ReplaceTile(newMaterial Material, flagFilters ...MaterialFlag) bool
	ReplaceSelf(newMaterial Material) bool
	AddTile(newMaterial Material, dirFilters ...pkg.Direction) bool
	AddTileGrassStyle(newMaterial Material) bool
}

type CollisionEnvironment interface {
	Actions() []Action

	IsFlagged(flag MaterialFlag) bool
	IsType(mType MaterialType) bool

	DampSourceForce(k float64) bool
	DampSourceHealth(step float64, flagFilters ...MaterialFlag) bool
	DampSelfHealth(step float64) bool
	DampSelfHealthByMassRate(step float64) bool

	ReflectSourceForce() bool
	ReflectSourceTargetForces(sourceForceDampK float64) bool

	MoveSandSource() bool
	MoveLiquidSource() bool

	SwapSourceTarget() bool
}
