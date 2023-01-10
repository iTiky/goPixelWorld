package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

type Environment struct {
	source     *types.Tile
	neighbours map[pkg.Direction]*types.Tile
	//
	actions []types.Action
}

func NewEnvironment(sourceTile *types.Tile) *Environment {
	return &Environment{
		source:     sourceTile,
		neighbours: make(map[pkg.Direction]*types.Tile, 8),
	}
}

func (e *Environment) SetNeighbour(dir pkg.Direction, tile *types.Tile) {
	if tile == nil {
		return
	}
	e.neighbours[dir] = tile
}

func (e *Environment) Health() float64 {
	return e.source.Particle.Health()
}

func (e *Environment) Position() types.Position {
	return e.source.Pos
}

func (e *Environment) Actions() []types.Action {
	return e.actions
}

func (e *Environment) getEmptyNeighbour(dir pkg.Direction) *types.Tile {
	neighbourTile := e.neighbours[dir]
	if neighbourTile == nil {
		return nil
	}

	if neighbourTile.HasParticle() {
		return nil
	}

	return neighbourTile
}

func (e *Environment) getNonEmptyNeighbour(dir pkg.Direction) *types.Tile {
	neighbourTile := e.neighbours[dir]
	if neighbourTile == nil {
		return nil
	}

	if !neighbourTile.HasParticle() {
		return nil
	}

	return neighbourTile
}

func (e *Environment) getNonEmptyNeighbourWithAndFlags(dir pkg.Direction, flagFilters ...types.MaterialFlag) *types.Tile {
	neighbourTile := e.getNonEmptyNeighbour(dir)
	if neighbourTile == nil {
		return nil
	}

	for _, flagFiler := range flagFilters {
		if !neighbourTile.Particle.Material().IsFlagged(flagFiler) {
			return nil
		}
	}

	return neighbourTile
}
