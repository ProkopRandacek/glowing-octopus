package main

import (
	"fmt"
	"math"
)

type RecipeDep struct {
	ItemName    string `json:"name"`
	Count       int    `json:"count"`
	MakeFactory bool
}

type Item struct {
	Name      string  `json:"name"`
	CraftTime float32 `json:"craftTime"` // craft_time / craft_amount
	Liquid    bool `json:"liquid"`
	Deps      []RecipeDep `json:"deps"`
}

const (
	up = iota
	down
	left
	right
)

type Building struct {
	Name      string   `json:"name"`
	Rotation  int      `json:"rotation"`
	CraftItem string   `json"craft_item"`
	Pos       Position `json:"pos"`
}

func (b *Bot) newFactory(item Item, ps float32) {
	bp := noFluidBp
	
	for _, d := range item.Deps {
		if d.MakeFactory {
			b.newFactory(d.ItemStruct, ps*float32(d.Count))
		}

		if d.ItemStruct.Liquid {
			bp = fluidBp
		}
	}

	asmCount := int(math.Ceil(float64(ps * item.CraftTime)))

	//fmt.Printf("asm count for %s: %d\n", item.Name, asmCount)

	out := make([]Building, asmCount * len(bp.Buildings))

	bCount := 0
	for i:=0; i < asmCount; i++ {
		for j, b := range bp.Buildings {
			b.Pos.Y += bp.Dims.Y * float64(i)
			out[bCount] = b

			if j == 0 {
				out[bCount].CraftItem = item.Name
			}

			bCount++
		}
	}

	fmt.Println(out)
}
