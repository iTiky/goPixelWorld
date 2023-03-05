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

// AddGravity adds the vertical gravity force Vector.
func (e *Environment) AddGravity() bool {
	e.actions = append(e.actions, types.NewAddForce(e.source.Pos, e.source.Particle.ID(), gravityDownVec))
	return true
}

// AddReverseGravity adds the reversed vertical gravity force Vector.
func (e *Environment) AddReverseGravity() bool {
	e.actions = append(e.actions, types.NewAddForce(e.source.Pos, e.source.Particle.ID(), gravityUpVec))
	return true
}

// FlipGravity flips the vertical gravity.
func FlipGravity() {
	gravityDownVec, gravityUpVec = gravityUpVec, gravityDownVec
}
