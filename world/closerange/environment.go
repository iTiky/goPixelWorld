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

// ForceVec returns the source Particle force Vector.
func (e *Environment) ForceVec() pkg.Vector {
	return e.source.Particle.ForceVector()
}

// StateParam returns the source Particle internal state param.
func (e *Environment) StateParam(key string) int {
	return e.source.Particle.GetStateParam(key)
}

// Actions returns the env output Actions.
func (e *Environment) Actions() []types.Action {
	return e.actions
}

// getNeighbours returns neighbour tiles matching criteria with corresponding direction (relative to the source).
func (e *Environment) getNeighbours(isEmpty *bool, dirs []pkg.Direction, dirsIn bool, mTypes []types.MaterialType, mTypesIn bool, mFlags []types.MaterialFlag, mFlagsIn bool) ([]*types.Tile, []pkg.Direction) {
	var tiles []*types.Tile
	var tilesDir []pkg.Direction
	for neighbourDir, neighbourTile := range e.neighbours {
		if neighbourTile == nil {
			continue
		}
		if isEmpty != nil {
			if neighbourTile.HasParticle() == *isEmpty {
				continue
			}
		}
		if len(dirs) > 0 {
			if pkg.SliceHasValue(dirs, neighbourDir) != dirsIn {
				continue
			}
		}
		if len(mTypes) > 0 {
			if pkg.SliceHasValue(mTypes, neighbourTile.Particle.Material().Type()) != mTypesIn {
				continue
			}
		}
		if len(mFlags) > 0 {
			if neighbourTile.Particle.Material().IsFlagged(mFlags...) != mFlagsIn {
				continue
			}
		}

		tiles = append(tiles, neighbourTile)
		tilesDir = append(tilesDir, neighbourDir)
	}

	return tiles, tilesDir
}

// getTilesInRange returns in a circle range tiles matching criteria with corresponding distances (relative to the source).
func (e *Environment) getTilesInRange(isEmpty *bool, distanceMax *float64, mTypes []types.MaterialType, mTypesIn bool, mFlags []types.MaterialFlag, mFlagsIn bool) ([]*types.Tile, []float64) {
	var tiles []*types.Tile
	var tilesDistance []float64
	for _, rangeTile := range e.tilesInRange {
		if rangeTile == nil {
			continue
		}
		if isEmpty != nil {
			if rangeTile.HasParticle() == *isEmpty {
				continue
			}
		}
		if len(mTypes) > 0 {
			if pkg.SliceHasValue(mTypes, rangeTile.Particle.Material().Type()) != mTypesIn {
				continue
			}
		}
		if len(mFlags) > 0 {
			if rangeTile.Particle.Material().IsFlagged(mFlags...) != mFlagsIn {
				continue
			}
		}

		rangeTileDistance := e.source.Pos.DistanceTo(rangeTile.Pos)
		if distanceMax != nil {
			if rangeTileDistance > *distanceMax {
				continue
			}
		}

		tiles = append(tiles, rangeTile)
		tilesDistance = append(tilesDistance, rangeTileDistance)
	}

	return tiles, tilesDistance
}
