package world

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/collision"
	"github.com/itiky/goPixelWorld/world/types"
)

func (m *Map) Update() {
	for i := 0; i < len(m.procActions); i++ {
		m.procActions[i] = m.procActions[i][:0]
	}

	m.iterateNonEmptyTiles(func(tile *types.Tile) {
		m.procTileWorkerWG.Add(1)
		m.procTileJobCh <- tile
	})
	m.procTileWorkerWG.Wait()

	m.processActions()
}

func (m *Map) tileWorker(id int) {
	tileEnv := closerange.NewEnvironment(nil)
	collisionEnv := collision.NewEnvironment(pkg.DirectionTop, nil, nil)

	for tile := range m.procTileJobCh {
		m.processTile(tile, tileEnv, collisionEnv, &m.procActions[id])
		m.procTileWorkerWG.Done()
	}
}

func (m *Map) processTile(tile *types.Tile, tileEnv *closerange.Environment, collisionEnv *collision.Environment, output *[]types.Action) {
	pushActions := func(actions ...types.Action) {
		*output = append(*output, actions...)
	}

	if tile.Particle.Material().Type() == types.MaterialTypeBorder {
		return
	}

	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.processTile")()
	}

	tile.Particle.UpdateState()

	if processTileEnv := m.buildTileEnv(tile, tileEnv); processTileEnv {
		tile.Particle.Material().ProcessInternal(tileEnv)
		pushActions(tileEnv.Actions()...)
	}

	targetEmptyTile, processCollisionEnv := m.buildCollisionEnv(tile, collisionEnv)
	if targetEmptyTile != nil {
		pushActions(types.NewMoveTile(tile.Pos, tile.Particle.ID(), targetEmptyTile.Pos))
	}
	if processCollisionEnv {
		collisionEnv.TargetMaterial().ProcessCollision(collisionEnv)
		pushActions(collisionEnv.Actions()...)
	}
}

func (m *Map) buildTileEnv(sourceTile *types.Tile, tileEnv *closerange.Environment) bool {
	envType := sourceTile.Particle.Material().CloseRangeType()
	if envType == types.MaterialCloseRangeTypeNone {
		return false
	}

	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.buildTileEnv")()
	}

	tileEnv.Reset(sourceTile)
	switch envType {
	case types.MaterialCloseRangeTypeSelfOnly:
	case types.MaterialCloseRangeTypeSurrounding:
		setNeighbor := func(dir pkg.Direction, dx, dy int) {
			x, y := sourceTile.Pos.X+dx, sourceTile.Pos.Y+dy
			if !m.isPositionValid(x, y) {
				return
			}
			tileEnv.SetNeighbour(dir, m.getTile(x, y))
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

			tileEnv.AddTileInRange(neighborTile)
		}
	}

	return true
}

func (m *Map) buildCollisionEnv(sourceTile *types.Tile, collisionEnv *collision.Environment) (*types.Tile, bool) {
	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.buildCollisionEnv")()
	}

	targetTile := sourceTile.TargetTile()
	if targetTile == nil {
		return nil, false
	}

	pathToTarget := sourceTile.Pos.CreatePathTo(targetTile.Pos, m.width, m.height)
	for _, pathPos := range pathToTarget {
		targetTile = m.getTile(pathPos.X, pathPos.Y)
		if targetTile.HasParticle() {
			break
		}
	}

	if !targetTile.HasParticle() {
		return targetTile, false
	}

	colDirection := pkg.NewDirectionFromCoords(sourceTile.Pos.X, sourceTile.Pos.Y, targetTile.Pos.X, targetTile.Pos.Y)

	collisionEnv.Reset(colDirection, sourceTile, targetTile)
	setNeighbor := func(dir pkg.Direction, dx, dy int) {
		x, y := targetTile.Pos.X+dx, targetTile.Pos.Y+dy
		if !m.isPositionValid(x, y) {
			return
		}
		collisionEnv.SetNeighbour(dir, m.getTile(x, y))
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

	return nil, true
}

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
