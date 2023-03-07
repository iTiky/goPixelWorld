package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Border{}

// Border is a static non-removable Particle.
type Border struct {
	base
}

func NewBorder() Border {
	return Border{
		base: newBase(
			types.MaterialTypeBorder,
			color.RGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF},
			withFlags(types.MaterialFlagIsUnremovable, types.MaterialFlagIsUnmovable),
			withCloseRangeType(types.MaterialCloseRangeTypeNone),
			withSourceDamping(0.2, 0.0),
		),
	}
}

func (m Border) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces(m.srcForceDamperK)
}
