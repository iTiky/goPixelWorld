package types

import (
	"fmt"
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
)

// lastParticleID is the last unique Particle ID.
var lastParticleID uint64 = 0

// Particle internal state parameter keys.
const (
	// ParticleStateParamSteady parameter defines a counter which is incremented if a Particle didn't move.
	ParticleStateParamSteady = "steady"
)

type (
	// Particle holds a single world object state (a pixel).
	Particle struct {
		id       uint64        // unique ID
		material Material      // Particle's Material
		forceVec pkg.Vector    // current force Vector
		health   float64       // current health state
		state    ParticleState // map of internal state counters
	}

	ParticleState map[string]int
)

// NewParticle creates a new Particle.
func NewParticle(material Material) *Particle {
	return &Particle{
		id:       nextParticleID(),
		material: material,
		forceVec: pkg.NewVector(0, 0),
		health:   material.InitialHealth(),
		state:    make(ParticleState),
	}
}

// ID returns the unique ID.
func (p *Particle) ID() uint64 {
	return p.id
}

// Material returns the Material.
func (p *Particle) Material() Material {
	return p.material
}

// ForceVector returns the current force Vector.
func (p *Particle) ForceVector() pkg.Vector {
	return p.forceVec
}

// Health returns the current health state.
func (p *Particle) Health() float64 {
	return p.health
}

// IsDestroyed returns true if the Particle should be destroyed.
func (p *Particle) IsDestroyed() bool {
	return p.health <= 0
}

// Color returns the current color adjusted by its health state.
// Not all Materials support that.
func (p *Particle) Color() color.Color {
	return p.material.ColorAdjusted(p.health)
}

// SetStateParam sets the internal state parameter by key.
func (p *Particle) SetStateParam(key string, value int) {
	p.state[key] = value
}

// GetStateParam returns the internal state parameter by key.
func (p *Particle) GetStateParam(key string) int {
	return p.state[key]
}

// IncStateParam increments the internal state parameter by key.
func (p *Particle) IncStateParam(key string) int {
	v := p.state[key]
	v++
	p.state[key] = v

	return v
}

// DecStateParam decrements the internal state parameter by key.
func (p *Particle) DecStateParam(key string) {
	p.state[key] = p.state[key] - 1
}

// OnMove updates the internal state, dropping its steady param to 0 (it is moving right now).
func (p *Particle) OnMove() {
	p.SetStateParam(ParticleStateParamSteady, 0)
}

// UpdateState updates the internal state parameters.
// Drops the force Vector if Particle is not moving (solves a huge accumulated force issue).
func (p *Particle) UpdateState() {
	steadyCnt := p.IncStateParam(ParticleStateParamSteady)
	if steadyCnt > 10 {
		p.SetStateParam(ParticleStateParamSteady, 0)
		p.forceVec = pkg.NewVector(0, 0)
	}
}

// AddForce adds a new Vector to the force Vector.
func (p *Particle) AddForce(force pkg.Vector) {
	p.forceVec = p.forceVec.Add(force)
	p.limitForce()
}

// MultiplyForce alters the force Vector magnitude.
func (p *Particle) MultiplyForce(k float64) {
	p.forceVec = p.forceVec.MultiplyByK(k)
	p.limitForce()
}

// RotateForce rotates the force Vector by angle.
func (p *Particle) RotateForce(angleRad float64) {
	p.forceVec = p.forceVec.Rotate(angleRad)
}

// ReflectForce rotates (180 degrees) the force Vector.
func (p *Particle) ReflectForce(horizontal, vertical bool) {
	p.forceVec = p.forceVec.Reflect(horizontal, vertical)
}

// SetForce sets the force Vector.
func (p *Particle) SetForce(force pkg.Vector) {
	p.forceVec = force
}

// ReduceHealth alters the current health state.
func (p *Particle) ReduceHealth(step float64) {
	p.health -= step
}

func (p *Particle) String() string {
	return fmt.Sprintf("Particle(id=%d, material=%T, force=%s, health=%f)", p.id, p.material, p.forceVec, p.health)
}

// limitForce limits the force Vector.
func (p *Particle) limitForce() {
	const maxForce = 10.0

	if m := p.forceVec.Magnitude(); m > maxForce {
		p.forceVec = p.forceVec.SetMagnitude(m / 2.0)
	}
}

// nextParticleID returns the next unique Particle ID.
func nextParticleID() uint64 {
	lastParticleID++

	return lastParticleID
}
