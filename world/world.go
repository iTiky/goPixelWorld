package world

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/itiky/goPixelWorld/monitor"
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/types"
)

// MapOption defines the Map constructor option.
type MapOption func(*Map) error

// Map keeps the world state and performs state changes in the background.
type Map struct {
	/* Grid state */
	// Grid size
	width, height int
	// ParticleID -> Tile mapping
	particles map[uint64]*types.Tile
	// Position (coordinates) -> Tile mapping
	grid [][]*types.Tile

	/* Processing state */
	// Tile workers input jobs queue (a Tile to process)
	procTileJobCh chan *types.Tile
	// procTileJobCh current size counter
	// Drops to 0 when workers did clean up the queue
	procTileWorkerWG sync.WaitGroup
	// Per tile worker output actions queue
	// Each worker pushes to its own queue Actions to alter the Map state
	procActions [][]types.Action
	// Processing start requests queue
	// The next Map state calculation is halted until the next request is received
	procRequestCh chan struct{}
	// Processing done acknowledgement queue
	// Ack is sent whenever procOutput is ready to be read
	procAckCh chan struct{}
	// The next Map state to export
	procOutput []types.Pixel

	/* External services */
	monitor *monitor.Keeper
}

// WithWidth options sets the grid width.
func WithWidth(width int) MapOption {
	return func(m *Map) error {
		if width <= 0 {
			return fmt.Errorf("invalid map width: %d", width)
		}

		m.width = width
		return nil
	}
}

// WithHeight options sets the grid height.
func WithHeight(height int) MapOption {
	return func(m *Map) error {
		if height <= 0 {
			return fmt.Errorf("invalid map height: %d", height)
		}

		m.height = height
		return nil
	}
}

// WithMonitor enables the external Monitor.
func WithMonitor(keeper *monitor.Keeper) MapOption {
	return func(m *Map) error {
		if keeper == nil {
			return fmt.Errorf("monitor keeper is nil")
		}

		m.monitor = keeper

		return nil
	}
}

// NewMap creates a new Map.
func NewMap(opts ...MapOption) (*Map, error) {
	// Grid init
	m := Map{
		width:     200,
		height:    200,
		particles: make(map[uint64]*types.Tile),
	}
	for _, opt := range opts {
		if err := opt(&m); err != nil {
			return nil, err
		}
	}
	m.initGrid(m.width, m.height)

	// Processing init
	m.initProcessing()

	return &m, nil
}

// Size returns the grid size.
func (m *Map) Size() (int, int) {
	return m.width, m.height
}

// ExportState exports the current Map state.
// Waits for the current processing round to end and starts a new one after the export is done.
func (m *Map) ExportState(fn func(pixel types.TileI)) {
	m.processingDone()
	defer m.processingStart()

	for i := 0; i < len(m.procOutput); i++ {
		if !m.procOutput[i].Ready {
			break
		}
		fn(m.procOutput[i])
	}
}

// CreateParticles adds new Particle(s) to the Map.
// If {radius} is GT 1, creates a set of Particles in a circle area.
// If {randomForce} is enabled, applies a random force to a new Particle(s).
func (m *Map) CreateParticles(x, y, radius int, materialBz types.MaterialI, randomForce bool) {
	m.processingDone()
	defer m.processingStart()

	for _, pos := range types.PositionsInCircle(x, y, radius, true) {
		// Skip out-of-range Positions
		if !m.isPositionValid(pos.X, pos.Y) {
			continue
		}

		// Skip an unknown Material
		material, ok := materialBz.(types.Material)
		if !ok {
			continue
		}

		// If a target Tile has a particle, remove it (if it is removable)
		tile := m.getTile(pos.X, pos.Y)
		if tile.HasParticle() {
			if mType := material.Type(); mType != types.MaterialTypeFire && mType != types.MaterialTypeAntiGraviton {
				continue
			}
			if !m.removeParticle(tile) {
				continue
			}
		}
		m.createParticle(tile, material)

		// Apply an initial random force
		if randomForce {
			forceVec := pkg.NewVector(
				float64(rand.Int31n(5)),
				pkg.RandomAngle(),
			)
			tile.Particle.SetForce(forceVec)
		}
	}
}

// RemoveParticles removes existing Particle(s) from the Map.
// If {radius} is GT 1, removes a set of Particles in a circle area.
// This method skips Particles that can't be removed (borders for ex.).
func (m *Map) RemoveParticles(x, y, radius int) {
	m.processingDone()
	defer m.processingStart()

	for _, pos := range types.PositionsInCircle(x, y, radius, true) {
		// Skip out-of-range Positions
		if !m.isPositionValid(pos.X, pos.Y) {
			continue
		}

		// Skip empty Tiles
		tile := m.getTile(pos.X, pos.Y)
		if !tile.HasParticle() {
			continue
		}

		m.removeParticle(tile)
	}
}

// FlipGravity flips the global vertical gravity vector.
func (m *Map) FlipGravity() {
	closerange.FlipGravity()
}

// isPositionValid checks if Position if within the grid.
func (m *Map) isPositionValid(x, y int) bool {
	return x >= 0 && y >= 0 && x < m.width && y < m.height
}
