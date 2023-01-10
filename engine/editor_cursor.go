package engine

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/colornames"

	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

const (
	cursorRadiusDef  = 17
	cursorRadiusMin  = 5
	cursorRadiusMax  = 80
	cursorRadiusStep = 3
)

type cursorTool struct {
	material           worldTypes.MaterialI
	radius             int
	applyForce         bool
	dotImage           *ebiten.Image
	circleColor        color.Color
	pendingWorldAction worldAction
}

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

func (t *cursorTool) GetPendingWorldAction() worldAction {
	if t.pendingWorldAction == nil {
		return nil
	}

	a := t.pendingWorldAction
	t.pendingWorldAction = nil

	return a
}

func (t *cursorTool) UpdateMaterial(material worldTypes.MaterialI) {
	var baseColor color.Color
	if material != nil {
		baseColor = material.Color()
	} else {
		baseColor = colornames.Black
	}

	circleColor := ColorToNRGBA(baseColor)
	circleColor.A = 0xA0

	t.dotImage.Fill(baseColor)
	t.material = material
	t.circleColor = circleColor
}

func (t *cursorTool) EnableCircle() {
	t.radius = cursorRadiusDef
}

func (t *cursorTool) EnableDot() {
	t.radius = 1
}

func (t *cursorTool) EnableRandomForce() {
	t.applyForce = true
}

func (t *cursorTool) DisableRandomForce() {
	t.applyForce = false
}

func (t *cursorTool) IncRadius() {
	if t.radius == 1 {
		return
	}

	t.radius += cursorRadiusStep
	if t.radius > cursorRadiusMax {
		t.radius = cursorRadiusMax
	}
}

func (t *cursorTool) DecRadius() {
	if t.radius == 1 {
		return
	}

	t.radius -= cursorRadiusStep
	if t.radius < cursorRadiusMin {
		t.radius = cursorRadiusMin
	}
}

func ColorToNRGBA(c color.Color) color.NRGBA {
	r, g, b, a := c.RGBA()

	return color.NRGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
