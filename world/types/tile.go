package types

import (
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"
)

var _ TileI = &Tile{}

// Tile represents a positioned Particle.
type Tile struct {
	Pos      Position
	Particle *Particle
}

// NewTile creates a new Tile.
func NewTile(pos Position, particle *Particle) *Tile {
	return &Tile{
		Pos:      pos,
		Particle: particle,
	}
}

// X returns the x coordinate.
func (t *Tile) X() int {
	return t.Pos.X
}

// Y returns the y coordinate.
func (t *Tile) Y() int {
	return t.Pos.Y
}

// Color return the current Particle color.
func (t *Tile) Color() color.Color {
	if !t.HasParticle() {
		return colornames.Wheat
	}

	return t.Particle.Color()
}

// HasParticle checks if the Tile has a Particle.
func (t *Tile) HasParticle() bool {
	return t.Particle != nil
}

// TargetTile calculates the target Tile based the Particle force Vector.
func (t *Tile) TargetTile() *Tile {
	if !t.HasParticle() {
		return nil
	}

	targetPos := NewPosition(
		t.Pos.X+int(t.Particle.forceVec.X()),
		t.Pos.Y+int(t.Particle.forceVec.Y()),
	)
	if targetPos.Equal(t.Pos) {
		return nil
	}

	return NewTile(targetPos, nil)
}

func (t *Tile) String() string {
	particleStr := "no particle"
	if t.Particle != nil {
		particleStr = t.Particle.String()
	}

	return fmt.Sprintf("Tile(%s, %s)", t.Pos, particleStr)
}
