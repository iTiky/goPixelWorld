package collision

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

// MoveSandSource moves the source Particle for sand-like Materials.
// Criteria:
//   - if target's left neighbour is empty, pick it;
//   - if target's right neighbour is empty, pick it;
//   - if target's top neighbour is empty, pick it (stack);
func (e *Environment) MoveSandSource() bool {
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionLeft); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, e.source.Particle.ID(), neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionRight); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, e.source.Particle.ID(), neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionTop); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, e.source.Particle.ID(), neighbourTile.Pos))
		return true
	}

	return false
}

// MoveLiquidSource moves the source Particle for liquid-like Materials.
// Criteria:
//   - if target's left neighbour is empty, pick it;
//   - if target's right neighbour is empty, pick it;
//   - if target's top-left neighbour is empty, pick it (spread);
//   - if target's top-right neighbour is empty, pick it (spread);
//   - if target's top neighbour is empty, pick it (stack);
func (e *Environment) MoveLiquidSource() bool {
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionLeft); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, e.source.Particle.ID(), neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionRight); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, e.source.Particle.ID(), neighbourTile.Pos))
		return true
	}

	var spreadDir1, spreadDir2 pkg.Direction
	spreadDir1, spreadDir2 = pkg.DirectionTopLeft, pkg.DirectionTopRight
	if neighbourTile := e.getEmptyNeighbour(spreadDir1); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, e.source.Particle.ID(), neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(spreadDir2); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, e.source.Particle.ID(), neighbourTile.Pos))
		return true
	}
	if neighbourTile := e.getEmptyNeighbour(pkg.DirectionTop); neighbourTile != nil {
		e.actions = append(e.actions, types.NewMoveTile(e.source.Pos, e.source.Particle.ID(), neighbourTile.Pos))
		return true
	}

	return false
}

// SwapSourceTarget swaps the source and target Particles.
func (e *Environment) SwapSourceTarget() bool {
	e.actions = append(e.actions, types.NewSwapTiles(e.source.Pos, e.source.Particle.ID(), e.target.Pos, e.target.Particle.ID()))
	return true
}
