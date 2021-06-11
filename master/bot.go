package main

import (
	"encoding/json"
	"fmt"
	"github.com/gtaylor/factorio-rcon"
	"io/ioutil"
)

type State struct { // Lua bot internal representation
	Pos     Position `json:"position"`
	Walking bool     `json:"walking_state"`
	Mining  bool     `json:"mining_state"`
}

type Bot struct {
	conn   *rcon.RCON
	mapper Mapper
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

	return bot, nil
}

func (b *Bot) refreshState() {
	dat, err := ioutil.ReadFile("./master/script-output/state.json")
	if err != nil {
		fmt.Println("Error reading state file: ", err)
		return
	}

	newState := State{}
	json.Unmarshal(dat, &newState)
	b.State = newState
}
