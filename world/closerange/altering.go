package closerange

import (
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) ReplaceTile(newMaterial types.Material, flagFilters ...types.MaterialFlag) bool {
	var tileCandidates []*types.Tile
	for _, dir := range pkg.AllDirections {
		tile := e.getNonEmptyNeighbourWithAndFlags(dir, flagFilters...)
		if tile == nil {
			continue
		}
		tileCandidates = append(tileCandidates, tile)
	}
	if len(tileCandidates) == 0 {
		return false
	}

	replacementTile := tileCandidates[rand.Intn(len(tileCandidates))]
	e.actions = append(e.actions, types.NewTileReplace(replacementTile.Pos, newMaterial))

	return true
}

func (e *Environment) AddTile(newMaterial types.Material) bool {
	for _, dir := range pkg.AllDirections {
		tile := e.getEmptyNeighbour(dir)
		if tile == nil {
			continue
		}
		e.actions = append(e.actions, types.NewTileAdd(tile.Pos, newMaterial))
		return true
	}

	return false
}
