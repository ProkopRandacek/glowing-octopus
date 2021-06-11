package main

import (
	"fmt"
	"os"
)

func main() {
	addr := "localhost:7000"
	pass := "123456"

	if len(os.Args) == 3 {
		addr = os.Args[1]
		pass = os.Args[2]
	}

	bot, err := newBot(addr, pass)
	if err != nil {
		fmt.Println("could not initialize bot")
		return
	}

	trees, _ := bot.getTrees(Box{Position{-100, -100}, Position{100, 100}})
	target := Position{trees[0][0], trees[0][1]}
	bot.walkTo(target)
	bot.waitForTaskDone()
	bot.mine(target)
}
