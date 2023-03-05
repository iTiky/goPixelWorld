package pkg

import (
	"fmt"
	"math"
)

// Vector defines a math vector which state is defined by magnitude and angle (normalized).
type Vector struct {
	mag float64
	ang float64
}

// NewVector builds a new Vector.
func NewVector(magnitude, angle float64) Vector {
	return Vector{
		mag: magnitude,
		ang: NormalizeAngle(angle),
	}
}

// NewVectorByCoordinates builds a new Vector from coordinates.
func NewVectorByCoordinates(magnitude, x0, y0, x1, y1 float64) Vector {
	return Vector{
		mag: magnitude,
		ang: NormalizeAngle(math.Atan2(y1-y0, x1-x0)),
	}
}

// Magnitude returns the magnitude.
func (v Vector) Magnitude() float64 {
	return v.mag
}

// Angle returns the angle.
func (v Vector) Angle() float64 {
	return v.ang
}

// X returns the X projection.
func (v Vector) X() float64 {
	return v.mag * math.Cos(v.ang)
}

// Y returns the Y projection.
func (v Vector) Y() float64 {
	return v.mag * math.Sin(v.ang)
}

// IsZero checks if the Vector is empty.
func (v Vector) IsZero() bool {
	if v.mag == 0.0 {
		return true
	}

	return false
}

// SetMagnitude sets the magnitude.
func (v Vector) SetMagnitude(mag float64) Vector {
	return Vector{
		mag: mag,
		ang: v.ang,
	}
}

// Add adds two Vectors.
func (v Vector) Add(v2 Vector) Vector {
	if v2.IsZero() {
		return v
	}

	v1x := v.X()
	v1y := v.Y()

	v2x := v2.X()
	v2y := v2.Y()

	v3x := v1x + v2x
	v3y := v1y + v2y

	v3m := math.Sqrt(math.Pow(v3x, 2) + math.Pow(v3y, 2))
	v3a := math.Atan2(v3y, v3x)

	return Vector{
		mag: v3m,
		ang: NormalizeAngle(v3a),
	}
}

// Reflect returns the reflected Vector.
func (v Vector) Reflect(horizontal, vertical bool) Vector {
	vx, vy := v.X(), v.Y()
	if horizontal {
		vx *= -1.0
	}
	if vertical {
		vy *= -1.0
	}

	return Vector{
		mag: v.mag,
		ang: math.Atan2(vy, vx),
	}
}

// MultiplyByK alters the Vector's magnitude.
func (v Vector) MultiplyByK(k float64) Vector {
	if k >= 0.0 {
		return Vector{
			mag: v.mag * k,
			ang: v.ang,
		}
	}

	return Vector{
		mag: v.mag * math.Abs(k),
		ang: NormalizeAngle(v.ang + math.Pi),
	}
}

// Rotate returns the rotated by angle Vector.
func (v Vector) Rotate(angleRad float64) Vector {
	return Vector{
		mag: v.mag,
		ang: NormalizeAngle(v.ang + angleRad),
	}
}

// DotProduct returns the dot product of two vectors (projection).
func (v Vector) DotProduct(v2 Vector) float64 {
	return v.X()*v2.X() + v.Y()*v2.Y()
}

func (v Vector) String() string {
	return fmt.Sprintf("Vector(mag=%0.2f, angle=%0.1f)", v.mag, RadToDegAngle(v.ang))
}
