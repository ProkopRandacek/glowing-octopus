package main

import (
	"fmt"
)

func (b *Bot) walkTo(x, y int) error {
	_, err := b.conn.Execute(fmt.Sprintf("/walkto {\"x\": %d, \"y\": %d}", x, y))
	return err
}
