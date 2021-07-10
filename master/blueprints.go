package main

type Blueprint struct {
	Dims      Position
	Buildings []Building
}

var noFluidBp = Blueprint{
	Position{8, 3},
	[]Building{
		//BBIASMIB
		//BBLASM B
		//BB ASM B
		{
			Name: "assembling-machine-%d",
			Pos:  Position{4, 1},
		},
		{
			Name: "small-electric-pole",
			Pos:  Position{2, 2},
		},
		{
			Name: "small-electric-pole",
			Pos:  Position{6, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 2},
		},
		{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      Position{2, 0},
		},
		{
			Name:     "long-handed-inserter",
			Rotation: dirEast,
			Pos:      Position{2, 1},
		},
		{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      Position{6, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{7, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{7, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{7, 2},
		},
	},
}

var fluidBp = Blueprint{
	Position{10, 3},
	/*
	     BBIASMIB
	   PUBBUASM B
	   M BBLASM B
	*/
	[]Building{
		{
			Name: "assembling-machine-%d",
			Pos:  Position{5, 0},
		},
		{
			Name: "small-electric-pole",
			Pos:  Position{1, 0},
		},
		{
			Name: "small-electric-pole",
			Pos:  Position{8, 1},
		},
		{
			Name:     "pipe",
			Rotation: dirSouth,
			Pos:      Position{0, 1},
		},
		{
			Name:     "underground-pipe",
			Rotation: dirWest,
			Pos:      Position{1, 1},
		},
		{
			Name:     "underground-pipe",
			Rotation: dirEast,
			Pos:      Position{4, 1},
		},
		{
			Name:     "pump",
			Rotation: dirSouth,
			Pos:      Position{0, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{2, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{9, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{2, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{9, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{2, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{9, 2},
		},
		{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      Position{3, 0},
		},
		{
			Name:     "long-inserter",
			Rotation: dirEast,
			Pos:      Position{3, 2},
		},
		{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      Position{8, 0},
		},
	},
}

var minerBp = Blueprint{
	Position{8, 3},
	/*
	   MMMBMMMP
	   MMMBMMM
	   MMMBMMM
	*/
	[]Building{
		{
			Name:     "electric-mining-drill",
			Rotation: dirEast,
			Pos:      Position{1, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{3, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{3, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{3, 2},
		},
		{
			Name:     "electric-mining-drill",
			Rotation: dirWest,
			Pos:      Position{4.5, 1},
		},
		{
			Name: "small-electric-pole",
			Pos:  Position{7, 0},
		},
	},
}

var smeltingHeaderBp = Blueprint{
	Position{4, 11},
	[]Building{
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 0},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{2, 0},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{3, 0},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 1},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 2},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 3},
		},
		{
			Name:     "underground-belt",
			Rotation: dirEast,
			Pos:      Position{0, 4},
			Ugbt:     "input",
		},
		{
			Name:     "underground-belt",
			Rotation: dirWest,
			Pos:      Position{2, 4},
			Ugbt:     "output",
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{3, 4},
		},
		{
			Name:     "belt",
			Rotation: dirWest,
			Pos:      Position{3, 5},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 4},
		},
		{
			Name:     "splitter",
			Rotation: dirEast,
			Pos:      Position{0, 5.5},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 5},
		},
		{
			Name:     "splitter",
			Rotation: dirWest,
			Pos:      Position{2.5, 5.5},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 6},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 7},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 8},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 9},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 10},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{2, 10},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{3, 10},
		},
	},
}

var smeltingBp = Blueprint{
	Position{2, 12},
	[]Building{
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{0, 0},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 0},
		},
		{
			Name:     "inserter",
			Rotation: dirNorth,
			Pos:      Position{0, 1},
		},
		{
			Name: "furnace",
			Pos:  Position{1, 3},
		},
		{
			Name:     "inserter",
			Rotation: dirNorth,
			Pos:      Position{0, 4},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{0, 5},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 5},
		},
		{
			Name:     "inserter",
			Rotation: dirSouth,
			Pos:      Position{0, 6},
		},
		{
			Name: "furnace",
			Pos:  Position{1, 7.5},
		},
		{
			Name:     "inserter",
			Rotation: dirSouth,
			Pos:      Position{0, 9},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{0, 10},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 10},
		},
	},
}

var smeltingFooterBp = Blueprint{
	Position{9, 12},
	[]Building{
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 0},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 1},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 2},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 3},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 4},
		},
		{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 5},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{0, 6},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 6},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{2, 6},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{3, 6},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{4, 6},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{5, 6},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{6, 6},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{7, 6},
		},
		{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{8, 6},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{0, 7},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{0, 8},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{0, 9},
		},
		{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{0, 10},
		},
	},
}
