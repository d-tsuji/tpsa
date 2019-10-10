package tpsa

import "math"

type Point struct {
	X float64
	Y float64
}

func dist(p1 Point, p2 Point) float64 {
	cost := math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
	return cost
}
