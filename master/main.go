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

	err := newBot(addr, pass) // bot is global
	if err != nil {
		fmt.Println("could not initialize bot")
		fmt.Println(err.Error())
		return
	}

	bot.allocWater(box(-224, -224, 224, 224))

	resrcs, err := bot.getResources(Box{Position{-700, -700}, Position{700, 700}}) // get the start area ores
	if err != nil {
		fmt.Println("could not get the resources")
		fmt.Println(err.Error())
		return
	}
	bot.Mapper.readResources(resrcs)

	fmt.Println(bot.newFactory("inserter", 4))

	bot.runShell()
}
