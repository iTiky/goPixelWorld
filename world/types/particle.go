package types

import (
	"fmt"
	"image/color"

	"github.com/itiky/goPixelWorld/pkg"
)

var lastParticleID uint64 = 0

type Particle struct {
	id        uint64
	material  Material
	forceVec  pkg.Vector
	health    float64
	steadyCnt int
}

func NewParticle(material Material) *Particle {
	return &Particle{
		id:       nextParticleID(),
		material: material,
		forceVec: pkg.NewVector(0, 0),
		health:   material.InitialHealth(),
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

func (p *Particle) HandleMoved() {
	p.steadyCnt = 0
}

func (p *Particle) UpdateState() {
	p.steadyCnt++
	p.limitForce()
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
