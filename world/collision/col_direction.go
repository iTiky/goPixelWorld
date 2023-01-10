package collision

import (
	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

type Direction int

const (
	DirectionN  Direction = iota
	DirectionNE Direction = iota
	DirectionE  Direction = iota
	DirectionSE Direction = iota
	DirectionS  Direction = iota
	DirectionSW Direction = iota
	DirectionW  Direction = iota
	DirectionNW Direction = iota
	//
	collisionDiagMarginDeg = 7.5
)

var (
	angleRadN  = pkg.DegToRadAngle(225.0 + collisionDiagMarginDeg)
	angleRadNE = pkg.DegToRadAngle(315.0 - collisionDiagMarginDeg)
	angleRadE  = pkg.DegToRadAngle(315.0 + collisionDiagMarginDeg)
	angleRadSE = pkg.DegToRadAngle(45.0 - collisionDiagMarginDeg)
	angleRadS  = pkg.DegToRadAngle(45.0 + collisionDiagMarginDeg)
	angleRadSW = pkg.DegToRadAngle(135.0 - collisionDiagMarginDeg)
	angleRadW  = pkg.DegToRadAngle(135.0 + collisionDiagMarginDeg)
	angleRadNW = pkg.DegToRadAngle(225.0 - collisionDiagMarginDeg)
)

func (d Direction) Rotate() Direction {
	return (d + 4) % 8
}

func GetDirection(pos1, pos2 types.Position) Direction {
	dirVec := pkg.NewVectorByCoordinates(0,
		float64(pos1.X), float64(pos1.Y),
		float64(pos2.X), float64(pos2.Y),
	)

	return GetVectorDirection(dirVec)
}

func GetVectorDirection(vec pkg.Vector) Direction {
	if vec.Angle() >= angleRadNE && vec.Angle() < angleRadE {
		return DirectionNE
	}
	if vec.Angle() >= angleRadE && vec.Angle() < angleRadSE {
		return DirectionE
	}
	if vec.Angle() >= angleRadSE && vec.Angle() < angleRadS {
		return DirectionSE
	}
	if vec.Angle() >= angleRadS && vec.Angle() < angleRadSW {
		return DirectionS
	}
	if vec.Angle() >= angleRadSW && vec.Angle() < angleRadW {
		return DirectionSW
	}
	if vec.Angle() >= angleRadW && vec.Angle() < angleRadNW {
		return DirectionW
	}
	if vec.Angle() >= angleRadNW && vec.Angle() < angleRadN {
		return DirectionNW
	}

	return DirectionN
}
