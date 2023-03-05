package collision

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

// ReflectSourceTargetForces reflects the source and the target Particle force Vectors based on collision angle and mass ratio.
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

	e.actions = append(e.actions, types.NewAlterForce(e.source.Pos, e.source.Particle.ID(), sourceVecAfter))
	if !e.target.Particle.Material().IsFlagged(types.MaterialFlagIsUnmovable) {
		e.actions = append(e.actions, types.NewAlterForce(e.target.Pos, e.target.Particle.ID(), targetVecAfter))
	}

	return true
}

// ReflectSourceForce reflects the source Particle force Vector based on collision direction.
// A simple logic with horizontal / vertical / both flips.
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
	e.actions = append(e.actions, types.NewReflectForce(e.source.Pos, e.source.Particle.ID(), flipV, flipH))

	return true
}
