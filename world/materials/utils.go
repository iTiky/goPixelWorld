package materials

import (
	"image/color"
	"math"

	"github.com/itiky/goPixelWorld/world/types"
)

// FindClosestMaterialByColor returns the closest by color Material.
func FindClosestMaterialByColor(c color.Color) types.Material {
	c1R, c1G, c1B, _ := c.RGBA()

	calcCorrelation := func(c2 color.Color) float64 {
		c2R, c2G, c2B, _ := c2.RGBA()

		r := math.Pow(float64(c1R-c2R), 2)
		g := math.Pow(float64(c1G-c2G), 2)
		b := math.Pow(float64(c1B-c2B), 2)

		return math.Sqrt(r + g + b)
	}

	bestType := types.MaterialTypeNone
	bestCorrelation := calcCorrelation(color.Black)

	for _, m := range AllMaterialsSet {
		correlation := calcCorrelation(m.Color())
		if correlation < bestCorrelation {
			bestType = m.Type()
			bestCorrelation = correlation
		}
	}

	if bestType != types.MaterialTypeNone {
		return AllMaterialsSet[bestType]
	}

	return nil
}
