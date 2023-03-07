package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) DampNeighboursHealthByFlag(step float64, typeFilters []types.MaterialType, flagFilters []types.MaterialFlag) int {
	tileCandidates, _ := e.getNeighbours(
		pkg.ValuePtr(false),
		pkg.AllDirections, true,
		typeFilters, true,
		flagFilters, true,
	)
	if len(tileCandidates) == 0 {
		return 0
	}

	for _, tile := range tileCandidates {
		e.actions = append(e.actions, types.NewReduceHealth(tile.Pos, tile.Particle.ID(), step))
	}

	return len(tileCandidates)
}
