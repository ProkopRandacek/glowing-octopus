package main

import (
	"fmt"
	"math"
	"time"
)

const (
	dirNorth     = 0
	dirNorthEast = 1
	dirEast      = 2
	dirSouthEast = 3
	dirSouth     = 4
	dirSouthWest = 5
	dirWest      = 6
	dirNorthWest = 7
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

type Measurable interface {
	getDims() Box
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

func (b *Box) Round() {
	b.Tl.X = math.Round(b.Tl.X)
	b.Tl.Y = math.Round(b.Tl.Y)
	b.Br.X = math.Round(b.Br.X)
	b.Br.Y = math.Round(b.Br.Y)
}

func box(a, b, c, d float64) Box { // a shortcut
	return Box{Position{a, b}, Position{c, d}}
}

func Find(slice []string, val string) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}

func getRawItemsFromItem(num float64, item string) map[string]float64 {
	stopItems := []string{"iron-plate", "copper-plate"}
	out := map[string]float64{}
	val, ok := recipes[item]

	if !ok || Find(stopItems, item) {
		return map[string]float64{item: num}
	}

	for _, d := range val.Deps {
		n := getRawItemsFromItem(float64(d.Count) * num / float64(val.CraftNum), d.Name)
		for n, c := range n {
			out[n] += c
		}
	}

	return out
}

func calcBPItems([]Building) map[string]int {
	buldingsCount := map[string]int{}
	for _, b := range []Building{} {
		buldingsCount[b.Name]++
	}

	items := map[string]float64{}

	for building, count := range buldingsCount {
		for n, c := range getRawItemsFromItem(float64(count), building,) {
			items[n] += c
		}
	}
	itemInts := map[string]int{}
	for k, v := range items {
		itemInts[k] = int(math.Ceil(v))
	}

	return itemInts
}
