package types

// InputActionType defines an input action type.
type InputActionType int

const (
	InputActionCreateParticles InputActionType = iota
	InputActionDeleteParticles
	InputActionFlipGravity
)

// InputAction defines a common input action interface.
type InputAction interface {
	Type() InputActionType
}

// CreateParticlesInputAction defines a request to creates a set of new Particles.
type CreateParticlesInputAction struct {
	X, Y       int       // Position to add new Particles
	Radius     int       // if GT 1, creates a set of Particles in a circle area
	Material   MaterialI // new Particle(s) Material
	ApplyForce bool      // if set, apply a random force to a new Particle(s).
}

func (a CreateParticlesInputAction) Type() InputActionType {
	return InputActionCreateParticles
}

// DeleteParticlesInputAction defines a request to delete a set of existing Particles.
type DeleteParticlesInputAction struct {
	X, Y   int // Position to remove existing Particles
	Radius int // if GT 1, remove a set of Particles in a circle area
}

func (a DeleteParticlesInputAction) Type() InputActionType {
	return InputActionDeleteParticles
}

// FlipGravityInputAction defines a request to flip the vertical gravity.
type FlipGravityInputAction struct{}

func (a FlipGravityInputAction) Type() InputActionType {
	return InputActionFlipGravity
}
