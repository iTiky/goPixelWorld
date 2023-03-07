package world

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/collision"
	"github.com/itiky/goPixelWorld/world/types"
)

// processingWorker is worker that iterates over all non-empty Tiles, processes them and
// handles all output Action events which updates the Map state.
// As a result it fills up the output ready to be collected buffer.
func (m *Map) processingWorker() {
	for {
		// Wait for the next processing roundrequest
		<-m.procRequestCh

		// Clean up output buffers
		for i := 0; i < len(m.procActions); i++ {
			m.procActions[i] = m.procActions[i][:0]
		}

		// Fill up the jobs queue and for it to be processed
		m.iterateNonEmptyTiles(func(tile *types.Tile) {
			m.procTileWorkerWG.Add(1)
			m.procTileJobCh <- tile
		})
		m.procTileWorkerWG.Wait()

		// Alter the Map state
		m.processActions()

		// Prepare the output buffer
		pixelIdx := 0
		m.iterateNonEmptyTiles(func(tile *types.Tile) {
			m.procOutput[pixelIdx].Ready = true
			m.procOutput[pixelIdx].PosX = tile.Pos.X
			m.procOutput[pixelIdx].PosY = tile.Pos.Y
			m.procOutput[pixelIdx].ParticleColor = tile.Color()
			pixelIdx++
		})
		m.procOutput[pixelIdx].Ready = false

		// Ack the processing round
		m.procAckCh <- struct{}{}
	}
}

// tileWorker processes a single Tile state update from the input job queue.
func (m *Map) tileWorker(id int) {
	tileEnv := closerange.NewEnvironment(nil)
	collisionEnv := collision.NewEnvironment(pkg.DirectionTop, nil, nil)

	for tile := range m.procTileJobCh {
		m.processTile(tile, tileEnv, collisionEnv, &m.procActions[id])
		m.procTileWorkerWG.Done()
	}
}

// processTile updates a singe Tile state based on its environment and movement path.
// All Tile / Map updates are not applied here, instead a series of Actions are pushed to the corresponding queue.
// Each Action can update this Particle state, surrounding Particle(s) or a Particle this Tile has collided with.
func (m *Map) processTile(tile *types.Tile, tileEnv *closerange.Environment, collisionEnv *collision.Environment, output *[]types.Action) {
	pushActions := func(actions ...types.Action) {
		*output = append(*output, actions...)
	}

	// Skip Borders processing
	if tile.Particle.Material().Type() == types.MaterialTypeBorder {
		return
	}

	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.processTile")()
	}

	// Self-update
	tile.Particle.UpdateState()

	// Build Tile surrounding environment and update it if required by the Tile's material
	if processTileEnv := m.buildTileEnv(tile, tileEnv); processTileEnv {
		tile.Particle.Material().ProcessInternal(tileEnv)
		pushActions(tileEnv.Actions()...)
	}

	// Process Tile movement based on path to a target Tile (where this one want to move to)
	targetEmptyTile, processCollisionEnv := m.buildCollisionEnv(tile, collisionEnv)
	if targetEmptyTile != nil {
		pushActions(types.NewMoveTile(tile.Pos, tile.Particle.ID(), targetEmptyTile.Pos))
	}
	if processCollisionEnv {
		collisionEnv.TargetMaterial().ProcessCollision(collisionEnv)
		pushActions(collisionEnv.Actions()...)
	}
}

// buildTileEnv builds a Tile surrounding state.
// Returns true, if state processing is required.
// Environment filling is based on Tile's Material requirements.
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
		// Build close neighbours environment
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
		// Build neighbours in a circle range environment
		r := sourceTile.Particle.Material().CloseRangeCircleRadius()
		for _, pos := range types.PositionsInCircle(sourceTile.Pos.X, sourceTile.Pos.Y, r, false) {
			if !m.isPositionValid(pos.X, pos.Y) {
				continue
			}

			neighborTile := m.getTile(pos.X, pos.Y)
			tileEnv.AddTileInRange(neighborTile)
		}
	}

	return true
}

// buildCollisionEnv builds a path to Tile's target position, walks through it and build a collision state if occurred.
// Returns non-nil target Tile if source Tile can freely move to it (no path collisions).
// Returns true in case of a collision meaning it should be processed by the source Material.
func (m *Map) buildCollisionEnv(sourceTile *types.Tile, collisionEnv *collision.Environment) (*types.Tile, bool) {
	if m.monitor != nil {
		defer m.monitor.TrackOpDuration("Map.buildCollisionEnv")()
	}

	// Get the target Tile based on the current Particle's force Vector
	targetTile := sourceTile.TargetTile()
	if targetTile == nil {
		// Particle "doesn't want to move"
		return nil, false
	}

	// Build a path to the target with the grid limitations
	pathToTarget := sourceTile.Pos.CreatePathTo(targetTile.Pos, m.width, m.height)
	for _, pathPos := range pathToTarget {
		targetTile = m.getTile(pathPos.X, pathPos.Y)
		if targetTile.HasParticle() {
			break
		}
	}

	// Check if the last path Tile is empty
	if !targetTile.HasParticle() {
		return targetTile, false
	}

	// Build collision direction (defines the target Particle neighbours order)
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
