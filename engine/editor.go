package engine

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

type editor struct {
	toolTiles []toolTile
	//
	mouseLeftInput  *mouseInput
	mouseRightInput *mouseInput
	keyboardInput   *keyboardInput
	//
	cursor *cursorTool
	//
	drawOpts *ebiten.DrawImageOptions
}

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

		e.keyboardInput.SetCallback(ebiten.KeyQ, func() {
			e.cursor.DecRadius()
		})
		e.keyboardInput.SetCallback(ebiten.KeyE, func() {
			e.cursor.IncRadius()
		})

		removeToggleTool := newRemoveToggleTile(e.cursor.UpdateMaterial)
		e.keyboardInput.SetCallback(ebiten.KeyD, func() {
			removeToggleTool.ToggleOn()
		})

		circleCursorTool := newGenericToggleTile(
			"Circle",
			func() {
				e.cursor.EnableCircle()
			},
			func() {
				e.cursor.EnableDot()
			},
		)
		e.keyboardInput.SetCallback(ebiten.KeyS, func() {
			circleCursorTool.Toggle()
		})

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

func (e *editor) Layout(screenWidth, screenHeight int) {
	for _, t := range e.toolTiles {
		t.Layout(screenWidth, screenHeight)
	}
}

func (e *editor) Draw(screen *ebiten.Image) {
	for _, t := range e.toolTiles {
		t.Draw(screen, e.drawOpts)
	}
	e.cursor.Draw(screen, e.drawOpts)
}

func (e *editor) HandleInput() {
	e.mouseLeftInput.Update()
	e.mouseRightInput.Update()
	e.keyboardInput.Update()

	e.handleLeftClick()
	e.handleRightClick()
}

func (e *editor) GetNextWorldAction() worldAction {
	return e.cursor.GetPendingWorldAction()
}

func (e *editor) handleLeftClick() {
	mouseX, mouseY, isPressed := e.mouseLeftInput.IsPressed()
	if !isPressed {
		return
	}

	e.cursor.OnPress(mouseX, mouseY)
}

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
