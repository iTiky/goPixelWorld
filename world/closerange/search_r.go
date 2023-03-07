package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) SearchTilesInRange(
	isEmpty *bool, maxDistance *float64,
	dirsFilter []pkg.Direction, dirsIn bool,
	typeFilters []types.MaterialType, typeIn bool,
	flagFilters []types.MaterialFlag, flagIn bool,
) ([]*types.Tile, []pkg.Direction, []float64) {

	tileCandidates, tilesDistance := e.getTilesInRange(
		isEmpty,
		maxDistance,
		typeFilters, typeIn,
		flagFilters, flagIn,
	)
	if len(tileCandidates) == 0 {
		return nil, nil, nil
	}

	outputTiles := make([]*types.Tile, 0, len(tileCandidates))
	outputDirs := make([]pkg.Direction, 0, len(tileCandidates))
	outputDistances := make([]float64, 0, len(tilesDistance))
	for tileIdx, tile := range tileCandidates {
		tileDir := pkg.NewDirectionFromCoords(e.source.Pos.X, e.source.Pos.Y, tile.Pos.X, tile.Pos.Y)
		if len(dirsFilter) > 0 {
			if pkg.SliceHasValue(dirsFilter, tileDir) != dirsIn {
				continue
			}
		}

		outputTiles = append(outputTiles, tile)
		outputDirs = append(outputDirs, tileDir)
		outputDistances = append(outputDistances, tilesDistance[tileIdx])
	}

	return outputTiles, outputDirs, outputDistances
}
