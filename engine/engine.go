package engine

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"github.com/itiky/goPixelWorld/monitor"
	"github.com/itiky/goPixelWorld/world"
	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

var _ ebiten.Game = &Runner{}

// RunnerOption defines the Runner constructor options.
type RunnerOption func(r *Runner) error

// Runner keeps the renderer state and implements the Ebiten Game interface.
type Runner struct {
	worldMap            *world.Map // the World
	editor              *editor    // the Editor
	screenWidthInitial  int        // initial screen width
	screenHeightInitial int        // initial screen height
	// Current state
	screenWidth  int                           // the current screen layout width
	screenHeight int                           // the current screen layout height
	tileSize     float64                       // the current Tile size relative to (screenWidth, screenHeight)
	tilesCache   map[color.Color]*ebiten.Image // cached pixels
	tileDrawOpts *ebiten.DrawImageOptions      // reused object to save some time on rendering
	// External services
	monitor *monitor.Keeper
}

// WithScreenSize defines the window size.
func WithScreenSize(width, height int) RunnerOption {
	return func(r *Runner) error {
		if width <= 0 || height <= 0 {
			return fmt.Errorf("invalid screen size: %dx%d", width, height)
		}

		r.screenWidthInitial, r.screenHeightInitial = width, height

		return nil
	}
}

// WithMonitor adds an external monitor.
func WithMonitor(monitor *monitor.Keeper) RunnerOption {
	return func(r *Runner) error {
		if monitor == nil {
			return fmt.Errorf("monitor is nil")
		}

		r.monitor = monitor

		return nil
	}
}

// NewRunner builds a new Runner.
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
	//ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	ebiten.SetTPS(60)

	return r, nil
}

// Run implements the ebiten.Game interface.
func (r *Runner) Run() error {
	return ebiten.RunGame(r)
}

// Update implements the ebiten.Game interface.
// Handles the mouse / keyboard inputs.
func (r *Runner) Update() error {
	if r.monitor != nil {
		defer r.monitor.TrackOpDuration("Runner.Update")()
	}

	// Handle the inputs (pass the action to the World)
	if r.editor != nil {
		r.editor.HandleInput()
		r.applyWorldAction(r.editor.GetNextWorldAction())
	}

	return nil
}

// Draw implements the ebiten.Game interface.
// Renders the grid and editor.
func (r *Runner) Draw(screen *ebiten.Image) {
	if r.monitor != nil {
		defer r.monitor.TrackOpDuration("Runner.Draw")()
		r.monitor.AddFrame()
	}

	// Render pixels
	drawnPixels := r.drawTiles(screen)

	// Render the editor
	if r.editor != nil {
		r.editor.Draw(screen)
	}

	// Render the debug text message
	mouseX, mouseY := ebiten.CursorPosition()
	fps := ebiten.ActualFPS()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("[%d, %d]\n%.1f\n%d",
		r.mouseCoordToWorld(mouseX), r.mouseCoordToWorld(mouseY),
		fps,
		drawnPixels,
	))
}

// Layout implements the ebiten.Game interface.
// Called on window resize.
// Calculates the new Tile size and drops Tiles caches (does the same for the Editor).
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

// drawTiles iterates over all non-empty Tiles and draws them.
// Utilizes the image cache to save some FPSs.
func (r *Runner) drawTiles(screen *ebiten.Image) int64 {
	drawnPixels := int64(0)

	r.worldMap.ExportState(func(tile worldTypes.TileI) {
		tileImage, found := r.tilesCache[tile.Color()]
		if !found {
			tileImage = ebiten.NewImage(int(r.tileSize), int(r.tileSize))
			tileImage.Fill(tile.Color())
			r.tilesCache[tile.Color()] = tileImage
		}

		tileDrawX, tileDrawY := float64(tile.X())*r.tileSize, float64(tile.Y())*r.tileSize
		r.tileDrawOpts.GeoM.Reset()
		r.tileDrawOpts.GeoM.Translate(tileDrawX, tileDrawY)

		screen.DrawImage(tileImage, r.tileDrawOpts)

		drawnPixels++
	})

	return drawnPixels
}

// applyWorldAction passes through editor input to the World by type.
// Method transforms the engine coordinates and sizes to the World's.
func (r *Runner) applyWorldAction(actionBz worldTypes.InputAction) {
	if actionBz == nil {
		return
	}

	switch action := actionBz.(type) {
	case worldTypes.CreateParticlesInputAction:
		action.X = int(float64(action.X) / r.tileSize)
		action.Y = int(float64(action.Y) / r.tileSize)
		action.Radius = int(float64(action.Radius) / r.tileSize)

		r.worldMap.PushInputAction(action)
	case worldTypes.DeleteParticlesInputAction:
		action.X = int(float64(action.X) / r.tileSize)
		action.Y = int(float64(action.Y) / r.tileSize)
		action.Radius = int(float64(action.Radius) / r.tileSize)

		r.worldMap.PushInputAction(action)
	case worldTypes.FlipGravityInputAction:
		r.worldMap.PushInputAction(action)
	}
}

// mouseCoordToWorld converts the mouse coordinate to the World coordinate.
func (r *Runner) mouseCoordToWorld(c int) int {
	return int(float64(c) / r.tileSize)
}
