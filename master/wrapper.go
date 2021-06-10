package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
)

type RawWorld struct {
	Water  [][]float32
	Iron   [][]float32 `json:"iron-ore"`
	Copper [][]float32 `json:"copper-ore"`
	Stone  [][]float32
	Coal   [][]float32
}

func (b *Bot) walkTo(x, y int) error {
	_, err := b.conn.Execute(fmt.Sprintf("/walkto {\"x\": %d, \"y\": %d}", x, y))
	return err
}

func (b *Bot) getWorld() (RawWorld, error) {
	_, err := b.conn.Execute("/writeworld")
	if (err != nil) { return RawWorld{}, err }

	dat, err := ioutil.ReadFile("./master/script-output/world.json")
	if (err != nil) { return RawWorld{}, err }

	var world RawWorld
	json.Unmarshal(dat, &world)

	return world, nil
}
