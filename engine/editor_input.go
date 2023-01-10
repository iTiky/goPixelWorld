package engine

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type mouseInputState int

const (
	mouseInputStateNone mouseInputState = iota
	mouseInputStatePressing
	mouseInputStateReleased
)

type mouseInput struct {
	btnType       ebiten.MouseButton
	state         mouseInputState
	mouseXOnPress int
	mouseYOnPress int
}

func newMouseInput(btnType ebiten.MouseButton) *mouseInput {
	return &mouseInput{
		btnType: btnType,
		state:   mouseInputStateNone,
	}
}

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

func (m *mouseInput) IsPressed() (int, int, bool) {
	if m.state != mouseInputStatePressing {
		return 0, 0, false
	}
	m.state = mouseInputStateNone

	mouseX, mouseY := ebiten.CursorPosition()

	return mouseX, mouseY, true
}

func (m *mouseInput) IsClicked() (int, int, bool) {
	if m.state != mouseInputStateReleased {
		return 0, 0, false
	}
	m.state = mouseInputStateNone

	return m.mouseXOnPress, m.mouseYOnPress, true
}

type keyboardInput struct {
	keyTimeout     time.Duration
	keysBuf        []ebiten.Key
	keyPressEvents map[ebiten.Key]time.Time
	callbacks      map[ebiten.Key]func()
}

func newKeyboardInput() *keyboardInput {
	return &keyboardInput{
		keyTimeout:     100 * time.Millisecond,
		keysBuf:        make([]ebiten.Key, 0, 3),
		keyPressEvents: make(map[ebiten.Key]time.Time),
		callbacks:      make(map[ebiten.Key]func()),
	}
}

func (m *keyboardInput) SetCallback(key ebiten.Key, callback func()) {
	m.callbacks[key] = callback
}

func (m *keyboardInput) Update() {
	now := time.Now()

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

	for key, pressedAt := range m.keyPressEvents {
		if now.Sub(pressedAt) < m.keyTimeout {
			continue
		}
		delete(m.keyPressEvents, key)

		m.callbacks[key]()
	}
}
