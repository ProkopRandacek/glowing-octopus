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

	A := Box{Position{0, 0}, Position{10, 10}}

	resrcs, err := bot.getResources(Box{Position{-70, -50}, Position{50, 50}}) // get the start area ores
	if err != nil {
		fmt.Println("could not get the resources")
		fmt.Println(err.Error())
		return
	}
	bot.Mapper.readResources(resrcs)

	fmt.Println(bot.Mapper.OrePatches)

	bot.clearAll(bot.Mapper.OrePatches["copper-ore"][0].Dims)
	bot.build(bot.newMiners(bot.Mapper.OrePatches["copper-ore"][0]))
}
