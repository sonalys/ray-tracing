package raytracing

import (
	"math"
	"math/rand"
)

type (
	Point3D [3]float64
)

func (self Point3D) Distance(other Point3D) float64 {
	dx := self[0] - other[0]
	dy := self[1] - other[1]
	dz := self[2] - other[2]
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (self Point3D) Length() float64 {
	return self.Distance(Point3D{})
}

func (self Point3D) LengthSquared() float64 {
	return self[0]*self[0] + self[1]*self[1] + self[2]*self[2]
}

func (self Point3D) Unit() Point3D {
	l := self.Length()
	return Point3D{
		self[0] / l,
		self[1] / l,
		self[2] / l,
	}
}

func (self Point3D) Multiply(v float64) Point3D {
	return Point3D{
		self[0] * v,
		self[1] * v,
		self[2] * v,
	}
}

func (self Point3D) Dot(other Point3D) float64 {
	return self[0]*other[0] +
		self[1]*other[1] +
		self[2]*other[2]

}

func (self Point3D) Cross(other Point3D) Point3D {
	return Point3D{
		self[1]*other[2] - self[2]*other[1],
		self[2]*other[0] - self[0]*other[2],
		self[0]*other[1] - self[1]*other[0],
	}
}

func (self Point3D) Add(other Point3D) Point3D {
	return Point3D{
		self[0] + other[0],
		self[1] + other[1],
		self[2] + other[2],
	}
}

func (self Point3D) Sub(other Point3D) Point3D {
	return Point3D{
		self[0] - other[0],
		self[1] - other[1],
		self[2] - other[2],
	}
}

func (self Point3D) Divide(v float64) Point3D {
	return Point3D{
		self[0] / v,
		self[1] / v,
		self[2] / v,
	}
}

func (self Point3D) NearZero() bool {
	const epsilon = 0.0001
	return self[0] < epsilon && self[1] < epsilon && self[2] < epsilon
}

func DegreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func RandomPoint(min, max float64) Point3D {
	return Point3D{
		RandFloat64Range(min, max),
		RandFloat64Range(min, max),
		RandFloat64Range(min, max),
	}
}

func RandomInUnitSphere() Point3D {
	for {
		p := RandomPoint(-1, 1)
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

func RandFloat64Range(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
