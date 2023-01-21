package closerange

import (
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) MoveGas() bool {
	possibleMoveDirs := make([]pkg.Direction, 0, 5)
	for _, dir := range pkg.DirectionTop.Sector(1) {
		if neighbourTile := e.getEmptyNeighbour(dir); neighbourTile != nil {
			possibleMoveDirs = append(possibleMoveDirs, dir)
		}
	}

	if len(possibleMoveDirs) == 0 {
		return false
	}

	moveDirAfter := possibleMoveDirs[rand.Intn(len(possibleMoveDirs))]
	e.actions = append(e.actions, types.NewRotateForce(e.source.Pos, e.source.Particle.ID(), moveDirAfter.Angle()))

	return true
}
