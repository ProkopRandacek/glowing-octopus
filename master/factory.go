package main

import (
	"fmt"
	"math"
)

type RecipeDep struct {
	ItemStruct Item
	Count int
	MakeFactory bool
}

type Item struct {
	Name string
	CraftTime float32 // craft_time / craft_amount
	Liquid bool
	Deps []RecipeDep
}

const (
	up = iota
	down
	left
	right
)

type Building struct {
	Name string `json:"name"`
	Rotation int `json:"rotation"`
	CraftItem string `json"craft_item"`
	Pos Position `json:"pos"`
}

func (b *Bot) newFactory(item Item, ps float32) {
	for _, d := range item.Deps {
		if d.MakeFactory {
			b.newFactory(d.ItemStruct, ps * float32(d.Count))
		}
	}

	asmCount := int(math.Ceil(float64(ps * item.CraftTime)))

	fmt.Printf("asm count for %s: %d\n", item.Name, asmCount)

	bp := twoItemBp

	out := make([]Building, int(float64(asmCount) * bp.Dims.X * bp.Dims.Y))

	for i:=0; i < asmCount; i++ {
		for j, b := range bp.Buildings {
			b.Pos.Y += bp.Dims.Y * float64(i)
			out[i*int(bp.Dims.X) + j] = b
			//fmt.Println(b)
		}
	}

	fmt.Println(out)
}
