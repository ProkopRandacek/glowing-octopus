package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type sharedDepLocation struct {
	Name string
	Pos  position
	Left float64
}

type recipeDep struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type item struct {
	CraftTime float64     `json:"craftTime"` // craft_time / craft_amount
	CraftNum  int         `json:"resultCount"`
	Liquid    bool        `json:"liquid"`
	Deps      []recipeDep `json:"deps"`
}

type building struct {
	Name      string   `json:"name"`
	Rotation  int      `json:"rotation"`
	CraftItem string   `json:"recipe"`
	Pos       position `json:"pos"`
	Ugbt      string   `json:"ugbt"` // underground belt type - "input" or "output"
}

// Returns an index in bot.SharedResources for `amount` of `item`.
// If it can't find any, it will try to create a new factory for said
// resource. If that fails, returns an error.
func (b *bot) findSharedResource(item string, amount float64) (int, error) {
	// go through available resources and check for a fitting one
	for i, r := range b.SharedResources[item] {
		if r.Left >= amount {
			b.SharedResources[item][i].Left -= amount
			return i, nil
		}
	}

	//// if the code gets here, it means there wasn't found any fitting location

	// try to create name of raw resource to `item`
	oreName := strings.Split(item, "-")[0]
	if item != "coal" { // coal is only coal
		oreName += "-ore"
	}

	// find the closest resource
	state, _ := b.state()
	botPos := state.Pos
	minDist := -1.0
	minIndex := -1
	for i, p := range b.Mapper.OrePatches[oreName] {
		if dist := math.Sqrt(math.Pow(botPos.X-p.Dims.Tl.X, 2) +
				math.Pow(botPos.Y-p.Dims.Tl.Y, 2));
		   dist < minDist || minIndex == -1 {

			minDist = dist
			minIndex = i
		}
	}

	// failed to find a resource
	if minIndex == -1 {
		return -1, errors.New("no patches ore resource locations found. Explore some more")
	}

	// build miner and smelters for the resource
	p := b.Mapper.OrePatches[oreName][minIndex]
	_, output := b.newMiners(p)
	pos, err := b.newSmelters(output)
	if err != nil {
		return -1, err
	}

	// add the resource location to bot and remove the or patch
	fmt.Printf("adding %f of %s at %v to shared resources\n", output, item, pos)
	b.SharedResources[item] = append(
		b.SharedResources[item],
		sharedDepLocation{Name: item, Pos: pos, Left: output})
	b.Mapper.OrePatches[oreName] = append(
		b.Mapper.OrePatches[oreName][:minIndex],
		b.Mapper.OrePatches[oreName][minIndex+1:]...)

	return b.findSharedResource(item, amount)
}

