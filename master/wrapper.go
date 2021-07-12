package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
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
		if !(s.Walking || s.Mining || s.ResourceMining || s.Placing || s.Puting || s.Taking || s.Clearing || s.Building) {
			break
		}
	}
	fmt.Println()
}

func (b *bot) collectItemsForBP(bp []building) {

	// TODO collecting from mall and stuff

	state, _ := b.state()
	items := howMuchMore(state.Inventory, calcBPItems(bp))

	fmt.Println("I need to gather %v", items)

	for item, count := range items {
		if count == 0 {
			continue
		}
		item = b.resolveBuildingName(item)
		if item == "iron-plate" || item == "copper-plate" { // things that we need to smelt
			fmt.Println(item, count)
			oreName := strings.Split(item, "-")[0] + "-ore"
			b.mineResource(b.Mapper.findPlaceToMine(oreName, count), count, oreName)
			b.waitForTaskDone()

			coalCount := int(math.Ceil(float64(count) / 13.0))

			b.mineResource(b.Mapper.findPlaceToMine("coal", coalCount), coalCount, "coal")
			b.waitForTaskDone()

			state, _ := b.state()
			furnacePos := state.Pos
			furnaceBox := makeBox(furnacePos.X-2, furnacePos.Y-2, furnacePos.X+2, furnacePos.Y+2)

			b.Mapper.findSpace(&furnaceBox)
			furnacePos = position{furnaceBox.Tl.X + 2, furnaceBox.Tl.Y + 2}
			b.clearAll(furnaceBox)
			b.waitForTaskDone()

			allocid, _ := b.Mapper.alloc(furnaceBox)

			b.place(furnacePos, "stone-furnace")
			b.waitForTaskDone()
			furnacePos.X -= 0.5
			b.put(furnacePos, oreName, count, 2)
			b.put(furnacePos, "coal", coalCount, 1)
			time.Sleep(2 * time.Second)
			b.take(furnacePos, item, count, 3)
			b.waitForTaskDone()
			b.mine(furnacePos)
			b.waitForTaskDone()

			b.Mapper.free(allocid)
		} else if item == "wood" {
			fmt.Println("need wood ", count)
		}
	}
	fmt.Println("Resource gathering done")

	for _, building := range bp {
		b.craft(building.Name, 1)
		b.waitForTaskDone()
	}
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
