package materials

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/itiky/goPixelWorld/pkg"
	"github.com/itiky/goPixelWorld/world/types"
)

var _ types.Material = Bug{}

// Bug ...
type Bug struct {
	base
	foodDampHealthStep   float64 // grass surrounding feed speed (the same amount of health is recovered)
	movementSpeedMagStep float64 // movement speed increment
	movementSpeedMagMax  float64 // max movement speed limit
	healthToSplit        float64 // the amount of health to split (create a new Bug)
	healthAfterSplit     float64 // the amount of health set after the split
}

const (
	BugStateParamMovementDir = "bug_move_dir"
)

func NewBug() Bug {
	return Bug{
		base: newBase(
			types.MaterialTypeBug,
			color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF},
			withFlags(types.MaterialFlagIsFlammable),
			withCloseRangeType(types.MaterialCloseRangeTypeInCircleRange),
			withCloseRangeCircleR(5),
			withMass(50.0),
			withSelfHealthReduction(100.0, 0.6),
			withSourceDamping(0.5, 0.0),
		),
		foodDampHealthStep:   1.0,
		movementSpeedMagStep: 0.2,
		movementSpeedMagMax:  1.8,
		healthToSplit:        1000.0,
		healthAfterSplit:     20.0,
	}
}

func (m Bug) ColorAdjusted(health float64) color.Color {
	if health < 20.0 {
		return color.RGBA{R: 0x9B, G: 0x71, B: 0x77, A: 0x35}
	} else if health < 40.0 {
		return color.RGBA{R: 0xC5, G: 0xAE, B: 0xA1, A: 0x5C}
	} else if health < 60.0 {
		return color.RGBA{R: 0xE1, G: 0xD0, B: 0xCD, A: 0x6F}
	} else if health < 80.0 {
		return color.RGBA{R: 0xEA, G: 0xEE, B: 0xDD, A: 0xFF}
	}

	return m.baseColor
}

