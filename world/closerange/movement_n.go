package closerange

import (
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) MoveTileWithNeighboursGasStyle() bool {
	tileCandidates, tilesDir := e.getNeighbours(
		pkg.ValuePtr(true),
		pkg.DirectionTop.Sector(1), true,
		nil, false,
		nil, false,
	)
	if len(tileCandidates) == 0 {
		return false
	}

	moveDir := tilesDir[rand.Intn(len(tilesDir))]
	e.actions = append(e.actions, types.NewRotateForce(e.source.Pos, e.source.Particle.ID(), moveDir.Angle()))

	return true
}
