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

	bot.walkTo(20, 20)
	bot.waitForWalkDone()
	//bot.craft("wooden-chest", 1)
	//bot.getWorld(Box{Position{-50, -50}, Position{50, 50}})
}
