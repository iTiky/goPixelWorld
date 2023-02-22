package types

import (
	"fmt"
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
)

var lastParticleID uint64 = 0

const (
	ParticleStateParamSteady = "steady"
)

type Particle struct {
	id        uint64
	material  Material
	forceVec  pkg.Vector
	health    float64
	steadyCnt int
	state     map[string]int
}

func NewParticle(material Material) *Particle {
	return &Particle{
		id:       nextParticleID(),
		material: material,
		forceVec: pkg.NewVector(0, 0),
		health:   material.InitialHealth(),
		state:    make(map[string]int),
	}
}

func (p *Particle) ID() uint64 {
	return p.id
}

func (p *Particle) Material() Material {
	return p.material
}

func (p *Particle) ForceVector() pkg.Vector {
	return p.forceVec
}

func (p *Particle) Health() float64 {
	return p.health
}

func (p *Particle) IsDestroyed() bool {
	return p.health <= 0
}

func (p *Particle) Color() color.Color {
	return p.material.ColorAdjusted(p.health)
}

func (p *Particle) SetStateParam(key string, value int) {
	p.state[key] = value
}

func (p *Particle) GetStateParam(key string) int {
	return p.state[key]
}

func (p *Particle) IncStateParam(key string) int {
	v := p.state[key]
	v++
	p.state[key] = v

	return v
}

func (p *Particle) DecStateParam(key string) {
	p.state[key] = p.state[key] - 1
}

func (p *Particle) OnMove() {
	p.SetStateParam(ParticleStateParamSteady, 0)
}

func (p *Particle) UpdateState() {
	steadyCnt := p.IncStateParam(ParticleStateParamSteady)
	if steadyCnt > 10 {
		p.SetStateParam(ParticleStateParamSteady, 0)
		p.forceVec = pkg.NewVector(0, 0)
	}
}

func (p *Particle) AddForce(force pkg.Vector) {
	p.forceVec = p.forceVec.Add(force)
	p.limitForce()
}

func (p *Particle) MultiplyForce(k float64) {
	p.forceVec = p.forceVec.MultiplyByK(k)
	p.limitForce()
}

func (p *Particle) RotateForce(angleRad float64) {
	p.forceVec = p.forceVec.Rotate(angleRad)
}

func (p *Particle) ReflectForce(horizontal, vertical bool) {
	p.forceVec = p.forceVec.Reflect(horizontal, vertical)
}

func (p *Particle) SetForce(force pkg.Vector) {
	p.forceVec = force
}

func (p *Particle) ReduceHealth(step float64) {
	p.health -= step
}

func (p *Particle) String() string {
	return fmt.Sprintf("Particle(id=%d, material=%T, force=%s, health=%f)", p.id, p.material, p.forceVec, p.health)
}

func (p *Particle) limitForce() {
	const maxForce = 10.0

	if m := p.forceVec.Magnitude(); m > maxForce {
		p.forceVec = p.forceVec.SetMagnitude(m / 2.0)
	}
}

func nextParticleID() uint64 {
	lastParticleID++

	return lastParticleID
}
