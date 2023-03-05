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

// Tool Tile rendering defaults
const (
	toolTileWidth  = 88 // tool tile width
	toolTileHeight = 32 // tool tile height
	//
	toolTileMargin      = 10 // margin between tools
	toolTileOffsetRight = 10 // offset from the right window side
	//
	toolTextOffsetLeft = 10 // tool text left offset (from tool's Tile)
	toolTestOffsetTop  = 23 // tool text top offset (from tool's Tile)
	//
	fontDPI = 72 // font rendering DPI
)

var (
	// The last unique tool ID
	lastToolID = -1
	// Toggle tools defaults
	toggleToolOffColor  = color.RGBA{R: 0x2B, G: 0x2B, B: 0x2B, A: 0xFF} // OFF color
	toggleToolOnColor   = color.RGBA{R: 0x00, G: 0x3C, B: 0x18, A: 0xFF} // ON color
	toggleToolFontColor = color.RGBA{R: 0xE5, G: 0xE5, B: 0xE5, A: 0xFF} // font color
	toggleToolFont      font.Face                                        // font
)

// toolTile defines an interface for all Tile tools.
type toolTile interface {
	// OnClick is a tool callback on mouse click.
	OnClick(cursorX, cursorY int) bool
	// Layout updates a tool rendering params on window resize.
	Layout(screenWidth, screenHeight int)
	// Draw draws a tool.
	Draw(screen *ebiten.Image, drawOpts *ebiten.DrawImageOptions)
}

// toolBase defines common fields for all tools.
type toolBase struct {
	id           int           // unique ID
	tileX, tileY float64       // top-left Tile coordinates
	textX, textY int           // top-left Tile's font coordinates
	image        *ebiten.Image // tool image
	text         string        // tool text
	toggled      bool          // is toggled flag
}

// Layout ...
func (t *toolBase) Layout(screenWidth, screenHeight int) {
	tileWidth, tileHeight := t.image.Size()

	t.tileX = float64(screenWidth - tileWidth - toolTileOffsetRight)
	t.tileY = float64(t.id*(tileHeight+toolTileMargin) + toolTileMargin)

	if t.text != "" {
		t.textX = int(t.tileX) + toolTextOffsetLeft
		t.textY = int(t.tileY) + toolTestOffsetTop
	}
}

// Draw ...
func (t *toolBase) Draw(screen *ebiten.Image, drawOpts *ebiten.DrawImageOptions) {
	drawOpts.GeoM.Reset()
	drawOpts.GeoM.Translate(t.tileX, t.tileY)

	screen.DrawImage(t.image, drawOpts)

	if t.text != "" {
		text.Draw(screen, t.text, toggleToolFont, t.textX, t.textY, toggleToolFontColor)
	}
}

// isCursorOver returns true if the mouse cursor is over this tool.
func (t *toolBase) isCursorOver(cursorX, cursorY int) bool {
	tileWidth, tileHeight := t.image.Size()

	return cursorX >= int(t.tileX) && cursorX <= int(t.tileX)+tileWidth &&
		cursorY >= int(t.tileY) && cursorY <= int(t.tileY)+tileHeight
}

// toggleOn changes a tool rendering and state on toggle ON.
func (t *toolBase) toggleOn() bool {
	if t.toggled {
		return false
	}

	t.image.Fill(toggleToolOnColor)
	t.toggled = true

	return true
}

// toggleOff changes a tool rendering and state on toggle OFF.
func (t *toolBase) toggleOff() bool {
	if !t.toggled {
		return false
	}

	t.image.Fill(toggleToolOffColor)
	t.toggled = false

	return true
}

// materialTile defines a Material pick tool.
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

// removeToggleTile defines a removal tool
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

// genericToggleTile defines a generic toggle ON/OFF tool.
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

// nextToolID returns the next unique tool ID.
func nextToolID() int {
	lastToolID++

	return lastToolID
}

func init() {
	// Load the font from Ebiten example
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