func (m Bug) ProcessInternal(env types.TileEnvironment) {
	var moveDir pkg.Direction

	m.commonProcessInternal(env)

	appendForce := func() {
		if env.ForceVec().Magnitude() >= m.movementSpeedMagMax {
			return
		}
		moveVec := pkg.NewVector(m.movementSpeedMagStep, moveDir.Angle())
		env.AddSelfForce(moveVec)
	}

	setForce := func() {
		moveVec := pkg.NewVector(env.ForceVec().Magnitude()+m.movementSpeedMagStep, moveDir.Angle())
		if moveVec.Magnitude() >= m.movementSpeedMagMax {
			moveVec = moveVec.SetMagnitude(m.movementSpeedMagMax)
		}
		env.SetSelfForce(moveVec)
	}

	dropForce := func() {
		moveVec := pkg.NewVector(0, 0)
		env.SetSelfForce(moveVec)
	}

	// Add gravity only if we are not jumping from a stacked position
	didJump := false
	defer func() {
		if didJump {
			return
		}
		env.AddGravity()
	}()

	// Reduce health
	env.DampSelfHealth(m.selfHealthDampStep)
	if env.Health() <= 0.0 {
		if pkg.FlipCoin() {
			env.ReplaceSelf(AllMaterialsSet[types.MaterialTypeWater])
		}
		return
	}

	// Time to split
	if curHealth := env.Health(); curHealth >= m.healthToSplit {
		env.DampSelfHealth(curHealth - m.healthAfterSplit)
		env.AddNewTileInRange(AllMaterialsSet[types.MaterialTypeBug])
	}

	// Feed if there is some food around
	foodTilesInCloseCnt := env.DampEnvHealthByTypeInRange(math.Sqrt2, m.foodDampHealthStep, []types.MaterialType{types.MaterialTypeGrass}, nil)
	if foodTilesInCloseCnt > 0 {
		// Drop the movement intention (everything is OK, food found, no need to move)
		moveDir = 0
		env.UpdateStateParam(BugStateParamMovementDir, pkg.DirectionNone.Int())
		dropForce()
		// Adjust our health
		env.DampSelfHealth(-m.foodDampHealthStep * float64(foodTilesInCloseCnt))
		return
	}
	// We need to find some

	// Get the current movement intention
	moveDir = pkg.Direction(env.StateParam(BugStateParamMovementDir))

	// We are willing to move
	if moveDir != pkg.DirectionNone {
		// Accelerate the current movement to oppose the gravity
		appendForce()

		// But we are blocked, try to adjust
		if steadyCnt := env.StateParam(types.ParticleStateParamSteady); steadyCnt > 5 {
			// Can we slightly change the movement direction? (excluding Top, no jumping)
			_, tilesDirs, _ := env.SearchTilesInRange(
				pkg.ValuePtr(true), pkg.ValuePtr(math.Sqrt2),
				moveDir.Sector(1), true,
				nil, false,
				nil, false,
			)
			tilesDirs = pkg.FilterSlice(tilesDirs, func(dir pkg.Direction) bool {
				return dir != pkg.DirectionTop
			})
			if len(tilesDirs) > 0 {
				emptyTileIdx := rand.Intn(len(tilesDirs))
				moveDir = tilesDirs[emptyTileIdx]
				env.UpdateStateParam(BugStateParamMovementDir, moveDir.Int())
				setForce()
				didJump = true
			} else {
				// Drop the movement intention
				moveDir = 0
				env.UpdateStateParam(BugStateParamMovementDir, pkg.DirectionNone.Int())
			}
		}
		return
	}

	// We are not moving, search for the next food target Tile to move in direction of
	foodTilesInDistance, _, _ := env.SearchTilesInRange(
		pkg.ValuePtr(false), pkg.ValuePtr(math.Sqrt2*3),
		nil, false,
		[]types.MaterialType{types.MaterialTypeGrass}, true,
		nil, false,
	)
	if len(foodTilesInDistance) > 0 {
		// Start moving to our new food target
		targetFoodTileIdx := rand.Intn(len(foodTilesInDistance))
		targetFoodTile := foodTilesInDistance[targetFoodTileIdx]
		moveDir = pkg.NewDirectionFromCoords(env.Position().X, env.Position().Y, targetFoodTile.Pos.X, targetFoodTile.Pos.Y)
		env.UpdateStateParam(BugStateParamMovementDir, moveDir.Int())
		setForce()
		return
	}

	// Start moving to a random direction (only if we are standing on smth hard)
	tiles, tilesDirs, _ := env.SearchTilesInRange(
		nil, pkg.ValuePtr(math.Sqrt2),
		[]pkg.Direction{pkg.DirectionTopRight, pkg.DirectionRight, pkg.DirectionBottomRight, pkg.DirectionBottom, pkg.DirectionBottomLeft, pkg.DirectionLeft, pkg.DirectionTopLeft}, true,
		nil, false,
		nil, false,
	)

	inAir := true
	var emptyDirs []pkg.Direction
	for i := 0; i < len(tiles); i++ {
		if tilesDirs[i] == pkg.DirectionBottom {
			if tiles[i].HasParticle() && !tiles[i].Particle.Material().IsFlagged(types.MaterialFlagIsLiquid) {
				inAir = false
			}
			continue
		}

		if !tiles[i].HasParticle() {
			emptyDirs = append(emptyDirs, tilesDirs[i])
		}
	}

	if !inAir && len(emptyDirs) > 0 {
		emptyDirIdx := rand.Intn(len(emptyDirs))
		moveDir = emptyDirs[emptyDirIdx]
		env.UpdateStateParam(BugStateParamMovementDir, moveDir.Int())
		appendForce()
	}
}

func (m Bug) ProcessCollision(env types.CollisionEnvironment) {
	env.ReflectSourceTargetForces(m.srcForceDamperK)
	env.DampSelfHealthByMassRate(m.selfHealthDampStep)
}
