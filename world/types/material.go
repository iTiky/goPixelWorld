package types

import (
	"image/color"
)

type MaterialType int

const (
	MaterialTypeBorder MaterialType = iota
	MaterialTypeWater
	MaterialTypeSand
	MaterialTypeWood
	MaterialTypeSmoke
	MaterialTypeFire
)

type MaterialFlag int

const (
	MaterialFlagIsSand MaterialFlag = iota
	MaterialFlagIsLiquid
	MaterialFlagIsGas
	MaterialFlagIsFlammable
	MaterialFlagIsUnremovable
)

type Material interface {
	MaterialI

	Type() MaterialType
	ColorAdjusted(health float64) color.Color
	IsFlagged(flag MaterialFlag) bool
	Mass() float64
	ProcessInternal(env TileEnvironment)
	ProcessCollision(env CollisionEnvironment)
}

type TileEnvironment interface {
	Actions() []Action

	Health() float64
	Position() Position

	ReduceHealth(step float64) bool

	AddGravity() bool
	AddReverseGravity() bool

	MoveGas() bool

	ReplaceTile(newMaterial Material, flagFilters ...MaterialFlag) bool
	AddTile(newMaterial Material) bool
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
