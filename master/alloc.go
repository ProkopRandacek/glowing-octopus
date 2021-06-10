package main

type Area struct {
	Dims Box
	Id int
}

var allocIdCounter = 0

func (b *Bot) canAlloc(dims Box) bool {
	for _, a := range b.Areas {
		d := a.Dims

		if d.Br.X >= dims.Tl.X && d.Br.Y >= dims.Tl.Y && d.Tl.X <= dims.Br.X && d.Tl.Y <= dims.Br.Y {
			return false
		}
	}
	return true
}

func (b *Bot) alloc(dims Box) int {
	if !b.canAlloc(dims) {
		return -1
	}

	b.Areas = append(b.Areas, Area{dims, allocIdCounter})
	allocIdCounter++

	b.drawBox(dims, Color{0, 1, 1})

	return allocIdCounter-1
}

func (b *Bot) free(id int) bool {
	for i, v := range b.Areas {
		if v.Id == id {
			b.Areas = append(b.Areas[:i], b.Areas[i+1:]...)
			return true
		}
	}

	return false
}
