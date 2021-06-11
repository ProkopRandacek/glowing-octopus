package main

type Color struct {
	R, G, B float32
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Box struct {
	Tl Position // top left
	Br Position // bottom right
}
