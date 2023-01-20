package world

import (
	"fmt"
	"math/rand"

	"github.com/itiky/goPixelWorld/monitor"
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/closerange"
	"github.com/itiky/goPixelWorld/world/collision"
	"github.com/itiky/goPixelWorld/world/materials"
	"github.com/itiky/goPixelWorld/world/types"
)

type MapOption func(*Map) error

type Map struct {
	width     int
	height    int
	particles map[uint64]*types.Tile
	grid      [][]*types.Tile
	//
	tileEnv      *closerange.Environment
	collisionEnv *collision.Environment
	//
	monitor *monitor.Keeper
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
		width:        200,
		height:       200,
		particles:    make(map[uint64]*types.Tile),
		tileEnv:      closerange.NewEnvironment(nil),
		collisionEnv: collision.NewEnvironment(pkg.DirectionTop, nil, nil),
	}
	for _, opt := range opts {
		if err := opt(&m); err != nil {
			return nil, err
		}
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
		}
	}

	return &m, nil
}

func (m *Map) Size() (int, int) {
	return m.width, m.height
}

func (m *Map) IterateTiles(fn func(tile types.TileI)) {
	m.iterateNonEmptyTiles(func(tile *types.Tile) {
		fn(tile)
	})
}

func (m *Map) Update() {
	m.iterateNonEmptyTiles(m.processTile)
}

func (m *Map) CreateParticles(x, y, radius int, materialBz types.MaterialI, randomForce bool) {
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
			if material.Type() != types.MaterialTypeFire {
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
