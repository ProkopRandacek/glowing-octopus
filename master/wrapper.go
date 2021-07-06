package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	scriptFolder = "./master/script-output/"
)

func (b *Bot) waitForTaskDone() { // Waits until task is done.
	for {
		log("Waiting for task done")
		time.Sleep(2 * time.Second)
		s := b.state()
		if !(s.Walking || s.Mining || s.ResourceMining || s.Placing || s.Clearing || s.Building) {
			break
		}
	}
	fmt.Println()
}

func (b *Bot) getResources(box Box) (map[string][]Position, error) {
	filename := scriptFolder + "resrc.json"
	os.Remove(filename)
	_, err := b.conn.Execute(fmt.Sprintf("/writeresrc [[%.2f,%.2f],[%.2f,%.2f]]", box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
	if err != nil {
		return nil, err
	}

	for { // wait until the file is written
		_, err := os.Stat(filename)
		if err == nil {
			break
		}

		log(fmt.Sprintf("Waiting for the %s to be generated", filename))
		time.Sleep(time.Second)
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dat, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var resrcs map[string][]Position
	err = json.Unmarshal(dat, &resrcs)
	if err != nil {
		return nil, err
	}

	b.Mapper.LoadedBoxes = append(b.Mapper.LoadedBoxes, box)

	return resrcs, nil
}

func (b *Bot) walkTo(p Position) {
	b.conn.Execute(fmt.Sprintf(`/walkto [%.2f,%.2f]`, p.X, p.Y))
}

func (b *Bot) drawBox(box Box, color Color) {
	b.conn.Execute(fmt.Sprintf(`/drawbox {"color":[%2.f, %2.f, %2.f, 0.2],"x1":%2.f,"y1":%2.f,"x2":%2.f,"y2":%2.f}`, color.R, color.G, color.B, box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
}

func (b *Bot) drawPoint(pos Position, color Color) {
	b.conn.Execute(fmt.Sprintf(`/drawpoint {"color":[%2.f, %2.f, %2.f, 0.2],"x":%2.f,"y":%2.f}`, color.R, color.G, color.B, pos.X, pos.Y))
}

func (b *Bot) craft(r string, c int) {
	b.conn.Execute(fmt.Sprintf(`/craft {"recipe":"%s","count":%d}`, r, c))
}

func (b *Bot) mine(p Position) {
	b.conn.Execute(fmt.Sprintf(`/mine [%.2f,%.2f]`, p.X, p.Y))
}

func (b *Bot) mineResource(p Position, amount int, name string) {
	b.conn.Execute(fmt.Sprintf(`/mineresource {"pos":[%.2f,%.2f],"amount":%d,"name":"%s"}`, p.X, p.Y, amount, name))
}

func (b *Bot) clearNature(box Box) {
	b.conn.Execute(fmt.Sprintf(`/cleararea {"area":[[%2.f, %2.f],[%2.f, %2.f]],"t":"nature"}`, box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
}

func (b *Bot) clearAll(box Box) {
	b.conn.Execute(fmt.Sprintf(`/cleararea {"area":[[%2.f, %2.f],[%2.f, %2.f]],"t":"all"}`, box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
}

func (b *Bot) place(p Position, item string) {
	b.conn.Execute(fmt.Sprintf(`/place {"pos":[%2.f,%2.f],"item":"%s"}`, p.X, p.Y, item))
}

func (b *Bot) placeRecipe(p Position, item string, recipe string) {
	b.conn.Execute(fmt.Sprintf(`/place {"pos":[%2.f,%2.f],"item":"%s","recipe":"%s"}`, p.X, p.Y, item, recipe))
}

func (b *Bot) placeDir(p Position, item string, dir int) {
	b.conn.Execute(fmt.Sprintf(`/place {"pos":[%2.f,%2.f],"item":"%s","dir":%d}`, p.X, p.Y, item, dir))
}

// https://lua-api.factorio.com/latest/defines.html#defines.inventory
// fuel = 1
// furnace_source = 2
// furnace result = 3
func (b *Bot) put(p Position, item string, amount int, slot int) {
	b.conn.Execute(fmt.Sprintf(`/put {"pos":[%2.f,%2.f],"item":"%s","amount":%d,"slot":%d}`, p.X, p.Y, item, amount, slot))
}

func (b *Bot) take(p Position, item string, amount int, slot int) {
	b.conn.Execute(fmt.Sprintf(`/take {"pos":[%2.f,%2.f],"item":"%s","amount":%d,"slot":%d}`, p.X, p.Y, item, amount, slot))
}

func (b *Bot) build(bs []Building) {
	s, _ := json.Marshal(bs)
	b.conn.Execute(fmt.Sprintf(`/build %s`, string(s)))
}
