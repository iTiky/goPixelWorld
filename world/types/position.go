package types

import (
	"fmt"
	"math"

	"github.com/itiky/goPixelWorld/pkg"
)

// Position wraps the coordinates.
type Position struct {
	X, Y int
}

// NewPosition creates a new Position.
func NewPosition(x, y int) Position {
	return Position{
		X: x,
		Y: y,
	}
}

// Equal checks equality.
func (p Position) Equal(p2 Position) bool {
	return p.X == p2.X && p.Y == p2.Y
}

// CreatePathTo creates a path to the target Position (discrete line approximation).
// A variation of Dijkstra's algo.
func (p Position) CreatePathTo(toPos Position, xMax, yMax int) (path []Position) {
	if toPos.Equal(p) {
		return nil
	}

	dx, dy := toPos.X-p.X, toPos.Y-p.Y
	dxAbs, dyAbs := pkg.AbsInt(dx), pkg.AbsInt(dy)
	xStep, yStep := pkg.CmpInt(p.X, toPos.X), pkg.CmpInt(p.Y, toPos.Y)

	appendToPathFailed := false
	defer func() {
		if appendToPathFailed {
			return
		}
		path = append(path, toPos)
	}()

	appendToPath := func(pos Position) bool {
		if pos.X < 0 || pos.X >= xMax {
			appendToPathFailed = true
			return false
		}
		if pos.Y < 0 || pos.Y >= yMax {
			appendToPathFailed = true
			return false
		}
		path = append(path, pos)

		return true
	}

	if dxAbs == 0 {
		path = make([]Position, 0, dyAbs)

		x := toPos.X
		for y := p.Y + yStep; pkg.AbsInt(toPos.Y-y) > 0; y += yStep {
			if !appendToPath(NewPosition(x, y)) {
				break
			}
		}
		return
	}

	if dyAbs == 0 {
		path = make([]Position, 0, dxAbs)

		y := toPos.Y
		for x := p.X + xStep; pkg.AbsInt(toPos.X-x) > 0; x += xStep {
			if !appendToPath(NewPosition(x, y)) {
				break
			}
		}
		return
	}

	if dxAbs == dyAbs {
		path = make([]Position, 0, dxAbs)

		for x, y := p.X+xStep, p.Y+yStep; pkg.AbsInt(toPos.X-x) > 0; x, y = x+xStep, y+yStep {
			if !appendToPath(NewPosition(x, y)) {
				break
			}
		}
		return
	}

	lineSlope := float64(dy) / float64(dx)
	lineIntercept := float64(p.Y) - lineSlope*float64(p.X)

	var xyNext func(i int) (int, int)
	if dxAbs > dyAbs {
		path = make([]Position, 0, dxAbs)
		xyNext = func(i int) (int, int) {
			x := p.X + xStep*i
			yFloat := lineSlope*float64(x) + lineIntercept
			return x, int(math.Round(yFloat))
		}
	} else {
		path = make([]Position, 0, dyAbs)
		xyNext = func(i int) (int, int) {
			y := p.Y + yStep*i
			xFloat := (float64(y) - lineIntercept) / lineSlope
			return int(math.Round(xFloat)), y
		}
	}

	for i := 1; i < cap(path); i++ {
		x, y := xyNext(i)
		if !appendToPath(NewPosition(x, y)) {
			break
		}
	}

	return
}

func (p Position) String() string {
	return fmt.Sprintf("Position(%d, %d)", p.X, p.Y)
}

// PositionsInCircle returns new Positions in the circle area.
func PositionsInCircle(x, y, r int, includeXY bool) []Position {
	if r <= 1 {
		return []Position{NewPosition(x, y)}
	}

	var positions []Position
	r -= 1
	for xC := x - r; xC <= x+r; xC++ {
		for yC := y - r; yC <= y+r; yC++ {
			posDistance := math.Sqrt(math.Pow(float64(xC-x), 2) + math.Pow(float64(yC-y), 2))
			if posDistance > float64(r) {
				continue
			}
			if !includeXY && xC == x && yC == y {
				continue
			}

			positions = append(positions, NewPosition(xC, yC))
		}
	}

	return positions
}
