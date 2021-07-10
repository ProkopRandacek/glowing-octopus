package main

import "math"

func (b *box) round() {
	b.Tl.X = math.Round(b.Tl.X)
	b.Tl.Y = math.Round(b.Tl.Y)
	b.Br.X = math.Round(b.Br.X)
	b.Br.Y = math.Round(b.Br.Y)
}

func find(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func contains(s []position, e position) bool {
	for _, a := range s {
		if a.X == e.X && a.Y == e.Y {
			return true
		}
	}
	return false
}

func makeBox(a, b, c, d float64) box {
	return box{position{a, b}, position{c, d}}
}

func removeVal(s []position, e position) []position {
	for i, a := range s {
		if a.X == e.X && a.Y == e.Y {
			return remove(s, i)
		}
	}
	return s
}

func remove(slice []position, s int) []position {
	return append(slice[:s], slice[s+1:]...)
}
