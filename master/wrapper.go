package main

import (
	"fmt"
)

func (b *Bot) walkTo(x, y int) error {
	_, err := b.conn.Execute(fmt.Sprintf("/walkto {\"x\": %d, \"y\": %d}", x, y))
	return err
}

func (b *Bot) getWorld() (string, error) {
	response, err := b.conn.Execute("/writeworld")
	if err != nil {
		return "", err
	}

	return response.Body, nil
}
