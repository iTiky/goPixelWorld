package world

import (
	"github.com/itiky/goPixelWorld/world/types"
)

// processActions iterates over each Tile worker output queue and applies pushed Actions modifying the Map state.
// Each Action "doesn't know" about its predecessors, so the apply operation must be idempotent.
func (m *Map) processActions() {
	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.processActions")()
	}

	getExistingTile := func(pos types.Position, pid uint64) *types.Tile {
		tile := m.getTile(pos.X, pos.Y)
		if !tile.HasParticle() || tile.Particle.ID() != pid {
			return nil
		}
		return tile
	}

	getEmptyTile := func(pos types.Position) *types.Tile {
		tile := m.getTile(pos.X, pos.Y)
		if tile.HasParticle() {
			return nil
		}
		return tile
	}

	for _, workerActions := range m.procActions {
		for _, aBz := range workerActions {
			switch a := aBz.(type) {
			case *types.MultiplyForce:
				tile := getExistingTile(a.TilePos, a.ParticleID)
				if tile == nil {
					break
				}
				tile.Particle.MultiplyForce(a.K)
			case *types.ReflectForce:
				tile := getExistingTile(a.TilePos, a.ParticleID)
				if tile == nil {
					break
				}
				tile.Particle.ReflectForce(a.Horizontal, a.Vertical)
			case *types.AddForce:
				tile := getExistingTile(a.TilePos, a.ParticleID)
				if tile == nil {
					break
				}
				tile.Particle.AddForce(a.ForceVec)
			case *types.AlterForce:
				tile := getExistingTile(a.TilePos, a.ParticleID)
				if tile == nil {
					break
				}
				tile.Particle.SetForce(a.NewForceVec)
			case *types.RotateForce:
				tile := getExistingTile(a.TilePos, a.ParticleID)
				if tile == nil {
					break
				}
				tile.Particle.RotateForce(a.Angle)
			case *types.MoveTile:
				tile1 := getExistingTile(a.TilePos, a.ParticleID)
				if tile1 == nil {
					break
				}
				tile2 := getEmptyTile(a.NewTilePos)
				if tile2 == nil {
					break
				}
				m.moveTile(tile1, a.NewTilePos)
			case *types.SwapTiles:
				tile1 := getExistingTile(a.TilePos, a.ParticleID)
				if tile1 == nil {
					break
				}
				tile2 := getExistingTile(a.SwapTilePos, a.SwapParticleID)
				if tile2 == nil {
					break
				}
				m.swapTiles(tile1, tile2)
			case *types.ReduceHealth:
				tile := getExistingTile(a.TilePos, a.ParticleID)
				if tile == nil {
					break
				}
				tile.Particle.ReduceHealth(a.HealthDelta)
				if tile.Particle.IsDestroyed() {
					m.removeParticle(tile)
				}
			case *types.TileReplace:
				tile := getExistingTile(a.TilePos, a.ParticleID)
				if tile == nil {
					break
				}
				m.removeParticle(tile)
				m.createParticle(tile, a.Material)
			case *types.UpdateStateParam:
				tile := getExistingTile(a.TilePos, a.ParticleID)
				if tile == nil {
					break
				}
				tile.Particle.SetStateParam(a.ParamKey, a.ParamValue)
			case *types.TileAdd:
				tile := getEmptyTile(a.TilePos)
				if tile == nil {
					break
				}
				m.createParticle(tile, a.Material)
			}
		}
	}
}
