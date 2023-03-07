package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) ReplaceSelf(newMaterial types.Material) bool {
	e.actions = append(e.actions, types.NewTileReplace(e.source.Pos, e.source.Particle.ID(), newMaterial))
	return true
}

func (e *Environment) UpdateStateParam(paramKey string, paramValue int) bool {
	if e.source.Particle.GetStateParam(paramKey) == paramValue {
		return false
	}

	e.actions = append(e.actions, types.NewUpdateStateParam(e.source.Pos, e.source.Particle.ID(), paramKey, paramValue))
	return true
}

func (e *Environment) AddSelfForce(vec pkg.Vector) bool {
	e.actions = append(e.actions, types.NewAddForce(e.source.Pos, e.source.Particle.ID(), vec))
	return true
}

func (e *Environment) SetSelfForce(vec pkg.Vector) bool {
	e.actions = append(e.actions, types.NewAlterForce(e.source.Pos, e.source.Particle.ID(), vec))
	return true
}
