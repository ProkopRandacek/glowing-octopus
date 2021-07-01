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

	bot.clearAll(Box{Position{-70,-20}, Position{20,20}}) // mine the ship

	resrcs, err := bot.getResources(Box{Position{-220, -220},Position{220, 220}}) // get the start area ores
	if err != nil {
		fmt.Println("could not get the resources")
		fmt.Println(err.Error())
		return
	}
	bot.Mapper.readResources(resrcs)

	bot.mineResource(bot.Mapper.Resrcs[2][0], 10, "coal")
	bot.waitForTaskDone()
	bot.mineResource(bot.Mapper.Resrcs[0][0], 10, "iron-ore")
	bot.waitForTaskDone()
	bot.place(bot.state().Pos, "stone-furnace") // place furnace
	bot.waitForTaskDone()
	bot.put(bot.state().Pos, "coal", 10, 1) // put coal in furnace
}
