package engine

import (
	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

// worldActionType defines a World input action type.
type worldActionType int

const (
	worldActionTypeCreateParticles worldActionType = iota
	worldActionDeleteParticles
	worldActionFlipGravity
)

// worldAction defines a common World input action interface.
type worldAction interface {
	Type() worldActionType
}

// createParticlesWorldAction defines a request to creates a set of new Particles.
type createParticlesWorldAction struct {
	mouseX, mouseY int
	mouseRadius    int
	material       worldTypes.MaterialI
	applyForce     bool
}

func (a createParticlesWorldAction) Type() worldActionType {
	return worldActionTypeCreateParticles
}

// deleteParticlesWorldAction defines a request to delete a set of existing Particles.
type deleteParticlesWorldAction struct {
	mouseX, mouseY int
	mouseRadius    int
}

func (a deleteParticlesWorldAction) Type() worldActionType {
	return worldActionDeleteParticles
}

// flipGravityWorldAction defines a request to flip the vertical gravity.
type flipGravityWorldAction struct{}

func (a flipGravityWorldAction) Type() worldActionType {
	return worldActionFlipGravity
}
