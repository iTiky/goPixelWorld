package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

// editor keeps the state of the right toolbar.
type editor struct {
	toolTiles []toolTile // material and toggle tiles
	//
	mouseLeftInput  *mouseInput    // left mouse button input state
	mouseRightInput *mouseInput    // right mouse button input state
	keyboardInput   *keyboardInput // keyboard input state
	//
	cursor *cursorTool // the current cursor tool state
	//
	drawOpts *ebiten.DrawImageOptions // reused object to save some time on rendering
}

// WithEditorUI enables the Editor toolbox.
func WithEditorUI(materials ...worldTypes.MaterialI) RunnerOption {
	return func(r *Runner) error {
		if len(materials) == 0 {
			return fmt.Errorf("no materials provided")
		}

		e := editor{
			mouseLeftInput:  newMouseInput(ebiten.MouseButtonLeft),
			mouseRightInput: newMouseInput(ebiten.MouseButtonRight),
			keyboardInput:   newKeyboardInput(),
			drawOpts:        &ebiten.DrawImageOptions{},
			cursor:          newCursorTool(),
		}

		// Register the "decrement circle cursor radius" keyboard callback
		e.keyboardInput.SetCallback(ebiten.KeyQ, func() {
			e.cursor.DecRadius()
		})
		// Register the "increment circle cursor radius" keyboard callback
		e.keyboardInput.SetCallback(ebiten.KeyE, func() {
			e.cursor.IncRadius()
		})
		// Register the "flip vertical gravity" keyboard callback
		e.keyboardInput.SetCallback(ebiten.KeyZ, func() {
			e.cursor.FlipGravity()
		})

		// Remove Particles tool
		removeToggleTool := newRemoveToggleTile(e.cursor.UpdateMaterial)
		e.keyboardInput.SetCallback(ebiten.KeyD, func() {
			removeToggleTool.ToggleOn()
		})

		// Circle / dot cursor type toggle tool
		circleCursorTool := newGenericToggleTile(
			"Circle",
			func() {
				e.cursor.EnableCircle()
			},
			func() {
				e.cursor.EnableDot()
			},
		)
		circleCursorTool.Toggle()
		e.keyboardInput.SetCallback(ebiten.KeyS, func() {
			circleCursorTool.Toggle()
		})

		// Apply random force to new Particles toggle tool
		randomForceTool := newGenericToggleTile(
			"RndF",
			func() {
				e.cursor.EnableRandomForce()
			},
			func() {
				e.cursor.DisableRandomForce()
			},
		)
		e.keyboardInput.SetCallback(ebiten.KeyF, func() {
			randomForceTool.Toggle()
		})

		// Create Material tools and assign 1..9 keyboard input callbacks to them
		for idx, m := range materials {
			materialTool := newMaterialTile(m, func(m worldTypes.MaterialI) {
				e.cursor.UpdateMaterial(m)
				removeToggleTool.ToggleOff()
			})

			e.keyboardInput.SetCallback(ebiten.KeyDigit1+ebiten.Key(idx), func() {
				materialTool.OnClick(-1, -1)
				removeToggleTool.ToggleOff()
			})
			e.toolTiles = append(e.toolTiles, materialTool)
		}
		e.toolTiles = append(e.toolTiles,
			removeToggleTool,
			circleCursorTool,
			randomForceTool,
		)
		e.toolTiles[0].OnClick(-1, -1)

		r.editor = &e

		return nil
	}
}

// Layout updates tool sizes on window resize.
func (e *editor) Layout(screenWidth, screenHeight int) {
	for _, t := range e.toolTiles {
		t.Layout(screenWidth, screenHeight)
	}
}

// Draw draws all the registered tools.
func (e *editor) Draw(screen *ebiten.Image) {
	for _, t := range e.toolTiles {
		t.Draw(screen, e.drawOpts)
	}
	e.cursor.Draw(screen, e.drawOpts)
}

// HandleInput updates the mouse / keyboard input states.
// That can create a new World input action.
func (e *editor) HandleInput() {
	e.mouseLeftInput.Update()
	e.mouseRightInput.Update()
	e.keyboardInput.Update()

	e.handleLeftClick()
	e.handleRightClick()
}

// GetNextWorldAction returns the next World input action (if any).
func (e *editor) GetNextWorldAction() worldTypes.InputAction {
	return e.cursor.GetPendingWorldAction()
}

// handleLeftClick passes through the left mouse click event to the Cursor tool.
func (e *editor) handleLeftClick() {
	mouseX, mouseY, isPressed := e.mouseLeftInput.IsPressed()
	if !isPressed {
		return
	}

	e.cursor.OnPress(mouseX, mouseY)
}

// handleRightClick passes through the right mouse click event to the Cursor tool.
func (e *editor) handleRightClick() {
	mouseX, mouseY, isClicked := e.mouseRightInput.IsClicked()
	if !isClicked {
		return
	}

	for _, t := range e.toolTiles {
		if t.OnClick(mouseX, mouseY) {
			return
		}
	}
}
