package collision

import (
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) DampSourceHealth(step float64, flagFilters ...types.MaterialFlag) bool {
	if len(flagFilters) == 0 || e.source.Particle.Material().IsFlagged(flagFilters...) {
		e.actions = append(e.actions, types.NewReduceHealth(e.source.Pos, step))
		return true
	}

	return false
}

func (e *Environment) DampSelfHealth(step float64) bool {
	e.actions = append(e.actions, types.NewReduceHealth(e.target.Pos, step))

	return true
}

func (e *Environment) DampSelfHealthByMassRate(step float64) bool {
	rate := e.source.Particle.Material().Mass() / e.target.Particle.Material().Mass()
	e.actions = append(e.actions, types.NewReduceHealth(e.target.Pos, step*rate))

	return true
}
