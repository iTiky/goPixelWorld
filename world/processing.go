package world

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/collision"
	"github.com/itiky/goPixelWorld/world/types"
)

func (m *Map) processTile(tile *types.Tile) {
	if tile.Particle.Material().Type() == types.MaterialTypeBorder {
		return
	}

	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.processTile")()
	}

	tile.Particle.UpdateState()

	tileEnv := m.buildTileEnv(tile)
	tile.Particle.Material().ProcessInternal(tileEnv)
	m.applyActions(tileEnv.Actions()...)

	if !tile.HasParticle() {
		return
	}

	targetEmptyTile, collisionEnv := m.buildCollisionEnv(tile)
	if targetEmptyTile != nil {
		m.applyActions(types.NewMoveTile(tile.Pos, targetEmptyTile.Pos))
	}
	if collisionEnv != nil {
		collisionEnv.TargetMaterial().ProcessCollision(collisionEnv)
		m.applyActions(collisionEnv.Actions()...)
	}
}

func (m *Map) buildTileEnv(sourceTile *types.Tile) *closerange.Environment {
	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.buildTileEnv")()
	}

	m.tileEnv.Reset(sourceTile)
	switch sourceTile.Particle.Material().CloseRangeType() {
	case types.MaterialCloseRangeTypeNone:
	case types.MaterialCloseRangeTypeSurrounding:
		setNeighbor := func(dir pkg.Direction, dx, dy int) {
			x, y := sourceTile.Pos.X+dx, sourceTile.Pos.Y+dy
			if !m.isPositionValid(x, y) {
				return
			}
			m.tileEnv.SetNeighbour(dir, m.getTile(x, y))
		}

		setNeighbor(pkg.DirectionTop, 0, -1)
		setNeighbor(pkg.DirectionTopRight, 1, -1)
		setNeighbor(pkg.DirectionRight, 1, 0)
		setNeighbor(pkg.DirectionBottomRight, 1, 1)
		setNeighbor(pkg.DirectionBottom, 0, 1)
		setNeighbor(pkg.DirectionBottomLeft, -1, 1)
		setNeighbor(pkg.DirectionLeft, -1, 0)
		setNeighbor(pkg.DirectionTopLeft, -1, -1)
	case types.MaterialCloseRangeTypeInCircleRange:
		r := sourceTile.Particle.Material().CloseRangeCircleRadius()
		for _, pos := range types.PositionsInCircle(sourceTile.Pos.X, sourceTile.Pos.Y, r, false) {
			if !m.isPositionValid(pos.X, pos.Y) {
				continue
			}

			neighborTile := m.getTile(pos.X, pos.Y)
			if !neighborTile.HasParticle() {
				continue
			}

			m.tileEnv.AddTileInRange(neighborTile)
		}
	}

	return m.tileEnv
}

