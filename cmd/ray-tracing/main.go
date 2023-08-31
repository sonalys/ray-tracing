package main

import raytracing "github.com/sonalys/ray-tracing"

func main() {
	scene := raytracing.Config{
		Width:           1920,
		Height:          1080,
		SamplesPerPixel: 100,
		MaxDepth:        500,
		Sky: &raytracing.Sky{
			TextureMaterial: raytracing.NewTexturedMaterial(raytracing.Pixel{
				0, 0, 1,
			}, "./data/beach.jpg", 0),
		},
		Objects: []raytracing.Sphere{
			// {
			// 	Center:      raytracing.Point3D{2, 4, 2},
			// 	Radius:      1,
			// 	Scatterable: raytracing.Light{},
			// },
			{
				Center: raytracing.Point3D{2, 2, 2},
				Radius: 1,
				Scatterable: raytracing.NewTexturedMaterial(raytracing.Pixel{
					0, 0, 0,
				}, "./data/earth.jpg", 0),
			},
			{
				Center: raytracing.Point3D{3.2, 3.2, 3.2},
				Radius: 0.25,
				Scatterable: raytracing.Glass{
					IndexOfRefraction: 1.217,
				},
			},
		},
		Camera: raytracing.NewCamera(raytracing.CameraParams{
			LookFrom:    raytracing.Point3D{4, 4, 4},
			LookAt:      raytracing.Point3D{0, 0, 0},
			VerticalFOV: 50,
			Aspect:      16. / 9.,
			VUP:         raytracing.Point3D{0, 1, 0},
		}),
	}

	raytracing.Render("out.png", scene)
}
