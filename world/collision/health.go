package collision

import (
	"github.com/itiky/goPixelWorld/world/types"
)

// DampSourceHealth alters the source Particle health state filtered.
func (e *Environment) DampSourceHealth(step float64, flagFilters ...types.MaterialFlag) bool {
	if len(flagFilters) == 0 || e.source.Particle.Material().IsFlagged(flagFilters...) {
		e.actions = append(e.actions, types.NewReduceHealth(e.source.Pos, e.source.Particle.ID(), step))
		return true
	}

	return false
}

// DampSelfHealth alters the target Particle health state.
func (e *Environment) DampSelfHealth(step float64) bool {
	e.actions = append(e.actions, types.NewReduceHealth(e.target.Pos, e.target.Particle.ID(), step))

	return true
}

// DampSelfHealthByMassRate alters the target Particle health state using the source/target mass ratio.
func (e *Environment) DampSelfHealthByMassRate(step float64) bool {
	rate := e.source.Particle.Material().Mass() / e.target.Particle.Material().Mass()
	e.actions = append(e.actions, types.NewReduceHealth(e.target.Pos, e.target.Particle.ID(), step*rate))

	return true
}
