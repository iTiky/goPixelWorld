package types

import (
	"github.com/itiky/goPixelWorld/pkg"
)

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

type Action interface {
	Type() ActionType
	GetTilePos() Position
	GetParticleID() uint64
}

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

type NoopAction struct {
	ActionBase
}

func (a NoopAction) Type() ActionType {
	return ActionTypeNone
}

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