func (m *Map) buildCollisionEnv(sourceTile *types.Tile) (*types.Tile, *collision.Environment) {
	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.buildCollisionEnv")()
	}

	targetTile := sourceTile.TargetTile()
	if targetTile == nil {
		return nil, nil
	}

	pathToTarget := sourceTile.Pos.CreatePathTo(targetTile.Pos, m.width, m.height)
	for _, pathPos := range pathToTarget {
		targetTile = m.getTile(pathPos.X, pathPos.Y)
		if targetTile.HasParticle() {
			break
		}
	}

	if !targetTile.HasParticle() {
		return targetTile, nil
	}

	colDirection := pkg.NewDirectionFromCoords(sourceTile.Pos.X, sourceTile.Pos.Y, targetTile.Pos.X, targetTile.Pos.Y)

	m.collisionEnv.Reset(colDirection, sourceTile, targetTile)
	setNeighbor := func(dir pkg.Direction, dx, dy int) {
		x, y := targetTile.Pos.X+dx, targetTile.Pos.Y+dy
		if !m.isPositionValid(x, y) {
			return
		}
		m.collisionEnv.SetNeighbour(dir, m.getTile(x, y))
	}

	switch colDirection {
	case pkg.DirectionTop:
		setNeighbor(pkg.DirectionTopLeft, 1, 1)
		setNeighbor(pkg.DirectionLeft, 1, 0)
		setNeighbor(pkg.DirectionTop, 0, 1)
		setNeighbor(pkg.DirectionRight, -1, 0)
		setNeighbor(pkg.DirectionTopRight, -1, 1)
	case pkg.DirectionTopRight:
		setNeighbor(pkg.DirectionTopLeft, 0, 1)
		setNeighbor(pkg.DirectionLeft, 1, 1)
		setNeighbor(pkg.DirectionTop, -1, 1)
		setNeighbor(pkg.DirectionRight, -1, -1)
		setNeighbor(pkg.DirectionTopRight, -1, 0)
	case pkg.DirectionRight:
		setNeighbor(pkg.DirectionTopLeft, -1, 1)
		setNeighbor(pkg.DirectionLeft, 0, 1)
		setNeighbor(pkg.DirectionTop, -1, 0)
		setNeighbor(pkg.DirectionRight, 0, -1)
		setNeighbor(pkg.DirectionTopRight, -1, -1)
	case pkg.DirectionBottomRight:
		setNeighbor(pkg.DirectionTopLeft, -1, 0)
		setNeighbor(pkg.DirectionLeft, -1, 1)
		setNeighbor(pkg.DirectionTop, -1, -1)
		setNeighbor(pkg.DirectionRight, 1, -1)
		setNeighbor(pkg.DirectionTopRight, 0, -1)
	case pkg.DirectionBottom:
		setNeighbor(pkg.DirectionTopLeft, -1, -1)
		setNeighbor(pkg.DirectionLeft, -1, 0)
		setNeighbor(pkg.DirectionTop, 0, -1)
		setNeighbor(pkg.DirectionRight, 1, 0)
		setNeighbor(pkg.DirectionTopRight, 1, -1)
	case pkg.DirectionBottomLeft:
		setNeighbor(pkg.DirectionTopLeft, 0, -1)
		setNeighbor(pkg.DirectionLeft, -1, -1)
		setNeighbor(pkg.DirectionTop, 1, -1)
		setNeighbor(pkg.DirectionRight, 1, 1)
		setNeighbor(pkg.DirectionTopRight, 1, 0)
	case pkg.DirectionLeft:
		setNeighbor(pkg.DirectionTopLeft, 1, -1)
		setNeighbor(pkg.DirectionLeft, 0, -1)
		setNeighbor(pkg.DirectionTop, 1, 0)
		setNeighbor(pkg.DirectionRight, 0, 1)
		setNeighbor(pkg.DirectionTopRight, 1, 1)
	case pkg.DirectionTopLeft:
		setNeighbor(pkg.DirectionTopLeft, -1, 0)
		setNeighbor(pkg.DirectionLeft, -1, 1)
		setNeighbor(pkg.DirectionTop, -1, -1)
		setNeighbor(pkg.DirectionRight, 1, -1)
		setNeighbor(pkg.DirectionTopRight, 0, -1)
	}

	return nil, m.collisionEnv
}

func (m *Map) applyActions(actions ...types.Action) {
	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.applyActions")()
	}

	for _, aBz := range actions {
		tile1Pos := aBz.GetTilePos()
		tile1 := m.getTile(tile1Pos.X, tile1Pos.Y)

		switch a := aBz.(type) {
		case types.MultiplyForce:
			if !tile1.HasParticle() {
				break
			}
			tile1.Particle.MultiplyForce(a.K)
		case types.ReflectForce:
			if !tile1.HasParticle() {
				break
			}
			tile1.Particle.ReflectForce(a.Horizontal, a.Vertical)
		case types.AddForce:
			if !tile1.HasParticle() {
				break
			}
			tile1.Particle.AddForce(a.ForceVec)
		case types.AlterForce:
			if !tile1.HasParticle() {
				break
			}
			tile1.Particle.SetForce(a.NewForceVec)
		case types.RotateForce:
			if !tile1.HasParticle() {
				break
			}
			tile1.Particle.RotateForce(a.Angle)
		case types.MoveTile:
			if !tile1.HasParticle() {
				break
			}
			m.moveTile(tile1, a.NewTilePos)
		case types.SwapTiles:
			if !tile1.HasParticle() {
				break
			}
			tile2 := m.getTile(a.SwapTilePos.X, a.SwapTilePos.Y)
			if !tile2.HasParticle() {
				break
			}
			m.swapTiles(tile1, tile2)
		case types.ReduceHealth:
			if !tile1.HasParticle() {
				break
			}
			tile1.Particle.ReduceHealth(a.HealthDelta)
			if tile1.Particle.IsDestroyed() {
				m.removeParticle(tile1)
			}
		case types.TileReplace:
			if !tile1.HasParticle() {
				break
			}
			m.removeParticle(tile1)
			m.createParticle(tile1, a.Material)
		case types.UpdateStateParam:
			if !tile1.HasParticle() {
				break
			}
			tile1.Particle.SetStateParam(a.ParamKey, a.ParamValue)
		case types.TileAdd:
			m.createParticle(tile1, a.Material)
		}
	}
}
