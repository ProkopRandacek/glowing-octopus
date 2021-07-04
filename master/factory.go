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


func (b *Bot) resolveBuildingName(curr string) string {
	switch curr {
	case "assembling-machine-%d":
		return fmt.Sprintf(curr, b.AssemblerLevel)
	case "inserter":
		if len(b.InserterLevel) > 0 {
			return b.InserterLevel + "-inserter"
		}
	case "belt":
		curr = ""
		if len(b.BeltLevel) > 0 {
			curr = b.BeltLevel + "-"
		}

		return curr + "transport-belt"
	case "underground-belt":
		curr = ""
		if len(b.BeltLevel) > 0 {
			curr = b.BeltLevel + "-"
		}

		return curr + "underground-belt"
	case "furnace":
		return b.FurnaceLevel + "-furnace"
	}

	return curr
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
			}

			out[bCount].Name = b.resolveBuildingName(out[bCount].Name)

			bCount++
		}
	}

	return out, nil
}

func (b *Bot) shouldBuildMiner(mPos Position) bool {
	for _, tiles := range b.Mapper.Resrcs {
		for _, tile := range tiles {
			if math.Abs(tile.X) >= math.Abs(mPos.X) - 1.5 &&
				math.Abs(tile.Y) >= math.Abs(mPos.Y) - 1.5 &&
				math.Abs(tile.X) <= math.Abs(mPos.X) + 1.5 &&
				math.Abs(tile.Y) >= math.Abs(mPos.Y) + 1.5 {
					return true
				}
		}
	}

	return false
}

func (b *Bot) newMiners(patch OrePatch) []Building {
	wcount := int(math.Abs(math.Ceil((patch.Dims.Tl.X - patch.Dims.Br.X) / minerBp.Dims.X)))
	hcount := int(math.Abs(math.Ceil((patch.Dims.Tl.Y - patch.Dims.Br.Y) / minerBp.Dims.Y)))

	out := make([]Building, wcount * hcount * len(minerBp.Buildings) + hcount + wcount) // count of bps * buildings in bp + additional poles
	bCount := 0

	for h:=0; h < hcount; h++ { // add poles on the left
		out[bCount] = Building{
			"small-electric-pole", dirNorth, "",
			Position{patch.Dims.Tl.X - 1, float64(h) * minerBp.Dims.Y + patch.Dims.Tl.Y}}

		bCount++
	}

	for w := 0; w < wcount; w++ { // add poles connecting all collumns
		out[bCount] = Building{
			"small-electric-pole", dirNorth, "",
			Position{float64(w)*minerBp.Dims.X + patch.Dims.Tl.X, patch.Dims.Tl.Y - 1}}

		bCount++
	}

	for w := 0; w < wcount; w++ {
		for h := 0; h < hcount; h++ {

			for _, building := range minerBp.Buildings {
				out[bCount] = building
				out[bCount].Pos.X += float64(w)*minerBp.Dims.X + patch.Dims.Tl.X
				out[bCount].Pos.Y += float64(h)*minerBp.Dims.Y + patch.Dims.Tl.Y

				if building.Name == "electric-mining-drill" && !b.shouldBuildMiner(out[bCount].Pos) {
					fmt.Println(out[bCount])
					bCount--
					continue
				}

				out[bCount].Name = b.resolveBuildingName(out[bCount].Name)

				bCount++
			}
		}
	}

	return out[:bCount]
}


func (b *Bot) newSmelters(maxInput float64) []Building {
	beltMax := 15.0
	if b.BeltLevel == "fast" {
		beltMax = 30.0
	}

	if maxInput > beltMax {
		return []Building{}
	}

	furnaceCount := int(math.Ceil(48.0 * (maxInput/beltMax) / 2.0)) // Divided by 2, since 2 are in 1 bp
	fmt.Println(furnaceCount)

	out := make([]Building, len(smeltingHeaderBp.Buildings) + furnaceCount * len(smeltingBp.Buildings))

	bCount := 0
	for _, building := range smeltingHeaderBp.Buildings {
		out[bCount] = building
		out[bCount].Name = b.resolveBuildingName(out[bCount].Name)
		bCount++
	}

	for i:=0; i < furnaceCount; i++ {
		for _, building := range smeltingBp.Buildings {
			out[bCount] = building
			out[bCount].Pos.X += float64(i) * smeltingBp.Dims.X + smeltingHeaderBp.Dims.X
			out[bCount].Name = b.resolveBuildingName(out[bCount].Name)
			bCount++
		}
	}

	return out
}
