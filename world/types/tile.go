package types

import (
	"fmt"
	"image/color"

	"golang.org/x/image/colornames"
)

var _ TileI = &Tile{}

type Tile struct {
	Pos      Position
	Particle *Particle
}

func NewTile(pos Position, particle *Particle) *Tile {
	return &Tile{
		Pos:      pos,
		Particle: particle,
	}
}

func (t *Tile) Position() Position {
	return t.Pos
}

func (t *Tile) Color() color.Color {
	if !t.HasParticle() {
		return colornames.Wheat
	}

	return t.Particle.Color()
}

func (t *Tile) HasParticle() bool {
	return t.Particle != nil
}

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
