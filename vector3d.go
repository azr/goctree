package goctree

import "fmt"

type Vector3D [3]float64

const (
	x = 0
	y = 1
	z = 2
)

func (v Vector3D) String() string {
	return fmt.Sprintf("[%5.2f,%5.2f,%5.2f]", v[x], v[y], v[z])
}

func (p Vector3D) Add(b Vector3D) Vector3D { return Vector3D{p[x] + b[x], p[y] + b[y], p[z] + b[z]} }

func (p Vector3D) Sub(b Vector3D) Vector3D { return Vector3D{p[x] - b[x], p[y] - b[y], p[z] - b[z]} }

func (p Vector3D) Imul(b float64) Vector3D { return Vector3D{p[x] * b, p[y] * b, p[z] * b} }

func (p Vector3D) Sqd(q Vector3D) float64 {
	var sum float64
	for dim, pCoord := range p {
		d := pCoord - q[dim]
		sum += d * d
	}
	return sum
}
