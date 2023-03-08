package world

import (
	"fmt"
	"sync"

	"github.com/itiky/goPixelWorld/monitor"
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

	/* Input state */
	inputActions []types.InputAction

	/* Nature state */
	natureEnabled           bool
	natureCloudsTimeout     int
	natureWindChangeTimeout int

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

// WithNatureEffects options enables nature effects like clouds, wind, etc.
func WithNatureEffects() MapOption {
	return func(m *Map) error {
		m.natureEnabled = true
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
	// Nature init
	m.initNatureEvents()

	return &m, nil
}

// Size returns the grid size.
func (m *Map) Size() (int, int) {
	return m.width, m.height
}

// ExportState exports the current Map state.
// Waits for the current processing round to end and starts a new one after the export is done.
func (m *Map) ExportState(fn func(pixel types.TileI)) {
	// Wait for the previous processing round to finish
	m.processingDone()
	defer m.processingStart()

	// Export
	for i := 0; i < len(m.procOutput); i++ {
		if !m.procOutput[i].Ready {
			break
		}
		fn(m.procOutput[i])
	}

	// Nature events
	if m.natureEnabled {
		m.inputActions = append(m.inputActions, m.handleNatureEvents()...)
	}

	// Handle input actions
	// That alters the map state, so we need to apply actions before the next processing round
	for _, actionBz := range m.inputActions {
		switch action := actionBz.(type) {
		case types.CreateParticlesInputAction:
			m.handleCreateParticlesInput(action)
		case types.DeleteParticlesInputAction:
			m.handleRemoveParticlesInput(action)
		case types.FlipGravityInputAction:
			m.handleFlipGravityInput()
		}
	}
	m.inputActions = m.inputActions[:0]
}

// isPositionValid checks if Position if within the grid.
func (m *Map) isPositionValid(x, y int) bool {
	return x >= 0 && y >= 0 && x < m.width && y < m.height
}
