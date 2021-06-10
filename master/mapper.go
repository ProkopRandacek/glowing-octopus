package main

import (
	"math"
)

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

type Area struct {
	Dims Box
	Id int
}

var allocIdCounter = 0

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
					move := Position{pos.X + ox, pos.Y + oy}            // move in that direction
					if contains(tiles, move) && !contains(seen, move) { // if that position is a ore and not visited already
						q = append(q, move)            // add it to queue
						tiles = removeVal(tiles, move) // and remove it from ore list
					}
				}
			}
			seen = append(seen, pos)
			maxx = math.Max(maxx, pos.X) // update box checks
			maxy = math.Max(maxy, pos.Y)
			minx = math.Min(minx, pos.X)
			miny = math.Min(miny, pos.Y)
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

// returns true, if dims is available
func (b *Bot) canAlloc(dims Box) bool {
	for _, a := range b.Areas {
		d := a.Dims

		if d.Br.X >= dims.Tl.X && d.Br.Y >= dims.Tl.Y && d.Tl.X <= dims.Br.X && d.Tl.Y <= dims.Br.Y {
			return false
		}
	}
	return true
}

// allocates area and returns it's id. Returns -1 if area not available.
func (b *Bot) alloc(dims Box) int {
	if !b.canAlloc(dims) {
		return -1
	}

	b.Areas = append(b.Areas, Area{dims, allocIdCounter})
	allocIdCounter++

	b.drawBox(dims, Color{0, 1, 1})

	return allocIdCounter-1
}

// frees area by id. Returns true, if successful
func (b *Bot) free(id int) bool {
	for i, v := range b.Areas {
		if v.Id == id {
			b.Areas = append(b.Areas[:i], b.Areas[i+1:]...)
			return true
		}
	}

	return false
}
