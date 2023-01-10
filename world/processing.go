package world

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/collision"
	"github.com/itiky/goPixelWorld/world/types"
)

func (m *Map) Update() {
	m.iterateNonEmptyTiles(m.processTile)
}

func (m *Map) processTile(tile *types.Tile) {
	if tile.Particle.Material().Type() == types.MaterialTypeBorder {
		return
	}

	tileEnv := m.buildTileEnv(tile)
	tile.Particle.Material().ProcessInternal(tileEnv)
	m.applyActions(tileEnv.Actions()...)

	tile = m.getTile(tile.Pos)
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
	state := closerange.NewEnvironment(sourceTile)
	setNeighbor := func(dir pkg.Direction, dx, dy int) {
		pos := types.NewPosition(sourceTile.Pos.X+dx, sourceTile.Pos.Y+dy)
		state.SetNeighbour(dir, m.getTile(pos))
	}

	setNeighbor(pkg.DirectionTop, 0, -1)
	setNeighbor(pkg.DirectionTopRight, 1, -1)
	setNeighbor(pkg.DirectionRight, 1, 0)
	setNeighbor(pkg.DirectionBottomRight, 1, 1)
	setNeighbor(pkg.DirectionBottom, 0, 1)
	setNeighbor(pkg.DirectionBottomLeft, -1, 1)
	setNeighbor(pkg.DirectionLeft, -1, 0)
	setNeighbor(pkg.DirectionTopLeft, -1, -1)

	return state
}

func (m *Map) buildCollisionEnv(sourceTile *types.Tile) (*types.Tile, *collision.Environment) {
	targetTile := sourceTile.TargetTile()
	if targetTile == nil {
		return nil, nil
	}

	pathToTarget := sourceTile.Pos.CreatePathTo(targetTile.Pos, m.width, m.height)
	for _, pathPos := range pathToTarget {
		targetTile = m.getTile(pathPos)
		if targetTile.HasParticle() {
			break
		}
	}

	if !targetTile.HasParticle() {
		return targetTile, nil
	}

	colDirection := pkg.NewDirectionFromCoords(sourceTile.Pos.X, sourceTile.Pos.Y, targetTile.Pos.X, targetTile.Pos.Y)

	var topLeftTile, leftTile, frontTile, rightTile, topRightTile *types.Tile
	getNeighbourTile := func(dx, dy int) *types.Tile {
		pos := types.NewPosition(targetTile.Pos.X+dx, targetTile.Pos.Y+dy)
		return m.getTile(pos)
	}

	switch colDirection {
	case pkg.DirectionTop:
		topLeftTile = getNeighbourTile(1, 1)
		leftTile = getNeighbourTile(1, 0)
		frontTile = getNeighbourTile(0, 1)
		rightTile = getNeighbourTile(-1, 0)
		topRightTile = getNeighbourTile(-1, 1)
	case pkg.DirectionTopRight:
		topLeftTile = getNeighbourTile(0, 1)
		leftTile = getNeighbourTile(1, 1)
		frontTile = getNeighbourTile(-1, 1)
		rightTile = getNeighbourTile(-1, -1)
		topRightTile = getNeighbourTile(-1, 0)
	case pkg.DirectionRight:
		topLeftTile = getNeighbourTile(-1, 1)
		leftTile = getNeighbourTile(0, 1)
		frontTile = getNeighbourTile(-1, 0)
		rightTile = getNeighbourTile(0, -1)
		topRightTile = getNeighbourTile(-1, -1)
	case pkg.DirectionBottomRight:
		topLeftTile = getNeighbourTile(-1, 0)
		leftTile = getNeighbourTile(-1, 1)
		frontTile = getNeighbourTile(-1, -1)
		rightTile = getNeighbourTile(1, -1)
		topRightTile = getNeighbourTile(0, -1)
	case pkg.DirectionBottom:
		topLeftTile = getNeighbourTile(-1, -1)
		leftTile = getNeighbourTile(-1, 0)
		frontTile = getNeighbourTile(0, -1)
		rightTile = getNeighbourTile(1, 0)
		topRightTile = getNeighbourTile(1, -1)
	case pkg.DirectionBottomLeft:
		topLeftTile = getNeighbourTile(0, -1)
		leftTile = getNeighbourTile(-1, -1)
		frontTile = getNeighbourTile(1, -1)
		rightTile = getNeighbourTile(1, 1)
		topRightTile = getNeighbourTile(1, 0)
	case pkg.DirectionLeft:
		topLeftTile = getNeighbourTile(1, -1)
		leftTile = getNeighbourTile(0, -1)
		frontTile = getNeighbourTile(1, 0)
		rightTile = getNeighbourTile(0, 1)
		topRightTile = getNeighbourTile(1, 1)
	case pkg.DirectionTopLeft:
		topLeftTile = getNeighbourTile(-1, 0)
		leftTile = getNeighbourTile(-1, 1)
		frontTile = getNeighbourTile(-1, -1)
		rightTile = getNeighbourTile(1, -1)
		topRightTile = getNeighbourTile(0, -1)
	}

	state := collision.NewEnvironment(colDirection, sourceTile, targetTile)
	state.SetNeighbour(pkg.DirectionTopLeft, topLeftTile)
	state.SetNeighbour(pkg.DirectionLeft, leftTile)
	state.SetNeighbour(pkg.DirectionTop, frontTile)
	state.SetNeighbour(pkg.DirectionRight, rightTile)
	state.SetNeighbour(pkg.DirectionTopRight, topRightTile)

	return nil, state
}

func (m *Map) applyActions(actions ...types.Action) {
	for _, aBz := range actions {
		tile1 := m.getTile(aBz.GetTilePos())
		if tile1 == nil {
			continue
		}

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
			tile2 := m.getTile(a.SwapTilePos)
			if tile2 == nil || !tile2.HasParticle() {
				break
			}
			m.swapTiles(tile1, tile2)
		case types.ReduceHealth:
			if !tile1.HasParticle() {
				break
			}
			tile1.Particle.ReduceHealth(a.HealthDelta)
			if tile1.Particle.IsDestroyed() {
				m.removeTile(tile1)
			}
		case types.TileReplace:
			if !tile1.HasParticle() {
				break
			}
			m.removeTile(tile1)
			m.createTile(tile1.Pos, a.Material)
		case types.TileAdd:
			if tile1.HasParticle() {
				break
			}
			m.createTile(tile1.Pos, a.Material)
		}
	}
}
