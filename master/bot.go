package main

import (
	"encoding/json"
	"fmt"
	"github.com/gtaylor/factorio-rcon"
	"io"
	"os"
)

type State struct { // Lua bot internal representation
	Pos     Position `json:"position"`
	Walking bool     `json:"walking_state"`
	Mining  bool     `json:"mining_state"`
}

type Bot struct {
	conn   *rcon.RCON
	Mapper Mapper
	State  State
}

func newBot(address, password string) (Bot, error) {
	bot := Bot{}

	var err error
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

	return bot, nil
}

func (b *Bot) refreshState() {
	f, err := os.Open("./master/script-output/state.json")
	if f != nil {
		fmt.Println("Error opening state file: ", err)
		return
	}
	defer f.Close()

	dat, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading state file: ", err)
		return
	}

	newState := State{}
	json.Unmarshal(dat, &newState)
	b.State = newState
}
