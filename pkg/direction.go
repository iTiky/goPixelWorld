package pkg

// Direction defines a relative to smth direction.
// It is used to build the CloseRange environment (neighbours relation to the source Particle) and
// the Collision environment (from which side the source Particle is colliding with the target one).
type Direction int

const (
	DirectionNone Direction = iota
	DirectionTop
	DirectionTopRight
	DirectionRight
	DirectionBottomRight
	DirectionBottom
	DirectionBottomLeft
	DirectionLeft
	DirectionTopLeft
)

var AllDirections = []Direction{DirectionTop, DirectionTopRight, DirectionRight, DirectionBottomRight, DirectionBottom, DirectionBottomLeft, DirectionLeft, DirectionTopLeft}

// Direction angle sectors ranges.
var (
	angleRadTop         = DegToRadAngle(270.0 - 45.0/2.0) // Top sector start
	angleRadTopRight    = DegToRadAngle(270.0 + 45.0/2.0) // Top sector end and Top-Right start
	angleRadRight       = DegToRadAngle(0.000 - 45.0/2.0) // etc
	angleRadBottomRight = DegToRadAngle(0.000 + 45.0/2.0)
	angleRadBottom      = DegToRadAngle(90.00 - 45.0/2.0)
	angleRadBottomLeft  = DegToRadAngle(90.00 + 45.0/2.0)
	angleRadLeft        = DegToRadAngle(180.0 - 45.0/2.0)
	angleRadTopLeft     = DegToRadAngle(180.0 + 45.0/2.0)
)

// NewDirectionFromAngle returns a new Direction based on angle.
func NewDirectionFromAngle(angleRad float64) Direction {
	angleRad = NormalizeAngle(angleRad)

	if IsRadAngleInRange(angleRad, angleRadTopRight, angleRadRight) {
		return DirectionTopRight
	} else if IsRadAngleInRange(angleRad, angleRadRight, angleRadBottomRight) {
		return DirectionRight
	} else if IsRadAngleInRange(angleRad, angleRadBottomRight, angleRadBottom) {
		return DirectionBottomRight
	} else if IsRadAngleInRange(angleRad, angleRadBottom, angleRadBottomLeft) {
		return DirectionBottom
	} else if IsRadAngleInRange(angleRad, angleRadBottomLeft, angleRadLeft) {
		return DirectionBottomLeft
	} else if IsRadAngleInRange(angleRad, angleRadLeft, angleRadTopLeft) {
		return DirectionLeft
	} else if IsRadAngleInRange(angleRad, angleRadTopLeft, angleRadTop) {
		return DirectionTopLeft
	}

	return DirectionTop
}

// NewDirectionFromCoords returns a new Direction based on Positions.
func NewDirectionFromCoords(fromX, fromY, toX, toY int) Direction {
	dirVec := NewVectorByCoordinates(0,
		float64(fromX), float64(fromY),
		float64(toX), float64(toY),
	)

	return NewDirectionFromAngle(dirVec.Angle())
}

// Angle returns the corresponding angle.
func (d Direction) Angle() (angleRad float64) {
	switch d {
	case DirectionTop:
		return Rad270
	case DirectionTopRight:
		return Rad315
	case DirectionRight:
		return Rad0
	case DirectionBottomRight:
		return Rad45
	case DirectionBottom:
		return Rad90
	case DirectionBottomLeft:
		return Rad135
	case DirectionLeft:
		return Rad180
	case DirectionTopLeft:
		return Rad225
	}

	return 0
}

// Sector returns a set of Directions in a sector (close to {d} within depth).
func (d Direction) Sector(depth uint) []Direction {
	if depth >= 4 {
		return AllDirections
	}

	normalizeDirection := func(dir Direction) Direction {
		if dir < 0 {
			return dir + 8
		}
		if dir > 7 {
			return dir - 8
		}
		return dir
	}

	neighbours := make([]Direction, 0, depth*2+1)
	neighbours = append(neighbours, d)
	for i := 1; i <= int(depth); i++ {
		neighbours = append(neighbours, normalizeDirection(d+Direction(i)))
		neighbours = append(neighbours, normalizeDirection(d-Direction(i)))
	}

	return neighbours
}

// Next returns the next Direction (clockwise).
func (d Direction) Next() Direction {
	dNext := d + 1
	if dNext > 7 {
		dNext = 0
	}

	return dNext
}

// Rotate180 returns the inverted Direction.
func (d Direction) Rotate180() Direction {
	return (d + 4) % 8
}

// Int ...
func (d Direction) Int() int {
	return int(d)
}

func (d Direction) String() string {
	switch d {
	case DirectionTop:
		return "Top"
	case DirectionTopRight:
		return "TopRight"
	case DirectionRight:
		return "Right"
	case DirectionBottomRight:
		return "BottomRight"
	case DirectionBottom:
		return "Bottom"
	case DirectionBottomLeft:
		return "BottomLeft"
	case DirectionLeft:
		return "Left"
	case DirectionTopLeft:
		return "TopLeft"
	}

	return "Unknown"
}
