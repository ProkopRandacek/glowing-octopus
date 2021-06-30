package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
)

var recipes map[string]Item

type RecipeDep struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
	MakeFactory bool
}

type Item struct {
	Name      string      `json:"name"`
	CraftTime float32     `json:"craftTime"` // craft_time / craft_amount
	Liquid    bool        `json:"liquid"`
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

func loadRecipes() error {
	f, err := os.Open("./master/recipes.json")
	if err != nil {
		return err
	}
	defer f.Close()

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

func (b *Bot) newFactory(itemStr string, ps float32) error {
	bp := noFluidBp
	item, exists := recipes[itemStr]
	if !exists {
		return errors.New("unknown item " + itemStr)
	}

	for _, d := range item.Deps {
		//if d.MakeFactory { TODO
		b.newFactory(d.Name, ps*float32(d.Count))
		//}

		if recipes[d.Name].Liquid {
			bp = fluidBp
		}
	}

	asmCount := int(math.Ceil(float64(ps * item.CraftTime))) // count of assemblers needed

	out := make([]Building, asmCount*len(bp.Buildings))

	bCount := 0 // count of building placed
	for i := 0; i < asmCount; i++ {
		for j, building := range bp.Buildings {
			building.Pos.Y += bp.Dims.Y * float64(i)
			out[bCount] = building

			if j == 0 {
				out[bCount].CraftItem = item.Name
				out[bCount].Name = fmt.Sprintf(out[bCount].Name, b.AssemblerLevel)
			}

			if out[bCount].Name == "inserter" && len(b.InserterLevel) > 0 {
				out[bCount].Name = b.InserterLevel + "-inserter"
			}

			if out[bCount].Name == "belt" {
				out[bCount].Name = ""
				if len(b.BeltLevel) > 0 {
					out[bCount].Name = b.BeltLevel + "-"
				}

				out[bCount].Name += "transport-belt"
			}

			bCount++
		}
	}

	t, _ := json.Marshal(out)
	fmt.Println(string(t))

	return nil
}
