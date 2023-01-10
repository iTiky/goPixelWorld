package engine

import (
	worldTypes "github.com/itiky/goPixelWorld/world/types"
)

type worldActionType int

const (
	worldActionTypeCreateParticles worldActionType = iota
	worldActionDeleteParticles
)

type worldAction interface {
	Type() worldActionType
}

type createParticlesWorldAction struct {
	mouseX, mouseY int
	mouseRadius    int
	material       worldTypes.MaterialI
	applyForce     bool
}

func (a createParticlesWorldAction) Type() worldActionType {
	return worldActionTypeCreateParticles
}

type deleteParticlesWorldAction struct {
	mouseX, mouseY int
	mouseRadius    int
}

func (a deleteParticlesWorldAction) Type() worldActionType {
	return worldActionDeleteParticles
}
