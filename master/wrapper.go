package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	worldFilename = "./master/script-output/world.json"
)

type RawWorld struct {
	Water  [][]float64 `json:"water"`
	Iron   [][]float64 `json:"iron-ore"`
	Copper [][]float64 `json:"copper-ore"`
	Stone  [][]float64 `json:"stone"`
	Coal   [][]float64 `json:"coal"`
}

func (b *Bot) walkTo(x, y int) error {
	b.State.Walking = true
	_, err := b.conn.Execute(fmt.Sprintf(`/walkto {"x":%d,"y":%d}`, x, y))
	return err
}

func (b *Bot) waitForWalkDone() {
	for {
		time.Sleep(2 * time.Second)
		b.refreshState()
		if b.State.Walking {
			fmt.Println("Waiting for walk done")
		} else {
			break
		}
	}
}

func (b *Bot) getWorld(box Box) (RawWorld, error) {
	os.Remove(worldFilename)
	_, err := b.conn.Execute(fmt.Sprintf("/writeworld [[%.2f,%.2f],[%.2f,%.2f]]", box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
	if err != nil {
		return RawWorld{}, err
	}

	for { // wait until the file is written
		_, err := os.Stat(worldFilename)
		if os.IsNotExist(err) {
			fmt.Println("Waiting for the world.json to be generated...")
			time.Sleep(time.Second)
		} else if err == nil {
			break
		}
	}

	dat, err := ioutil.ReadFile(worldFilename)
	if err != nil {
		return RawWorld{}, err
	}

	var world RawWorld
	json.Unmarshal(dat, &world)

	return world, nil
}

func (b *Bot) drawBox(box Box, color Color) {
	b.conn.Execute(fmt.Sprintf(`/drawbox {"color":[%2.f, %2.f, %2.f, 0.2],"x1":%2.f,"y1":%2.f,"x2":%2.f,"y2":%2.f}`, color.R, color.G, color.B, box.Tl.X, box.Tl.Y, box.Br.X, box.Br.Y))
}

func (b *Bot) craft(r string, c int) {
	b.conn.Execute(fmt.Sprintf(`/craft {"recipe":"%s","count":%d}`, r, c))
}
