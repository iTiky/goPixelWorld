package engine

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/itiky/goPixelWorld/world"
	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

var _ ebiten.Game = &Runner{}

type RunnerOption func(r *Runner) error

type Runner struct {
	worldMap            *world.Map
	editor              *editor
	screenWidthInitial  int
	screenHeightInitial int
	//
	screenWidth  int
	screenHeight int
	//
	tileSize     float64
	tilesCache   map[color.Color]*ebiten.Image
	tileDrawOpts *ebiten.DrawImageOptions
}

func WithScreenSize(width, height int) RunnerOption {
	return func(r *Runner) error {
		if width <= 0 || height <= 0 {
			return fmt.Errorf("invalid screen size: %dx%d", width, height)
		}

		r.screenWidthInitial, r.screenHeightInitial = width, height

		return nil
	}
}

func NewRunner(worldMap *world.Map, opts ...RunnerOption) (*Runner, error) {
	const (
		screenWidthInitial  = 800
		screenHeightInitial = 600
	)

	if worldMap == nil {
		return nil, fmt.Errorf("worldMap is nil")
	}

	r := &Runner{
		worldMap:            worldMap,
		screenWidthInitial:  screenWidthInitial,
		screenHeightInitial: screenHeightInitial,
		tilesCache:          make(map[color.Color]*ebiten.Image),
		tileDrawOpts:        &ebiten.DrawImageOptions{},
	}
	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, fmt.Errorf("applying option: %w", err)
		}
	}

	ebiten.SetWindowTitle("Go Pixel World")
	ebiten.SetWindowSize(r.screenWidthInitial, r.screenHeightInitial)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	return r, nil
}

func (r *Runner) Run() error {
	return ebiten.RunGame(r)
}

func (r *Runner) Update() error {
	if r.editor != nil {
		r.editor.HandleInput()
		r.applyWorldAction(r.editor.GetNextWorldAction())
	}

	r.worldMap.Update()

	return nil
}

func (r *Runner) Draw(screen *ebiten.Image) {
	r.drawTiles(screen)

	if r.editor != nil {
		r.editor.Draw(screen)
	}
}

func (r *Runner) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if r.screenWidth != outsideWidth || r.screenHeight != outsideHeight {
		mapWidth, mapHeight := r.worldMap.Size()

		tileSize := float64(outsideWidth) / float64(mapWidth)
		if v := float64(outsideHeight) / float64(mapHeight); v < tileSize {
			tileSize = v
		}

		r.screenWidth, r.screenHeight = outsideWidth, outsideHeight
		r.tileSize = tileSize
		r.tilesCache = make(map[color.Color]*ebiten.Image)

		if r.editor != nil {
			r.editor.Layout(outsideWidth, outsideHeight)
		}
	}

	return outsideWidth, outsideHeight
}

func (r *Runner) drawTiles(screen *ebiten.Image) {
	r.worldMap.IterateTiles(func(tile worldTypes.TileI) {
		tilePos := tile.Position()
		tileColor := tile.Color()

		tileImage, found := r.tilesCache[tileColor]
		if !found {
			tileImage = ebiten.NewImage(int(r.tileSize), int(r.tileSize))
			tileImage.Fill(tileColor)
			r.tilesCache[tileColor] = tileImage
		}

		tileDrawX, tileDrawY := float64(tilePos.X)*r.tileSize, float64(tilePos.Y)*r.tileSize
		r.tileDrawOpts.GeoM.Reset()
		r.tileDrawOpts.GeoM.Translate(tileDrawX, tileDrawY)

		screen.DrawImage(tileImage, r.tileDrawOpts)
	})
}

func (r *Runner) applyWorldAction(actionBz worldAction) {
	if actionBz == nil {
		return
	}

	switch action := actionBz.(type) {
	case createParticlesWorldAction:
		x := int(float64(action.mouseX) / r.tileSize)
		y := int(float64(action.mouseY) / r.tileSize)
		radius := int(float64(action.mouseRadius) / r.tileSize)

		r.worldMap.CreateParticles(x, y, radius, action.material, action.applyForce)
	case deleteParticlesWorldAction:
		x := int(float64(action.mouseX) / r.tileSize)
		y := int(float64(action.mouseY) / r.tileSize)
		radius := int(float64(action.mouseRadius) / r.tileSize)

		r.worldMap.RemoveParticles(x, y, radius)
	}
}