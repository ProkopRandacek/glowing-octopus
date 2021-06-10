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
	rawWorld, err := bot.getWorld()

	bot.mapper.readRawWorld(rawWorld)

	bot.drawBox(bot.mapper.resMaps.ironPatches[0], Color{1, 0, 0})
	id1 := bot.mapper.alloc(Box{Position{1, 5}, Position{9, 12}})
	id2 := bot.mapper.alloc(Box{Position{25, 5}, Position{30, 12}})
	id3 := bot.mapper.alloc(Box{Position{1, 5}, Position{9, 12}})
	fmt.Printf("allocated: %d %d %d\n", id1, id2, id3)
	fmt.Printf("can alloc: %v\n", bot.mapper.canAlloc(Box{Position{2, 2}, Position{20, 9}}))
	f1 := bot.mapper.free(0)
	f2 := bot.mapper.free(1)
	f3 := bot.mapper.free(2)
	fmt.Printf("freed: %v %v %v\n", f1, f2, f3)

	fmt.Println(bot.mapper.resMaps.ironPatches)
}
