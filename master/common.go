package main

import (
	"fmt"
	"time"
)

type Color struct {
	R, G, B float32
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Box struct {
	Tl Position // top left
	Br Position // bottom right
}

var lastLog string

func log(msg string) {
	bar := []string{"[.  ]", "[ . ]", "[  .]"}

	if msg != lastLog {
		fmt.Print("\n")
	}

	fps := 2.6
	fmt.Printf("\x1b[2K\x1b[0G%s %s", bar[int(time.Now().UnixNano()/int64(1000000000.0/fps))%len(bar)], msg)
	lastLog = msg
}
