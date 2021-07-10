package main

type blueprint struct {
	Dims      position
	Buildings []building
}

var noFluidBp = blueprint{
	position{8, 3},
	[]building{
		//BBIASMIB
		//BBLASM B
		//BB ASM B
		{
			Name: "assembling-machine-%d",
			Pos:  position{4, 1},
		},
		{
			Name: "small-electric-pole",
			Pos:  position{2, 2},
		},
		{
			Name: "small-electric-pole",
			Pos:  position{6, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 2},
		},
		{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      position{2, 0},
		},
		{
			Name:     "long-handed-inserter",
			Rotation: dirEast,
			Pos:      position{2, 1},
		},
		{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      position{6, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{7, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{7, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{7, 2},
		},
	},
}

var fluidBp = blueprint{
	position{10, 3},
	/*
	     BBIASMIB
	   PUBBUASM B
	   M BBLASM B
	*/
	[]building{
		{
			Name: "assembling-machine-%d",
			Pos:  position{5, 0},
		},
		{
			Name: "small-electric-pole",
			Pos:  position{1, 0},
		},
		{
			Name: "small-electric-pole",
			Pos:  position{8, 1},
		},
		{
			Name:     "pipe",
			Rotation: dirSouth,
			Pos:      position{0, 1},
		},
		{
			Name:     "underground-pipe",
			Rotation: dirWest,
			Pos:      position{1, 1},
		},
		{
			Name:     "underground-pipe",
			Rotation: dirEast,
			Pos:      position{4, 1},
		},
		{
			Name:     "pump",
			Rotation: dirSouth,
			Pos:      position{0, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{2, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{9, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{2, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{9, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{2, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{9, 2},
		},
		{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      position{3, 0},
		},
		{
			Name:     "long-inserter",
			Rotation: dirEast,
			Pos:      position{3, 2},
		},
		{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      position{8, 0},
		},
	},
}

var minerBp = blueprint{
	position{8, 3},
	/*
	   MMMBMMMP
	   MMMBMMM
	   MMMBMMM
	*/
	[]building{
		{
			Name:     "electric-mining-drill",
			Rotation: dirEast,
			Pos:      position{1, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{3, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{3, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{3, 2},
		},
		{
			Name:     "electric-mining-drill",
			Rotation: dirWest,
			Pos:      position{4.5, 1},
		},
		{
			Name: "small-electric-pole",
			Pos:  position{7, 0},
		},
	},
}

var smeltingHeaderBp = blueprint{
	position{4, 11},
	[]building{
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{1, 0},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{2, 0},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{3, 0},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{1, 1},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{1, 2},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{1, 3},
		},
		{
			Name:     "underground-belt",
			Rotation: dirEast,
			Pos:      position{0, 4},
			Ugbt:     "input",
		},
		{
			Name:     "underground-belt",
			Rotation: dirWest,
			Pos:      position{2, 4},
			Ugbt:     "output",
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{3, 4},
		},
		{
			Name:     "belt",
			Rotation: dirWest,
			Pos:      position{3, 5},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{1, 4},
		},
		{
			Name:     "splitter",
			Rotation: dirEast,
			Pos:      position{0, 5.5},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{1, 5},
		},
		{
			Name:     "splitter",
			Rotation: dirWest,
			Pos:      position{2.5, 5.5},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 6},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 7},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 8},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{1, 9},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{1, 10},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{2, 10},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{3, 10},
		},
	},
}

var smeltingBp = blueprint{
	position{2, 12},
	[]building{
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{0, 0},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{1, 0},
		},
		{
			Name:     "inserter",
			Rotation: dirNorth,
			Pos:      position{0, 1},
		},
		{
			Name: "furnace",
			Pos:  position{1, 3},
		},
		{
			Name:     "inserter",
			Rotation: dirNorth,
			Pos:      position{0, 4},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{0, 5},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{1, 5},
		},
		{
			Name:     "inserter",
			Rotation: dirSouth,
			Pos:      position{0, 6},
		},
		{
			Name: "furnace",
			Pos:  position{1, 7.5},
		},
		{
			Name:     "inserter",
			Rotation: dirSouth,
			Pos:      position{0, 9},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{0, 10},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{1, 10},
		},
	},
}

var smeltingFooterBp = blueprint{
	position{9, 12},
	[]building{
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 3},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 4},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      position{0, 5},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      position{0, 6},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{0, 7},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{0, 8},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{0, 9},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      position{0, 10},
		},
	},
}
