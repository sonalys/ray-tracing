package raytracing

import "math"

type (
	Ray struct {
		Origin    Point3D
		Direction Point3D
	}

	HitRecord struct {
		T         float64
		Point     Point3D
		Normal    Point3D
		FrontFace bool
		Texture   Scatterable
		U, V      float64
	}

	Hittable interface {
		Hit(r Ray, tMin, tMax float64) (hit HitRecord, ok bool)
	}
)

func (self Ray) At(t float64) Point3D {
	return self.Origin.Add(self.Direction.Multiply(t))
}

func Reflect(v, n Point3D) Point3D {
	return v.Sub(n.Multiply(2 * v.Dot(n)))
}

func Refract(uv, n Point3D, ETAIOverETAT float64) Point3D {
	cosTheta := math.Min(uv.Multiply(-1).Dot(n), 1.0)
	rOutPerp := uv.Add(n.Multiply(cosTheta)).Multiply(ETAIOverETAT)
	rOutParallel := n.Multiply(-1 * math.Sqrt(math.Abs(1-rOutPerp.Length())))
	return rOutPerp.Add(rOutParallel)
}

func Reflectance(cosine, refIndex float64) float64 {
	r0 := (1 - refIndex) / (1 + refIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
