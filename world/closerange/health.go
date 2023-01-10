package closerange

import (
	"github.com/itiky/goPixelWorld/world/types"
)

func (e *Environment) ReduceHealth(step float64) bool {
	e.actions = append(e.actions, types.NewReduceHealth(e.source.Pos, step))
	return true
}
