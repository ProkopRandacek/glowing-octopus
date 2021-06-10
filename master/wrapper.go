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
	_, err := b.conn.Execute(fmt.Sprintf("/walkto {\"x\": %d, \"y\": %d}", x, y))
	return err
}

func (b *Bot) getWorld() (RawWorld, error) {
	os.Remove(worldFilename)
	_, err := b.conn.Execute("/writeworld")
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

func (b *Bot) drawBox(box Box) {
	b.conn.Execute(fmt.Sprintf(`/drawbox {"color":[1, 0, 0, 0.2],"x1":%2.f,"y1":%2.f,"x2":%2.f,"y2":%2.f}`, box.tl.x, box.tl.y, box.br.x, box.br.y))
}
