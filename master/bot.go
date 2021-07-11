package main

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"math"

	rcon "github.com/gtaylor/factorio-rcon"
)

var fbot = bot{}

type state struct { // Lua fbot internal representation
	Pos            position       `json:"position"`
	Walking        bool           `json:"walking_state"`
	Mining         bool           `json:"mining_state"`
	ResourceMining bool           `json:"mining_resource_state"`
	Placing        bool           `json:"placing_state"`
	Puting         bool           `json:"puting_state"`
	Taking         bool           `json:"taking_state"`
	Clearing       bool           `json:"clearing_state"`
	Building       bool           `json:"building_state"`
	Inventory      map[string]int `json:"inv"`
}

type task func(*bot) error

type bot struct {
	conn            *rcon.RCON
	Mapper          mapper
	TaskList        *list.List
	InserterLevel   string
	AssemblerLevel  int
	BeltLevel       string
	FurnaceLevel    string
	SharedResources map[string][]sharedDepLocation
}

func newBot(address, password string) error {
	err := loadRecipes()
	if err != nil {
		return err
	}

	fbot.conn, err = rcon.Dial(address)
	if err != nil {
		return err
	}

	err = fbot.conn.Authenticate(password)
	if err != nil {
		return err
	}

	fbot.Mapper = mapper{}
	fbot.Mapper.Resources = make(map[string][]position, 6)
	fbot.Mapper.OrePatches = make(map[string][]orePatch, 6)

	fbot.Mapper.Resources = map[string][]position{
		"iron-ore":    {},
		"copper-ore":  {},
		"coal":        {},
		"stone":       {},
		"uranium-ore": {},
		"crude-oil":   {},
	}

	fbot.Mapper.OrePatches = map[string][]orePatch{
		"iron-ore":    {},
		"copper-ore":  {},
		"coal":        {},
		"stone":       {},
		"uranium-ore": {},
		"crude-oil":   {},
	}

	fbot.TaskList = list.New()

	fbot.InserterLevel = "fast"
	fbot.AssemblerLevel = 1
	fbot.BeltLevel = ""
	fbot.FurnaceLevel = "stone"

	fbot.SharedResources = map[string][]sharedDepLocation{}

	return nil
}

func (b *bot) safeExecute(command string) {
	_, err := b.conn.Execute(command)
	if err != nil {
		panic("error while calling factorio command: " + err.Error())
	}
}

func (b *bot) state() (state, error) {
	f, err := os.Open("./master/script-output/state.json")
	if f == nil {
		fmt.Println("Error opening state file: ", err)
		return state{}, err
	}
	defer f.Close() // TODO: handle error from f.Close()

	dat, err := io.ReadAll(f)

	if err != nil {
		fmt.Println("Error reading state file: ", err)
		return state{}, err
	}

	s := state{}
	err = json.Unmarshal(dat, &s)
	if err != nil {
		return state{}, nil
	}

	return s, nil
}

func (b *bot) addTask(t task) {
	b.TaskList.PushBack(t)
}

func (b *bot) pushTask(t task) {
	b.TaskList.PushFront(t)
}

func (b *bot) doTask() error {
	if b.TaskList.Len() == 0 {
		return errors.New("no tasks left")
	}

	err := b.TaskList.Front().Value.(task)(b)

	b.TaskList.Remove(b.TaskList.Front())

	return err
}

func (b *bot) findPlaceToMine(item string, count int) (position, error) {
	index := -1
	minDist := -1.0
	state, _ := b.state()
	pos := state.Pos

	for i, t := range b.Mapper.Resources[item] {
		if dist := math.Pow(t.X - pos.X, 2) + math.Pow(t.Y - pos.Y, 2); (dist < minDist || index == -1) && b.Mapper.ResourceAmounts[item][i] >= count {
			index = i
			minDist = dist
		}
	}

	if index == -1 {
		return position{}, errors.New("no tile found")
	}

	b.Mapper.ResourceAmounts[item][index] -= count

	return b.Mapper.Resources[item][index], nil
}
