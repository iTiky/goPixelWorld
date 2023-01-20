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
}

type ActionBase struct {
	TilePos Position
}

func (b ActionBase) GetTilePos() Position {
	return b.TilePos
}

type MoveTile struct {
	ActionBase
	NewTilePos Position
}

func NewMoveTile(tilePos, newTilePos Position) MoveTile {
	return MoveTile{
		ActionBase: ActionBase{
			TilePos: tilePos,
		},
		NewTilePos: newTilePos,
	}
}

func (a MoveTile) Type() ActionType {
	return ActionTypeMoveTile
}

type SwapTiles struct {
	ActionBase
	SwapTilePos Position
}

func NewSwapTiles(tile1Pos, tile2Pos Position) SwapTiles {
	return SwapTiles{
		ActionBase: ActionBase{
			TilePos: tile1Pos,
		},
		SwapTilePos: tile2Pos,
	}
}

func (a SwapTiles) Type() ActionType {
	return ActionTypeSwapTiles
}

type MultiplyForce struct {
	ActionBase
	K float64
}

func NewMultiplyForce(tilePos Position, k float64) MultiplyForce {
	return MultiplyForce{
		ActionBase: ActionBase{
			TilePos: tilePos,
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

func NewReflectForce(tilePos Position, vertical, horizontal bool) ReflectForce {
	return ReflectForce{
		ActionBase: ActionBase{
			TilePos: tilePos,
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

func NewAlterForce(tilePos Position, newForceVec pkg.Vector) AlterForce {
	return AlterForce{
		ActionBase: ActionBase{
			TilePos: tilePos,
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

func NewAddForce(tilePos Position, forceVec pkg.Vector) AddForce {
	return AddForce{
		ActionBase: ActionBase{
			TilePos: tilePos,
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

func NewRotateForce(tilePos Position, angle float64) RotateForce {
	return RotateForce{
		ActionBase: ActionBase{
			TilePos: tilePos,
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

func NewReduceHealth(tilePos Position, healthDelta float64) ReduceHealth {
	return ReduceHealth{
		ActionBase: ActionBase{
			TilePos: tilePos,
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

func NewTileReplace(tilePos Position, material Material) TileReplace {
	return TileReplace{
		ActionBase: ActionBase{
			TilePos: tilePos,
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

func NewTileAdd(tilePos Position, material Material) TileAdd {
	return TileAdd{
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

func NewUpdateStateParam(tilePos Position, paramKey string, paramValue int) UpdateStateParam {
	return UpdateStateParam{
		ActionBase: ActionBase{
			TilePos: tilePos,
		},
		ParamKey:   paramKey,
		ParamValue: paramValue,
	}
}

func (a UpdateStateParam) Type() ActionType {
	return ActionTypeUpdateStateParam
}
