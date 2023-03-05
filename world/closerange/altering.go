package closerange

import (
	"math"
	"math/rand"
	"sort"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

// ReplaceTile replaces a neighbour Tile with a new Particle.
// Candidate is selected randomly from non-empty neighbours matching the filter.
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
	e.removeHealthReductions(replacementTile.Pos)
	e.actions = append(e.actions, types.NewTileReplace(replacementTile.Pos, replacementTile.Particle.ID(), newMaterial))

	return true
}

// AddTile adds a new neighbour Particle.
// Candidate is selected randomly from empty neighbours matching the filter.
func (e *Environment) AddTile(newMaterial types.Material, dirFilters ...pkg.Direction) bool {
	var dirs []pkg.Direction
	if len(dirFilters) > 0 {
		dirs = dirFilters
	} else {
		dirs = make([]pkg.Direction, len(pkg.AllDirections))
		copy(dirs, pkg.AllDirections)
	}
	rand.Shuffle(len(dirs), func(i, j int) {
		dirs[i], dirs[j] = dirs[j], dirs[i]
	})

	for _, dir := range dirs {
		tile := e.getEmptyNeighbour(dir)
		if tile == nil {
			continue
		}
		e.actions = append(e.actions, types.NewTileAdd(tile.Pos, newMaterial))
		return true
	}

	return false
}

// ReplaceSelf replaces the source Particle with a new one.
func (e *Environment) ReplaceSelf(newMaterial types.Material) bool {
	e.actions = append(e.actions, types.NewTileReplace(e.source.Pos, e.source.Particle.ID(), newMaterial))
	return true
}

// AddTileGrassStyle adds a new grass-like neighbour.
// Candidate select criteria:
//   - three close empty Tiles (for ex.: Top-Left, Top and Top-Right);
//   - random;
func (e *Environment) AddTileGrassStyle(newMaterial types.Material) bool {
	var dirs []pkg.Direction
	for dir, tile := range e.neighbours {
		if tile.HasParticle() {
			continue
		}
		dirs = append(dirs, dir)
	}
	if len(dirs) < 3 {
		return false
	}

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i] < dirs[j]
	})

	idx := rand.Intn(len(dirs))
	nextIdx := func() int {
		i := idx + 1
		if i >= len(dirs) {
			return 0
		}
		return i
	}
	prevIdx := func() int {
		i := idx - 1
		if i < 0 {
			return len(dirs) - 1
		}
		return i
	}
	for i := 0; i < len(dirs); i++ {
		idxPrev, idxNext := prevIdx(), nextIdx()
		if dirs[idxPrev].Next() != dirs[idx] || dirs[idx].Next() != dirs[idxNext] {
			idx = idxNext
			continue

		}

		tile := e.neighbours[dirs[idx]]
		e.actions = append(e.actions, types.NewTileAdd(tile.Pos, newMaterial))
		return true
	}

	return false
}

// UpdateStateParam updates the internal source Particle state param.
func (e *Environment) UpdateStateParam(paramKey string, paramValue int) bool {
	if e.source.Particle.GetStateParam(paramKey) == paramValue {
		return false
	}

	e.actions = append(e.actions, types.NewUpdateStateParam(e.source.Pos, e.source.Particle.ID(), paramKey, paramValue))
	return true
}

// AddForceInRange adds a force Vector in a circle range.
// If {mag} is LT 0, force is reflected.
func (e *Environment) AddForceInRange(mag float64, notFlagFilters ...types.MaterialFlag) bool {
	magAbs := math.Abs(mag)

	rotateForceVec := false
	if mag < 0.0 {
		rotateForceVec = true
	}

	added := false
	for _, tile := range e.tilesInRange {
		if tile.Particle.Material().IsFlagged(notFlagFilters...) {
			continue
		}

		added = true
		forceVec := pkg.NewVectorByCoordinates(magAbs, float64(tile.Pos.X), float64(tile.Pos.Y), float64(e.source.Pos.X), float64(e.source.Pos.Y))
		if rotateForceVec {
			forceVec = forceVec.Rotate(math.Pi)
		}
		e.actions = append(e.actions, types.NewAddForce(tile.Pos, tile.Particle.ID(), forceVec))
	}

	return added
}
