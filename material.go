package raytracing

import (
	"math"
	"math/rand"
)

type (
	Scatterable interface {
		Scatter(r Ray, hr HitRecord) (Ray, bool, Pixel)
	}

	Light      struct{}
	Lambertian struct {
		Albedo Pixel
	}
	Metal struct {
		Albedo Pixel
		Fuzz   float64
	}
	Glass struct {
		IndexOfRefraction float64
	}
	TextureMaterial struct {
		Albedo Pixel
		Texture
		HOffset float64
	}
)

func (self Light) Scatter(r Ray, hr HitRecord) (Ray, bool, Pixel) {
	return Ray{}, true, Pixel{1, 1, 1}
}

func (self Lambertian) Scatter(r Ray, hr HitRecord) (Ray, bool, Pixel) {
	scatterDirection := hr.Normal.Add(RandomInUnitSphere())
	if scatterDirection.NearZero() {
		scatterDirection = hr.Normal
	}
	target := hr.Point.Add(scatterDirection)
	scattered := Ray{
		Origin:    hr.Point,
		Direction: target.Sub(hr.Point),
	}
	return scattered, true, self.Albedo
}

func (self Metal) Scatter(r Ray, hr HitRecord) (Ray, bool, Pixel) {
	reflected := Reflect(r.Direction, hr.Normal)
	scattered := Ray{
		Origin:    hr.Point,
		Direction: reflected.Add(RandomInUnitSphere()).Multiply(self.Fuzz),
	}
	if scattered.Direction.Dot(hr.Normal) > 0 {
		return scattered, true, self.Albedo
	}
	return Ray{}, false, Pixel{}
}

func (self Glass) Scatter(r Ray, hr HitRecord) (Ray, bool, Pixel) {
	attenuation := Pixel{1, 1, 1}
	var refractionRatio float64
	if hr.FrontFace {
		refractionRatio = 1 / self.IndexOfRefraction
	} else {
		refractionRatio = self.IndexOfRefraction
	}
	unitDirection := r.Direction.Unit()
	cosTheta := math.Min(unitDirection.Multiply(-1).Dot(hr.Normal), 1)
	sinTheta := math.Sqrt((1 - cosTheta*cosTheta))
	cannotRefract := refractionRatio*sinTheta > 1
	if cannotRefract || Reflectance(cosTheta, refractionRatio) > rand.Float64() {
		reflected := Reflect(unitDirection, hr.Normal)
		scattered := Ray{
			Origin:    hr.Point,
			Direction: reflected,
		}
		return scattered, true, attenuation
	} else {
		direction := Refract(unitDirection, hr.Normal, refractionRatio)
		scattered := Ray{
			Origin:    hr.Point,
			Direction: direction,
		}
		return scattered, true, attenuation
	}
}

func NewTexturedMaterial(albedo Pixel, path string, rot float64) TextureMaterial {
	texture := LoadTexture(path)
	return TextureMaterial{
		Albedo:  albedo,
		Texture: texture,
		HOffset: rot,
	}
}

func (self TextureMaterial) GetAlbedo(u, v float64) Pixel {
	rot := u + self.HOffset
	if rot > 1 {
		rot = rot - 1
	}

	uu := rot * float64(self.Width)
	vv := (1. - v) * float64(self.Height-1.)

	x := int(math.Floor(uu))
	y := int(math.Floor(vv))
	r, g, b, a := self.Image.At(x, y).RGBA()
	return Pixel{
		float64(r) / float64(a),
		float64(g) / float64(a),
		float64(b) / float64(a),
	}
}

func (self TextureMaterial) Scatter(r Ray, hr HitRecord) (Ray, bool, Pixel) {
	scatterDirection := hr.Normal.Add(RandomInUnitSphere())
	if scatterDirection.NearZero() {
		scatterDirection = hr.Normal
	}
	target := hr.Point.Add(scatterDirection)
	scattered := Ray{
		Origin:    hr.Point,
		Direction: target.Sub(hr.Point),
	}
	return scattered, true, self.GetAlbedo(hr.U, hr.V)
}
