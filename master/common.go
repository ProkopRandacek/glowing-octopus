package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
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

var recipes map[string]item

type color struct {
	R, G, B float32
}

type position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type box struct {
	Tl position // top left
	Br position // bottom right
}

var lastLog string

func log(msg string) {
	bar := []string{"[.  ]", "[ . ]", "[  .]"}

	if msg != lastLog {
		fmt.Println()
	}

	fps := 2.6
	fmt.Printf("\x1b[2K\x1b[0G%s %s", bar[int(time.Now().UnixNano()/int64(1000000000.0/fps))%len(bar)], msg)
	lastLog = msg
}

func loadRecipes() error {
	f, err := os.Open("./master/recipes.json")
	if err != nil {
		return err
	}
	defer f.Close() // TODO: Handle error from f.Close()

	dat, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dat, &recipes)
	if err != nil {
		return err
	}

	return nil
}

func getRawItemsFromItem(num float64, item string) map[string]float64 {
	stopItems := []string{"iron-plate", "copper-plate"}
	out := map[string]float64{}
	val, ok := recipes[item]

	if !ok || find(stopItems, item) {
		return map[string]float64{item: num}
	}

	for _, d := range val.Deps {
		n := getRawItemsFromItem(float64(d.Count)*num/float64(val.CraftNum), d.Name)
		for n, c := range n {
			out[n] += c
		}
	}

	return out
}

func calcBPItems(bs []building) map[string]int {
	buildingsCount := map[string]int{}
	for _, b := range bs {
		buildingsCount[fbot.resolveBuildingName(b.Name)]++
	}

	items := map[string]float64{}

	for building, count := range buildingsCount {
		for n, c := range getRawItemsFromItem(float64(count), building) {
			items[n] += c
		}
	}
	itemInts := map[string]int{}
	for k, v := range items {
		itemInts[k] = int(math.Ceil(v))
	}

	return itemInts
}

func tileArrayToBP(tiles []tile) []building {
	var bp []building
	wasLastUg := false
	for _, t := range tiles {
		typ := "input"
		name := "transport-belt"
		if t.Ug {
			name = "underground-belt"
			wasLastUg = true
		} else if wasLastUg {
			name = "underground-belt"
			wasLastUg = false
			typ = "output"
		}
		bp = append(bp, building{
			Name:     name,
			Rotation: t.Dir,
			Pos:      t.Pos,
			Ugbt:     typ,
		})
	}
	return bp
}

func findComponents(tilesOrig []position) (boxes []orePatch) { // divide graph into components but only store component bounds
	pos := position{}

	tiles := make([]position, len(tilesOrig))

	for i := range tilesOrig {
		tiles[i] = tilesOrig[i]
	}

	for len(tiles) > 0 { // iterate over all positions
		pos, tiles = tiles[0], tiles[1:]                                               // pop first value
		seen := []position{pos}                                                        // make this a set?
		q := []position{pos}                                                           // queue
		maxx, maxy, minx, miny := math.Inf(-1), math.Inf(-1), math.Inf(1), math.Inf(1) // component bounds
		for {
			pos, q = q[0], q[1:] // pop first value
			for _, ox := range [3]float64{-1.0, 0.0, 1.0} {
				for _, oy := range [3]float64{-1.0, 0.0, 1.0} { // for each direction from our position
					move := position{pos.X + ox, pos.Y + oy}            // move in that direction
					if contains(tiles, move) && !contains(seen, move) { // if that position is a ore and not visited already
						q = append(q, move)            // add it to queue
						tiles = removeVal(tiles, move) // and remove it from ore list
					}
				}
			}
			seen = append(seen, pos)
			maxx = math.Max(maxx, pos.X) // update box checks
			maxy = math.Max(maxy, pos.Y)
			minx = math.Min(minx, pos.X)
			miny = math.Min(miny, pos.Y)
			if len(q) == 0 {
				break
			}
		}
		boxes = append(boxes, orePatch{box{position{minx, miny}, position{maxx, maxy}}, true})
	}
	return
}

func isBorderSharedWithBox(box box, boxes []box, ignore int) bool {
	box.round()

	for i, b := range boxes {
		if i == ignore {
			continue
		}

		b.round()

		if box.Tl.X == b.Tl.X ||
			box.Tl.Y == b.Tl.Y ||
			box.Br.X == b.Br.Y ||
			box.Br.Y == b.Br.Y ||
			box.Tl.X == b.Br.X ||
			box.Tl.Y == b.Br.Y {
			return true
		}
	}

	return false
}

// takes bot inventory and item list and returns how many more of each item it needs.
func howMuchMore(have, need map[string]int) map[string]int {
	reallyNeed := map[string]int{}
	for item, count := range need {
		reallyNeed[item] = int(math.Max(float64(count-have[item]), 0))
	}
	return reallyNeed
}
