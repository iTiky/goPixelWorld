package engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// mouseInputState defines an enum for the mouse click state.
type mouseInputState int

const (
	mouseInputStateNone mouseInputState = iota
	mouseInputStatePressing
	mouseInputStateReleased
)

// mouseInput keeps the mouse click input state.
type mouseInput struct {
	btnType       ebiten.MouseButton // left / right
	state         mouseInputState    // current state
	mouseXOnPress int                // X coordinate on press
	mouseYOnPress int                // Y coordinate on press
}

// newMouseInput creates a new mouseInput.
func newMouseInput(btnType ebiten.MouseButton) *mouseInput {
	return &mouseInput{
		btnType: btnType,
		state:   mouseInputStateNone,
	}
}

// Update updates the input state machine.
func (m *mouseInput) Update() {
	isPressed := ebiten.IsMouseButtonPressed(m.btnType)

	switch m.state {
	case mouseInputStateNone:
		if isPressed {
			m.state = mouseInputStatePressing
			m.mouseXOnPress, m.mouseYOnPress = ebiten.CursorPosition()
		}
	case mouseInputStatePressing:
		if !isPressed {
			m.state = mouseInputStateReleased
		}
	}
}

// IsPressed return the mouse cursor coordinates and isPressed flag (button is being pressed).
func (m *mouseInput) IsPressed() (int, int, bool) {
	if m.state != mouseInputStatePressing {
		return 0, 0, false
	}
	m.state = mouseInputStateNone

	mouseX, mouseY := ebiten.CursorPosition()

	return mouseX, mouseY, true
}

// IsClicked return the mouse cursor coordinates and isClicked flag (button was pressed and released).
func (m *mouseInput) IsClicked() (int, int, bool) {
	if m.state != mouseInputStateReleased {
		return 0, 0, false
	}
	m.state = mouseInputStateNone

	return m.mouseXOnPress, m.mouseYOnPress, true
}

// keyboardInput keeps the keyboard press input state.
type keyboardInput struct {
	keyTimeout     time.Duration            // press timeout to avoid "drift"
	keysBuf        []ebiten.Key             // key codes we are interested in
	keyPressEvents map[ebiten.Key]time.Time // key onPress timestamps by code
	callbacks      map[ebiten.Key]func()    // key callbacks
}

// newKeyboardInput creates a new keyboardInput.
func newKeyboardInput() *keyboardInput {
	return &keyboardInput{
		keyTimeout:     100 * time.Millisecond,
		keysBuf:        make([]ebiten.Key, 0, 3),
		keyPressEvents: make(map[ebiten.Key]time.Time),
		callbacks:      make(map[ebiten.Key]func()),
	}
}

// SetCallback sets a new callback for the code.
func (m *keyboardInput) SetCallback(key ebiten.Key, callback func()) {
	m.callbacks[key] = callback
}

// Update updates the input state machine.
func (m *keyboardInput) Update() {
	now := time.Now()

	// Get currently pressed keys buffer and store the press event timestamp (if not stored already)
	m.keysBuf = inpututil.AppendPressedKeys(m.keysBuf[:0])
	for _, key := range m.keysBuf {
		callback := m.callbacks[key]
		if callback == nil {
			continue
		}

		if _, found := m.keyPressEvents[key]; !found {
			m.keyPressEvents[key] = now
		}
	}

	// Check the timeout and emit callbacks
	for key, pressedAt := range m.keyPressEvents {
		if now.Sub(pressedAt) < m.keyTimeout {
			continue
		}
		delete(m.keyPressEvents, key)

		m.callbacks[key]()
	}
}
