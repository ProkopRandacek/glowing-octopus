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
		return
	}

	bot.addTask(func(b *Bot) error {
		b.walkTo(Position{10, 10})
		b.waitForTaskDone()
		return nil
	})
	bot.addTask(func(b *Bot) error {
		b.walkTo(Position{-10, -10})
		b.waitForTaskDone()
		return nil
	})
	bot.pushTask(func(b *Bot) error {
		fmt.Println("start")
		return nil
	})
	bot.addTask(func(b *Bot) error {
		fmt.Println("done")
		return nil
	})

	for bot.doTask() == nil {
	}
}
