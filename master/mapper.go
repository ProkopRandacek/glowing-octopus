package main

import (
	"math"
)

type OrePatch struct {
	Dims Box
	Unsafe bool
}

type ResourceMaps struct {
	Iron          []Position
	Copper        []Position
	Stone         []Position
	Coal          []Position
	IronPatches   []OrePatch
	CopperPatches []OrePatch
	StonePatches  []OrePatch
	CoalPatches   []OrePatch
}

type Area struct {
	Dims Box
	Id   int
}

type Mapper struct {
	Areas          []Area // marks what areas area allocated
	AllocIdCounter int    // for area id generation
	MallMap        map[string]Position
	AllocMap       []Box
	WaterTiles     []Position // the exact tiles that are in water
	WaterBox       Box        // water boxes for fast intersection check
	ResMaps        ResourceMaps
	OrePatches     map[string][]OrePatch
	LoadedBoxes []Box
}

func isBorderSharedWithBox(box Box, boxes []Box, ignore int) bool {
	for i, b := range boxes {
		if i == ignore {
			continue
		}

		if box.Tl.X == b.Tl.X ||
				box.Tl.Y == b.Tl.Y ||
				box.Br.X == b.Br.Y ||
				box.Br.Y == b.Br.Y {
			return true
		}
	}

	return false
}

func (o *OrePatch) isUnsafe(boxes []Box) {
	for i, b := range boxes {
		if isBorderSharedWithBox(o.Dims, []Box{b}, -1) && !isBorderSharedWithBox(b, boxes, i) {
			o.Unsafe = true
			return
		}
	}

	o.Unsafe = false
}

func CalcBoxes(tiles []Position) (boxes []OrePatch) {
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
		boxes = append(boxes, OrePatch{Box{Position{maxx, miny}, Position{minx, maxy}}, false})
	}
	return
}

func (rm *ResourceMaps) CalcPatches() {
	rm.IronPatches = CalcBoxes(rm.Iron)
	rm.CopperPatches = CalcBoxes(rm.Copper)
	rm.StonePatches = CalcBoxes(rm.Stone)
	rm.CoalPatches = CalcBoxes(rm.Coal)
}

func (m *Mapper) calcUnsafe() {
	
}

func (m *Mapper) readRawWorld(rw RawWorld) {
	m.WaterTiles = make([]Position, len(rw.Water))
	m.ResMaps.Iron = make([]Position, len(rw.Iron))
	m.ResMaps.Copper = make([]Position, len(rw.Copper))
	m.ResMaps.Stone = make([]Position, len(rw.Stone))
	m.ResMaps.Coal = make([]Position, len(rw.Coal))
	for i, t := range rw.Water {
		m.WaterTiles[i] = Position{t[0], t[1]}
	}
	for i, t := range rw.Iron {
		m.ResMaps.Iron[i] = Position{t[0], t[1]}
	}
	for i, t := range rw.Copper {
		m.ResMaps.Copper[i] = Position{t[0], t[1]}
	}
	for i, t := range rw.Stone {
		m.ResMaps.Stone[i] = Position{t[0], t[1]}
	}
	for i, t := range rw.Coal {
		m.ResMaps.Coal[i] = Position{t[0], t[1]}
	}
	m.ResMaps.CalcPatches()
}

// for when you get new world chunk and want to add the infomation to the mapper
func (m *Mapper) addRawWorld(rw RawWorld) {
	for _, t := range rw.Water {
		m.WaterTiles = append(m.WaterTiles, Position{t[0], t[1]})
	}
	for _, t := range rw.Iron {
		m.ResMaps.Iron = append(m.ResMaps.Iron, Position{t[0], t[1]})
	}
	for _, t := range rw.Copper {
		m.ResMaps.Copper = append(m.ResMaps.Copper, Position{t[0], t[1]})
	}
	for _, t := range rw.Stone {
		m.ResMaps.Stone = append(m.ResMaps.Stone, Position{t[0], t[1]})
	}
	for _, t := range rw.Coal {
		m.ResMaps.Coal = append(m.ResMaps.Coal, Position{t[0], t[1]})
	}
	m.ResMaps.CalcPatches() // patches need to be recalculated
}

// returns true, if dims is available
func (m *Mapper) canAlloc(dims Box) bool {
	for _, a := range m.Areas {
		d := a.Dims

		if d.Br.X >= dims.Tl.X && d.Br.Y >= dims.Tl.Y && d.Tl.X <= dims.Br.X && d.Tl.Y <= dims.Br.Y {
			return false
		}
	}
	return true
}

// allocates area and returns it's id. Returns -1 if area not available.
func (m *Mapper) alloc(dims Box) int {
	if !m.canAlloc(dims) {
		return -1
	}

	m.Areas = append(m.Areas, Area{dims, m.AllocIdCounter})
	m.AllocIdCounter++

	return m.AllocIdCounter - 1
}

// frees area by id. Returns true, if successful
func (m *Mapper) free(id int) bool {
	for i, v := range m.Areas {
		if v.Id == id {
			m.Areas = append(m.Areas[:i], m.Areas[i+1:]...)
			return true
		}
	}
	return false
}
