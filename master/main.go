package main

import (
	"fmt"
	"os"
	//time"
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
		fmt.Println(err.Error())
		return
	}

	bp, err := bot.newFactory("inserter", 4)
	if err != nil {
		fmt.Println("could not build the factory")
		fmt.Println(err.Error())
		return
	}

	bot.build(bp)
}
