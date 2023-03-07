package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) DampEnvHealthByTypeInRange(distance, step float64, typeFilters []types.MaterialType, flagFilters []types.MaterialFlag) int {
	tileCandidates, _ := e.getTilesInRange(
		pkg.ValuePtr(false),
		pkg.ValuePtr(distance),
		typeFilters, true,
		flagFilters, true,
	)

	for _, tile := range tileCandidates {
		e.actions = append(e.actions, types.NewReduceHealth(tile.Pos, tile.Particle.ID(), step))
	}

	return len(tileCandidates)
}
