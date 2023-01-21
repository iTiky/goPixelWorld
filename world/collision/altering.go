package collision

import (
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) DampSourceForce(k float64) bool {
	e.actions = append(e.actions, types.NewMultiplyForce(e.source.Pos, e.source.Particle.ID(), k))

	return true
}
