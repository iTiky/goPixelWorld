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
	//
	windVec = pkg.NewVector(0, 0)
)

func (e *Environment) AddGravity() bool {
	e.actions = append(e.actions, types.NewAddForce(e.source.Pos, e.source.Particle.ID(), gravityDownVec))
	return true
}

func (e *Environment) AddReverseGravity() bool {
	e.actions = append(e.actions, types.NewAddForce(e.source.Pos, e.source.Particle.ID(), gravityUpVec))
	return true
}

func (e *Environment) AddWind() bool {
	if windVec.IsZero() {
		return false
	}
	if e.source.Particle.Material().IsFlagged(types.MaterialFlagIsUnmovable) {
		return false
	}

	e.actions = append(e.actions, types.NewAddForce(e.source.Pos, e.source.Particle.ID(), windVec))
	return true
}

// FlipGravity flips the vertical gravity.
func FlipGravity() {
	gravityDownVec, gravityUpVec = gravityUpVec, gravityDownVec
}

// SetWind adds the global wind force Vector for movable particles only.
func SetWind(mag float64, left bool) {
	var angle float64
	if left {
		angle = pkg.DirectionLeft.Angle()
	} else {
		angle = pkg.DirectionRight.Angle()
	}

	windVec = pkg.NewVector(mag, angle)
}

// GetWind returns the current global wind Vector.
func GetWind() pkg.Vector {
	return windVec
}
