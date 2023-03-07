package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) SearchNeighbours(
	isEmpty *bool,
	dirsFilter []pkg.Direction, dirsIn bool,
	typeFilters []types.MaterialType, typeIn bool,
	flagFilters []types.MaterialFlag, flagIn bool,
) ([]*types.Tile, []pkg.Direction) {

	return e.getNeighbours(
		isEmpty,
		dirsFilter, dirsIn,
		typeFilters, typeIn,
		flagFilters, flagIn,
	)
}
