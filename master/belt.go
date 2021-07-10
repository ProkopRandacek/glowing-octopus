package main

// Belt pathfinding start at the goal and pathfinds towards the start
// Tiles have direction which is the final direction of the belt

import (
	"fmt"
	"math"

	"github.com/beefsack/go-astar"
)

type tile struct {
	Pos   position
	Ug    bool // underground belt flags
	UgLen int  // only set for exit underground belts
	Dir   int
}

func (t *tile) direct() []tile {
	return []tile{
		{position{t.Pos.X, t.Pos.Y - 1}, false, 0, dirSouth},
		{position{t.Pos.X + 1, t.Pos.Y}, false, 0, dirWest},
		{position{t.Pos.X, t.Pos.Y + 1}, false, 0, dirNorth},
		{position{t.Pos.X - 1, t.Pos.Y}, false, 0, dirEast},
	}
}

// underground generates underground enters that exit on this tile with this direction
func (t *tile) underground(maxBeltLen int) []tile {
	var out []tile

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

	for i := 2; i < maxBeltLen; i++ {
		out = append(out, tile{position{t.Pos.X + xo*float64(i), t.Pos.Y + yo*float64(i)}, true, i, t.Dir})
	}
	return out
}

func (t *tile) isEmpty() bool {
	return fbot.Mapper.isTileEmpty(t.Pos)
}

//PathNeighborCost calculates the cost of the move from tile to its neighbour
func (t tile) PathNeighborCost(to astar.Pather) float64 {
	if t.Dir != to.(tile).Dir { // prefer straight lines
		return 2
	}
	return 1
}

//PathEstimatedCost estimates the distance to goal
func (t tile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(tile)
	absX := math.Abs(toT.Pos.X - t.Pos.X)
	absY := math.Abs(toT.Pos.Y - t.Pos.Y)
	return absX + absY
}

func (t tile) validNeighbor(to tile) bool {
	if t.Ug || to.Ug { // if the belt is supposed to be after an exit from underground
		if t.Dir != to.Dir { // its direction has to be the same
			return false
		}
	}
	if t.Ug && to.Ug { // cant go from ug straight into another ug
		return false
	}
	if !to.isEmpty() { // if the position is not empty. This is the most expensive check so it's at the end
		return false
	}
	return true
}

//PathNeighbors finds neighbors from given position
func (t tile) PathNeighbors() []astar.Pather {
	var out []astar.Pather

	maxBeltLen := 6
	if fbot.BeltLevel == "fast" {
		maxBeltLen = 8
	} else if fbot.BeltLevel == "express" {
		maxBeltLen = 10
	}

	for _, n := range append(t.underground(maxBeltLen), t.direct()...) {
		if t.validNeighbor(n) {
			out = append(out, n)
		}
	}
	return out
}

func (m *mapper) findBeltPath(from, to position) []tile {
	fromTile := tile{from, false, 0, 0}
	toTile := tile{to, false, 0, 0}
	path, _, found := astar.Path(toTile, fromTile)
	if !found {
		fmt.Println("Could not find a path")
		return nil
	}
	var tilePath []tile
	for _, v := range path[1 : len(path)-1] {
		tilePath = append(tilePath, v.(tile))
	}
	return tilePath
}

func (m *mapper) allocBelt(tiles []tile) {
	for _, t := range tiles {
		space := makeBox(t.Pos.X, t.Pos.Y, t.Pos.X+1, t.Pos.Y)
		m.forceAlloc(space) // We don't care if we alloc over (shouldn't ever happen actually) something or too close to something
		fbot.clearAll(space)
		fbot.waitForTaskDone()
	}
}
