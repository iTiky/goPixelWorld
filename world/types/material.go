package types

import (
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
)

// MaterialType defines Material type.
type MaterialType int

const (
	MaterialTypeNone MaterialType = iota
	MaterialTypeBorder
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
	MaterialTypeAntiGraviton
)

// MaterialFlag defines Material property.
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

// MaterialCloseRangeType defines Material surrounding environment filling that is required to self-process a Particle.
// Used during the Tile processing step.
type MaterialCloseRangeType int

const (
	MaterialCloseRangeTypeNone MaterialCloseRangeType = iota
	MaterialCloseRangeTypeSelfOnly
	MaterialCloseRangeTypeSurrounding
	MaterialCloseRangeTypeInCircleRange
)

// Material defines the contract for all Material types.
type Material interface {
	MaterialI

	// Type returns a Material type key.
	Type() MaterialType
	// ColorAdjusted returns a Particle color adjusted by its current health (dimmed or brightened base color).
	ColorAdjusted(health float64) color.Color
	// IsFlagged checks if a Material has a set of properties.
	IsFlagged(flags ...MaterialFlag) bool
	// Mass returns a Material mass.
	Mass() float64
	// InitialHealth returns a Particle initial health.
	InitialHealth() float64

	// CloseRangeType returns a surrounding environment type that must be considered while building the closerange.Environment object.
	CloseRangeType() MaterialCloseRangeType
	// CloseRangeCircleRadius returns the circle area radius used for the MaterialCloseRangeTypeInCircleRange MaterialCloseRangeType.
	CloseRangeCircleRadius() int

	// ProcessInternal defines a handler to self-process a Particle.
	// That could change Particle's parameters and change the surrounding Particles.
	ProcessInternal(env TileEnvironment)
	// ProcessCollision defines a handler to process a collision with that Material.
	ProcessCollision(env CollisionEnvironment)
}

// TileEnvironment defines a contract each Material must implement in order to self-process a Particle.
// Environment doesn't change the Map state, instead it builds a set of Actions to apply.
// Environment holds the source Particle data alongside its neighbours.
// Each "action" method returns true, if the required Action(s) was added (it might fail, for ex., if a proper neighbour wasn't found).
type TileEnvironment interface {
	// Actions returns output Actions that should be applied for this environment.
	Actions() []Action

	// Health returns the Particle current health.
	Health() float64
	// Position returns the Particle current Position.
	Position() Position

	// StateParam returns the Particle internal state parameter by key.
	StateParam(key string) int
	// UpdateStateParam adds an Action to set the internal Particle state parameter by key.
	UpdateStateParam(paramKey string, paramValue int) bool

	// DampSelfHealth adds an Action to alter the Particle's health (increase / decrease).
	DampSelfHealth(step float64) bool
	// DampEnvHealthByFlag adds Actions to alter the Particle's neighbours' health.
	// {flagFilters} filters source neighbours by MaterialFlag (AND, including).
	DampEnvHealthByFlag(step float64, flagFilters ...MaterialFlag) int
	// DampEnvHealthByType adds Actions to alter the Particle's neighbours' health.
	// {typeFilters} filters source neighbours by MaterialType (AND, including).
	DampEnvHealthByType(step float64, typeFilters ...MaterialType) int
	// RemoveSelfHealthDamps removes any previously added environment Actions that reduces the Particle health.
	RemoveSelfHealthDamps() bool

	// AddGravity adds an Action which add the vertical gravity to the Particle's force Vector.
	AddGravity() bool
	// AddReverseGravity adds an Action which add the inverted vertical gravity to the Particle's force Vector.
	AddReverseGravity() bool
	// AddForceInRange adds Actions that add the provided magnitude to neighbours' force Vectors.
	// {notFlagFilters} filters source neighbours by MaterialType (AND, NOT including).
	AddForceInRange(mag float64, notFlagFilters ...MaterialFlag) bool

	// MoveGas adds an Action that implements gas-like Material movement.
	MoveGas() bool

	// ReplaceTile adds Actions that replaces a particular Tile's Particle with a new one.
	// {flagFilters} filters source neighbours by MaterialFlag (AND, including).
	ReplaceTile(newMaterial Material, flagFilters ...MaterialFlag) bool
	// ReplaceSelf adds an Action that replaces the Particle with a new one.
	ReplaceSelf(newMaterial Material) bool
	// AddTile adds Actions that creates a new Particle(s).
	// {dirFilters} filters source neighbours by their relative position.
	AddTile(newMaterial Material, dirFilters ...pkg.Direction) bool
	// AddTileGrassStyle adds an Action that creates a new Particle for grass-like Material.
	AddTileGrassStyle(newMaterial Material) bool
}

// CollisionEnvironment defines a contract each Material must implement in order to process Particles collision.
// Environment doesn't change the Map state, instead it builds a set of Actions to apply.
// Environment holds the target and source Particle data alongside target neighbours.
// Source Particle is the one which "want" to collide with the target one.
// Collision logic is defined by the target's Material and it can alter both (source and target).
// Each "action" method returns true, if the required Action(s) was added (it might fail, for ex., if a proper neighbour wasn't found).
type CollisionEnvironment interface {
	// Actions returns output Actions that should be applied for this environment.
	Actions() []Action

	// IsFlagged checks if the source Particle has the specified MaterialFlag.
	IsFlagged(flag MaterialFlag) bool
	// IsType checks if the source Particle has the specified MaterialFlag.
	IsType(mType MaterialType) bool

	// DampSourceForce adds an Action that alters the source Particle force Vector.
	DampSourceForce(k float64) bool
	// DampSourceHealth adds an Action that alters the source Particle health.
	DampSourceHealth(step float64, flagFilters ...MaterialFlag) bool
	// DampSelfHealth adds an Action that alters the target Particle health.
	DampSelfHealth(step float64) bool
	// DampSelfHealthByMassRate adds an Action that alters the target Particle health.
	// An actual {step} value is adjusted according to source / target mass ratio.
	DampSelfHealthByMassRate(step float64) bool

	// ReflectSourceForce adds an Action that modifies the source Particle force Vector.
	// This is a simple reflection logic.
	ReflectSourceForce() bool
	// ReflectSourceTargetForces adds Actions that modifies both source and target Particle force Vector.
	// This is an advanced reflection logic that takes in consideration angles and masses.
	ReflectSourceTargetForces(sourceForceDampK float64) bool

	// MoveSandSource adds Actions that performs reflection logic for sand-like Material.
	MoveSandSource() bool
	// MoveLiquidSource adds Actions that performs reflection logic for liquid-like Material.
	MoveLiquidSource() bool

	// SwapSourceTarget adds Actions that swaps source and target Particles.
	SwapSourceTarget() bool
}
