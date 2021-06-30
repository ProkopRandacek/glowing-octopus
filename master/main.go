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

	A := Box{Position{0, 0}, Position{10, 10}}

	for i := 0; i < 100; i++ {
		B := A
		bot.Mapper.findSpace(&B)
		bot.drawBox(B, Color{0, 0, 1})
		bot.Mapper.alloc(B)
	}
}
