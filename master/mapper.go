package main

import (
	"fmt"
	"math"
)

type OrePatch struct {
	Dims   Box
	Unsafe bool
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
	WaterTiles     []Position            // the exact tiles that are in water
	WaterBox       Box                   // water boxes for fast intersection check
	Resrcs         map[string][]Position // all the individual ore tiles
	OrePatches     map[string][]OrePatch // ore tiles grouped together into patches
	LoadedBoxes    []Box                 // all the area boxes that we requested from the game
}

// Ore types
const (
	iron = iota
	copper
	stone
	coal
	uran
)

func isBorderSharedWithBox(box Box, boxes []Box, ignore int) bool {
	box.Round()

	for i, b := range boxes {
		if i == ignore {
			continue
		}

		b.Round()

		if box.Tl.X == b.Tl.X ||
			box.Tl.Y == b.Tl.Y ||
			box.Br.X == b.Br.Y ||
			box.Br.Y == b.Br.Y ||
			box.Tl.X == b.Br.X ||
			box.Tl.Y == b.Br.Y {
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

func findComponents(tiles []Position) (boxes []OrePatch) { // divide graph into components but only store component bounds
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
		boxes = append(boxes, OrePatch{Box{Position{minx, miny}, Position{maxx, maxy}}, true})
	}
	return
}

func (m *Mapper) calcPatches() {
	for t := range m.OrePatches { // for each resource type
		m.OrePatches[t] = findComponents(m.Resrcs[t])
	}
}

func (m *Mapper) calcUnsafe() {
	for i := range m.OrePatches {
		for j := range m.OrePatches[i] {
			m.OrePatches[i][j].isUnsafe(m.LoadedBoxes)
		}
	}
}

func (m *Mapper) readResources(r map[string][]Position) {
	for t := range m.Resrcs { // for each resource type
		if len(r[t]) > 0 {
			m.Resrcs[t] = append(m.Resrcs[t], r[t]...) // append the ores to it
			fmt.Println("read", t, "len:", len(r[t]))
		}
	}
	m.calcPatches()
	m.calcUnsafe()
}

// returns true, if dims is available
func (m *Mapper) canAlloc(b Box) bool {
	const gap = 5 // how much space to keep between allocated areas (for belts and stuff)
	for _, t := range m.Areas {
		a := t.Dims

		if !((b.Tl.X - a.Br.X) > gap ||
			(b.Br.X - a.Tl.X) < -gap ||
			(b.Tl.Y - a.Br.Y) > gap ||
			(b.Br.Y - a.Tl.Y) < -gap) {
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

// Moves the box to a near place where it fits.
func (m *Mapper) findSpace(dims *Box) {
	side := 1

	for {
		for i := 0; i < side; i++ { // move right
			if m.canAlloc(*dims) {
				return
			}
			dims.Tl.X += 1
			dims.Br.X += 1
		}
		for i := 0; i < side; i++ { // move down
			if m.canAlloc(*dims) {
				return
			}
			dims.Tl.Y += 1
			dims.Br.Y += 1
		}
		side += 1
		for i := 0; i < side; i++ { // move left
			if m.canAlloc(*dims) {
				return
			}
			dims.Tl.X -= 1
			dims.Br.X -= 1
		}
		for i := 0; i < side; i++ { // move up
			if m.canAlloc(*dims) {
				return
			}
			dims.Tl.Y -= 1
			dims.Br.Y -= 1
		}
		side += 1
	}
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
