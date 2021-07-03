package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

var recipes map[string]Item

type RecipeDep struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
	MakeFactory bool
}

type Item struct {
	CraftTime float32     `json:"craftTime"` // craft_time / craft_amount
	Liquid    bool        `json:"liquid"`
	Deps      []RecipeDep `json:"deps"`
}

type Building struct {
	Name      string   `json:"name"`
	Rotation  int      `json:"rotation"`
	CraftItem string   `json:"recipe"`
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

func (b *Bot) newFactory(itemStr string, ps float32) ([]Building, error) {
	bp := noFluidBp
	item, exists := recipes[itemStr]
	if !exists {
		return nil, errors.New("unknown item " + itemStr)
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
		for _, building := range bp.Buildings {
			building.Pos.Y += bp.Dims.Y * float64(i)
			out[bCount] = building

			if strings.HasPrefix(out[bCount].Name, "assembling-machine") {
				out[bCount].CraftItem = itemStr
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

	return out, nil
}

func (b *Bot) newMiners(patch OrePatch) []Building {
	wcount := int(math.Abs(math.Ceil((patch.Dims.Tl.X - patch.Dims.Br.X) / minerBp.Dims.X)))
	hcount := int(math.Abs(math.Ceil((patch.Dims.Tl.Y - patch.Dims.Br.Y) / minerBp.Dims.Y)))

	out := make([]Building, wcount * hcount * len(minerBp.Buildings) + hcount + wcount) // count of bps * buildings in bp + additional poles
	bCount := 0

	for h:=0; h < hcount; h++ { // add poles on the left
		out[bCount] = Building{
			"small-electric-pole", up, "",
			Position{patch.Dims.Tl.X - 1, float64(h) * minerBp.Dims.Y + patch.Dims.Tl.Y}}

		bCount++
	}

	for w:=0; w < wcount; w++ { // add poles connecting all collumns
		out[bCount] = Building{
			"small-electric-pole", up, "",
			Position{float64(w) * minerBp.Dims.X + patch.Dims.Tl.X, patch.Dims.Tl.Y - 1}}

		bCount++
	}

	for w:=0; w < wcount; w++ {
		for h:=0; h < hcount; h++ {

			for _, building := range minerBp.Buildings {
				out[bCount] = building
				out[bCount].Pos.X += float64(w) * minerBp.Dims.X + patch.Dims.Tl.X
				out[bCount].Pos.Y += float64(h) * minerBp.Dims.Y + patch.Dims.Tl.Y

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
	}

	return out
}
