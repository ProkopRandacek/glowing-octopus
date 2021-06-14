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
		//return
	}

	greenCircuit := Item{
		"green-circuit",
		0.5,
		false,
		[]RecipeDep{
			RecipeDep{
				Item{
					"copper-cable",
					0.25,
					false,
					[]RecipeDep{},
				},
				3,
				true,
			},
			RecipeDep{
				Item{
					Name: "iron-plate",
				},
				1,
				false,
			},
		},
	}

	bot.newFactory(greenCircuit, 2)

	/*r, _ := bot.getResources(Box{Position{-100, -100}, Position{100, 100}})
	bot.Mapper.readResources(r)

	bot.drawBox(bot.Mapper.LoadedBoxes[0], Color{0, 0, 1})
	for _, t := range bot.Mapper.OrePatches {
		for _, o := range t {
			c := Color{0, 1, 0}
			if o.Unsafe {
				c = Color{1, 0, 0}
			}
			bot.drawBox(o.Dims, c)
		}
	}

	fmt.Println()
	for {
		bot.refreshState()
		fmt.Println("\x1b[A", bot.State)
		time.Sleep(1 * time.Second)
	}*/
}
