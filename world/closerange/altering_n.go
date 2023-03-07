package closerange

import (
	"math/rand"
	"sort"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) AddNewNeighbourTile(newMaterial types.Material, dirFilters []pkg.Direction) bool {
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

	tileCandidates, _ := e.getNeighbours(
		pkg.ValuePtr(true),
		dirs, true,
		nil, false,
		nil, false,
	)
	if len(tileCandidates) == 0 {
		return false
	}
	e.actions = append(e.actions, types.NewTileAdd(tileCandidates[0].Pos, newMaterial))

	return true
}

func (e *Environment) ReplaceNeighbourTile(newMaterial types.Material, flagFilters []types.MaterialFlag) bool {
	tileCandidates, _ := e.getNeighbours(
		pkg.ValuePtr(false),
		pkg.AllDirections, true,
		nil, false,
		flagFilters, true,
	)
	if len(tileCandidates) == 0 {
		return false
	}

	replacementTile := tileCandidates[rand.Intn(len(tileCandidates))]
	e.removeHealthReductions(replacementTile.Pos)
	e.actions = append(e.actions, types.NewTileReplace(replacementTile.Pos, replacementTile.Particle.ID(), newMaterial))

	return true
}

func (e *Environment) AddNewNeighbourTileGrassStyle(newMaterial types.Material) bool {
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
