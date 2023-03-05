package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"

	"github.com/itiky/goPixelWorld/pkg"
	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

const (
	cursorRadiusDef  = 22 // circle cursor default radius
	cursorRadiusMin  = 5  // circle cursor min radius
	cursorRadiusMax  = 80 // circle cursor max radius
	cursorRadiusStep = 3  // circle cursor radius inc/dec step
)

// cursorTool keeps the currently selected Material tool data and pending World input actions.
type cursorTool struct {
	material           worldTypes.MaterialI // current Material selected (nil if not)
	radius             int                  // current circle type radius
	applyForce         bool                 // if "apply random force" is toggled
	dotImage           *ebiten.Image        // dot type image
	circleColor        color.Color          // current tools color
	pendingWorldAction worldAction          // the next World input action to apply (nil if none)
}

// newCursorTool creates a new Cursor tool.
func newCursorTool() *cursorTool {
	const (
		tileWidth  = 5
		tileHeight = 5
	)

	return &cursorTool{
		radius:     1,
		applyForce: false,
		dotImage:   ebiten.NewImage(tileWidth, tileHeight),
	}
}

// Draw ...
func (t *cursorTool) Draw(screen *ebiten.Image, drawOpts *ebiten.DrawImageOptions) {
	mouseX, mouseY := ebiten.CursorPosition()

	if t.radius == 1 {
		toolWidth, toolHeight := t.dotImage.Size()

		x := float64(mouseX - toolWidth)
		y := float64(mouseY - toolHeight)

		drawOpts.GeoM.Reset()
		drawOpts.GeoM.Translate(x, y)

		screen.DrawImage(t.dotImage, drawOpts)
	} else {
		ebitenutil.DrawCircle(screen, float64(mouseX), float64(mouseY), float64(t.radius), t.circleColor)
	}
}

// OnPress generates a new World input action.
func (t *cursorTool) OnPress(mouseX, mouseY int) {
	if t.material != nil {
		t.pendingWorldAction = createParticlesWorldAction{
			mouseX:      mouseX,
			mouseY:      mouseY,
			mouseRadius: t.radius,
			material:    t.material,
			applyForce:  t.applyForce,
		}
	} else {
		t.pendingWorldAction = deleteParticlesWorldAction{
			mouseX:      mouseX,
			mouseY:      mouseY,
			mouseRadius: t.radius,
		}
	}
}

// GetPendingWorldAction ...
func (t *cursorTool) GetPendingWorldAction() worldAction {
	if t.pendingWorldAction == nil {
		return nil
	}

	a := t.pendingWorldAction
	t.pendingWorldAction = nil

	return a
}

// UpdateMaterial updates the cursor color on Material change.
func (t *cursorTool) UpdateMaterial(material worldTypes.MaterialI) {
	var baseColor color.Color
	if material != nil {
		baseColor = material.Color()
	} else {
		baseColor = colornames.Black
	}

	circleColor := pkg.ColorToNRGBA(baseColor)
	circleColor.A = 0xA0

	t.dotImage.Fill(baseColor)
	t.material = material
	t.circleColor = circleColor
}

// EnableCircle switches the cursor to the circle mode.
func (t *cursorTool) EnableCircle() {
	t.radius = cursorRadiusDef
}

// EnableDot switches the cursor to the dot mode.
func (t *cursorTool) EnableDot() {
	t.radius = 1
}

// EnableRandomForce enables the "apply random force" mode.
func (t *cursorTool) EnableRandomForce() {
	t.applyForce = true
}

// DisableRandomForce disables the "apply random force" mode.
func (t *cursorTool) DisableRandomForce() {
	t.applyForce = false
}

// IncRadius increments the circle radius for the circle mode.
func (t *cursorTool) IncRadius() {
	if t.radius == 1 {
		return
	}

	t.radius += cursorRadiusStep
	if t.radius > cursorRadiusMax {
		t.radius = cursorRadiusMax
	}
}

// DecRadius decrements the circle radius for the circle mode.
func (t *cursorTool) DecRadius() {
	if t.radius == 1 {
		return
	}

	t.radius -= cursorRadiusStep
	if t.radius < cursorRadiusMin {
		t.radius = cursorRadiusMin
	}
}

// FlipGravity generates a new World input action.
func (t *cursorTool) FlipGravity() {
	t.pendingWorldAction = flipGravityWorldAction{}
}
