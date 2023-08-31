package raytracing

type (
	Config struct {
		Width, Height   uint
		SamplesPerPixel uint32
		MaxDepth        uint
		*Sky
		Camera
		Objects []Sphere
	}
)
