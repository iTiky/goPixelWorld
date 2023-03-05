package world

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/itiky/goPixelWorld/monitor"
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/materials"
	"github.com/itiky/goPixelWorld/world/types"
)

const (
	workersNum          = 32
	tileWorkerJobChSize = 50000
)

type MapOption func(*Map) error

type Map struct {
	width     int
	height    int
	particles map[uint64]*types.Tile
	grid      [][]*types.Tile
	//
	monitor *monitor.Keeper
	//
	procTileWorkerWG sync.WaitGroup
	procTileJobCh    chan *types.Tile
	procActions      [][]types.Action
	//
	processingRequestCh chan struct{}
	processingAckCh     chan struct{}
	processingOutput    []types.Pixel
}

func WithWidth(width int) MapOption {
	return func(m *Map) error {
		if width <= 0 {
			return fmt.Errorf("invalid map width: %d", width)
		}

		m.width = width
		return nil
	}
}

func WithHeight(height int) MapOption {
	return func(m *Map) error {
		if height <= 0 {
			return fmt.Errorf("invalid map height: %d", height)
		}

		m.height = height
		return nil
	}
}

func WithMonitor(keeper *monitor.Keeper) MapOption {
	return func(m *Map) error {
		if keeper == nil {
			return fmt.Errorf("monitor keeper is nil")
		}

		m.monitor = keeper

		return nil
	}
}

func NewMap(opts ...MapOption) (*Map, error) {
	m := Map{
		width:     200,
		height:    200,
		particles: make(map[uint64]*types.Tile),
		//
		procTileJobCh: make(chan *types.Tile, tileWorkerJobChSize),
		//
		processingRequestCh: make(chan struct{}),
		processingAckCh:     make(chan struct{}),
	}
	for _, opt := range opts {
		if err := opt(&m); err != nil {
			return nil, err
		}
	}

	m.initGrid(m.width, m.height)

	for i := 0; i < workersNum; i++ {
		m.procActions = append(m.procActions, make([]types.Action, 0, tileWorkerJobChSize))
		go m.tileWorker(i)
	}

	go m.processingWorker()
	m.processingStart()

	return &m, nil
}

func (m *Map) Size() (int, int) {
	return m.width, m.height
}

func (m *Map) ExportState(fn func(pixel types.Pixel)) {
	m.processingDone()
	defer m.processingStart()

	for i := 0; i < len(m.processingOutput); i++ {
		if !m.processingOutput[i].Ready {
			break
		}
		fn(m.processingOutput[i])
	}
}

func (m *Map) CreateParticles(x, y, radius int, materialBz types.MaterialI, randomForce bool) {
	m.processingDone()
	defer m.processingStart()

	for _, pos := range types.PositionsInCircle(x, y, radius, true) {
		if !m.isPositionValid(pos.X, pos.Y) {
			continue
		}

		material, ok := materialBz.(types.Material)
		if !ok {
			continue
		}

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

		if randomForce {
			forceVec := pkg.NewVector(
				float64(rand.Int31n(5)),
				pkg.RandomAngle(),
			)
			tile.Particle.SetForce(forceVec)
		}
	}
}

func (m *Map) RemoveParticles(x, y, radius int) {
	m.processingDone()
	defer m.processingStart()

	for _, pos := range types.PositionsInCircle(x, y, radius, true) {
		if !m.isPositionValid(pos.X, pos.Y) {
			continue
		}

		tile := m.getTile(pos.X, pos.Y)
		if !tile.HasParticle() {
			continue
		}

		m.removeParticle(tile)
	}
}

func (m *Map) FlipGravity() {
	closerange.FlipGravity()
}

func (m *Map) isPositionValid(x, y int) bool {
	return x >= 0 && y >= 0 && x < m.width && y < m.height
}

func (m *Map) initGrid(width, height int) {
	m.width = width
	m.height = height

	for pID := range m.particles {
		delete(m.particles, pID)
	}

	m.grid = make([][]*types.Tile, m.width)
	for x := 0; x < m.width; x++ {
		m.grid[x] = make([]*types.Tile, m.height)
		for y := 0; y < m.height; y++ {
			var material types.Material
			if x == 0 || y == 0 || x == m.width-1 || y == m.height-1 {
				material = materials.NewBorder()
			}

			m.createTile(types.NewPosition(x, y), material)
			m.processingOutput = append(m.processingOutput, types.Pixel{})
		}
	}
}
