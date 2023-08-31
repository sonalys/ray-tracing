package raytracing

import "math"

type (
	CameraParams struct {
		LookFrom    Point3D
		LookAt      Point3D
		VUP         Point3D
		VerticalFOV float64
		Aspect      float64
	}

	Camera struct {
		Origin          Point3D
		LowerLeftCorner Point3D
		FocalLength     float64
		Horizontal      Point3D
		Vertical        Point3D
		CameraParams
	}
)

func NewCamera(params CameraParams) Camera {
	theta := DegreesToRadians(params.VerticalFOV)
	halfHeight := math.Tan(theta / 2)
	halfWidth := params.Aspect * halfHeight

	w := params.LookFrom.Sub(params.LookAt).Unit()
	u := params.VUP.Cross(w).Unit()
	v := w.Cross(u)
	return Camera{
		Origin:          params.LookFrom,
		LowerLeftCorner: params.LookFrom.Sub(u.Multiply(halfWidth)).Sub(v.Multiply(halfHeight)).Sub(w),
		Horizontal:      u.Multiply(2 * halfWidth),
		Vertical:        v.Multiply(2 * halfHeight),
		FocalLength:     params.LookFrom.Sub(params.LookAt).Length(),
		CameraParams:    params,
	}
}

func (self Camera) GetRay(u, v float64) Ray {
	return Ray{
		Origin:    self.Origin,
		Direction: self.LowerLeftCorner.Add(self.Horizontal.Multiply(u)).Add(self.Vertical.Multiply(v)).Sub(self.Origin),
	}
}
