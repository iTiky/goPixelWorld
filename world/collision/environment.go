package collision

import (
	"strings"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.CollisionEnvironment = &Environment{}

type Environment struct {
	direction  pkg.Direction
	source     *types.Tile
	target     *types.Tile
	neighbours map[pkg.Direction]*types.Tile
	//
	actions []types.Action
}

func NewEnvironment(direction pkg.Direction, sourceTile, targetTile *types.Tile) *Environment {
	e := Environment{
		direction:  direction,
		source:     sourceTile,
		target:     targetTile,
		neighbours: make(map[pkg.Direction]*types.Tile, 5),
	}

	return &e
}

func (e *Environment) TargetMaterial() types.Material {
	return e.target.Particle.Material()
}

func (e *Environment) SetNeighbour(dir pkg.Direction, tile *types.Tile) {
	if tile == nil {
		return
	}
	e.neighbours[dir] = tile
}

func (e *Environment) IsFlagged(flag types.MaterialFlag) bool {
	return e.source.Particle.Material().IsFlagged(flag)
}

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
