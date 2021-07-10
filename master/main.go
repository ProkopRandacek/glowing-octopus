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

	err := newBot(addr, pass) // fbot is global
	if err != nil {
		fmt.Println("could not initialize fbot")
		fmt.Println(err.Error())
		return
	}

	fbot.allocWater(makeBox(-224, -224, 224, 224))

	resources, err := fbot.getResources(box{position{-700, -700}, position{700, 700}}) // get the start area ores
	if err != nil {
		fmt.Println("could not get the resources")
		fmt.Println(err.Error())
		return
	}
	fbot.Mapper.readResources(resources)

	fmt.Println(fbot.newFactory("inserter", 4))

	fbot.runShell()
}