// returns the building name with current available level in mind
func (b *bot) resolveBuildingName(curr string) string {
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

// Generates the blueprint for a factory at position `pos`.
func (b *bot) genFactory(pos position, itemName string, asmCount int, bp blueprint) []building {
	out := make([]building, asmCount*len(bp.Buildings))

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

// Creates a new factory that manufactures `ps` of `itemStr` per second.
// Returns the position, where the factory was built.
// It automatically builds it's dependencies.
func (b *bot) newFactory(itemStr string, ps float64) (position, error) {
	bp := noFluidBp
	item, exists := recipes[itemStr]
	if !exists {
		return position{}, errors.New("unknown item " + itemStr)
	}

	asmCount := int(math.Ceil(ps * item.CraftTime)) // count of assemblers needed

	// allocate space
	space := box{Tl: position{0, 0}, Br: position{bp.Dims.X, bp.Dims.Y * float64(asmCount)}}
	b.Mapper.findSpace(&space)
	_, err := b.Mapper.alloc(space)
	if err != nil {
		return position{}, errors.New("could not find space to allocate")
	}

	// resolve deps
	resources := map[string]position{} // map of places where a path needs to be built from
	for i, d := range item.Deps {
		if d.Name != "iron-plate" && d.Name != "copper-plate" {
			fmt.Printf("%s needs to be built.\n", d.Name)
			var err error
			resources[d.Name], err = b.newFactory(d.Name, ps*float64(d.Count))
			if err != nil {
				return position{}, err
			}

		} else {
			fmt.Printf("%s is shared. I will look for it.\n", d.Name)
			index, err := b.findSharedResource(d.Name, ps*float64(d.Count))
			if err != nil {
				return position{}, err
			}

			resources[d.Name] = b.SharedResources[d.Name][index].Pos
			b.SharedResources[d.Name][index].Pos = position{space.Tl.X + float64(i), space.Br.Y + 1}
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
	b.clearAll(space)
	b.waitForTaskDone()
	fmt.Println("cleared space for factory")

	buildings := b.genFactory(space.Tl, itemStr, asmCount, bp)

	// generate and build paths for all resources
	for i:=0; key := range resources; i++ {
		fmt.Printf(
			"building path for resource %s from %v to %v\n",
			key, resources[key],
			position{space.Tl.X + float64(i), space.Tl.Y - 1})

		path := b.Mapper.findBeltPath(resources[key], position{space.Tl.X + float64(i), space.Tl.Y - 2})
		b.Mapper.allocBelt(path)
		pathBp := tileArrayToBP(path)
		b.build(pathBp)
		b.waitForTaskDone()
	}

	fmt.Println("building factory")
	b.build(buildings)
	b.waitForTaskDone()

	return position{space.Br.X, space.Br.Y + 2}, nil
}

// TODO: very slow
func (b *bot) shouldBuildMiner(mPos position) bool {
	for _, tiles := range b.Mapper.Resources {
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

// Build new miners mining `patch`.
// Returns position where they are built and how much
// they mine per second.
func (b *bot) newMiners(patch orePatch) (position, float64) {
	wCount := int(math.Abs(math.Ceil((patch.Dims.Tl.X - patch.Dims.Br.X) / minerBp.Dims.X)))
	hCount := int(math.Abs(math.Ceil((patch.Dims.Tl.Y - patch.Dims.Br.Y) / minerBp.Dims.Y)))

	space := box{position{patch.Dims.Tl.X - 1, patch.Dims.Tl.Y - 1}, position{patch.Dims.Br.X, patch.Dims.Br.Y + 1}}
	b.clearAll(space)
	b.waitForTaskDone()

	out := make([]building, wCount*hCount*len(minerBp.Buildings)+hCount+wCount+int(math.Abs(patch.Dims.Tl.X-patch.Dims.Br.X))) // count of bps * buildings in bp + additional poles
	bCount := 0

	for h := 0; h < hCount; h++ { // add poles on the left
		out[bCount] = building{
			Name:     "small-electric-pole",
			Rotation: dirNorth,
			Pos:      position{patch.Dims.Tl.X - 1, float64(h)*minerBp.Dims.Y + patch.Dims.Tl.Y}}

		bCount++
	}

	for w := 0; w < wCount; w++ { // add poles connecting all columns
		out[bCount] = building{
			Name:     "small-electric-pole",
			Rotation: dirNorth,
			Pos:      position{float64(w)*minerBp.Dims.X + patch.Dims.Tl.X, patch.Dims.Tl.Y - 1}}

		bCount++
	}

	minerCount := 0
	for w := 0; w < wCount; w++ {
		for h := 0; h < hCount; h++ {
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
		out[bCount] = building{Name: b.resolveBuildingName("belt"), Rotation: dirEast, Pos: position{w, patch.Dims.Br.Y}}
		bCount++
	}

	b.build(out[:bCount])
	b.waitForTaskDone()

	return space.Br, float64(minerCount) / 0.5 // rate is 0.5/s except for uranium TODO
}

func (b *bot) newSmelters(maxInput float64) (position, error) {
	beltMax := 15.0
	if b.BeltLevel == "fast" {
		beltMax = 30.0
	}

	if maxInput > beltMax {
		//return box{}, errors.New("Too much input for one smelting array")
		maxInput = beltMax
	}

	furnaceCount := int(math.Ceil(48.0 * (maxInput / beltMax) / 2.0)) // Divided by 2, since 2 are in 1 bp

	space := box{position{0, 0},
		position{
			smeltingBp.Dims.X*float64(furnaceCount) + smeltingHeaderBp.Dims.X,
			smeltingBp.Dims.Y}}

	b.Mapper.findSpace(&space)
	_, err := b.Mapper.alloc(space)
	if err != nil {
		return position{}, errors.New("could not find space to allocate")
	}
	b.clearAll(space)
	b.waitForTaskDone()

	out := make([]building, 0,
		len(smeltingHeaderBp.Buildings) + 
		furnaceCount*len(smeltingBp.Buildings) +
		len(smeltingFooterBp.Buildings))

	bCount := 0
	out = append(out, b.repeatBp(
		smeltingHeaderBp.Buildings,
		1, space.Tl.X, space.Tl.Y))

	for i := 0; i < furnaceCount; i++ {
		out = append(out, b.repeatBp(
			smeltingBp.Buildings,
			1,
			space.Tl.X + float64(i)*smeltingBp.Dims.X +
				smeltingHeaderBp.Dims.X,
				space.Tl.Y))
	}

	out = append(out, b.repeatBp(
		smeltingFooterBp.Buildings,
		1,
		space.Tl.X + smeltingBp.Dims.X +
			smeltingHeaderBp.Dims.X,
			space.Tl.Y))

	b.build(out)
	b.waitForTaskDone()

	return position{
		smeltingHeaderBp.Dims.X + smeltingBp.Dims.X*float64(furnaceCount) + 2,
		6,
	}, nil
}
