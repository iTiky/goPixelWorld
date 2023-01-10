package collision

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) MoveSandSource() bool {
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionLeft); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionRight); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionTop); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, neighbourTile.Pos))
		return true
	}

	return false
}

func (e *Environment) MoveLiquidSource() bool {
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionLeft); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionRight); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, neighbourTile.Pos))
		return true
	}

	var spreadDir1, spreadDir2 pkg.Direction
	//if pkg.FlipCoin() {
	spreadDir1, spreadDir2 = pkg.DirectionTopLeft, pkg.DirectionTopRight
	//} else {
	//	spreadDir1, spreadDir2 = pkg.DirectionTopRight, pkg.DirectionTopLeft
	//}
	if neighbourTile := e.getEmptyNeighbour(spreadDir1); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(spreadDir2); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionTop); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, neighbourTile.Pos))
		return true
	}

	return false
}

func (e *Environment) SwapSourceTarget() bool {
	e.actions = append(e.actions, types.NewSwapTiles(e.source.Pos, e.target.Pos))
	return true
}
