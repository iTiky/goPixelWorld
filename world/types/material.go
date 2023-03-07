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
	MaterialTypeBug
)

func (t MaterialType) String() string {
	switch t {
	case MaterialTypeBorder:
		return "Border"
	case MaterialTypeWater:
		return "Water"
	case MaterialTypeSand:
		return "Sand"
	case MaterialTypeWood:
		return "Wood"
	case MaterialTypeSmoke:
		return "Smoke"
	case MaterialTypeFire:
		return "Fire"
	case MaterialTypeSteam:
		return "Steam"
	case MaterialTypeGrass:
		return "Grass"
	case MaterialTypeMetal:
		return "Metal"
	case MaterialTypeRock:
		return "Rock"
	case MaterialTypeGraviton:
		return "Graviton"
	case MaterialTypeAntiGraviton:
		return "Anti-Graviton"
	case MaterialTypeBug:
		return "Bug"
	}

	return ""
}

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
	// ForceVec returns the Particle current force Vector.
	ForceVec() pkg.Vector
	// StateParam return the internal Particle state parameter value.
	StateParam(key string) int

	// AddGravity adds the vertical gravity force Vector to the Particle.
	AddGravity() (isApplied bool)
	// AddReverseGravity adds the reversed vertical gravity force Vector to the Particle.
	AddReverseGravity() (isApplied bool)

	// ReplaceSelf replaces the Particle with a new one.
	ReplaceSelf(newMaterial Material) (flagIn bool)
	// UpdateStateParam updates the internal Particle state param.
	UpdateStateParam(paramKey string, paramValue int) (flagIn bool)
	// AddSelfForce adds a new force to the Particle.
	AddSelfForce(vec pkg.Vector) (flagIn bool)
	// SetSelfForce sets a new force to the Particle.
	SetSelfForce(vec pkg.Vector) (flagIn bool)
	// DampSelfHealth alters the Particle health.
	DampSelfHealth(step float64) (flagIn bool)
	// RemoveSelfHealthDamps removes previously added Particle health reduction Actions.
	RemoveSelfHealthDamps() (flagIn bool)

	// AddNewNeighbourTile adds a new neighbour Particle.
	// Candidate is selected randomly from empty neighbours matching the filter.
	// {dirFilters} filter includes candidates IN directions.
	AddNewNeighbourTile(newMaterial Material, dirFilters []pkg.Direction) (isApplied bool)
	// ReplaceNeighbourTile replaces a neighbour Tile with a new Particle.
	// Candidate is selected randomly from non-empty neighbours matching the filter.
	// {flagFilters} filter includes candidates WITH flags.
	ReplaceNeighbourTile(newMaterial Material, flagFilters []MaterialFlag) (isApplied bool)
	// AddNewNeighbourTileGrassStyle adds a new grass-like neighbour.
	// Candidate select criteria:
	//   - three close empty Tiles (for ex.: Top-Left, Top and Top-Right);
	//   - random;
	AddNewNeighbourTileGrassStyle(newMaterial Material) (isApplied bool)
	// DampNeighboursHealthByFlag alters neighbour(s) health.
	// Non-empty candidates are selected by AND filters.
	// {typeFilters} filter includes candidates MATCHING types.
	// {flagFilters} filter includes candidates WITH flags.
	DampNeighboursHealthByFlag(step float64, typeFilters []MaterialType, flagFilters []MaterialFlag) (neighboursAffectedCnt int)
	// MoveTileWithNeighboursGasStyle moves the Particle up for gas-like Materials.
	// Move to a randomly selected empty Tile in the upper direction sector (Top-Left + Top + Top-Right).
	MoveTileWithNeighboursGasStyle() (isApplied bool)
	// SearchNeighbours performs a full neighbours search.
	// {isEmpty} if not nil, defines empty/non-empty Tile criteria.
	// {dirsFilter, dirsIn} filter includes particles IN direction relative to the source ({dirsIn} = false, inverts the filter making it NOT).
	// {typeFilters, typeIn} filter includes candidates MATCHING types ({typeIn} = false, inverts the filter making it NOT).
	// {flagFilters, flagIn} filter includes candidates WITH flags ({flagIn} = false, inverts the filter making it NOT).
	SearchNeighbours(
		isEmpty *bool,
		dirsFilter []pkg.Direction, dirsIn bool,
		typeFilters []MaterialType, typeIn bool,
		flagFilters []MaterialFlag, flagIn bool,
	) (tiles []*Tile, tileDirs []pkg.Direction)

	// AddNewTileInRange adds a new Particle in a circle range.
	// Candidate is selected randomly from empty tiles.
	AddNewTileInRange(newMaterial Material) bool
	// AddForceInRange adds a force Vector in a circle range.
	// If {mag} is LT 0, force is reflected.
	// Non-empty candidates are selected by NOT filter.
	// {notFlagFilters} filters OUT particles WITH flags.
	AddForceInRange(mag float64, notFlagFilters []MaterialFlag) (isApplied bool)
	// DampEnvHealthByTypeInRange alters tiles in a circle range health.
	// Non-empty candidates are selected by AND filters and distance limit.
	// {distance} limits a candidate distance to the source.
	// {typeFilters} filter includes candidates MATCHING types.
	// {flagFilters} filter includes candidates WITH flags.
	DampEnvHealthByTypeInRange(distance, step float64, typeFilters []MaterialType, flagFilters []MaterialFlag) (tilesAffectedCnt int)
	// SearchTilesInRange performs a full in a circle area search.
	// {isEmpty} if not nil, defines empty/non-empty Tile criteria.
	// {dirsFilter, dirsIn} filter includes particles IN direction relative to the source ({dirsIn} = false, inverts the filter making it NOT).
	// {typeFilters, typeIn} filter includes candidates MATCHING types ({typeIn} = false, inverts the filter making it NOT).
	// {flagFilters, flagIn} filter includes candidates WITH flags ({flagIn} = false, inverts the filter making it NOT).
	SearchTilesInRange(
		isEmpty *bool, maxDistance *float64,
		dirsFilter []pkg.Direction, dirsIn bool,
		typeFilters []MaterialType, typeIn bool,
		flagFilters []MaterialFlag, flagIn bool,
	) (tiles []*Tile, sourceToTileDirs []pkg.Direction, sourceToTileDistances []float64)
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
