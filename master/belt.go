package main

// Belt pathfinding start at the goal and pathfinds towards the start
// Tiles have direction which is the final direction of the belt

import (
	"fmt"
	"github.com/beefsack/go-astar"
	"math"
)

type Tile struct {
	Pos   Position
	Ug    bool // undereground belt flags
	UgLen int  // only set for exit underground belts
	Dir   int
}

const (
	UgBeltLen = 6
)

func (t *Tile) Up() Tile    { return Tile{Position{t.Pos.X, t.Pos.Y - 1}, false, 0, dirSouth} }
func (t *Tile) Right() Tile { return Tile{Position{t.Pos.X + 1, t.Pos.Y}, false, 0, dirWest} }
func (t *Tile) Down() Tile  { return Tile{Position{t.Pos.X, t.Pos.Y + 1}, false, 0, dirNorth} }
func (t *Tile) Left() Tile  { return Tile{Position{t.Pos.X - 1, t.Pos.Y}, false, 0, dirEast} }

// Generates underground enters that exit on this tile with this direction
func (t *Tile) Underground() []Tile {
	out := []Tile{}

	xo := 0.0
	yo := 0.0

	if t.Dir == dirSouth {
		yo = -1
	} else if t.Dir == dirWest {
		xo = 1
	} else if t.Dir == dirNorth {
		yo = 1
	} else if t.Dir == dirEast {
		xo = -1
	}

	for i := 2; i < UgBeltLen; i++ {
		out = append(out, Tile{Position{t.Pos.X + xo*float64(i), t.Pos.Y + yo*float64(i)}, true, i, t.Dir})
	}
	return out
}

func (t *Tile) IsEmpty() bool {
	return bot.Mapper.isTileEmpty(t.Pos)
}

func (t Tile) PathNeighborCost(to astar.Pather) float64 {
	if t.Dir != to.(Tile).Dir { // prefer straight lines
		return 2
	}
	return 1
}

func (t Tile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(Tile)
	absX := math.Abs(toT.Pos.X - t.Pos.X)
	absY := math.Abs(toT.Pos.Y - t.Pos.Y)
	return float64(absX + absY)
}

func (from Tile) ValidNeighbor(to Tile) bool {
	if from.Ug || to.Ug { // if the belt is supposed to be after an exit from undergound
		if from.Dir != to.Dir { // its direction has to be the same
			return false
		}
	}
	if from.Ug && to.Ug { // cant go from ug straight into another ug
		return false
	}
	if !to.IsEmpty() { // if the position is not empty. This is the most expensive check so it's at the end
		return false
	}
	return true
}

func (t Tile) PathNeighbors() []astar.Pather {
	out := []astar.Pather{}
	for _, n := range t.Underground() {
		if t.ValidNeighbor(n) {
			out = append(out, n)
		}
	}
	for _, n := range []astar.Pather{t.Up(), t.Right(), t.Down(), t.Left()} {
		if t.ValidNeighbor(n.(Tile)) {
			out = append(out, n)
		}
	}
	return out
}

func (m *Mapper) FindBeltPath(from, to Position) []Tile {
	fromTile := Tile{from, false, 0, 0}
	toTile := Tile{to, false, 0, 0}
	path, _, found := astar.Path(toTile, fromTile)
	if !found {
		fmt.Println("Could not find a path")
		return nil
	}
	tilePath := []Tile{}
	for _, v := range path[1:len(path)-1] {
		tilePath = append(tilePath, v.(Tile))
	}
	return tilePath
}

func (m *Mapper) TileArrayToBP(tiles []Tile) []Building {
	bp := []Building{}
	wasLastUg := false
	for _, t := range tiles {
		typ := "input"
		name := "transport-belt"
		if t.Ug {
			name = "underground-belt"
			wasLastUg = true
		} else if wasLastUg {
			name = "underground-belt"
			wasLastUg = false
			typ = "output"
		}
		bp = append(bp, Building{
			Name:     name,
			Rotation: t.Dir,
			Pos:      t.Pos,
			Ugbt:     typ,
		})
	}
	return bp
}
