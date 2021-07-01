package main

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gtaylor/factorio-rcon"
	"io"
	"os"
)

type State struct { // Lua bot internal representation
	Pos      Position `json:"position"`
	Walking  bool     `json:"walking_state"`
	Mining   bool     `json:"mining_state"`
	Placing  bool     `json:"placing_state"`
	Clearing bool     `json:"clearing_state"`
	Building bool     `json:"building_state"`
}

type Task func(*Bot) error

type Bot struct {
	conn           *rcon.RCON
	Mapper         Mapper
	TaskList       *list.List
	InserterLevel  string
	AssemblerLevel int
	BeltLevel      string
}

func newBot(address, password string) (Bot, error) {
	bot := Bot{}

	err := loadRecipes()
	if err != nil {
		return bot, err
	}

	bot.conn, err = rcon.Dial(address)
	if err != nil {
		return bot, err
	}

	err = bot.conn.Authenticate(password)
	if err != nil {
		return bot, err
	}

	bot.Mapper = Mapper{}
	bot.Mapper.Resrcs = make([][]Position, 4)
	bot.Mapper.OrePatches = make([][]OrePatch, 4)

	bot.TaskList = list.New()

	bot.InserterLevel = "fast"
	bot.AssemblerLevel = 1
	bot.BeltLevel = ""

	return bot, nil
}

func (b *Bot) state() State {
	f, err := os.Open("./master/script-output/state.json")
	if f == nil {
		fmt.Println("Error opening state file: ", err)
		return State{}
	}
	defer f.Close()

	dat, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading state file: ", err)
		return State{}
	}

	state := State{}
	json.Unmarshal(dat, &state)

	return state
}

func (b *Bot) addTask(t Task) {
	b.TaskList.PushBack(t)
}

func (b *Bot) pushTask(t Task) {
	b.TaskList.PushFront(t)
}

func (b *Bot) doTask() error {
	if b.TaskList.Len() == 0 {
		return errors.New("no tasks left")
	}

	err := b.TaskList.Front().Value.(Task)(b)

	b.TaskList.Remove(b.TaskList.Front())

	return err
}
