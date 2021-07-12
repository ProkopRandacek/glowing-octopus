package main

import (
	"errors"
	"fmt"
)

type orePatch struct {
	Dims   box
	Unsafe bool
}

type area struct {
	Dims box
	Id   int
}

type mapper struct {
	Areas          []area // marks what areas area allocated
	AllocIdCounter int    // for area id generation
	MallMap        map[string]position
	AllocMap       []box
	WaterTiles     []position            // the exact tiles that are in water
	WaterBox       box                   // water boxes for fast intersection check
	Resources      map[string][]position // all the individual ore tiles
	OrePatches     map[string][]orePatch // ore tiles grouped together into patches
	LoadedBoxes    []box                 // all the area boxes that we requested from the game
}

func (o *orePatch) isUnsafe(boxes []box) {
	for i, b := range boxes {
		if isBorderSharedWithBox(o.Dims, []box{b}, -1) && !isBorderSharedWithBox(b, boxes, i) {
			o.Unsafe = true
			return
		}
	}

	o.Unsafe = false
}

func (m *mapper) calcPatches() {
	for t := range m.OrePatches { // for each resource type
		m.OrePatches[t] = findComponents(m.Resources[t])
	}
}

func (m *mapper) calcUnsafe() {
	for i := range m.OrePatches {
		for j := range m.OrePatches[i] {
			m.OrePatches[i][j].isUnsafe(m.LoadedBoxes)
		}
	}
}

func (m *mapper) allocatePatches() error {
	for i := range m.OrePatches {
		for j := range m.OrePatches[i] {
			if m.OrePatches[i][j].Unsafe {
				continue
			}

			_, err := m.alloc(m.OrePatches[i][j].Dims)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *mapper) readResources(r map[string][]position) error {
	for t := range m.Resources { // for each resource type
		if len(r[t]) > 0 {
			m.Resources[t] = append(m.Resources[t], r[t]...) // append the ores to it
			fmt.Println("read", t, "len:", len(r[t]))
		}
	}
	m.calcPatches()
	m.calcUnsafe()
	err := m.allocatePatches()
	if err != nil {
		return err
	}
	return nil
}

func (m *mapper) isTileEmpty(p position) bool {
	for _, t := range m.Areas {
		a := t.Dims

		if !(p.X < a.Tl.X || a.Br.X < p.X ||
			p.Y < a.Tl.Y || a.Br.Y < p.Y) {
			return false
		}
	}
	return true
}

// returns true, if dims is available
func (m *mapper) canAlloc(b box) bool {
	const gap = 7 // how much space to keep between allocated areas (for belts and stuff)
	for _, t := range m.Areas {
		a := t.Dims

		if !((b.Tl.X-a.Br.X) > gap ||
			(b.Br.X-a.Tl.X) < -gap ||
			(b.Tl.Y-a.Br.Y) > gap ||
			(b.Br.Y-a.Tl.Y) < -gap) {
			return false
		}
	}

	return true
}

// allocates area and returns it's id. Returns -1 if area not available.
func (m *mapper) alloc(dims box) (int, error) {
	if !m.canAlloc(dims) {
		return 0, errors.New("trying to allocate over another allocated position")
	}

	return m.forceAlloc(dims), nil
}

func (m *mapper) forceAlloc(dims box) int {
	m.Areas = append(m.Areas, area{dims, m.AllocIdCounter})
	m.AllocIdCounter++

	fbot.drawBox(dims, color{1, 0, 0})

	return m.AllocIdCounter - 1
}

// Moves the box to a near place where it fits.
func (m *mapper) findSpace(dims *box) {
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
func (m *mapper) free(id int) bool {
	for i, v := range m.Areas {
		if v.Id == id {
			m.Areas = append(m.Areas[:i], m.Areas[i+1:]...)
			return true
		}
	}
	return false
}

func (m *mapper) findPlaceToMine(item string, count int) position {
	// TODO: can run out of the resource, is often far away...
	return m.Resources[item][0]
}
