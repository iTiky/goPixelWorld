package collision

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

// https://www.vobarian.com/collisions/2dcollisions2.pdf
func (e *Environment) ReflectSourceTargetForces(sourceForceDampK float64) bool {
	sourceVecBefore, targetVecBefore := e.source.Particle.ForceVector(), e.target.Particle.ForceVector()
	sourceMass, targetMass := e.source.Particle.Material().Mass(), e.target.Particle.Material().Mass()

	if sourceForceDampK > 0.0 {
		sourceVecBefore = sourceVecBefore.MultiplyByK(sourceForceDampK)
	}

	normalVec := pkg.NewVectorByCoordinates(1.0, targetVecBefore.X(), targetVecBefore.Y(), sourceVecBefore.X(), sourceVecBefore.Y())
	tangentVec := normalVec.Rotate(pkg.Rad90)

	sourceVecNormalPrBefore := normalVec.DotProduct(sourceVecBefore)
	sourceVecTangentPrAfter := tangentVec.DotProduct(sourceVecBefore)
	targetVecNormalPrBefore := normalVec.DotProduct(targetVecBefore)
	targetVecTangentPrAfter := tangentVec.DotProduct(targetVecBefore)

	sourceVecNormalPrAfter := (sourceVecNormalPrBefore*(sourceMass-targetMass) + 2.0*targetMass*targetVecNormalPrBefore) / (sourceMass + targetMass)
	targetVecNormalPrAfter := (targetVecNormalPrBefore*(targetMass-sourceMass) + 2.0*sourceMass*sourceVecNormalPrBefore) / (sourceMass + targetMass)

	sourceVecAfter := normalVec.MultiplyByK(sourceVecNormalPrAfter).Add(tangentVec.MultiplyByK(sourceVecTangentPrAfter))
	targetVecAfter := normalVec.MultiplyByK(targetVecNormalPrAfter).Add(tangentVec.MultiplyByK(targetVecTangentPrAfter))

	e.actions = append(e.actions, []types.Action{
		types.NewAlterForce(e.source.Pos, sourceVecAfter),
		types.NewAlterForce(e.target.Pos, targetVecAfter),
	}...)

	return true
}

func (e *Environment) ReflectSourceForce() bool {
	flipV, flipH := false, false
	switch e.direction {
	case pkg.DirectionTop, pkg.DirectionBottom:
		flipV = true
	case pkg.DirectionLeft, pkg.DirectionRight:
		flipH = true
	case pkg.DirectionTopRight, pkg.DirectionBottomRight, pkg.DirectionBottomLeft, pkg.DirectionTopLeft:
		flipV, flipH = true, true
	}
	e.actions = append(e.actions, types.NewReflectForce(e.source.Pos, flipV, flipH))

	return true
}
