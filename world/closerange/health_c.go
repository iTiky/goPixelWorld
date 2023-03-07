package closerange

import (
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) DampSelfHealth(step float64) bool {
	e.sourceHealth -= step
	e.actions = append(e.actions, types.NewReduceHealth(e.source.Pos, e.source.Particle.ID(), step))
	return true
}

func (e *Environment) RemoveSelfHealthDamps() bool {
	return e.removeHealthReductions(e.source.Pos)
}

func (e *Environment) removeHealthReductions(targetPos types.Position) bool {
	removed := false

	n := 0
	for _, a := range e.actions {
		if a.Type() == types.ActionTypeReduceHealth && a.GetTilePos().Equal(targetPos) {
			removed = true
			continue
		}

		e.actions[n] = a
		n++
	}
	e.actions = e.actions[:n]

	return removed
}
