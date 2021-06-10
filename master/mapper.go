package main

import (
	"math"
)

type Position struct {
	x float64
	y float64
}

type Box struct {
	tl Position // top left
	br Position // bottom right
}

type ResourceMaps struct {
	iron          []Position
	copper        []Position
	stone         []Position
	coal          []Position
	ironPatches   []Box
	copperPatches []Box
	stonePatches  []Box
	coalPatches   []Box
}

func CalcBoxes(tiles []Position) (boxes []Box) {
	pos := Position{}
	for len(tiles) > 0 { // iterate over all positions
		pos, tiles = tiles[0], tiles[1:]                                               // pop firts value
		seen := []Position{pos}                                                        // make this a set?
		q := []Position{pos}                                                           // queue
		maxx, maxy, minx, miny := math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1) // component bounds
		for {
			pos, q = q[0], q[1:] // pop firts value
			for _, ox := range [3]float64{-1.0, 0.0, 1.0} {
				for _, oy := range [3]float64{-1.0, 0.0, 1.0} { // for each direction from our position
					move := Position{pos.x + ox, pos.y + oy}            // move in that direction
					if contains(tiles, move) && !contains(seen, move) { // if that position is a ore and not visited already
						q = append(q, move)            // add it to queue
						tiles = removeVal(tiles, move) // and remove it from ore list
					}
				}
			}
			seen = append(seen, pos)
			maxx = math.Max(maxx, pos.x) // update box checks
			maxy = math.Max(maxy, pos.y)
			minx = math.Min(minx, pos.x)
			miny = math.Min(miny, pos.y)
			if len(q) == 0 {
				break
			}
		}
		boxes = append(boxes, Box{Position{maxx, miny}, Position{minx, maxy}})
	}
	return
}

func (rm *ResourceMaps) CalcPatches() {
	rm.ironPatches = CalcBoxes(rm.iron)
	rm.copperPatches = CalcBoxes(rm.copper)
	rm.stonePatches = CalcBoxes(rm.stone)
	rm.coalPatches = CalcBoxes(rm.coal)
}

type Mapper struct {
	mallMap    map[string]Position
	allocMap   []Box
	waterTiles []Position // the exact tiles that are in water
	waterBox   Box        // water boxes for fast intersection check
	resMaps    ResourceMaps
	orePatches map[string][]Box
}

func (m *Mapper) readRawWorld(rw RawWorld) {
	m.waterTiles = make([]Position, len(rw.Water))
	m.resMaps.iron = make([]Position, len(rw.Iron))
	m.resMaps.copper = make([]Position, len(rw.Copper))
	m.resMaps.stone = make([]Position, len(rw.Stone))
	m.resMaps.coal = make([]Position, len(rw.Coal))
	for i, t := range rw.Water {
		m.waterTiles[i] = Position{t[0], t[1]}
	}
	for i, t := range rw.Iron {
		m.resMaps.iron[i] = Position{t[0], t[1]}
	}
	for i, t := range rw.Copper {
		m.resMaps.copper[i] = Position{t[0], t[1]}
	}
	for i, t := range rw.Stone {
		m.resMaps.stone[i] = Position{t[0], t[1]}
	}
	for i, t := range rw.Coal {
		m.resMaps.coal[i] = Position{t[0], t[1]}
	}
	m.resMaps.CalcPatches()
}
