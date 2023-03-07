package closerange

import (
	"math"
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) AddForceInRange(mag float64, notFlagFilters []types.MaterialFlag) bool {
	magAbs := math.Abs(mag)

	rotateForceVec := false
	if mag < 0.0 {
		rotateForceVec = true
	}

	tileCandidates, _ := e.getTilesInRange(
		pkg.ValuePtr(false),
		nil,
		nil, false,
		notFlagFilters, false,
	)
	if len(tileCandidates) == 0 {
		return false
	}

	for _, tile := range tileCandidates {
		forceVec := pkg.NewVectorByCoordinates(magAbs, float64(tile.Pos.X), float64(tile.Pos.Y), float64(e.source.Pos.X), float64(e.source.Pos.Y))
		if rotateForceVec {
			forceVec = forceVec.Rotate(math.Pi)
		}
		e.actions = append(e.actions, types.NewAddForce(tile.Pos, tile.Particle.ID(), forceVec))
	}

	return true
}

func (e *Environment) AddNewTileInRange(newMaterial types.Material) bool {
	tileCandidates, _ := e.getTilesInRange(
		pkg.ValuePtr(true), nil,
		nil, false,
		nil, false,
	)
	if len(tileCandidates) == 0 {
		return false
	}

	tileCandidate := tileCandidates[rand.Intn(len(tileCandidates))]
	e.actions = append(e.actions, types.NewTileAdd(tileCandidate.Pos, newMaterial))

	return true
}
