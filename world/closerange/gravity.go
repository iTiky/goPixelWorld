package closerange

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

const (
	gravityMag = 0.15
)

var (
	gravityDownVec = pkg.NewVector(gravityMag, pkg.Rad90)
	gravityUpVec   = pkg.NewVector(gravityMag, pkg.Rad270)
)

func (e *Environment) AddGravity() bool {
	e.actions = append(e.actions, types.NewAddForce(e.source.Pos, gravityDownVec))
	return true
}

func (e *Environment) AddReverseGravity() bool {
	e.actions = append(e.actions, types.NewAddForce(e.source.Pos, gravityUpVec))
	return true
}
