package pkg

import (
	"math"
	"math/rand"
	"time"
)

const (
	Rad0   = 0
	Rad45  = Rad90 / 2.0
	Rad90  = math.Pi / 2.0
	Rad135 = Rad90 + Rad45
	Rad180 = math.Pi
	Rad225 = Rad180 + Rad45
	Rad270 = 3.0 * math.Pi / 2.0
	Rad315 = Rad270 + Rad45
)

// AbsInt ...
func AbsInt(v int) int {
	if v < 0 {
		return -v
	}

	return v
}

// CmpInt ...
func CmpInt(v1, v2 int) int {
	if v2 > v1 {
		return 1
	}
	if v2 < v1 {
		return -1
	}

	return 0
}

// DegToRadAngle converts angle in degrees to rads.
func DegToRadAngle(v float64) float64 {
	return NormalizeAngle(v * math.Pi / 180.0)
}

// RadToDegAngle converts angle in rads to degrees.
func RadToDegAngle(v float64) float64 {
	return NormalizeAngle(v) * 180.0 / math.Pi
}

// NormalizeAngle returns the angle in the [-pi, pi] range.
func NormalizeAngle(angleRad float64) float64 {
	return math.Atan2(math.Sin(angleRad), math.Cos(angleRad))
}

// IsRadAngleInRange checks if the angle is in range (must be normalized).
func IsRadAngleInRange(angle, min, max float64) bool {
	if angle < 0.0 {
		angle += 2 * math.Pi
	}
	if min < 0.0 {
		min += 2 * math.Pi
	}
	if max < 0.0 {
		max += 2 * math.Pi
	}

	return angle >= min && angle < max
}

// RandomAngle ...
func RandomAngle() (angleRad float64) {
	return DegToRadAngle(float64(rand.Int31n(360)))
}

// FlipCoin ...
func FlipCoin() bool {
	return rand.Intn(2) == 0
}

// RollDice ...
func RollDice(size int) bool {
	if size < 2 {
		size = 2
	}
	return rand.Intn(size) == 0
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
