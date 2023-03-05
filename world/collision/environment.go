package collision

import (
	"strings"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.CollisionEnvironment = &Environment{}

// Environment defines the state for a Particle-Particle collision processing.
type Environment struct {
	direction  pkg.Direction                 // collision direction (relative to the target, from which side the source is flying from)
	source     *types.Tile                   // source Particle (the one who wants to move to the target's Position)
	target     *types.Tile                   // target Particle
	neighbours map[pkg.Direction]*types.Tile // neighbour Particles by a relative to source -> target direction
	//
	actions []types.Action // output Actions
}

// NewEnvironment creates a new empty Environment.
func NewEnvironment(direction pkg.Direction, sourceTile, targetTile *types.Tile) *Environment {
	e := Environment{
		direction:  direction,
		source:     sourceTile,
		target:     targetTile,
		neighbours: make(map[pkg.Direction]*types.Tile, 5),
	}

	return &e
}

// Reset the Environment state.
// Used to reduce the amount of allocations.
func (e *Environment) Reset(direction pkg.Direction, sourceTile, targetTile *types.Tile) {
	e.direction = direction
	e.source = sourceTile
	e.target = targetTile
	e.actions = e.actions[:0]

	for dir := range e.neighbours {
		e.neighbours[dir] = nil
	}
}

// TargetMaterial returns the target Particle Material.
func (e *Environment) TargetMaterial() types.Material {
	return e.target.Particle.Material()
}

// SetNeighbour sets the target neighbor by direction.
func (e *Environment) SetNeighbour(dir pkg.Direction, tile *types.Tile) {
	if tile == nil {
		return
	}
	e.neighbours[dir] = tile
}

// IsFlagged checks if the source Particle has MaterialFlag.
func (e *Environment) IsFlagged(flag types.MaterialFlag) bool {
	return e.source.Particle.Material().IsFlagged(flag)
}

// IsType checks if the source Particle Material is MaterialType.
func (e *Environment) IsType(mType types.MaterialType) bool {
	return e.source.Particle.Material().Type() == mType
}

// Actions returns the current output Actions.
func (e *Environment) Actions() []types.Action {
	return e.actions
}

func (e *Environment) String() string {
	str := strings.Builder{}

	str.WriteString("CollisionEnvironment{\n")
	str.WriteString("  direction: " + e.direction.String() + "\n")
	str.WriteString("  source: " + e.source.String() + "\n")
	str.WriteString("  target: " + e.target.String() + "\n")
	for dir, tile := range e.neighbours {
		str.WriteString("  neighbour(" + dir.String() + ": " + tile.String() + "\n")
	}
	str.WriteString("}")

	return str.String()
}

// getEmptyNeighbour returns an empty target neighbour by direction.
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
