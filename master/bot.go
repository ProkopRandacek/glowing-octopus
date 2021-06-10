package main

import (
	"github.com/gtaylor/factorio-rcon"
)

type Bot struct {
	conn   *rcon.RCON
	mapper Mapper
}

func newBot(address, password string) (Bot, error) {
	out := Bot{}

	var err error
	out.conn, err = rcon.Dial(address)
	if err != nil {
		return out, err
	}

	err = out.conn.Authenticate(password)
	if err != nil {
		return out, err
	}

	return out, nil
}
