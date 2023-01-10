package engine

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

const (
	toolTileWidth  = 88
	toolTileHeight = 32
	//
	toolTileMargin      = 10
	toolTileOffsetRight = 10
	//
	toolTextOffsetLeft = 10
	toolTestOffsetTop  = 23
	//
	fontDPI = 72
)

var (
	lastToolID          = -1
	toggleToolOffColor  = color.RGBA{R: 0x2B, G: 0x2B, B: 0x2B, A: 0xFF}
	toggleToolOnColor   = color.RGBA{R: 0x00, G: 0x3C, B: 0x18, A: 0xFF}
	toggleToolFontColor = color.RGBA{R: 0xE5, G: 0xE5, B: 0xE5, A: 0xFF}
	toggleToolFont      font.Face
)

type toolTile interface {
	OnClick(cursorX, cursorY int) bool
	Layout(screenWidth, screenHeight int)
	Draw(screen *ebiten.Image, drawOpts *ebiten.DrawImageOptions)
}

type toolBase struct {
	id           int
	tileX, tileY float64
	textX, textY int
	image        *ebiten.Image
	text         string
	toggled      bool
}

func (t *toolBase) Layout(screenWidth, screenHeight int) {
	tileWidth, tileHeight := t.image.Size()

	t.tileX = float64(screenWidth - tileWidth - toolTileOffsetRight)
	t.tileY = float64(t.id*(tileHeight+toolTileMargin) + toolTileMargin)

	if t.text != "" {
		t.textX = int(t.tileX) + toolTextOffsetLeft
		t.textY = int(t.tileY) + toolTestOffsetTop
	}
}

func (t *toolBase) Draw(screen *ebiten.Image, drawOpts *ebiten.DrawImageOptions) {
	drawOpts.GeoM.Reset()
	drawOpts.GeoM.Translate(t.tileX, t.tileY)

	screen.DrawImage(t.image, drawOpts)

	if t.text != "" {
		text.Draw(screen, t.text, toggleToolFont, t.textX, t.textY, toggleToolFontColor)
	}
}

func (t *toolBase) isCursorOver(cursorX, cursorY int) bool {
	tileWidth, tileHeight := t.image.Size()

	return cursorX >= int(t.tileX) && cursorX <= int(t.tileX)+tileWidth &&
		cursorY >= int(t.tileY) && cursorY <= int(t.tileY)+tileHeight
}

func (t *toolBase) toggleOn() bool {
	if t.toggled {
		return false
	}

	t.image.Fill(toggleToolOnColor)
	t.toggled = true

	return true
}

func (t *toolBase) toggleOff() bool {
	if !t.toggled {
		return false
	}

	t.image.Fill(toggleToolOffColor)
	t.toggled = false

	return true
}

type materialTile struct {
	toolBase
	material worldTypes.MaterialI
	callback func(m worldTypes.MaterialI)
}

func newMaterialTile(material worldTypes.MaterialI, callback func(m worldTypes.MaterialI)) *materialTile {
	t := materialTile{
		toolBase: toolBase{
			id:    nextToolID(),
			image: ebiten.NewImage(toolTileWidth, toolTileHeight),
		},
		material: material,
		callback: callback,
	}
	t.toolBase.image.Fill(t.material.Color())

	return &t
}

func (t *materialTile) OnClick(cursorX, cursorY int) bool {
	if !t.isCursorOver(cursorX, cursorY) && cursorX != -1 && cursorY != -1 {
		return false
	}

	t.callback(t.material)

	return true
}

type removeToggleTile struct {
	toolBase
	callback func(m worldTypes.MaterialI)
}

func newRemoveToggleTile(callback func(m worldTypes.MaterialI)) *removeToggleTile {
	t := removeToggleTile{
		toolBase: toolBase{
			id:      nextToolID(),
			image:   ebiten.NewImage(toolTileWidth, toolTileHeight),
			text:    "Remove",
			toggled: true,
		},
		callback: callback,
	}
	t.ToggleOff()

	return &t
}

func (t *removeToggleTile) OnClick(cursorX, cursorY int) bool {
	if !t.isCursorOver(cursorX, cursorY) {
		return false
	}

	t.ToggleOn()

	return true
}

func (t *removeToggleTile) ToggleOn() {
	if !t.toggleOn() {
		return
	}
	t.callback(nil)
}

func (t *removeToggleTile) ToggleOff() {
	t.toggleOff()
}

type genericToggleTile struct {
	toolBase
	callbackOn  func()
	callbackOff func()
}

func newGenericToggleTile(text string, callbackOn, callbackOff func()) *genericToggleTile {
	t := genericToggleTile{
		toolBase: toolBase{
			id:      nextToolID(),
			image:   ebiten.NewImage(toolTileWidth, toolTileHeight),
			text:    text,
			toggled: true,
		},
		callbackOn:  callbackOn,
		callbackOff: callbackOff,
	}
	t.Toggle()

	return &t
}

func (t *genericToggleTile) OnClick(cursorX, cursorY int) bool {
	if !t.isCursorOver(cursorX, cursorY) {
		return false
	}

	t.Toggle()

	return true
}

func (t *genericToggleTile) Toggle() {
	if t.toggleOn() {
		t.callbackOn()
		return
	}
	if t.toggleOff() {
		t.callbackOff()
	}
}

func nextToolID() int {
	lastToolID++

	return lastToolID
}

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal("Parsing font:", err)
	}

	toggleToolFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    18,
		DPI:     fontDPI,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("Creating font face:", err)
	}
}
