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
			"assembling-machine-%d",
			dirNorth,
			"",
			Position{4, 1},
		},
		Building{
			"small-electric-pole",
			dirNorth,
			"",
			Position{2, 2},
		},
		Building{
			"small-electric-pole",
			dirNorth,
			"",
			Position{6, 2},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{0, 0},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{0, 1},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{0, 2},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 0},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 1},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 2},
		},
		Building{
			"inserter",
			dirEast,
			"",
			Position{2, 0},
		},
		Building{
			"long-handed-inserter",
			dirEast,
			"",
			Position{2, 1},
		},
		Building{
			"inserter",
			dirEast,
			"",
			Position{6, 0},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{7, 0},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{7, 1},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{7, 2},
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
			"assembling-machine-%d",
			dirWest,
			"",
			Position{5, 0},
		},
		Building{
			"small-electric-pole",
			dirNorth,
			"",
			Position{1, 0},
		},
		Building{
			"small-electric-pole",
			dirNorth,
			"",
			Position{8, 1},
		},
		Building{
			"pipe",
			dirSouth,
			"",
			Position{0, 1},
		},
		Building{
			"underground-pipe",
			dirWest,
			"",
			Position{1, 1},
		},
		Building{
			"underground-pipe",
			dirEast,
			"",
			Position{4, 1},
		},
		Building{
			"pump",
			dirSouth,
			"",
			Position{0, 2},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 0},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{2, 0},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{9, 0},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 1},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{2, 1},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{9, 1},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 2},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{2, 2},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{9, 2},
		},
		Building{
			"inserter",
			dirEast,
			"",
			Position{3, 0},
		},
		Building{
			"long-inserter",
			dirEast,
			"",
			Position{3, 2},
		},
		Building{
			"inserter",
			dirEast,
			"",
			Position{8, 0},
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
			"electric-mining-drill",
			dirEast,
			"",
			Position{1, 1},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{3, 0},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{3, 1},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{3, 2},
		},
		Building{
			"electric-mining-drill",
			dirWest,
			"",
			Position{4.5, 1},
		},
		Building{
			"small-electric-pole",
			dirNorth,
			"",
			Position{7, 0},
		},
	},
}

var smeltingHeaderBp = Blueprint{
	Position{4, 11},
	[]Building{
		Building{
			"belt",
			dirEast,
			"",
			Position{1, 0},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{2, 0},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{3, 0},
		},
		Building{
			"belt",
			dirNorth,
			"",
			Position{1, 1},
		},
		Building{
			"belt",
			dirNorth,
			"",
			Position{1, 2},
		},
		Building{
			"belt",
			dirNorth,
			"",
			Position{1, 3},
		},
		Building{
			"underground-belt",
			dirEast,
			"",
			Position{0, 4},
		},
		Building{
			"underground-belt",
			dirWest,
			"",
			Position{2, 4},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{3, 4},
		},
		Building{
			"belt",
			dirWest,
			"",
			Position{3, 5},
		},
		Building{
			"belt",
			dirNorth,
			"",
			Position{1, 4},
		},
		Building{
			"splitter",
			dirEast,
			"",
			Position{0, 5.5},
		},
		Building{
			"belt",
			dirNorth,
			"",
			Position{1, 5},
		},
		Building{
			"splitter",
			dirWest,
			"",
			Position{2.5, 5.5},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 6},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 7},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 8},
		},
		Building{
			"belt",
			dirSouth,
			"",
			Position{1, 9},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{1, 10},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{2, 10},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{3, 10},
		},
	},
}

var smeltingBp = Blueprint {
	Position{2, 12},
	[]Building{
		Building{
			"belt",
			dirEast,
			"",
			Position{0, 0},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{1, 0},
		},
		Building{
			"inserter",
			dirNorth,
			"",
			Position{0, 1},
		},
		Building{
			"furnace",
			dirNorth,
			"",
			Position{1, 3},
		},
		Building{
			"inserter",
			dirNorth,
			"",
			Position{0, 4},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{0, 5},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{1, 5},
		},
		Building{
			"inserter",
			dirSouth,
			"",
			Position{0, 6},
		},
		Building{
			"furnace",
			dirNorth,
			"",
			Position{1, 7.5},
		},
		Building{
			"inserter",
			dirSouth,
			"",
			Position{0, 9},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{0, 10},
		},
		Building{
			"belt",
			dirEast,
			"",
			Position{1, 10},
		},
	},
}
