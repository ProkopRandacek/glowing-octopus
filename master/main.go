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

	err := newBot(addr, pass) // octopus is global
	if err != nil {
		fmt.Println("could not initialize octopus")
		fmt.Println(err.Error())
		return
	}

	octopus.allocWater(makeBox(-224, -224, 224, 224)) // the starting area

	resources, err := octopus.getResources(box{position{-700, -700}, position{700, 700}}) // get the start area ores
	if err != nil {
		fmt.Println("could not get the resources")
		fmt.Println(err.Error())
		return
	}
	octopus.Mapper.readResources(resources)

	octopus.clearAll(makeBox(-50, -50, 50, 50))
	octopus.waitForTaskDone()

	octopus.collectItemsForBP(noFluidBp.Buildings)
}
