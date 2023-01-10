package materials

import (
	"image/color"

	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Border{}

type Border struct {
	base
}

func NewBorder() Border {
	return Border{
		base: newBase(
			color.RGBA{R: 0x1A, G: 0x1A, B: 0x1A, A: 0xFF},
			withForceDamperK(0.1),
			withFlags(types.MaterialFlagIsUnremovable),
		),
	}
}

func (m Border) Type() types.MaterialType {
	return types.MaterialTypeBorder
}

func (m Border) ProcessCollision(env types.CollisionEnvironment) {
	env.DampSourceForce(m.forceDamperK)
	env.ReflectSourceForce()
}
