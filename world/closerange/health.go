package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) ReduceHealth(step float64) bool {
	e.sourceHealth -= step
	e.actions = append(e.actions, types.NewReduceHealth(e.source.Pos, step))
	return true
}

func (e *Environment) RemoveHealthSelfReductions() bool {
	return e.removeHealthReductions(e.source.Pos)
}

func (e *Environment) ReduceEnvHealthByFlag(step float64, flagFilters ...types.MaterialFlag) int {
	actionsCreated := 0
	for _, dir := range pkg.AllDirections {
		neighbourTile := e.getNonEmptyNeighbourWithAndFlags(dir, flagFilters...)
		if neighbourTile == nil {
			continue
		}

		actionsCreated++
		e.actions = append(e.actions, types.NewReduceHealth(neighbourTile.Pos, step))
	}

	return actionsCreated
}

func (e *Environment) ReduceEnvHealthByType(step float64, typeFilters ...types.MaterialType) int {
	actionsCreated := 0
	for _, dir := range pkg.AllDirections {
		neighbourTile := e.getNonEmptyNeighbourWithAndTypes(dir, typeFilters...)
		if neighbourTile == nil {
			continue
		}

		actionsCreated++
		e.actions = append(e.actions, types.NewReduceHealth(neighbourTile.Pos, step))
	}

	return actionsCreated
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
