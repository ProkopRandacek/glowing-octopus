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

type SharedDepLocation struct {
	Name string
	Pos Position
	Left float64
}

type RecipeDep struct {
	Name        string `json:"name"`
	Count       int    `json:"count"`
}

type Item struct {
	CraftTime float64     `json:"craftTime"` // craft_time / craft_amount
	CraftNum  int         `json:"resultCount"`
	Liquid    bool        `json:"liquid"`
	Deps      []RecipeDep `json:"deps"`
}

type Building struct {
	Name      string   `json:"name"`
	Rotation  int      `json:"rotation"`
	CraftItem string   `json:"recipe"`
	Pos       Position `json:"pos"`
	Ugbt      string   `json:"ugbt"` // underground belt type - "input" or "output"
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

// returns the index in b.SharedResources[item] or error
func (b *Bot) findSharedResource(item string, amount float64) (int, error) {
	for i, r := range b.SharedResources[item] {
		if r.Left >= amount {
			b.SharedResources[item][i].Left -= amount
			return i, nil
		}
	}

	// if the code gets here, it means there wasn't found any fitting location
	oreName := strings.Split(item, "-")[0]
	if item != "coal" { // coal is only coal
		oreName += "-ore"
	}

	for i, p := range b.Mapper.OrePatches[oreName] {
		pos, output := b.newMiners(p)
		pos, err := b.newSmelters(pos, output)
		if err != nil {
			return -1, err
		}

		fmt.Printf("adding %f of %s at %v to shared resources\n", output, item, pos)
		b.SharedResources[item] = append(b.SharedResources[item], SharedDepLocation{Name: item, Pos: pos, Left: output}) // add the pos to shared resources
		b.Mapper.OrePatches[oreName] = append(b.Mapper.OrePatches[oreName][:i], b.Mapper.OrePatches[oreName][i+1:]...) // remove the patch from mapper, as it no longer can be used
		return b.findSharedResource(item, amount)
	}

	return -1, errors.New("No patches ore resource locations found. Explore some more")
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

func (b *Bot) genFactory(pos Position, itemName string, asmCount int, bp Blueprint) []Building {
	out := make([]Building, asmCount*len(bp.Buildings))

	bCount := 0 // count of building placed
	for i := 0; i < asmCount; i++ {
		for _, building := range bp.Buildings {
			building.Pos.Y += bp.Dims.Y*float64(i) + pos.Y
			building.Pos.X += pos.X
			out[bCount] = building

			if strings.HasPrefix(out[bCount].Name, "assembling-machine") {
				out[bCount].CraftItem = itemName
			}

			out[bCount].Name = b.resolveBuildingName(out[bCount].Name)

			bCount++
		}
	}

	return out
}

func (b *Bot) newFactory(itemStr string, ps float64) (Position, error) {
	bp := noFluidBp
	item, exists := recipes[itemStr]
	if !exists {
		return Position{}, errors.New("unknown item " + itemStr)
	}

	resources := map[string]Position{}

	asmCount := int(math.Ceil(float64(ps * item.CraftTime))) // count of assemblers needed
	space := Box{Tl: Position{0, 0}, Br: Position{bp.Dims.X, bp.Dims.Y * float64(asmCount)}}
	b.Mapper.findSpace(&space)
	if b.Mapper.alloc(space) < 0 {
		return Position{}, errors.New("Could not find space to allocate")
	}

	for i, d := range item.Deps {
		if d.Name != "iron-plate" && d.Name != "copper-plate" {
			fmt.Printf("%s needs to be built.\n", d.Name)
			var err error
			resources[d.Name], err = b.newFactory(d.Name, ps*float64(d.Count))
			if err != nil {
				return Position{}, err
			}

		} else {
			fmt.Printf("%s is shared. I will look for it.\n", d.Name)
			index, err := b.findSharedResource(d.Name, ps*float64(d.Count))
			if err != nil {
				return Position{}, err
			}

			resources[d.Name] = b.SharedResources[d.Name][index].Pos
			b.SharedResources[d.Name][index].Pos = Position{space.Tl.X + float64(i), space.Br.Y + 1}
			if i == 1 {
				// TODO add splitter
			}
		}

		fmt.Printf("resolved dependency for %s\n", d.Name)
		if recipes[d.Name].Liquid {
			bp = fluidBp
		}
	}

	fmt.Println("all dependencies resolved")
	fmt.Println("allocated space for factory")
	b.clearAll(space)
	b.waitForTaskDone()
	fmt.Println("cleared space for factory")

	buildings := b.genFactory(space.Tl, itemStr, asmCount, bp)

	i := 0
	for key := range resources {
		fmt.Printf("building path for resource %s from %v to %v\n", key, resources[key], Position{space.Tl.X + float64(i), space.Tl.Y - 1})
		path := b.Mapper.FindBeltPath(resources[key], Position{space.Tl.X + float64(i), space.Tl.Y - 2})
		b.Mapper.AllocBelt(path)
		pathBp := b.Mapper.TileArrayToBP(path)
		b.build(pathBp)
		b.waitForTaskDone()
		i++
	}

	fmt.Println("building factory")
	b.build(buildings)
	b.waitForTaskDone()

	return Position{space.Br.X, space.Br.Y + 2}, nil
}

func (b *Bot) shouldBuildMiner(mPos Position) bool {
	for _, tiles := range b.Mapper.Resrcs {
		for _, tile := range tiles {
			if math.Abs(tile.X) >= math.Abs(mPos.X)-1.5 &&
				math.Abs(tile.Y) >= math.Abs(mPos.Y)-1.5 &&
				math.Abs(tile.X) <= math.Abs(mPos.X)+1.5 &&
				math.Abs(tile.Y) >= math.Abs(mPos.Y)+1.5 {
				return true
			}
		}
	}

	return false
}

func (b *Bot) newMiners(patch OrePatch) (Position, float64) {
	wcount := int(math.Abs(math.Ceil((patch.Dims.Tl.X - patch.Dims.Br.X) / minerBp.Dims.X)))
	hcount := int(math.Abs(math.Ceil((patch.Dims.Tl.Y - patch.Dims.Br.Y) / minerBp.Dims.Y)))

	space := Box{Position{patch.Dims.Tl.X - 1, patch.Dims.Tl.Y - 1}, Position{patch.Dims.Br.X, patch.Dims.Br.Y + 1}}
	b.clearAll(space)
	b.waitForTaskDone()

	out := make([]Building, wcount*hcount*len(minerBp.Buildings)+hcount+wcount+int(math.Abs(patch.Dims.Tl.X-patch.Dims.Br.X))) // count of bps * buildings in bp + additional poles
	bCount := 0

	for h := 0; h < hcount; h++ { // add poles on the left
		out[bCount] = Building{
			Name:     "small-electric-pole",
			Rotation: dirNorth,
			Pos:      Position{patch.Dims.Tl.X - 1, float64(h)*minerBp.Dims.Y + patch.Dims.Tl.Y}}

		bCount++
	}

	for w := 0; w < wcount; w++ { // add poles connecting all collumns
		out[bCount] = Building{
			Name:     "small-electric-pole",
			Rotation: dirNorth,
			Pos:      Position{float64(w)*minerBp.Dims.X + patch.Dims.Tl.X, patch.Dims.Tl.Y - 1}}

		bCount++
	}

	minerCount := 0
	for w := 0; w < wcount; w++ {
		for h := 0; h < hcount; h++ {

			for _, building := range minerBp.Buildings {
				out[bCount] = building
				out[bCount].Pos.X += float64(w)*minerBp.Dims.X + patch.Dims.Tl.X
				out[bCount].Pos.Y += float64(h)*minerBp.Dims.Y + patch.Dims.Tl.Y

				if building.Name == "electric-mining-drill" && !b.shouldBuildMiner(out[bCount].Pos) {
					bCount--
					continue
				}

				if building.Name == "electric-mining-drill" {
					minerCount++
				}

				out[bCount].Name = b.resolveBuildingName(out[bCount].Name)

				bCount++
			}
		}
	}

	for w := patch.Dims.Tl.X; w < patch.Dims.Br.X; w++ {
		out[bCount] = Building{Name: b.resolveBuildingName("belt"), Rotation: dirEast, Pos: Position{w, patch.Dims.Br.Y}}
		bCount++
	}

	b.build(out[:bCount])
	b.waitForTaskDone()

	return space.Br, float64(minerCount) / 0.5 // rate is 0.5/s except for uranium TODO
}

func (b *Bot) newSmelters(resPos Position, maxInput float64) (Position, error) {
	beltMax := 15.0
	if b.BeltLevel == "fast" {
		beltMax = 30.0
	}

	if maxInput > beltMax {
		//return Box{}, errors.New("Too much input for one smelting array")
		maxInput = beltMax
	}

	furnaceCount := int(math.Ceil(48.0 * (maxInput / beltMax) / 2.0)) // Divided by 2, since 2 are in 1 bp

	space := Box{Position{0, 0},
		Position{
			smeltingBp.Dims.X*float64(furnaceCount) + smeltingHeaderBp.Dims.X,
			smeltingBp.Dims.Y}}

	b.Mapper.findSpace(&space)
	if b.Mapper.alloc(space) < 0 {
		return Position{}, errors.New("Could not find space to allocate")
	}
	b.clearAll(space)
	b.waitForTaskDone()

	out := make([]Building, len(smeltingHeaderBp.Buildings)+furnaceCount*len(smeltingBp.Buildings)+len(smeltingFooterBp.Buildings))

	bCount := 0
	for _, building := range smeltingHeaderBp.Buildings {
		building.Pos.X += space.Tl.X
		building.Pos.Y += space.Tl.Y
		out[bCount] = building
		out[bCount].Name = b.resolveBuildingName(out[bCount].Name)
		bCount++
	}

	for i := 0; i < furnaceCount; i++ {
		for _, building := range smeltingBp.Buildings {
			out[bCount] = building
			out[bCount].Pos.X += space.Tl.X + float64(i)*smeltingBp.Dims.X + smeltingHeaderBp.Dims.X
			out[bCount].Pos.Y += space.Tl.Y
			out[bCount].Name = b.resolveBuildingName(out[bCount].Name)
			bCount++
		}
	}

	for _, building := range smeltingFooterBp.Buildings {
		building.Pos.X += space.Tl.X + smeltingHeaderBp.Dims.X + smeltingBp.Dims.X*float64(furnaceCount)
		building.Pos.Y += space.Tl.Y
		out[bCount] = building
		out[bCount].Name = b.resolveBuildingName(out[bCount].Name)
		bCount++
	}

	b.build(out)
	b.waitForTaskDone()

	return Position{
		smeltingHeaderBp.Dims.X + smeltingBp.Dims.X*float64(furnaceCount) + 1,
		6,
	}, nil
}
