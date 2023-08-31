package raytracing

import "math"

type (
	Sphere struct {
		Center Point3D
		Radius float64
		Scatterable
	}
)

func UVHitpointOnSphere(hitpoint Point3D) (u, v float64) {
	n := hitpoint.Unit()
	return math.Atan2(n[0], n[2])/(2*math.Pi) + 0.5, n[1]*0.5 + 0.5
}

func (self Sphere) Hit(r Ray, tMin, tMax float64) (hit HitRecord, ok bool) {
	oc := r.Origin.Sub(self.Center)
	a := r.Direction.LengthSquared()
	halfB := oc.Dot(r.Direction)
	c := oc.LengthSquared() - self.Radius*self.Radius
	discriminant := halfB*halfB - (a * c)

	if discriminant >= 0.0 {
		sqrtD := math.Sqrt(discriminant)
		rootA := ((-halfB) - sqrtD) / a
		rootB := ((-halfB) + sqrtD) / a
		for i := rootA; i <= rootB; i++ {
			if i >= tMax || i <= tMin {
				continue
			}
			p := r.At(i)
			hp := p.Sub(self.Center)
			normal := hp.Divide(self.Radius)
			frontFace := r.Direction.Dot(normal) < 0.0
			u, v := UVHitpointOnSphere(hp)

			var hNormal Point3D
			if frontFace {
				hNormal = normal
			} else {
				hNormal = normal.Multiply(-1)
			}
			return HitRecord{
				T:         i,
				Point:     p,
				Normal:    hNormal,
				FrontFace: frontFace,
				Texture:   self.Scatterable,
				U:         u,
				V:         v,
			}, true
		}
	}
	return HitRecord{}, false
}
