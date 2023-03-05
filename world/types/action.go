package types

import (
	"github.com/itiky/goPixelWorld/pkg"
)

// ActionType defines Action type.
type ActionType int

const (
	ActionTypeNone ActionType = iota
	ActionTypeMultiplyForce
	ActionTypeReflectForce
	ActionTypeAlterForce
	ActionTypeMoveTile
	ActionTypeSwapTiles
	ActionTypeAddForce
	ActionTypeRotateForce
	ActionTypeReduceHealth
	ActionTypeTileReplace
	ActionTypeTileAdd
	ActionTypeUpdateStateParam
)

// Action defines the contract for all Action types.
// Each Action is applied to a single Tile (target Tile).
// Since each Action is idempotent, it doesn't include the target Tile object itself, but rather keeps its Position and ParticleID.
// Action apply operation must check if the target Tile (or Particle) is "still there".
type Action interface {
	// Type returns the Action type.
	Type() ActionType
	// GetTilePos returns the target Tile Position.
	GetTilePos() Position
	// GetParticleID returns the target Tile Particle's ID.
	GetParticleID() uint64
}

// ActionBase defines a common Action fields.
type ActionBase struct {
	TilePos    Position
	ParticleID uint64
}

func (b ActionBase) GetTilePos() Position {
	return b.TilePos
}

func (b ActionBase) GetParticleID() uint64 {
	return b.ParticleID
}

// NoopAction defines an empty Action (not used ATM).
type NoopAction struct {
	ActionBase
}

func (a NoopAction) Type() ActionType {
	return ActionTypeNone
}

// MoveTile defines an Action which moves a Tile to a new Position.
type MoveTile struct {
	ActionBase
	NewTilePos Position
}

func NewMoveTile(tilePos Position, tilePID uint64, newTilePos Position) *MoveTile {
	return &MoveTile{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		NewTilePos: newTilePos,
	}
}

func (a MoveTile) Type() ActionType {
	return ActionTypeMoveTile
}

// SwapTiles defines an Action which swaps Particles between two Tiles.
type SwapTiles struct {
	ActionBase
	SwapTilePos    Position
	SwapParticleID uint64
}

func NewSwapTiles(tile1Pos Position, tile1PID uint64, tile2Pos Position, tile2PID uint64) *SwapTiles {
	return &SwapTiles{
		ActionBase: ActionBase{
			TilePos:    tile1Pos,
			ParticleID: tile1PID,
		},
		SwapTilePos:    tile2Pos,
		SwapParticleID: tile2PID,
	}
}

func (a SwapTiles) Type() ActionType {
	return ActionTypeSwapTiles
}

// MultiplyForce defines an Action which modifies a Particle's force Vector magnitude.
type MultiplyForce struct {
	ActionBase
	K float64
}

func NewMultiplyForce(tilePos Position, tilePID uint64, k float64) *MultiplyForce {
	return &MultiplyForce{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		K: k,
	}
}

func (a MultiplyForce) Type() ActionType {
	return ActionTypeMultiplyForce
}

// ReflectForce defines an Action which reflects a Particle's force Vector angle.
type ReflectForce struct {
	ActionBase
	Vertical   bool
	Horizontal bool
}

func NewReflectForce(tilePos Position, tilePID uint64, vertical, horizontal bool) *ReflectForce {
	return &ReflectForce{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		Vertical:   vertical,
		Horizontal: horizontal,
	}
}

func (a ReflectForce) Type() ActionType {
	return ActionTypeReflectForce
}

// AlterForce defines an Action which sets a new Particle's force Vector value.
type AlterForce struct {
	ActionBase
	NewForceVec pkg.Vector
}

func NewAlterForce(tilePos Position, tilePID uint64, newForceVec pkg.Vector) *AlterForce {
	return &AlterForce{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		NewForceVec: newForceVec,
	}
}

func (a AlterForce) Type() ActionType {
	return ActionTypeAlterForce
}

// AddForce defines an Action which adds a new Vector to the Particle's force Vector.
type AddForce struct {
	ActionBase
	ForceVec pkg.Vector
}

func NewAddForce(tilePos Position, tilePID uint64, forceVec pkg.Vector) *AddForce {
	return &AddForce{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		ForceVec: forceVec,
	}
}

func (a AddForce) Type() ActionType {
	return ActionTypeAddForce
}

// RotateForce defines an Action which rotates the Particle's force Vector.
type RotateForce struct {
	ActionBase
	Angle float64
}

func NewRotateForce(tilePos Position, tilePID uint64, angle float64) *RotateForce {
	return &RotateForce{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		Angle: angle,
	}
}

func (a RotateForce) Type() ActionType {
	return ActionTypeRotateForce
}

// ReduceHealth defines an Action which modifies the Particle's health (increase / decrease).
type ReduceHealth struct {
	ActionBase
	HealthDelta float64
}

func NewReduceHealth(tilePos Position, tilePID uint64, healthDelta float64) *ReduceHealth {
	return &ReduceHealth{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		HealthDelta: healthDelta,
	}
}

func (a ReduceHealth) Type() ActionType {
	return ActionTypeReduceHealth
}

// TileReplace defines an Action which replaces the Tile's Particle with a new one.
type TileReplace struct {
	ActionBase
	Material Material
}

func NewTileReplace(tilePos Position, tilePID uint64, material Material) *TileReplace {
	return &TileReplace{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		Material: material,
	}
}

func (a TileReplace) Type() ActionType {
	return ActionTypeTileReplace
}

// TileAdd defines an Action which adds a new Particle.
type TileAdd struct {
	ActionBase
	Material Material
}

func NewTileAdd(tilePos Position, material Material) *TileAdd {
	return &TileAdd{
		ActionBase: ActionBase{
			TilePos: tilePos,
		},
		Material: material,
	}
}

func (a TileAdd) Type() ActionType {
	return ActionTypeTileAdd
}

// UpdateStateParam defines an Action which modifies the Tile internal state.
type UpdateStateParam struct {
	ActionBase
	ParamKey   string
	ParamValue int
}

func NewUpdateStateParam(tilePos Position, tilePID uint64, paramKey string, paramValue int) *UpdateStateParam {
	return &UpdateStateParam{
		ActionBase: ActionBase{
			TilePos:    tilePos,
			ParticleID: tilePID,
		},
		ParamKey:   paramKey,
		ParamValue: paramValue,
	}
}

func (a UpdateStateParam) Type() ActionType {
	return ActionTypeUpdateStateParam
}
