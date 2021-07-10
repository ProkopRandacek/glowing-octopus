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

func (b *bot) waitForTaskDone() { // Waits until task is done.
	for {
		log("Waiting for task done")
		time.Sleep(2 * time.Second)
		s, err := b.state()
		if err != nil {
			panic("error while waiting for task: " + err.Error())
		}
		if !(s.Walking || s.Mining || s.ResourceMining || s.Placing || s.Clearing || s.Building) {
			break
		}
	}
	fmt.Println()
}

func (b *bot) mineRawResources(map[string]int) {
}

func (b *bot) getResources(box box) (map[string][]position, error) {
	filename := scriptFolder + "resrc.json"
	err := os.Remove(filename)
	if err != nil {
		return nil, err
	}
	_, err = b.conn.Execute(fmt.Sprintf("/writeresrc [[%.2f,%.2f],[%.2f,%.2f]]", box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
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
	defer f.Close() // TODO: handle error from f.Close()

	dat, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var resources map[string][]position
	err = json.Unmarshal(dat, &resources)
	if err != nil {
		return nil, err
	}

	b.Mapper.LoadedBoxes = append(b.Mapper.LoadedBoxes, box)

	return resources, nil
}

func (b *bot) allocWater(box box) error {
	filename := scriptFolder + "water.json"
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	_, err = b.conn.Execute(fmt.Sprintf("/writewater [[%.2f,%.2f],[%.2f,%.2f]]", box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
	if err != nil {
		return err
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
		return err
	}
	defer f.Close() // TODO: Handle error from f.Close()

	dat, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var water []position
	err = json.Unmarshal(dat, &water)
	if err != nil {
		return err
	}

	fmt.Println(water)

	for _, pos := range water {
		b.Mapper.forceAlloc(makeBox(pos.X, pos.Y, pos.X+32, pos.Y+32))
	}
	return nil
}

func (b *bot) walkTo(p position) {
	b.safeExecute(fmt.Sprintf(`/walkto [%.2f,%.2f]`, p.X, p.Y))
}

func (b *bot) drawBox(box box, color color) {
	b.safeExecute(fmt.Sprintf(`/drawbox {"color":[%2.f, %2.f, %2.f, 0.2],"x1":%2.f,"y1":%2.f,"x2":%2.f,"y2":%2.f}`, color.R, color.G, color.B, box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
}

func (b *bot) drawPoint(pos position, color color) {
	b.safeExecute(fmt.Sprintf(`/drawpoint {"color":[%2.f, %2.f, %2.f, 0.2],"x":%2.f,"y":%2.f}`, color.R, color.G, color.B, pos.X, pos.Y))
}

func (b *bot) craft(r string, c int) {
	b.safeExecute(fmt.Sprintf(`/craft {"recipe":"%s","count":%d}`, r, c))
}

func (b *bot) mine(p position) {
	b.safeExecute(fmt.Sprintf(`/mine [%.2f,%.2f]`, p.X, p.Y))
}

func (b *bot) mineResource(p position, amount int, name string) {
	b.safeExecute(fmt.Sprintf(`/mineresource {"pos":[%.2f,%.2f],"amount":%d,"name":"%s"}`, p.X, p.Y, amount, name))
}

func (b *bot) clearNature(box box) {
	b.safeExecute(fmt.Sprintf(`/cleararea {"area":[[%2.f, %2.f],[%2.f, %2.f]],"t":"nature"}`, box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
}

func (b *bot) clearAll(box box) {
	b.safeExecute(fmt.Sprintf(`/cleararea {"area":[[%2.f, %2.f],[%2.f, %2.f]],"t":"all"}`, box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
}

func (b *bot) place(p position, item string) {
	b.safeExecute(fmt.Sprintf(`/place {"pos":[%2.f,%2.f],"item":"%s"}`, p.X, p.Y, item))
}

func (b *bot) placeRecipe(p position, item string, recipe string) {
	b.safeExecute(fmt.Sprintf(`/place {"pos":[%2.f,%2.f],"item":"%s","recipe":"%s"}`, p.X, p.Y, item, recipe))
}

func (b *bot) placeDir(p position, item string, dir int) {
	b.safeExecute(fmt.Sprintf(`/place {"pos":[%2.f,%2.f],"item":"%s","dir":%d}`, p.X, p.Y, item, dir))
}

// https://lua-api.factorio.com/latest/defines.html#defines.inventory
// fuel = 1
// furnace_source = 2
// furnace result = 3
func (b *bot) put(p position, item string, amount int, slot int) {
	b.safeExecute(fmt.Sprintf(`/put {"pos":[%2.f,%2.f],"item":"%s","amount":%d,"slot":%d}`, p.X, p.Y, item, amount, slot))
}

func (b *bot) take(p position, item string, amount int, slot int) {
	b.safeExecute(fmt.Sprintf(`/take {"pos":[%2.f,%2.f],"item":"%s","amount":%d,"slot":%d}`, p.X, p.Y, item, amount, slot))
}

func (b *bot) build(bs []building) {
	s, _ := json.Marshal(bs)
	b.safeExecute(fmt.Sprintf(`/build %s`, string(s)))
}
