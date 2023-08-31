package raytracing

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

func WriteImage(path string, img image.Image) {
	output, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	if err := png.Encode(output, img); err != nil {
		panic(err)
	}
}

func HitWorld(world []Sphere, r Ray, tMin, tMax float64) (HitRecord, bool) {
	closestSoFar := tMax
	var hitRecord HitRecord
	var wasHit bool
	for i := range world {
		sphere := world[i]
		if hit, ok := sphere.Hit(r, tMin, closestSoFar); ok {
			closestSoFar = hit.T
			hitRecord = hit
			wasHit = true
		}
	}
	return hitRecord, wasHit
}

func Clamp[T ~float32 | ~float64 | ~int](v T) T {
	switch {
	case v < 0:
		return 0
	case v > 1:
		return 1
	default:
		return v
	}
}

func RayColor(r Ray, scene Config, lights []Sphere, depth, maxDepth uint) Pixel {
	if depth <= 0 {
		return Pixel{}
	}
	hit, ok := HitWorld(scene.Objects, r, 0.001, math.MaxFloat64)
	if ok {
		scatteredRay, ok, albedo := hit.Texture.Scatter(r, hit)
		if ok {
			r, g, b, prob := 0., 0., 0., 0.1
			switch hit.Texture.(type) {
			case Light:
				prob = 0.05
			}
			lightsLen := float64(len(lights))
			if len(lights) > 0 &&
				rand.Float64() > (1-lightsLen*prob) &&
				depth > (maxDepth-2) {
				for i := range lights {
					light := lights[i]
					lightRay := Ray{
						Origin:    hit.Point,
						Direction: light.Center.Sub(hit.Point),
					}
					targetColor := RayColor(lightRay, scene, lights, 2, 5)
					r += albedo[0] * targetColor[0]
					g += albedo[1] * targetColor[1]
					b += albedo[2] * targetColor[2]
				}
				r /= lightsLen
				g /= lightsLen
				b /= lightsLen
			}
			targetColor := RayColor(scatteredRay, scene, lights, maxDepth, depth-1)
			return Pixel{
				Clamp(r + albedo[0]*targetColor[0]),
				Clamp(g + albedo[1]*targetColor[1]),
				Clamp(b + albedo[2]*targetColor[2]),
			}
		}
		return albedo
	}
	t := Clamp(0.5*r.Direction.Unit()[1] + 1)
	u := Clamp(0.5*r.Direction.Unit()[0] + 1)
	if sky := scene.Sky; sky != nil {
		x := u * float64(sky.Width-1)
		y := ((1 - t) * float64(sky.Height-1))
		r, g, b, a := sky.Texture.Image.At(int(x), int(y)).RGBA()
		return Pixel{
			float64(r) / float64(a),
			float64(g) / float64(a),
			float64(b) / float64(a),
		}
	} else {
		return Pixel{
			(1-t)*1 + t*0.5,
			(1-t)*1 + t*0.7,
			(1-t)*1 + t*1.0,
		}
	}
}

func RenderLine(img *image.RGBA, scene Config, lights []Sphere, y int) {
	width, height := scene.Width, scene.Height
	scale := 1. / float64(scene.SamplesPerPixel)
	fy := float64(y)
	for x := 0; x < int(width); x++ {
		fx := float64(x)
		var px Pixel
		for sample := 0; sample < int(scene.SamplesPerPixel); sample++ {
			u := (fx + rand.Float64()) / float64(width-1)
			v := (float64(height) - (fy + rand.Float64())) / float64(height-1)
			r := scene.Camera.GetRay(u, v)
			c := RayColor(r, scene, lights, scene.MaxDepth, scene.MaxDepth)
			px[0] += c[0]
			px[1] += c[1]
			px[2] += c[2]
		}
		color := color.RGBA{
			R: uint8(px[0] * scale * 255.),
			G: uint8(px[1] * scale * 255.),
			B: uint8(px[2] * scale * 255.),
			A: 255,
		}
		img.Set(x, y, color)
	}
}

func Render(path string, scene Config) {
	var lights []Sphere
	var objects []Sphere
	for i := range scene.Objects {
		if _, ok := scene.Objects[i].Scatterable.(Light); ok {
			lights = append(lights, scene.Objects[i])
		} else {
			objects = append(objects, scene.Objects[i])
		}
	}
	scene.Objects = objects
	img := image.NewRGBA(image.Rect(0, 0, int(scene.Width), int(scene.Height)))
	var wg sync.WaitGroup
	t1 := time.Now()
	for y := 0; y < int(scene.Height); y++ {
		wg.Add(1)
		go func(y int) {
			defer wg.Done()
			RenderLine(img, scene, lights, y)
		}(y)
	}
	wg.Wait()
	fmt.Printf("Frame time: %s", time.Since(t1))
	WriteImage(path, img)
}
