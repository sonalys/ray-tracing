package raytracing

import (
	"image"
	"image/jpeg"
	"os"
)

type (
	Pixel   [3]float64
	Texture struct {
		Image  image.Image
		Width  uint
		Height uint
		Path   string
	}

	Sky struct {
		TextureMaterial
	}
)

func LoadTexture(path string) Texture {
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()
	image.RegisterFormat("jpg", "?", jpeg.Decode, jpeg.DecodeConfig)
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err.Error())
	}
	bounds := img.Bounds()
	return Texture{
		Width:  uint(bounds.Dx()),
		Height: uint(bounds.Dy()),
		Path:   path,
		Image:  img,
	}
}
