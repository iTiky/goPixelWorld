package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

// Environment defines the state for a Particle self-processing (altering itself and surrounding neighbours based on Material logic).
type Environment struct {
	source       *types.Tile                   // source Tile
	sourceHealth float64                       // source health (the current value since output Actions can rely on it)
	neighbours   map[pkg.Direction]*types.Tile // neighbour tiles by a relative to source direction
	//
	tilesInRange []*types.Tile // tiles in a circle range
	//
	actions []types.Action // processing output
}

// NewEnvironment creates a new empty Environment.
func NewEnvironment(sourceTile *types.Tile) *Environment {
	env := &Environment{
		source:       sourceTile,
		sourceHealth: 0.0,
		neighbours:   make(map[pkg.Direction]*types.Tile, 8),
		tilesInRange: make([]*types.Tile, 0, 8),
	}
	if sourceTile != nil {
		env.sourceHealth = sourceTile.Particle.Health()
	}

	return env
}

// Reset the Environment state.
// Used to reduce the amount of allocations.
func (e *Environment) Reset(sourceTile *types.Tile) {
	e.source = sourceTile
	e.sourceHealth = sourceTile.Particle.Health()
	e.actions = e.actions[:0]
	e.tilesInRange = e.tilesInRange[:0]

	for dir := range e.neighbours {
		e.neighbours[dir] = nil
	}
}

// SetNeighbour sets the neighbour by relative direction.
func (e *Environment) SetNeighbour(dir pkg.Direction, tile *types.Tile) {
	if tile == nil {
		return
	}
	e.neighbours[dir] = tile
}

// AddTileInRange adds a neighbour in a circle range.
func (e *Environment) AddTileInRange(tile *types.Tile) {
	e.tilesInRange = append(e.tilesInRange, tile)
}

// Health returns the current source Particle health.
func (e *Environment) Health() float64 {
	return e.sourceHealth
}

// Position returns the source Particle position.
func (e *Environment) Position() types.Position {
	return e.source.Pos
}

// StateParam returns the source Particle internal state param.
func (e *Environment) StateParam(key string) int {
	return e.source.Particle.GetStateParam(key)
}

// Actions returns the env output Actions.
func (e *Environment) Actions() []types.Action {
	return e.actions
}

// getEmptyNeighbour returns an empty neighbour Tile by direction.
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

// getNonEmptyNeighbour returns a non-empty neighbour Tile by direction.
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

// getNonEmptyNeighbourWithAndFlags returns a non-empty neighbour Tile by direction matching MaterialFlag.
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

// getNonEmptyNeighbourWithAndTypes returns a non-empty neighbour Tile by direction matching MaterialType.
func (e *Environment) getNonEmptyNeighbourWithAndTypes(dir pkg.Direction, typeFilters ...types.MaterialType) *types.Tile {
	neighbourTile := e.getNonEmptyNeighbour(dir)
	if neighbourTile == nil {
		return nil
	}

	for _, typeFilter := range typeFilters {
		if neighbourTile.Particle.Material().Type() != typeFilter {
			return nil
		}
	}

	return neighbourTile
}
