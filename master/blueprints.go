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
		Building{
			Name: "assembling-machine-%d",
			Pos:  Position{4, 1},
		},
		Building{
			Name: "small-electric-pole",
			Pos:  Position{2, 2},
		},
		Building{
			Name: "small-electric-pole",
			Pos:  Position{6, 2},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 1},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{0, 2},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 1},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 2},
		},
		Building{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      Position{2, 0},
		},
		Building{
			Name:     "long-handed-inserter",
			Rotation: dirEast,
			Pos:      Position{2, 1},
		},
		Building{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      Position{6, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{7, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{7, 1},
		},
		Building{
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
		Building{
			Name: "assembling-machine-%d",
			Pos:  Position{5, 0},
		},
		Building{
			Name: "small-electric-pole",
			Pos:  Position{1, 0},
		},
		Building{
			Name: "small-electric-pole",
			Pos:  Position{8, 1},
		},
		Building{
			Name:     "pipe",
			Rotation: dirSouth,
			Pos:      Position{0, 1},
		},
		Building{
			Name:     "underground-pipe",
			Rotation: dirWest,
			Pos:      Position{1, 1},
		},
		Building{
			Name:     "underground-pipe",
			Rotation: dirEast,
			Pos:      Position{4, 1},
		},
		Building{
			Name:     "pump",
			Rotation: dirSouth,
			Pos:      Position{0, 2},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{2, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{9, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 1},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{2, 1},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{9, 1},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 2},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{2, 2},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{9, 2},
		},
		Building{
			Name:     "inserter",
			Rotation: dirEast,
			Pos:      Position{3, 0},
		},
		Building{
			Name:     "long-inserter",
			Rotation: dirEast,
			Pos:      Position{3, 2},
		},
		Building{
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
		Building{
			Name:     "electric-mining-drill",
			Rotation: dirEast,
			Pos:      Position{1, 1},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{3, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{3, 1},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{3, 2},
		},
		Building{
			Name:     "electric-mining-drill",
			Rotation: dirWest,
			Pos:      Position{4.5, 1},
		},
		Building{
			Name: "small-electric-pole",
			Pos:  Position{7, 0},
		},
	},
}

var smeltingHeaderBp = Blueprint{
	Position{4, 11},
	[]Building{
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{2, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{3, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 1},
		},
		Building{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 2},
		},
		Building{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 3},
		},
		Building{
			Name:     "underground-belt",
			Rotation: dirEast,
			Pos:      Position{0, 4},
			Ugbt: "input",
		},
		Building{
			Name:     "underground-belt",
			Rotation: dirWest,
			Pos:      Position{2, 4},
			Ugbt: "output",
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{3, 4},
		},
		Building{
			Name:     "belt",
			Rotation: dirWest,
			Pos:      Position{3, 5},
		},
		Building{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 4},
		},
		Building{
			Name:     "splitter",
			Rotation: dirEast,
			Pos:      Position{0, 5.5},
		},
		Building{
			Name:     "belt",
			Rotation: dirNorth,
			Pos:      Position{1, 5},
		},
		Building{
			Name:     "splitter",
			Rotation: dirWest,
			Pos:      Position{2.5, 5.5},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 6},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 7},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 8},
		},
		Building{
			Name:     "belt",
			Rotation: dirSouth,
			Pos:      Position{1, 9},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 10},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{2, 10},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{3, 10},
		},
	},
}

var smeltingBp = Blueprint{
	Position{2, 12},
	[]Building{
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{0, 0},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 0},
		},
		Building{
			Name:     "inserter",
			Rotation: dirNorth,
			Pos:      Position{0, 1},
		},
		Building{
			Name: "furnace",
			Pos:  Position{1, 3},
		},
		Building{
			Name:     "inserter",
			Rotation: dirNorth,
			Pos:      Position{0, 4},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{0, 5},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 5},
		},
		Building{
			Name:     "inserter",
			Rotation: dirSouth,
			Pos:      Position{0, 6},
		},
		Building{
			Name: "furnace",
			Pos:  Position{1, 7.5},
		},
		Building{
			Name:     "inserter",
			Rotation: dirSouth,
			Pos:      Position{0, 9},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{0, 10},
		},
		Building{
			Name:     "belt",
			Rotation: dirEast,
			Pos:      Position{1, 10},
		},
	},
}

var smeltingFooterBp = Blueprint{
	Position{9, 12},
	[]Building {
		Building{
			Name: "belt",
			Rotation: dirSouth,
			Pos: Position{0, 0},
		},
		Building{
			Name: "belt",
			Rotation: dirSouth,
			Pos: Position{0, 1},
		},
		Building{
			Name: "belt",
			Rotation: dirSouth,
			Pos: Position{0, 2},
		},
		Building{
			Name: "belt",
			Rotation: dirSouth,
			Pos: Position{0, 3},
		},
		Building{
			Name: "belt",
			Rotation: dirSouth,
			Pos: Position{0, 4},
		},
		Building{
			Name: "belt",
			Rotation: dirSouth,
			Pos: Position{0, 5},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{0, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{1, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{2, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{3, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{4, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{5, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{6, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{7, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirEast,
			Pos: Position{8, 6},
		},
		Building{
			Name: "belt",
			Rotation: dirNorth,
			Pos: Position{0, 7},
		},
		Building{
			Name: "belt",
			Rotation: dirNorth,
			Pos: Position{0, 8},
		},
		Building{
			Name: "belt",
			Rotation: dirNorth,
			Pos: Position{0, 9},
		},
		Building{
			Name: "belt",
			Rotation: dirNorth,
			Pos: Position{0, 10},
		},
	},
}
